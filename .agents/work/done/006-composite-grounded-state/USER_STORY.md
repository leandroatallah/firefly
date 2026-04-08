# US-006 — Composite Grounded State (Sub-State Machine)

**Branch:** `006-composite-grounded-state`
**Bounded Context:** Entity

## Story

As a game developer,
I want the player's grounded behaviour to be expressed as a composite state with named sub-states,
So that grounded transitions (idle → run → duck → lock) are explicit and independently testable rather than nested inside a single `handleState` switch.

## Acceptance Criteria

- AC1: A `Grounded` composite state owns sub-states: `Idle`, `Walking`, `Ducking`, `AimLock`.
- AC2: Each sub-state implements `OnStart(currentCount int)`, `OnFinish()`, and state transition logic via return value, and is independently constructable.
- AC3: `Grounded` implements the `ActorState` interface and delegates state transitions to the active sub-state.
- AC4: Sub-state transitions are driven by input and do not bypass the parent `Grounded` state's `OnStart()`/`OnFinish()` logic.
- AC5: The existing flat `handleState` switch in `Character` is not broken — the composite state plugs in as a single `ActorStateEnum` value.
- AC6: Each sub-state transition is covered by a unit test asserting the correct sub-state is active after the triggering input.

## Edge Cases

- Transitioning from `Grounded/Ducking` to `Falling`: parent `Grounded.OnFinish()` is called, sub-state `Ducking.OnFinish()` is also called.
- Re-entering `Grounded` from `Falling`: sub-state defaults to `Idle` and `Idle.OnStart()` is called.

## Notes

- Depends on US-003 (Ducking sub-state) and US-004 (dash exits grounded).
- Lives in `internal/game/entity/actors/states/`.
