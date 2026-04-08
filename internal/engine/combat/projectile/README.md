# Projectile Management

The `projectile` package provides a high-performance system for managing projectiles in the game world.

## Projectile Manager

The `Manager` struct implements the `combat.ProjectileManager` interface and handles the spawning, updating, and removal of projectiles.

### Lifecycle

1.  **Spawn**: A new `projectile` is created and registered with the `BodiesSpace`.
2.  **Update**:
    - Projectiles move by adding their `speedX16` and `speedY16` to their current position.
    - Physics collisions are resolved using the `BodiesSpace`.
    - Bounds checking removes the projectile if it leaves the tilemap.
3.  **Removal**:
    - Automatically removed if it touches a `Collidable` that isn't its owner.
    - Automatically removed if it hits a blocking wall (`OnBlock`).
    - Automatically removed if it goes out of bounds.

## Projectile Implementation

The internal `projectile` struct manages its own state and implements the `body.Touchable` interface.

- **Collision Resolution**: Uses the engine's physics space to resolve overlaps.
- **Bounds Checking**: Uses the `TilemapDimensionsProvider` from the space to determine if the projectile has left the playable area.

### Usage Example

```go
import (
    "github.com/boilerplate/ebiten-template/internal/engine/combat/projectile"
)

// Create a manager (requires a BodiesSpace)
manager := projectile.NewManager(mySpace)

// Spawn a projectile (usually called by a Weapon)
manager.SpawnProjectile("bullet", x16, y16, vx16, vy16, playerEntity)

// In the main game loop
manager.Update()
manager.Draw(screen)
```

## Internal Details

- **Fixed-Point Arithmetic**: Positions and velocities use `fp16` units (1 pixel = 16 units).
- **Batch Processing**: Removals are queued and processed at the end of the update cycle to avoid concurrent modification issues.
- **Rendering**: Default bullets are rendered as a 2x1 white rectangle.
