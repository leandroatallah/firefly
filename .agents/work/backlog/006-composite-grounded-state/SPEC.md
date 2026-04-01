# SPEC 006 — Composite Grounded State (Sub-State Machine)

**Branch:** `006-composite-grounded-state`
**Bounded Context:** Entity
**Package:** `internal/game/entity/actors/states/`

## Technical Requirements

### Sub-state interface (local to package)

```go
type groundedSubState interface {
    Enter()
    Update() groundedSubStateEnum
    Exit()
}

type groundedSubStateEnum int

const (
    SubStateIdle groundedSubStateEnum = iota
    SubStateWalking
    SubStateDucking
    SubStateAimLock
)
```

### Sub-state implementations

Each is a separate type in its own file:

| Type | File |
|---|---|
| `IdleSubState` | `idle_sub_state.go` |
| `WalkingSubState` | `walking_sub_state.go` |
| `DuckingSubState` | wraps `DuckingState` from SPEC-003 |
| `AimLockSubState` | `aim_lock_sub_state.go` |

Each implements `Enter()`, `Update() groundedSubStateEnum`, `Exit()`.

### `GroundedState` — composite state

```go
type GroundedState struct { /* injected deps, active sub-state */ }

func NewGroundedState(deps GroundedDeps) *GroundedState
func (g *GroundedState) Enter()                // sets active sub-state to IdleSubState, calls sub.Enter()
func (g *GroundedState) Update() ActorStateEnum // delegates to active sub-state; handles sub-state transitions
func (g *GroundedState) Exit()                 // calls active sub-state.Exit(), then parent cleanup
```

- `Update()` calls `activeSub.Update()` → if returned sub-state differs from current, calls `current.Exit()`, sets new sub, calls `new.Enter()`.
- If a transition out of `Grounded` is needed (e.g. jump → `Falling`, dash → `Dashing`), `Update()` returns the appropriate `ActorStateEnum` to the parent state machine.
- Plugs into the existing `Character` `handleState` switch as a single `ActorStateEnum` value (e.g. `StateGrounded`).
- Parent `Grounded.Exit()` calls `activeSub.Exit()` before its own cleanup.
- Re-entering `Grounded` from `Falling`: `Enter()` always resets sub-state to `Idle`.

### Transition rules (inside `GroundedState.Update`)

| Condition | Sub-state transition |
|---|---|
| no input, grounded | → `Idle` |
| horizontal input | → `Walking` |
| duck input held | → `Ducking` |
| aim-lock input | → `AimLock` |
| jump input | exit `Grounded` → return `StateFalling`/`StateJumping` |
| dash input | exit `Grounded` → return `StateDashing` |

## Pre-conditions

- Actor is grounded.
- `GroundedState` is registered as `StateGrounded` in the parent state machine.

## Post-conditions

- Active sub-state is always valid (never nil).
- Sub-state `Exit()` is always called before switching sub-states.
- Parent `Grounded.Exit()` always calls active sub-state's `Exit()`.
- Re-entry always starts at `Idle`.

## Integration Points

- `DuckingSubState` wraps `DuckingState` from SPEC-003.
- `DashState` from SPEC-004 is triggered by returning `StateDashing` from `GroundedState.Update()`.
- Existing `Character.handleState` switch: replace grounded case(s) with single `StateGrounded` delegation.
- No changes to `body` contracts required.

## Red Phase — Failing Test Scenario

File: `internal/game/entity/actors/states/grounded_state_test.go`

`TestGroundedSubStateTransitions` (table-driven):

| case | input | initial sub-state | expected sub-state after Update() |
|---|---|---|---|
| no input | none | Idle | Idle |
| horizontal input | right | Idle | Walking |
| duck input | down | Walking | Ducking |
| duck released + clearance | none | Ducking | Idle |
| jump input | jump | Idle | (returns StateFalling to parent) |

`TestGroundedStateExitCallsSubExit`:
- Enter `GroundedState`, transition to `Ducking` sub-state, call `GroundedState.Exit()` → assert `DuckingSubState.Exit()` was called.

`TestGroundedStateReEntryResetsToIdle`:
- Enter → transition to `Walking` → Exit → Enter again → assert active sub-state is `Idle`.

Test must fail (types do not exist yet) → implement → test passes.
