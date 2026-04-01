# Contracts

This package defines **interfaces (contracts)** used throughout the engine.

## Why No Tests?

This package contains **only interface definitions** — no implementation logic. Interfaces are validated at compile-time when concrete types implement them. Testing happens at the implementation level (e.g., `../physics/body/`, `../physics/space/`).

## Key Contracts

- `body/`: Physical body interfaces (`Movable`, `Collidable`, `MovableCollidable`, `BodiesSpace`, `OneWayPlatform`).
  - `OneWayPlatform`: Extends `Body` with `IsOneWay()`, `SetPassThrough(actor, frames)`, and `IsPassThrough(actor)` for drop-through support.
- `scene/`: Scene-level contracts.
  - `Freezable`: `FreezeFrame(durationFrames int)` and `IsFrozen() bool` — injectable interface for the hit-stop freeze effect.
