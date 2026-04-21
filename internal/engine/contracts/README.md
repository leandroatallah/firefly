# Contracts

This package defines **interfaces (contracts)** used throughout the engine.

## Why No Tests?

This package contains **only interface definitions** — no implementation logic. Interfaces are validated at compile-time when concrete types implement them. Testing happens at the implementation level (e.g., `../physics/body/`, `../physics/space/`, `../combat/projectile/`).

## Key Contracts

- `animation/`: Animation abstractions (`FacingDirectionEnum`, frame queries).
- `body/`: Physical body interfaces.
  - `Movable`, `Collidable`, `MovableCollidable`, `BodiesSpace` — core physics.
  - `OneWayPlatform` — extends `Body` with `IsOneWay()`, `SetPassThrough(actor, frames)`, `IsPassThrough(actor)` for drop-through support.
  - `Passthrough` — short-lived "ignore me" flag used by projectiles during spawn to avoid self-collision.
  - `StateTransitionHandler` — callback used by skills (e.g., `ShootingSkill`) to request a state change on the owning body.
  - `ShootDirection` — enum for straight/up/down/diagonal fire directions.
- `combat/`: Combat interfaces.
  - `Weapon`, `Inventory`, `ProjectileManager` — firing pipeline.
  - `Damageable` (`TakeDamage(amount)`) and `Destructible` (adds `IsDestroyed()`).
  - `Factioned` + `Faction` constants — side identification for damage gating.
  - `EnemyShooter`, `EnemyShooterFactory`, `ShootMode`, `TargetBody` — automatic enemy firing.
- `config/`: Engine-side configuration structures.
- `context/`: App/scene context hand-off interfaces.
- `navigation/`: Scene-navigation hooks (`NavigateTo`, `NavigateBack`).
- `projectile/`: `Manager` interface (spawn/update/draw/clear) — kept separate from `combat` to avoid a circular import between `combat` and `projectile`.
- `scene/`: Scene-level contracts.
  - `Freezable` — `FreezeFrame(durationFrames int)` and `IsFrozen() bool`, the injectable interface for the hit-stop freeze effect.
- `sequences/`: `Sequence` and `Command` interfaces used by the scripted event player.
- `tilemaplayer/`: Tilemap layer access for scene rendering and collision.
- `vfx/`: `Manager` interface — particles, screen shake, floating text, directional puffs (muzzle flash), screen flash.
