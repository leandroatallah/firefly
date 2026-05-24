# SPEC — 066 Beat-Em-Up Airborne State Transitions

## 1. Scope

Extend `beatemupMovementTransitions` to drive `Jumping`, `Falling`, `Landing` states from the altitude axis (`VAltitude16()`, `Altitude()`). Ground-plane velocity (`vx`,`vy`) drives `Walking`/`Idle` only while grounded. No new contracts.

## 2. Files

- Modify: `internal/kit/actors/beatemup/beatemup_character.go` (function `beatemupMovementTransitions`)
- Modify: `internal/kit/actors/beatemup/beatemup_character_test.go` (add table-driven tests)

No engine, contracts, or game changes.

## 3. Altitude API (engine, read-only) [AC-1, AC-2, AC-3, AC-6]

Available on `*actors.Character` via embedded `MovableBody`:

```
c.Altitude() int        // current altitude (>= 0)
c.VAltitude16() int     // signed altitude velocity (fp16)
                        //   < 0  ascending
                        //   > 0  descending
                        //   == 0 at apex or grounded
c.State() actors.ActorStateEnum
c.IsAnimationFinished() bool
c.Velocity() (vx, vy int)        // ground-plane velocity (fp16)
c.SetNewStateFatal(s actors.ActorStateEnum)
```

Constants from `config.Get().Physics.DownwardGravity` used as ground-plane motion threshold (existing pattern).

## 4. Handler Pseudocode [AC-1..AC-6]

```
func beatemupMovementTransitions(c *actors.Character):
    vx, vy   := c.Velocity()
    threshold := config.Get().Physics.DownwardGravity
    isMovingGround := |vx| > threshold || |vy| > threshold
    state    := c.State()
    vAlt16   := c.VAltitude16()
    altitude := c.Altitude()
    airborne := altitude > 0
    set      := func(s) { c.SetNewStateFatal(s) }

    switch:
      # Landing: lock until animation completes [AC-4]
      case state == Landing:
          if !c.IsAnimationFinished(): return
          if isMovingGround: set(Walking)
          else:               set(Idle)
          return

      # Apex / descent during a Jumping state [AC-6]
      case state == Jumping && vAlt16 >= 0 && airborne:
          set(Falling); return

      # Ascending [AC-1]
      case vAlt16 < 0:
          if state != Jumping: set(Jumping)
          return

      # Descending while still in the air [AC-2]
      case vAlt16 > 0 && airborne:
          if state != Falling: set(Falling)
          return

      # Touchdown after a fall [AC-3]
      case state == Falling && !airborne:
          set(Landing); return

      # Airborne guard: never apply ground transitions [AC-5]
      case airborne:
          return

      # Ground-plane transitions (existing behaviour)
      case state != Walking && isMovingGround: set(Walking)
      case state != Idle && !isMovingGround:   set(Idle)
```

### Ordering rationale (one line)
The `Landing` lock comes first to satisfy AC-4. The `Jumping → Falling` apex check (AC-6) comes before the generic ascend/descend branches so the existing `Jumping` state is preserved across the first descending frame test. The airborne guard (AC-5) prevents ground transitions even when `vAlt16 == 0` mid-flight (impossible at the apex unless captured exactly — defensive).

## 5. Pre/Post Conditions

| Rule | Pre | Post |
|---|---|---|
| R1 [AC-1] | `state != Jumping`, `VAltitude16() < 0` | `state == Jumping` |
| R2 [AC-2] | `state ∉ {Jumping, Falling, Landing}` (or `state == Jumping` and not apex-rule), `VAltitude16() > 0`, `Altitude() > 0` | `state == Falling` |
| R3 [AC-3] | `state == Falling`, `Altitude() == 0` | `state == Landing` |
| R4 [AC-4 a] | `state == Landing`, `!IsAnimationFinished()` | `state == Landing` (unchanged) |
| R5 [AC-4 b] | `state == Landing`, `IsAnimationFinished()`, ground-moving | `state == Walking` |
| R6 [AC-4 c] | `state == Landing`, `IsAnimationFinished()`, ground-still | `state == Idle` |
| R7 [AC-5] | `state ∈ {Jumping, Falling}`, ground-moving, airborne | `state` unchanged (no Walking) |
| R8 [AC-6] | `state == Jumping`, `VAltitude16() >= 0`, `Altitude() > 0` | `state == Falling` |
| R9 (edge) | Jump-frame `Altitude() == 0`, `VAltitude16() < 0` | `state == Jumping` (R1 fires; never lands spuriously) |
| R10 (edge) | `state == Hurted` | handler is allowed to overwrite only via the normal switch; existing call-order guarantees Hurted is set after contributors. Not exercised here. |

## 6. Red Phase — Table-Driven Tests

File: `internal/kit/actors/beatemup/beatemup_character_test.go`

Test function: `TestBeatemupMovementTransitions_AirborneStates(t *testing.T)`.

