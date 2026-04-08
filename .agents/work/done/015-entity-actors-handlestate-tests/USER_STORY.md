# US-015 — Entity Actors handleState Test Coverage

**Branch:** `015-entity-actors-handlestate-tests`
**Bounded Context:** Entity

## Story

As a developer, I want `handleState` in `internal/engine/entity/actors/character.go` to have comprehensive test coverage, so that state machine regressions are caught automatically.

## Context

`handleState` is the core of the actor state machine. It currently sits at **19.5% coverage** (the overall `entity/actors` package is at 67.1%). The function contains 15+ conditional branches covering transitions between Idle, Walking, Jumping, Falling, Landing, Hurted, Dying, Dead, Exiting, and Ducking states — most are untested.

## Acceptance Criteria

- **AC1:** Every state transition branch in `handleState` is covered by at least one test.
- **AC2:** Tests cover the `StateTransitionHandler` override path (returns `true` skips default logic).
- **AC3:** Tests cover the invulnerability timer decrement and `SetInvulnerability(false)` at zero.
- **AC4:** Tests cover `Dying → Dead` transition when animation finishes.
- **AC5:** Tests cover early-exit for `Exiting`, `Dying`, `Dead` states.
- **AC6:** Tests cover `Health <= 0` forcing `Dying` from any non-dying state.
- **AC7:** Coverage for `internal/engine/entity/actors` package reaches **≥ 60%** (up from 19.5% on `handleState`).
- **AC8:** All tests are deterministic, table-driven where applicable, and use no `time.Sleep`.
