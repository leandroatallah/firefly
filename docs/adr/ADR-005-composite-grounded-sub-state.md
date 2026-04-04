# ADR-005 — Composite Grounded Sub-State Machine

## Status
Accepted

## Context
The grounded state encompasses several distinct behaviours: idle, walking, ducking, and aim-lock. Implementing each as a separate top-level actor state would require the parent state machine to manage transitions between them, duplicating the "is grounded" guard in every state and scattering grounded-specific input logic across multiple files. Flat states also make it harder to share the animation frame counter across the grounded family.

## Decision
`GroundedState` is a composite actor state that owns an inner `groundedSubState` (an unexported interface). Sub-states (`idleSubState`, `walkingSubState`, `duckingSubState`, `aimLockSubState`) implement `transitionTo(input)` to return the next `GroundedSubStateEnum`. `GroundedState.Update()` evaluates jump/dash inputs first (exiting to the parent machine), then delegates to the active sub-state's `transitionTo` and swaps sub-states when the key changes. The animation frame counter (`count`) is owned by `GroundedState` and shared across all sub-states via `OnStart`.

## Consequences
- All grounded-family transitions are contained in one file; the parent state machine only sees a single `StateGrounded` enum value.
- Adding a new grounded sub-state requires implementing `groundedSubState` and adding a case to `newSubState` — no changes to the parent machine.
- `ForceSubState` and `ActiveSubState` are exposed for testing without requiring reflection or exported internals.
- The two-level machine adds a layer of indirection; debugging a grounded transition requires inspecting both `GroundedState.Update` and the active sub-state's `transitionTo`.
