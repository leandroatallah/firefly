# SPEC — US-035 — VFX Manager Integration for Projectile Manager

**Branch:** `035-vfx-projectile-manager`

## Technical Requirements

### Manager Modification
**File:** `internal/engine/combat/projectile/manager.go`

```go
type Manager struct {
    projectiles []*projectile
    space       body.BodiesSpace
    counter     int
    vfxManager  vfx.Manager // NEW
}

func (m *Manager) SetVFXManager(manager vfx.Manager) {
    m.vfxManager = manager
}
```

### Spawn Method Update
**File:** `internal/engine/combat/projectile/manager.go`

In `Spawn()` method, pass VFX manager to projectile:

```go
p := &projectile{
    movable:  movableBody,
    body:     collidableBody,
    space:    m.space,
    speedX16: vx16,
    speedY16: vy16,
    vfxManager: m.vfxManager, // NEW
}
```

### Projectile Struct Update
**File:** `internal/engine/combat/projectile/projectile.go`

```go
type projectile struct {
    movable    body.Movable
    body       body.Collidable
    space      body.BodiesSpace
    speedX16   int
    speedY16   int
    vfxManager vfx.Manager // NEW
}
```

## Pre-conditions
- `Manager` exists without VFX integration
- `projectile` struct exists without VFX manager field
- `vfx.Manager` contract exists in `internal/engine/contracts/vfx/`

## Post-conditions
- `Manager` has `vfxManager` field and setter
- Projectiles receive VFX manager reference during spawn
- Nil VFX manager is safe (no panics)
- Existing functionality unchanged

## Integration Points
- **Contract:** `internal/engine/contracts/vfx/vfx.go` → `Manager` interface
- **Consumers:** US-031 (impact VFX), US-032 (despawn VFX)
- **Injection:** Game code calls `SetVFXManager()` after creating projectile manager

## Red Phase

### Test File
`internal/engine/combat/projectile/manager_test.go`

### Failing Test Scenario
```go
func TestManager_SetVFXManager(t *testing.T) {
    space := &mockBodiesSpace{}
    mgr := NewManager(space)
    vfxMgr := &mockVFXManager{}
    
    mgr.SetVFXManager(vfxMgr)
    
    // Assert manager stores VFX manager
    // Spawn a projectile
    // Assert projectile receives VFX manager reference
}

func TestManager_NilVFXManager(t *testing.T) {
    space := &mockBodiesSpace{}
    mgr := NewManager(space)
    
    // Don't set VFX manager (nil)
    // Spawn projectile
    // Assert no panic
}
```

**Expected failure:** `SetVFXManager` method does not exist.
