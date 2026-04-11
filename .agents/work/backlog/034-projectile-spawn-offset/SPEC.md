# SPEC — US-034 — Projectile Spawn Offset Configuration

**Branch:** `034-projectile-spawn-offset`

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
    muzzleEffectType string,
    spawnOffsetX16 int, // NEW - fp16 units
    spawnOffsetY16 int, // NEW - fp16 units
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
    muzzleEffectType string
    vfxManager       vfx.Manager
    spawnOffsetX16   int // NEW
    spawnOffsetY16   int // NEW
}
```

### Fire Method Update
```go
func (w *ProjectileWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection) {
    // Calculate spawn position with offset
    offsetX16 := w.spawnOffsetX16
    if faceDir == animation.FaceDirectionLeft {
        offsetX16 = -offsetX16 // Negate X offset when facing left
    }
    
    spawnX16 := x16 + offsetX16
    spawnY16 := y16 + w.spawnOffsetY16
    
    // Spawn muzzle flash VFX at offset position
    if w.vfxManager != nil && w.muzzleEffectType != "" {
        x := float64(spawnX16) / 16.0
        y := float64(spawnY16) / 16.0
        w.vfxManager.SpawnPuff(w.muzzleEffectType, x, y, 1, 0.0)
    }
    
    // Spawn projectile at offset position
    vx16, vy16 := w.calculateVelocity(direction, faceDir)
    w.manager.SpawnProjectile(w.projectileType, spawnX16, spawnY16, vx16, vy16, w.owner)
    w.currentCooldown = w.cooldownFrames
}
```

## Pre-conditions
- `ProjectileWeapon` has 6-parameter constructor (from US-030)
- `Fire()` spawns at entity position without offset
- Muzzle flash VFX implemented (US-030)

## Post-conditions
- Constructor accepts 8 parameters (adds offset fields)
- Offset applied to both projectile and VFX spawn positions
- X offset negated when facing left
- Offset (0, 0) maintains current behavior (backward compatible)

## Integration Points
- **Package:** `internal/engine/combat/weapon/`
- **Resolves:** TODO in `internal/engine/physics/skill/skill_shooting.go`
- **Usage:** Game code provides offset values matching sprite alignment

## Red Phase

### Test File
`internal/engine/combat/weapon/weapon_test.go`

### Failing Test Scenario
```go
func TestProjectileWeapon_SpawnOffset_FacingRight(t *testing.T) {
    mockProjMgr := &mockProjectileManager{}
    mockVFXMgr := &mockVFXManager{}
    
    // Offset: 8 pixels right, 4 pixels down (in fp16: 128, 64)
    weapon := NewProjectileWeapon("gun", 10, "bullet", 160, mockProjMgr, "muzzle_flash", 128, 64)
    weapon.SetVFXManager(mockVFXMgr)
    
    weapon.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight)
    
    // Assert projectile spawned at: x16=448 (320+128), y16=544 (480+64)
    // Assert VFX spawned at: x=28.0 (448/16), y=34.0 (544/16)
}

func TestProjectileWeapon_SpawnOffset_FacingLeft(t *testing.T) {
    mockProjMgr := &mockProjectileManager{}
    mockVFXMgr := &mockVFXManager{}
    
    weapon := NewProjectileWeapon("gun", 10, "bullet", 160, mockProjMgr, "muzzle_flash", 128, 64)
    weapon.SetVFXManager(mockVFXMgr)
    
    weapon.Fire(320, 480, animation.FaceDirectionLeft, body.ShootDirectionStraight)
    
    // Assert projectile spawned at: x16=192 (320-128), y16=544 (480+64)
    // Assert VFX spawned at: x=12.0 (192/16), y=34.0 (544/16)
}

func TestProjectileWeapon_ZeroOffset(t *testing.T) {
    mockProjMgr := &mockProjectileManager{}
    
    weapon := NewProjectileWeapon("gun", 10, "bullet", 160, mockProjMgr, "", 0, 0)
    
    weapon.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight)
    
    // Assert projectile spawned at: x16=320, y16=480 (no offset)
}
```

**Expected failure:** Constructor signature mismatch, offset not applied.
