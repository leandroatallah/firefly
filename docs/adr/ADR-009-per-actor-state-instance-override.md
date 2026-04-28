# ADR-009 — Per-Actor State Instance Override

## Status
Accepted

## Context
The registry-based state pattern (ADR-002) assigns a `StateConstructor` factory per `ActorStateEnum`. `NewState(actor, enum)` calls the factory fresh on every invocation. This works well for stateless or lightly-parameterized states.

Melee swing state breaks the pattern in two ways:

1. **Complex dependencies** — `melee.State` needs a pre-wired weapon reference, VFX spawner, and owner interface that the narrow `StateConstructor(actor) ActorState` signature cannot carry.
2. **Shared instance across enums** — combo-step animation states (step-0, step-1, …) all delegate to the *same* `melee.State` so its frame counter is consistent across the whole swing. A fresh instance per enum would reset the counter mid-combo.

## Decision
`Character` holds an optional `perActorInstances map[ActorStateEnum]ActorState`. `Character.NewState(enum)` checks this map first; if a pre-built instance is registered for that enum it is returned directly, bypassing the registry factory entirely.

The install helper `melee.InstallState(char, ...)` constructs one `*State` and calls `char.SetStateInstance(enum, st)` for the parent `meleeAttackEnum` and every step-state enum. All enums resolve to the same pointer.

## Consequences
- States that require constructor-time injection or shared identity escape the registry without modifying the registry contract or the `StateConstructor` signature.
- The registry remains the default path; the override is opt-in and localized to the `InstallState` call site.
- `perActorInstances` is keyed by enum, so two actors of the same type each get their own instance — no cross-actor state leakage.
- The override map is populated at wiring time (before the game loop) and never mutated during play, so no concurrency concerns arise.
