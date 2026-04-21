# Combat Weapon

The `weapon` package provides concrete implementations of the `combat.Weapon` interface. It handles firing logic, cooldown management, and projectile velocity calculation based on the owner's facing and aiming direction.

## Projectile Weapon

The primary implementation is the `ProjectileWeapon`, which spawns a projectile using a `ProjectileManager`.

### Firing Directions

The `Fire` method takes a `ShootDirection` and calculates the appropriate velocity (including diagonal scaling):

- `ShootDirectionStraight`: Fires forward.
- `ShootDirectionUp / Down`: Fires vertically.
- `ShootDirectionDiagonalUpForward / DiagonalDownForward`: Fires at a 45-degree angle.

Velocity calculations use a 0.707 scaling factor for diagonal speeds to maintain consistent magnitude.

### Cooldown Management

Each weapon has a `cooldownFrames` setting. After firing, the weapon becomes "unready" until its `Update()` method has been called enough times to reduce the cooldown to zero.

### Per-State Spawn Offset

`ProjectileWeapon.SetStateSpawnOffsets(map[int][2]int)` registers spawn offsets keyed by the owner's current actor state (fp16 units). When `Fire(..., state)` is called, the weapon looks up the state in the map and overrides the default `spawnOffsetX16/Y16`. Useful when the barrel origin shifts across animation states (e.g., crouching, jumping). Pass a nil or empty map to clear overrides.

### Muzzle Flash

When a `vfx.Manager` is wired via `SetVFXManager`, `Fire` calls `SpawnDirectionalPuff(muzzleEffectType, x, y, faceRight, ...)` so the flash sprite extends outward from the barrel regardless of facing.

### Enemy Shooting

`enemy_shooting.go` defines `EnemyShooting`, which implements `combat.EnemyShooter`. It wraps a `ProjectileWeapon` and applies three gates per frame — state, range + `ShootMode`, and cooldown — before firing. See the [combat overview](../README.md#enemy-shooting) for the gate chain.

## Weapon Factory

A JSON factory is available to create weapons from configuration data:

```json
{
    "id": "laser_blaster",
    "type": "projectile",
    "cooldown_frames": 10,
    "projectile": {
        "type": "bullet",
        "speed": 327680,
        "damage": 1
    }
}
```

### Usage Example

```go
import (
    "github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
)

// data is a []byte with the JSON content
w, err := weapon.NewWeaponFromJSON(data, myProjectileManager)
if err != nil {
    log.Fatal(err)
}

// Add to inventory or use directly
w.Fire(x16, y16, animation.FaceDirectionRight, body.ShootDirectionStraight)
```
