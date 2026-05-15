# SPEC — 056-eight-dir-movement-skill

## 1. Scope

Add `EightDirectionalMovementSkill` to `internal/kit/skills/` as a genre-agnostic, 8-direction input → body bridge. Mirrors the architecture of `HorizontalMovementSkill` (`platform_move.go`) but drives both X (`OnMoveLeft`/`OnMoveRight`) and ground-plane Y (`OnMoveUp`/`OnMoveDown`). Altitude (`VAltitude16`) is never touched.

No new contracts. No mock generation required. Factory wiring (`FromConfig` mode dispatch) is **deferred to story 058** — this story only adds the skill type, its constructor, and tests.

## 2. File Layout

| File | Purpose |
|---|---|
| `internal/kit/skills/eight_dir_move.go` | Production code (new). |
| `internal/kit/skills/eight_dir_move_test.go` | Table-driven unit tests (new). |
| `internal/kit/skills/package_surface_test.go` | Append assertion `var _ skill.Skill = (*kitskills.EightDirectionalMovementSkill)(nil)`. |

## 3. Type & Constructor [AC-1, AC-5, AC-6]

```go
package kitskills

import (
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
    "github.com/boilerplate/ebiten-template/internal/engine/input"
    physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
    "github.com/boilerplate/ebiten-template/internal/engine/skill"
    "github.com/hajimehoshi/ebiten/v2"
)

// EightDirectionalMovementSkill drives a body in 8 directions on the X/Y
// ground plane. Genre-agnostic: imports neither beat-em-up, platformer, nor
// top-down packages. Y is ground-plane depth; altitude is never written.
type EightDirectionalMovementSkill struct {
    skill.SkillBase
    activationKey ebiten.Key
}

// NewEightDirectionalMovementSkill creates a ready-state skill.
func NewEightDirectionalMovementSkill() *EightDirectionalMovementSkill
```

Post-construct invariants:
- `s.State() == skill.StateReady`
- `s.ActivationKey() == ebiten.Key(0)` (zero value; movement skill has no single activation key)

The `physicsmovement` import is required only because the engine `ActiveSkill` interface (`internal/engine/skill/skill.go`) takes `*physicsmovement.PlatformMovementModel`. This is the same constraint that already binds `HorizontalMovementSkill`; it is treated as pass-through and only used for `IsInputBlocked()`. The skill does **not** import any genre package.

## 4. Method Signatures

```go
func (s *EightDirectionalMovementSkill) Update(
    b body.MovableCollidable,
    model *physicsmovement.PlatformMovementModel,
)

func (s *EightDirectionalMovementSkill) ActivationKey() ebiten.Key

func (s *EightDirectionalMovementSkill) HandleInput(
    b body.MovableCollidable,
    model *physicsmovement.PlatformMovementModel,
    space body.BodiesSpace,
)
```

Interface satisfaction: `skill.ActiveSkill` (and therefore `skill.Skill`).

## 5. Behavior Pseudocode [AC-2, AC-3, AC-4]

```
Update(b, model):
    SkillBase.Update(b, model)        // no-op base
    return                            // AC-5

ActivationKey():
    return s.activationKey            // zero ebiten.Key

HandleInput(b, model, _space):
    # Guard 1: input-blocked guard (matches HorizontalMovementSkill order).
    if model != nil && model.IsInputBlocked():
        return                        # AC-4

    # Guard 2: immobile guard. Zero both axes (X & Y ground-plane),
    # leave altitude untouched.
    if b.Immobile():
        _, _ = b.Velocity()           # read both for clarity (optional)
        _, _ = b.Acceleration()
        b.SetVelocity(0, 0)
        b.SetAcceleration(0, 0)
        return                        # AC-3 — no OnMove* calls

    cmds := input.CommandsReader()
    speed := b.Speed()

    if cmds.Left:  b.OnMoveLeft(speed)
    if cmds.Right: b.OnMoveRight(speed)
    if cmds.Up:    b.OnMoveUp(speed)
    if cmds.Down:  b.OnMoveDown(speed)
    # AC-2. Diagonals: both calls fire (e.g. Left+Up → OnMoveLeft & OnMoveUp).
    # Conflict resolution (Left+Right, Up+Down) deferred to movement model.
```

**Guard order rationale (resolved):** `IsInputBlocked` first, then `Immobile`. Matches the existing `HorizontalMovementSkill` ordering and is safe for all callers: a blocked input must take precedence regardless of mobility state (e.g. cutscene playing while actor is also immobile).

**No `axis` smoothing / inertia path:** Unlike `HorizontalMovementSkill`, which has an `input.HorizontalAxis` and a `HorizontalInertia` branch, this skill always calls `OnMove*` directly. Beat-em-up / top-down inertia is the movement model's responsibility (see story 057).

## 6. Pre/Post-Conditions (Checkable)

