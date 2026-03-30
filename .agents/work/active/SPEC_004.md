# SPEC 004 — Tween-Based Dash Deceleration

**Branch:** `004-tween-dash-deceleration`
**Bounded Context:** Physics / Entity
**Package:** `internal/game/entity/actors/states/` (DashState) + `internal/engine/physics/tween/` (InOutSine tween)

## Technical Requirements

### 1. `InOutSineTween` — `internal/engine/physics/tween/`

```go
type InOutSineTween struct { /* unexported */ }

func NewInOutSineTween(from, to float64, durationFrames int) *InOutSineTween
func (t *InOutSineTween) Tick() float64   // advances one frame, returns current value
func (t *InOutSineTween) Done() bool
func (t *InOutSineTween) Reset()
```

- Formula: `value = from + (to-from) * (1 - cos(π * progress)) / 2` where `progress = currentFrame / durationFrames`.
- `Tick()` increments internal frame counter then returns the interpolated value.
- `Done()` returns true when `currentFrame >= durationFrames`.

### 2. `DashState` — `internal/game/entity/actors/states/`

```go
type DashState struct { /* injected deps */ }

func NewDashState(body body.MovableCollidable, space body.BodiesSpace, cfg DashConfig) *DashState
func (s *DashState) Enter()
func (s *DashState) Update() ActorStateEnum
func (s *DashState) Exit()

type DashConfig struct {
    Speed            int     // DashSpeed (x16 units)
    DurationFrames   int     // e.g. 18
    BlockDistance    int     // DashBlockDistance pixels
    Cooldown         int     // frames
    DuckHeight       int
}
```

- `Enter()`:
  - Checks `space.Query(...)` in facing direction within `BlockDistance` → if blocked, abort (return immediately, transition to `Idle`/`Falling`).
  - Zeroes vertical velocity, sets `body.SetFreeze(true)` (gravity suspended).
  - Resizes hitbox via `ResizeFixedBottom` to `DuckHeight`.
  - Initialises `InOutSineTween(DashSpeed, 0, DurationFrames)`.
  - Marks air-dash used if not grounded.
- `Update()`:
  - If tween not done: apply `tween.Tick()` as horizontal velocity in facing direction → return `StateDashing`.
  - If tween done: `body.SetFreeze(false)`, restore hitbox → return `StateFalling` or `StateIdle` based on grounded check.
  - If wall collision detected mid-dash: abandon tween, `body.SetFreeze(false)`, restore hitbox → return `StateFalling`/`StateIdle`.
- `Exit()`: ensures freeze is off, hitbox restored, cooldown timer started.
- Air-dash allowance: tracked as `airDashUsed bool`; reset in `Enter()` when grounded, set when dash starts airborne.
- Cooldown: `Update()` returns `StateDashing` (blocks re-trigger) until cooldown expires; managed by the state itself or the parent state machine.

### State Machine Transitions

```
Any grounded/air state → DashState   (dash input, not on cooldown, not already dashing, air-dash available if airborne)
DashState → Falling                  (tween complete, airborne)
DashState → Idle                     (tween complete, grounded)
DashState → Falling/Idle             (wall collision mid-dash)
```

## Pre-conditions

- Dash input pressed.
- Not already in `DashState`.
- Cooldown elapsed.
- If airborne: air-dash not yet used this jump.
- No solid body within `BlockDistance` in facing direction.

## Post-conditions

- Horizontal velocity follows `InOutSine` curve from `DashSpeed` to `0`.
- Vertical velocity == 0 and gravity suspended for full duration.
- Hitbox is duck-height for full duration.
- On completion: gravity restored, hitbox restored, correct next state entered.

## Integration Points

- `ResizeFixedBottom` from SPEC-001.
- `body.MovableCollidable` and `body.BodiesSpace` from `internal/engine/contracts/body`.
- `body.SetFreeze(bool)` / `body.Freeze()` already on `Movable` contract.
- `body.BodiesSpace.Query(rect)` for block-distance check.

## Red Phase — Failing Test Scenario

File: `internal/engine/physics/tween/inoutsine_test.go`

`TestInOutSineTween`:
- After `durationFrames` ticks, `Done() == true`.
- First tick value is close to `from` (progress near 0).
- Middle tick value is near midpoint.
- Final tick value is close to `to`.

File: `internal/game/entity/actors/states/dash_state_test.go`

`TestDashStateUpdate`:

| case | setup | expected state after Update() |
|---|---|---|
| tween in progress | frame < duration | StateDashing |
| tween complete, grounded | frame == duration, grounded | StateIdle |
| tween complete, airborne | frame == duration, not grounded | StateFalling |
| wall collision | blocked mid-dash | StateFalling or StateIdle |
| dash while already dashing | Enter() called twice | second Enter() is no-op |

Test must fail (types do not exist yet) → implement → test passes.
