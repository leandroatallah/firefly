# SPEC 003 — Duck State with Reduced Hitbox

**Branch:** `003-player-duck-state`
**Bounded Context:** Entity / Physics
**Package:** `internal/game/entity/actors/states/`

## Technical Requirements

New state `DuckingState` implementing the Actor state interface (`Enter()`, `Update()`, `Exit()`):

```go
type DuckingState struct { /* injected deps */ }

func NewDuckingState(body body.MovableCollidable, input InputSource) *DuckingState
func (s *DuckingState) Enter()
func (s *DuckingState) Update() ActorStateEnum
func (s *DuckingState) Exit()
```

- `InputSource` is a minimal interface (new, local to the states package or `internal/engine/input/`):
  ```go
  type InputSource interface {
      DuckHeld() bool
      HasCeilingClearance() bool  // or injected as a func
  }
  ```
- `Enter()`: calls `ResizeFixedBottom(body.Position(), duckHeight)` → `body.SetSize(...)`, sets `body.SetVelocity(0, vy)` (zeroes horizontal only).
- `Update()`: if duck input released AND ceiling clearance → return `StateIdle`; else return `StateDucking`.
- `Exit()`: restores full hitbox via `ResizeFixedTop` or direct `body.SetSize(fullWidth, fullHeight)`.
- Jump input is ignored while in `DuckingState` (no transition to jump state from here).
- `Ducking` is a grounded-only state; falling while ducking is not handled by this state.

### State Machine Transitions

```
Idle/Walking → DuckingState  (duck input pressed, grounded)
DuckingState → Idle          (duck released AND ceiling clearance)
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

File: `internal/game/entity/actors/states/ducking_state_test.go`

Table-driven test `TestDuckingStateTransitions`:

| case | setup | Update() result |
|---|---|---|
| duck held, no clearance check needed | duckHeld=true | StateDucking |
| duck released, has clearance | duckHeld=false, clearance=true | StateIdle |
| duck released, no clearance (ceiling) | duckHeld=false, clearance=false | StateDucking |

Additional assertions in `TestDuckingStateEnter`:
- After `Enter()`: body height == duckHeight, horizontal velocity == 0.

Additional assertions in `TestDuckingStateExit`:
- After `Exit()`: body height == fullHeight.

Test must fail (type does not exist yet) → implement → test passes.