| Case | Pre | Post |
|---|---|---|
| New | — | `s.State() == StateReady` |
| Blocked | `model.IsInputBlocked()==true`, any cmds | No `OnMove*` called; velocity unchanged |
| Immobile | `b.Immobile()==true`, velocity=(5,3), accel=(2,1) | velocity=(0,0); accel=(0,0); no `OnMove*` called; altitude unchanged |
| Left only | `cmds={Left:true}`, mobile, unblocked | `OnMoveLeft(speed)` called exactly once; `OnMoveRight/Up/Down` not called |
| Right only | `cmds={Right:true}` | `OnMoveRight(speed)` only |
| Up only | `cmds={Up:true}` | `OnMoveUp(speed)` only |
| Down only | `cmds={Down:true}` | `OnMoveDown(speed)` only |
| Diagonal Left+Up | `cmds={Left:true, Up:true}` | `OnMoveLeft(speed)` and `OnMoveUp(speed)` both called once; others not called |
| No input | `cmds={}` | No `OnMove*` calls; velocity unchanged |

## 7. Red Phase Test Scenarios [AC-7]

Table-driven test `TestEightDirectionalMovementSkill_HandleInput` in `eight_dir_move_test.go`. Each case stubs `input.CommandsReader` via swap-and-restore.

```
T-1: move_left_only
  pre:  cmds.Left=true; b.Speed()=200; b.Immobile()=false; model unblocked
  act:  s.HandleInput(b, model, nil)
  post: b.OnMoveLeft called once with arg=200; OnMoveRight/Up/Down not called

T-2: move_right_only
  pre:  cmds.Right=true; speed=200
  act:  HandleInput
  post: OnMoveRight(200) called once; others not called

T-3: move_up_only
  pre:  cmds.Up=true; speed=200
  act:  HandleInput
  post: OnMoveUp(200) called once; others not called

T-4: move_down_only
  pre:  cmds.Down=true; speed=200
  act:  HandleInput
  post: OnMoveDown(200) called once; others not called

T-5: diagonal_left_up
  pre:  cmds.Left=true, cmds.Up=true; speed=200
  act:  HandleInput
  post: OnMoveLeft(200) and OnMoveUp(200) each called once; OnMoveRight/Down not called

T-6: immobile_guard
  pre:  b.Immobile()=true; b.SetVelocity(fp16.To16(5), fp16.To16(3)); b.SetAcceleration(fp16.To16(2), fp16.To16(1)); cmds.Left=true
  act:  HandleInput
  post: vx==0, vy==0; accX==0, accY==0; no OnMove* called; b.Altitude16() unchanged

T-7: input_blocked_guard
  pre:  model.IsInputBlocked()==true (via mockPlayerMovementBlocker{blocked:true}); cmds.Left=true; b.SetVelocity(fp16.To16(7), fp16.To16(4))
  act:  HandleInput
  post: velocity unchanged (7,4); no OnMove* called

T-8: no_input
  pre:  cmds = zero-value PlayerCommands{}; mobile; unblocked
  act:  HandleInput
  post: no OnMove* called; velocity unchanged
```

Additional non-table tests:

```
T-9: TestEightDirectionalMovementSkill_New
  post: s != nil; s.State() == skill.StateReady

T-10: TestEightDirectionalMovementSkill_Update_NoOp
  pre:  s := NewEightDirectionalMovementSkill(); pre-state captured
  act:  s.Update(b, model)
  post: state unchanged; no OnMove* called

T-11: TestEightDirectionalMovementSkill_ActivationKey
  post: s.ActivationKey() == ebiten.Key(0)

T-12 (package_surface_test.go addition):
  compile-time: var _ skill.Skill = (*kitskills.EightDirectionalMovementSkill)(nil)
```

## 8. Test Harness Inventory

All harness types already exist in `internal/kit/skills/`:
- `newMockMovableCollidable()` — mock body with `OnMoveLeft/Right/Up/Down` call recording. (Verify it records all four; extend if Up/Down are missing.)
- `mockPlayerMovementBlocker{blocked bool}` — used by `movement.NewPlatformMovementModel`.
- Pattern for stubbing `input.CommandsReader`: save original, override, defer restore. See `coverage_test.go` for examples.

**Action for TDD Specialist:** If the existing `newMockMovableCollidable` does not record `OnMoveUp`/`OnMoveDown` calls, extend it to do so (track `onMoveUpCalls`, `onMoveDownCalls` counters with last-arg). Do **not** introduce a new mock file — keep harness changes local to the test file or shared mocks already in `mocks_test.go`.

## 9. Out of Scope (explicit)

- Schema change `MovementConfig.Mode` — story 058.
- `FromConfig` dispatch on mode — story 058.
- `BeatEmUpMovementModel` integration — story 057.
- Diagonal speed normalization (1/√2) — movement model's responsibility (story 057).
- Altitude handling, jump, dash — not part of this skill.

## 10. Acceptance Criteria → Section Map

- AC-1, AC-5, AC-6 → §3 (type, package, imports, no genre coupling)
- AC-2 → §5 (HandleInput pseudocode), T-1..T-5
- AC-3 → §5 (immobile guard), T-6
- AC-4 → §5 (blocked guard), T-7
- AC-5 → §4, §5 (Update no-op), T-10
- AC-7 → §7 (T-1..T-8 table coverage)
