# ADR-008 — StateContributor Hook for Extensible State Transitions

## Status
Accepted

## Context
`Character.handleState` owns the per-frame actor state machine: after an incoming `StateTransitionHandler` chance, it either keeps the current state or selects one of a handful of built-in movement transitions (`Walking`, `Jumping`, `Falling`, `Ducking`, `Landing`, `Hurted`, `Idle`).

As the game grew beyond basic platforming (dash, shooting, and their compound shooting-while-jumping/walking/falling states), new game-specific transitions kept needing a home:

- Wedging them into the engine's `handleState` `switch` would couple the engine to game concepts (shooting, dashing) it otherwise knows nothing about.
- Making each a full registered state (per [ADR-002](ADR-002-registry-based-state-pattern.md)) doesn't help — the states already exist; the missing piece is **who decides when to enter them** once a skill goes active.
- Overloading `StateTransitionHandler` (a pre-existing, eager-override callback) would blur its role and make it harder to reason about ordering.

## Decision
Introduce `StateContributor` — a narrow optional interface polled by `Character.handleState` **after** the `StateTransitionHandler` check and **before** the built-in movement transitions:

```go
type StateContributor interface {
    ContributeState(current ActorStateEnum) (ActorStateEnum, bool)
}
```

- `Character.AddStateContributor(sc)` registers contributors. They are polled in insertion order; the first to return `(target, true)` wins and is applied as the new state.
- Contributors are **skipped** during animation-critical states (`Hurted`, `Landing`, `Jumping`) so scripted recovery animations are never cut short.
- Contributors live in the **game layer** and adapt engine-side skills into game-specific states. Example: `internal/game/entity/actors/player/state_contributors.go` defines `dashContributor` (wraps `engineskill.DashSkill` → `StateDashing`) and `shootingContributor` (wraps `engineskill.ShootingSkill` → `IdleShooting` / `WalkingShooting` / `JumpingShooting` / `FallingShooting`).

The `ClimberPlayer` wires its contributors via `WireStateContributors(character, movementChecker)`, decoupling the engine-level `Character` from the game-level compound states it will take on.

## Consequences
- The engine stays agnostic of game-specific skills. Adding a new skill-driven state is a game-layer change: implement a contributor, wire it in `WireStateContributors`.
- Ordering is deterministic and insertion-based. Contributors that can fire simultaneously (e.g., dash + shoot) must be ordered by the caller; in practice, dash takes precedence because it is registered first.
- Animation-critical states are protected centrally; contributors do not need to individually check for `Hurted` / `Landing` / `Jumping`.
- The two hooks (`StateTransitionHandler` and `StateContributor`) have distinct roles: the former is a hard-override callback (e.g., forced transitions from physics events), the latter is a cooperative, skill-driven suggestion mechanism.
- Testing is straightforward: a stub contributor returning a fixed `(state, true)` pair exercises the override path without needing a full skill stack. See `character_test.go` (`TestCharacter_handleState_StateContributorWins` / `Defers`).
