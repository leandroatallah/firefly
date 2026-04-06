# SPEC — 022-weapon-system

**Branch:** `022-weapon-system`
**Bounded Context:** Engine (`internal/engine/combat/weapon/`)

## Technical Requirements

### New Contracts

**`internal/engine/contracts/combat/weapon.go`**

```go
type Weapon interface {
    ID() string
    Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection)
    CanFire() bool
    Update()
    Cooldown() int
    SetCooldown(frames int)
}
```

**`internal/engine/contracts/combat/projectile_manager.go`**

```go
type ProjectileManager interface {
    SpawnProjectile(projectileType string, x16, y16, vx16, vy16 int, owner interface{})
}
```

### New Package Structure

```
internal/engine/combat/weapon/
├── weapon.go           // ProjectileWeapon implementation
├── weapon_test.go      // Unit tests
├── factory.go          // Factory for JSON → Weapon
└── factory_test.go     // Factory tests
```

### Implementation Details

**`ProjectileWeapon` struct:**
- Fields: `id string`, `cooldownFrames int`, `currentCooldown int`, `projectileType string`, `projectileSpeed int`, `projectileDamage int`, `manager ProjectileManager`
- `Fire()` calculates velocity from direction/facing, delegates spawn to `manager.SpawnProjectile()`
- `Update()` decrements `currentCooldown` if > 0
- `CanFire()` returns `currentCooldown == 0`

**JSON Schema (embedded in factory):**
```json
{
  "id": "basic_blaster",
  "type": "projectile",
  "cooldown_frames": 15,
  "projectile": {
    "type": "bullet",
    "speed": 327680,
    "damage": 1
  }
}
```

**Factory:**
- `NewWeaponFromJSON(data []byte, manager ProjectileManager) (Weapon, error)`
- Validates `type == "projectile"` (only supported type for this story)
- Returns error if `projectile` sub-object missing or invalid

## Pre-conditions

- `internal/engine/contracts/body/shooter.go` exists (current scene-level contract)
- `internal/engine/physics/skill/skill_shooting.go` uses `body.Shooter`

## Post-conditions

- `internal/engine/contracts/combat/` package created with `Weapon` and `ProjectileManager` interfaces
- `internal/engine/combat/weapon/` package created with `ProjectileWeapon` and factory
- No imports of `internal/game/` or scene types in `internal/engine/combat/`
- 80%+ test coverage for weapon package

## Integration Points

- **Decoupling:** `ShootingSkill` will later be refactored to hold a `Weapon` instead of calling `body.Shooter` directly (out of scope for this story)
- **Projectile spawning:** `ProjectileManager` contract allows scene or phase to inject spawn logic without weapon knowing about `BodiesSpace`
- **Velocity calculation:** Reuses diagonal math from `skill_shooting.go` (707/1000 for 45° angles)

## Red Phase Scenario

**Test:** `TestProjectileWeapon_Fire_SpawnsProjectileWithCorrectVelocity`

**Setup:**
- Mock `ProjectileManager` recording calls to `SpawnProjectile()`
- `ProjectileWeapon` with `id="test"`, `cooldownFrames=10`, `projectileSpeed=100`, `projectileType="bullet"`

**Action:**
- Call `Fire(1000, 2000, FaceDirectionRight, ShootDirectionDiagonalUpForward)`

**Expected (failing initially):**
- Mock received call: `SpawnProjectile("bullet", 1000, 2000, 70, -70, nil)` (70 ≈ 100*707/1000)
- `CanFire()` returns `false` after fire
- `Cooldown()` returns `10`

**Test:** `TestWeaponFactory_InvalidType_ReturnsError`

**Setup:**
- JSON with `"type": "melee"`

**Action:**
- Call `NewWeaponFromJSON(data, mockManager)`

**Expected (failing initially):**
- Returns error containing "unsupported weapon type"

## Design Decisions

- **Cooldown in frames:** Matches engine's frame-based timing (no `time.Duration`)
- **Fixed-point positions:** `x16`/`y16` consistent with constitution
- **Projectile type as string:** Allows data-driven projectile definitions (future story)
- **No melee/charge:** Deferred to keep this story minimal
