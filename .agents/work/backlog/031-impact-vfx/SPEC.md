# SPEC — US-031 — Impact VFX on Projectile Hit

**Branch:** `031-impact-vfx`

## Technical Requirements

### Struct Modification
**File:** `internal/engine/combat/projectile/projectile.go`

```go
type projectile struct {
    movable          contractsbody.Movable
    body             contractsbody.Collidable
    space            contractsbody.BodiesSpace
    speedX16         int
    speedY16         int
    impactEffectType string      // NEW
    vfxManager       vfx.Manager  // NEW
}
```

### OnTouch Method Update
```go
func (p *projectile) OnTouch(other contractsbody.Collidable) {
    if other != p.body.Owner() {
        // Spawn impact VFX before removal (NEW)
        if p.vfxManager != nil && p.impactEffectType != "" {
            x16, y16 := p.body.GetPosition16()
            x := float64(x16) / 16.0
            y := float64(y16) / 16.0
            p.vfxManager.SpawnPuff(p.impactEffectType, x, y, 3, 1.0)
        }
        
        p.space.QueueForRemoval(p.body)
    }
}
```

### OnBlock Method Update
```go
func (p *projectile) OnBlock(other contractsbody.Collidable) {
    // Spawn impact VFX before removal (NEW)
    if p.vfxManager != nil && p.impactEffectType != "" {
        x16, y16 := p.body.GetPosition16()
        x := float64(x16) / 16.0
        y := float64(y16) / 16.0
        p.vfxManager.SpawnPuff(p.impactEffectType, x, y, 3, 1.0)
    }
    
    p.space.QueueForRemoval(p.body)
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
        movable:          movableBody,
        body:             collidableBody,
        space:            m.space,
        speedX16:         vx16,
        speedY16:         vy16,
        vfxManager:       m.vfxManager,
        impactEffectType: config.ImpactEffect, // NEW
    }
    
    // ... rest of method ...
}
```

## Pre-conditions
- `projectile` struct exists without VFX fields
- `OnTouch()` and `OnBlock()` queue removal without VFX
- US-035 implemented (VFX manager in Manager)
- US-036 implemented (`ImpactEffect` field in config)

## Post-conditions
- Projectile stores `impactEffectType` and `vfxManager`
- Impact VFX spawns before removal in both collision callbacks
- Empty effect type or nil manager = no VFX (backward compatible)
- VFX spawns at projectile position with count=3, randRange=1.0

## Integration Points
- **Contract:** `vfx.Manager.SpawnPuff(typeKey, x, y, count, randRange)`
- **Particle type:** References `bullet_impact` from US-033
- **Config:** Uses `ProjectileConfig.ImpactEffect` from US-036
- **Manager:** Uses `Manager.vfxManager` from US-035

## Red Phase

### Test File
`internal/engine/combat/projectile/projectile_test.go`

### Failing Test Scenario
```go
func TestProjectile_ImpactVFX_OnTouch(t *testing.T) {
    mockSpace := &mockBodiesSpace{}
    mockVFXMgr := &mockVFXManager{}
    mockBody := &mockCollidableBody{}
    mockBody.SetPosition16(320, 480)
    
    p := &projectile{
        body:             mockBody,
        space:            mockSpace,
        impactEffectType: "bullet_impact",
        vfxManager:       mockVFXMgr,
    }
    
    otherBody := &mockCollidableBody{}
    p.OnTouch(otherBody)
    
    // Assert SpawnPuff called with:
    // - typeKey: "bullet_impact"
    // - x: 20.0 (320/16)
    // - y: 30.0 (480/16)
    // - count: 3
    // - randRange: 1.0
    // Assert QueueForRemoval called
}

func TestProjectile_ImpactVFX_OnBlock(t *testing.T) {
    mockSpace := &mockBodiesSpace{}
    mockVFXMgr := &mockVFXManager{}
    mockBody := &mockCollidableBody{}
    mockBody.SetPosition16(320, 480)
    
    p := &projectile{
        body:             mockBody,
        space:            mockSpace,
        impactEffectType: "bullet_impact",
        vfxManager:       mockVFXMgr,
    }
    
    otherBody := &mockCollidableBody{}
    p.OnBlock(otherBody)
    
    // Assert SpawnPuff called
    // Assert QueueForRemoval called
}

func TestProjectile_NoVFXWhenManagerNil(t *testing.T) {
    mockSpace := &mockBodiesSpace{}
    mockBody := &mockCollidableBody{}
    
    p := &projectile{
        body:             mockBody,
        space:            mockSpace,
        impactEffectType: "bullet_impact",
        vfxManager:       nil, // No VFX manager
    }
    
    otherBody := &mockCollidableBody{}
    p.OnTouch(otherBody)
    
    // Assert no panic
}
```

**Expected failure:** Fields do not exist, VFX not spawned.
