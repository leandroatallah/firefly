# Combat System

The `internal/engine/combat` package provides a modular and extensible combat system for the engine. It is designed to handle weapon management, firing logic, and projectile lifecycles in a way that decouples the game logic from the underlying physics and rendering.

## Architecture Overview

The system is built around three core concepts defined in `internal/engine/contracts/combat`:

1.  **Inventory**: Manages a collection of weapons and tracks ammo.
2.  **Weapon**: Defines how a weapon fires, its cooldown, and its identity.
3.  **ProjectileManager**: Handles the spawning and management of projectiles in the game world.

### Component Interaction

1.  An **Entity** (e.g., Player or Enemy) holds an `Inventory`.
2.  The `Inventory` contains one or more `Weapon` objects.
3.  When a `Weapon` is fired, it uses a `ProjectileManager` to spawn a projectile.
4.  The `ProjectileManager` registers the projectile in the `BodiesSpace` for physics and collision handling.
5.  The `Projectile` automatically removes itself when it hits something or goes out of bounds.

## Sub-packages

- **[inventory](./inventory)**: Implementation of the weapon collection and ammo tracking.
- **[weapon](./weapon)**: Common weapon implementations (like `ProjectileWeapon`) and a JSON-based weapon factory.
- **[projectile](./projectile)**: A high-performance projectile manager that handles the lifecycle of active projectiles.

## Core Interfaces

The system is driven by interfaces to ensure modularity:

- `combat.Inventory`: Methods for adding, switching, and updating weapons.
- `combat.Weapon`: Methods for firing and managing cooldowns.
- `combat.ProjectileManager`: A simple interface for spawning projectiles by type and position.
- `combat.Damageable` / `combat.Destructible`: Implemented by anything that can receive damage. `Destructible` adds `IsDestroyed()` for lifecycle queries.
- `combat.Factioned`: Reports an entity's `Faction`. See [Faction System](#faction-system).
- `combat.EnemyShooter`: Encapsulates automatic firing gates (state, range, cooldown) for enemies. Implemented by `weapon.EnemyShooting`.

## Faction System

Factions prevent projectiles from damaging their own side. `faction.go` defines three constants:

- `FactionNeutral` (zero value — safe default, damages everyone)
- `FactionPlayer`
- `FactionEnemy`

The projectile's `applyDamage` resolves a `Damageable` target, reads its `Faction` (via the `Factioned` interface on the body or its owner), and **skips damage only when both sides are non-neutral and equal**. This is why player bullets don't hurt players and enemy bullets don't hurt enemies, but a neutral hazard can hurt anyone.

## Damage & Lifetime

`projectile.ProjectileConfig` carries per-type combat data:

- `Damage` — applied to `Damageable` targets on `OnTouch` / `OnBlock`.
- `LifetimeFrames` — when `> 0`, the projectile despawns automatically once the frame counter reaches zero (`0` = infinite, for backward compatibility).
- `Faction` — used by the faction gate described above.
- `ImpactEffect` — VFX key sprayed at the collision point.
- `DespawnEffect` — VFX key sprayed when `LifetimeFrames` expires.

## VFX Integration

Weapons and projectiles trigger VFX through the `vfx.Manager` contract:

- **Muzzle flash** — `ProjectileWeapon` calls `vfx.Manager.SpawnDirectionalPuff(muzzleEffectType, x, y, faceRight, ...)` on fire. The directional puff anchors to the correct edge so the flash extends outward from the barrel regardless of facing.
- **Impact / Despawn** — the projectile calls `SpawnPuff(impactEffect, ...)` on hit and `SpawnPuff(despawnEffect, ...)` on lifetime expiry.

## Enemy Shooting

`combat/weapon/enemy_shooting.go` implements `EnemyShooter` and wraps a `ProjectileWeapon` with three gates:

1. **State gate** (optional) — require the owner's current actor state to match a configured shoot state.
2. **Range + mode gate** — `ShootModeOnSight` requires a `TargetBody` within `Range()` pixels and flips the owner's face toward the target; `ShootModeAlways` fires regardless of target position.
3. **Cooldown gate** — standard weapon cooldown check.

Enemies typically compose one `EnemyShooting` per weapon in their `Update` loop.