Use the existing `newTestFixtures()` helper. Construct via `NewBeatEmUpCharacter`, then manipulate via:
- `c.SetVAltitude16(v16 int)` — sets altitude velocity (exposed on Character via embedded body).
- `c.SetAltitude(alt int)` — sets altitude (existing setter on body; confirm during TDD; if absent use `SetAltitude16(fp16.To16(alt))`).
- `c.SetVx16(v16)` / `c.SetVy16(v16)` — set ground velocity for the moving/still distinction.
- `c.SetNewStateFatal(s)` — seed the initial state.

Then invoke the handler indirectly: `c.MovementTransitionHandler(c.Character)`.

### Cases

```
T-A1 [AC-1] ascend from Idle → Jumping
  pre:  state=Idle, vAlt16=-1000, alt=0, vx=0, vy=0
  act:  handler(c)
  post: state==Jumping

T-A2 [AC-1] ascend with ground velocity → still Jumping (no Walking)
  pre:  state=Idle, vAlt16=-1000, alt=10, vx=2000, vy=0
  act:  handler(c)
  post: state==Jumping

T-A3 [AC-6] apex/descent during Jumping → Falling
  pre:  state=Jumping, vAlt16=0, alt=20
  act:  handler(c)
  post: state==Falling

T-A4 [AC-6] descent during Jumping → Falling
  pre:  state=Jumping, vAlt16=500, alt=20
  act:  handler(c)
  post: state==Falling

T-A5 [AC-2] descend airborne from Idle → Falling
  pre:  state=Idle, vAlt16=500, alt=20
  act:  handler(c)
  post: state==Falling

T-A6 [AC-3] touchdown from Falling → Landing
  pre:  state=Falling, vAlt16=0, alt=0
  act:  handler(c)
  post: state==Landing

T-A7 [AC-4] Landing locked while animation playing
  pre:  state=Landing, IsAnimationFinished=false, vx=0, vy=0, alt=0
  act:  handler(c)
  post: state==Landing

T-A8 [AC-4] Landing → Idle when finished and still
  pre:  state=Landing, IsAnimationFinished=true, vx=0, vy=0, alt=0
  act:  handler(c)
  post: state==Idle

T-A9 [AC-4] Landing → Walking when finished and ground-moving
  pre:  state=Landing, IsAnimationFinished=true, vx=2000, vy=0, alt=0
  act:  handler(c)
  post: state==Walking

T-A10 [AC-5] Jumping + ground velocity → stay Jumping
  pre:  state=Jumping, vAlt16=-500, alt=20, vx=2000
  act:  handler(c)
  post: state==Jumping

T-A11 [AC-5] Falling + ground velocity → stay Falling
  pre:  state=Falling, vAlt16=500, alt=20, vx=2000
  act:  handler(c)
  post: state==Falling

T-A12 (edge) same-frame land: jump impulse w/ alt=0 → Jumping (no spurious Landing)
  pre:  state=Idle, vAlt16=-1000, alt=0, vx=0
  act:  handler(c)
  post: state==Jumping

T-A13 (edge) buffered jump on landing frame: Falling + alt=0 → Landing first
  pre:  state=Falling, vAlt16=0, alt=0
  act:  handler(c)
  post: state==Landing   (jump skill will re-fire next frame; out of scope here)

T-A14 (regression) idle on ground, no motion → Idle (unchanged)
  pre:  state=Idle, vAlt16=0, alt=0, vx=0, vy=0
  act:  handler(c)
  post: state==Idle

T-A15 (regression) ground-moving on ground → Walking
  pre:  state=Idle, vAlt16=0, alt=0, vx=2000
  act:  handler(c)
  post: state==Walking
```

### Helper for IsAnimationFinished

`*actors.Character.IsAnimationFinished()` reads from the active state instance. For Landing tests, force the answer by:

```
// Option A (preferred): SetStateInstance(actors.Landing, fakeFinishedState)
// Option B: drive the animation by calling c.Update enough frames.
```

If a fake state is needed, define `fakeLandingState` in the test file implementing `actors.ActorState` with a controllable `IsAnimationFinished()`. TDD specialist chooses the lowest-friction path; both are acceptable.

## 7. Mock / Contract Inventory

- No new contracts.
- No new shared mocks.
- Optional package-local fake actor state in `beatemup_character_test.go` to control `IsAnimationFinished` in Landing tests.

## 8. Acceptance Mapping

- AC-1 → R1, T-A1, T-A2
- AC-2 → R2, T-A5
- AC-3 → R3, T-A6, T-A13
- AC-4 → R4–R6, T-A7, T-A8, T-A9
- AC-5 → R7, T-A2, T-A10, T-A11
- AC-6 → R8, T-A3, T-A4
- AC-7 → All T-A* cases are table-driven.

## 9. Non-Goals

- No new contracts or interfaces.
- No changes to `BeatEmUpMovementModel`, jump skill, or contributors.
- No changes to platformer handler.
- No new sprite assets; `cody.json` already declares `jump`/`fall`/`land`.
