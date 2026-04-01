# SPEC 003 — Duck State with Reduced Hitbox

**Branch:** `003-player-duck-state`
**Bounded Context:** Entity / Physics
**Package:** `internal/game/entity/actors/states/`

## Technical Requirements

New state `DuckingState` implementing the `ActorState` interface (`OnStart()`, `OnFinish()`, plus `State()`, `GetAnimationCount()`, `IsAnimationFinished()`), registered via `RegisterState`:

```go
type DuckingState struct {
    BaseState
    input InputSource
}

func (s *DuckingState) OnStart(currentCount int)
func (s *DuckingState) OnFinish()
```

- `InputSource` is a minimal interface local to the `actors` package:
  ```go
  type InputSource interface {
      DuckHeld() bool
      HasCeilingClearance() bool
  }
  ```
- `OnStart()`: calls `bodyphysics.ResizeFixedBottom(body.Position(), duckHeight)` → `body.SetSize(...)`, zeroes horizontal velocity via `body.SetVelocity(0, vy)`.
- `OnFinish()`: restores full hitbox via `body.SetSize(fullWidth, fullHeight)`.
- Transition logic (duck released + clearance → `Idle`) is driven by the character's `StateTransitionHandler`, not inside the state itself — consistent with how other states work.
- Jump input is ignored while in `DuckingState`.
- Jump input is ignored while in `DuckingState` (no transition to jump state from here).
- `Ducking` is a grounded-only state; falling while ducking is not handled by this state.

### State Machine Transitions

```
Idle/Walking → DuckingState  (duck input pressed, grounded)
DuckingState → Idle          (duck released AND ceiling clearance, via StateTransitionHandler)
DuckingState → [no jump]     (jump input ignored)
```

## Pre-conditions

- Actor is grounded.
- Duck input is held.
- `ResizeFixedBottom` (US-001) is available.

## Post-conditions

- While in `DuckingState`: `body.Height() == duckHeight`, horizontal velocity == 0.
- On exit: body dimensions restored to full height.
- Jump is not triggered from this state.

## Integration Points

- Depends on `internal/engine/physics/body.ResizeFixedBottom` (SPEC-001).
- Depends on `internal/engine/input.HorizontalAxis` or equivalent (SPEC-002) for clean directional input.
- Plugs into `Grounded` composite state (SPEC-006) as a sub-state.
- Uses `body.MovableCollidable` from `internal/engine/contracts/body`.

## Red Phase — Failing Test Scenario

File: `internal/engine/entity/actors/ducking_state_test.go`

Table-driven test `TestDuckingStateTransitions` (via `StateTransitionHandler` or direct `OnStart`/`OnFinish` assertions):

| case | setup | assertion |
|---|---|---|
| `OnStart()` called | body at full height | body height == duckHeight, vx == 0 |
| `OnFinish()` called | body at duck height | body height == fullHeight |

Additional: `DuckingState` must be registered and constructable via `NewState(actor, Ducking)`.
