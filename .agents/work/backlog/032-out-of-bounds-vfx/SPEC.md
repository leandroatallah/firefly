# SPEC — US-032 — Projectile Lifetime and Despawn VFX

**Branch:** `032-projectile-lifetime-despawn`

## Technical Requirements

### Config Modification
**File:** `internal/engine/combat/projectile/config.go`

```go
type ProjectileConfig struct {
    Width          int
    Height         int
    Damage         int
    ImpactEffect   string
    DespawnEffect  string
    LifetimeFrames int // NEW - 0 = infinite
}
```

### Struct Modification
**File:** `internal/engine/combat/projectile/projectile.go`

```go
type projectile struct {
    movable           contractsbody.Movable
    body              contractsbody.Collidable
    space             contractsbody.BodiesSpace
    speedX16          int
    speedY16          int
    impactEffectType  string
    despawnEffectType string      // NEW
    vfxManager        vfx.Manager
    lifetimeFrames    int         // NEW
    currentLifetime   int         // NEW
}
```

### Update Method Modification
```go
func (p *projectile) Update() {
    // Decrement lifetime (NEW)
    if p.lifetimeFrames > 0 {
        p.currentLifetime--
        if p.currentLifetime <= 0 {
            // Spawn despawn VFX
            if p.vfxManager != nil && p.despawnEffectType != "" {
                x16, y16 := p.body.GetPosition16()
                x := float64(x16) / 16.0
                y := float64(y16) / 16.0
                p.vfxManager.SpawnPuff(p.despawnEffectType, x, y, 5, 1.5)
            }
            p.space.QueueForRemoval(p.body)
            return
        }
    }
    
    // Existing movement logic
    x, y := p.body.GetPosition16()
    x += p.speedX16
    y += p.speedY16
    p.body.SetPosition16(x, y)

    p.space.ResolveCollisions(p.body)

    // Existing bounds check (no VFX)
    provider := p.space.GetTilemapDimensionsProvider()
    if provider == nil {
        return
    }
    w := provider.GetTilemapWidth()
    h := provider.GetTilemapHeight()

    if x < 0 || y < 0 || x > w<<4 || y > h<<4 {
        p.space.QueueForRemoval(p.body)
    }
}
```

### Manager Spawn Update
**File:** `internal/engine/combat/projectile/manager.go`

```go
func (m *Manager) Spawn(cfg interface{}, x16, y16, vx16, vy16 int, owner interface{}) {
    config, ok := cfg.(ProjectileConfig)
    if !ok {
        return
    }
    
    // ... existing body creation ...
    
    p := &projectile{
        movable:           movableBody,
        body:              collidableBody,
        space:             m.space,
        speedX16:          vx16,
        speedY16:          vy16,
        vfxManager:        m.vfxManager,
        impactEffectType:  config.ImpactEffect,
        despawnEffectType: config.DespawnEffect, // NEW
        lifetimeFrames:    config.LifetimeFrames, // NEW
        currentLifetime:   config.LifetimeFrames, // NEW
    }
    
    // ... rest of method ...
}
```

## Pre-conditions
- `projectile` has `Update()` without lifetime logic
- US-035 implemented (VFX manager)
- US-036 implemented (`DespawnEffect` field)
- US-031 implemented (impact VFX)

## Post-conditions
- `ProjectileConfig` has `LifetimeFrames` field
- Projectile decrements lifetime each frame
- Despawn VFX spawns at projectile position when lifetime expires
- Lifetime 0 = infinite (backward compatible)
- Out-of-bounds check remains (no VFX, safety fallback)

## Integration Points
- **Contract:** `vfx.Manager.SpawnPuff(typeKey, x, y, count, randRange)`
- **Particle type:** References `bullet_despawn` from US-033
- **Config:** Uses `ProjectileConfig.DespawnEffect` and `LifetimeFrames`

## Red Phase

### Test File
`internal/engine/combat/projectile/projectile_test.go`

### Failing Test Scenario
```go
func TestProjectile_LifetimeDespawn(t *testing.T) {
    mockSpace := &mockBodiesSpace{}
    mockVFXMgr := &mockVFXManager{}
    mockBody := &mockCollidableBody{}
    mockBody.SetPosition16(320, 480)
    
    p := &projectile{
        body:              mockBody,
        space:             mockSpace,
        despawnEffectType: "bullet_despawn",
        vfxManager:        mockVFXMgr,
        lifetimeFrames:    3,
        currentLifetime:   3,
    }
    
    // Frame 1
    p.Update()
    // Assert currentLifetime = 2, not removed
    
    // Frame 2
    p.Update()
    // Assert currentLifetime = 1, not removed
    
    // Frame 3
    p.Update()
    // Assert currentLifetime = 0
    // Assert SpawnPuff called with:
    // - typeKey: "bullet_despawn"
    // - x: 20.0, y: 30.0
    // - count: 5
    // - randRange: 1.5
    // Assert QueueForRemoval called
}

func TestProjectile_InfiniteLifetime(t *testing.T) {
    mockSpace := &mockBodiesSpace{}
    mockBody := &mockCollidableBody{}
    
    p := &projectile{
        body:            mockBody,
        space:           mockSpace,
        lifetimeFrames:  0, // Infinite
        currentLifetime: 0,
    }
    
    // Update many times
    for i := 0; i < 100; i++ {
        p.Update()
    }
    
    // Assert never queued for removal
}

func TestProjectile_OutOfBoundsNoVFX(t *testing.T) {
    mockSpace := &mockBodiesSpace{}
    mockVFXMgr := &mockVFXManager{}
    mockBody := &mockCollidableBody{}
    mockBody.SetPosition16(-100, 0) // Out of bounds
    
    p := &projectile{
        body:              mockBody,
        space:             mockSpace,
        despawnEffectType: "bullet_despawn",
        vfxManager:        mockVFXMgr,
        lifetimeFrames:    100,
        currentLifetime:   100,
    }
    
    p.Update()
    
    // Assert QueueForRemoval called
    // Assert SpawnPuff NOT called (out-of-bounds = no VFX)
}
```

**Expected failure:** Lifetime fields do not exist, logic not implemented.
