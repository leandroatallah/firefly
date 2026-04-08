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
