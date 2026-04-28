# ADR-002 — Registry-Based State Pattern with `init()` Registration

## Status
Accepted

## Context
Actor states (idle, falling, grounded, dashing, etc.) need to be extensible: game-specific states must be addable without modifying engine code. A classic enum switch would require editing the engine every time a new state is introduced. Dependency injection of a state factory would require wiring every state at construction time, coupling the actor builder to all possible states.

## Decision
States are registered into a package-level map via `RegisterState(name, constructor)`. Each state file calls `RegisterState` inside an `init()` function, which Go guarantees runs before `main`. The registry maps a string name to a unique `ActorStateEnum` integer and a `StateConstructor` factory. `NewState(actor, enum)` looks up the constructor and builds the state on demand.

## Consequences
- New states are added by creating a new file with an `init()` call — no engine changes required.
- The registry is global package-level state (intentional, documented with `//nolint:gochecknoglobals`).
- State enum values are assigned at runtime in registration order, so they must never be persisted or compared across process restarts.
- Tests must register states before exercising state transitions; the registry is shared across tests in the same process, so registration order can affect test isolation.

**Note**: For states that need complex constructor arguments unavailable to a registry factory (e.g., a pre-wired weapon reference shared across multiple enums), see [ADR-009](ADR-009-per-actor-state-instance-override.md) — per-actor state instance override.
