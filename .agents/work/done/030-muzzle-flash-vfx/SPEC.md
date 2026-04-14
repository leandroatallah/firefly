# SPEC — US-030 — Muzzle Flash VFX on Weapon Fire

**Branch:** `030-muzzle-flash-vfx`

## Technical Requirements

### Constructor Signature Change
**File:** `internal/engine/combat/weapon/weapon.go`

```go
func NewProjectileWeapon(
    id string,
    cooldownFrames int,
    projectileType string,
    projectileSpeed int,
    manager combat.ProjectileManager,
    muzzleEffectType string, // NEW
) *ProjectileWeapon
```

### Struct Modification
```go
type ProjectileWeapon struct {
    id               string
    cooldownFrames   int
    currentCooldown  int
    projectileType   string
    projectileSpeed  int
    manager          combat.ProjectileManager
    owner            interface{}
    muzzleEffectType string      // NEW
    vfxManager       vfx.Manager  // NEW
}
```

### VFX Manager Setter
```go
func (w *ProjectileWeapon) SetVFXManager(manager vfx.Manager) {
    w.vfxManager = manager
}
```

### Fire Method Update
```go
func (w *ProjectileWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection) {
    // Spawn muzzle flash VFX (NEW)
    if w.vfxManager != nil && w.muzzleEffectType != "" {
        x := float64(x16) / 16.0
        y := float64(y16) / 16.0
        w.vfxManager.SpawnPuff(w.muzzleEffectType, x, y, 1, 0.0)
    }
    
    // Existing projectile spawn logic
    vx16, vy16 := w.calculateVelocity(direction, faceDir)
    w.manager.SpawnProjectile(w.projectileType, x16, y16, vx16, vy16, w.owner)
    w.currentCooldown = w.cooldownFrames
}
```

## Pre-conditions
- `ProjectileWeapon` exists with 5-parameter constructor
- `Fire()` method spawns projectiles without VFX
- `vfx.Manager` contract exists

## Post-conditions
- Constructor accepts 6 parameters (adds `muzzleEffectType`)
- `SetVFXManager()` method available
- `Fire()` spawns VFX before projectile (if configured)
- Empty `muzzleEffectType` or nil `vfxManager` = no VFX (backward compatible)

## Integration Points
- **Contract:** `vfx.Manager.SpawnPuff(typeKey, x, y, count, randRange)`
- **Particle type:** References `muzzle_flash` from US-033
- **Injection:** Game code calls `SetVFXManager()` after weapon creation

## Red Phase

### Test File
`internal/engine/combat/weapon/weapon_test.go`

### Failing Test Scenario
```go
func TestProjectileWeapon_MuzzleFlashVFX(t *testing.T) {
    mockProjMgr := &mockProjectileManager{}
    mockVFXMgr := &mockVFXManager{}
    
    weapon := NewProjectileWeapon("gun", 10, "bullet", 160, mockProjMgr, "muzzle_flash")
    weapon.SetVFXManager(mockVFXMgr)
    
    weapon.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight)
    
    // Assert SpawnPuff called with:
    // - typeKey: "muzzle_flash"
    // - x: 20.0 (320/16)
    // - y: 30.0 (480/16)
    // - count: 1
    // - randRange: 0.0
}

func TestProjectileWeapon_NoVFXWhenManagerNil(t *testing.T) {
    mockProjMgr := &mockProjectileManager{}
    
    weapon := NewProjectileWeapon("gun", 10, "bullet", 160, mockProjMgr, "muzzle_flash")
    // Don't set VFX manager
    
    weapon.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight)
    
    // Assert no panic
}
```

**Expected failure:** Constructor signature mismatch, `SetVFXManager` does not exist.
