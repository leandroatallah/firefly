# SPEC -- US-035 -- VFX Manager Integration for Projectile Manager

**Branch:** `035-vfx-projectile-manager`
**Bounded Context:** Engine / Combat (`internal/engine/combat/projectile/`)

## Technical Requirements

### 1. Manager struct -- VFX fields (AC1, AC7)

**File:** `internal/engine/combat/projectile/manager.go`

The `Manager` struct holds:

- `vfxManager contractsvfx.Manager` -- optional VFX dependency (nil by default).
- `impactEffect string` -- default effect name for collision impacts (initialized to `"bullet_impact"` in `NewManager`).
- `despawnEffect string` -- default effect name for out-of-bounds despawn (initialized to `"bullet_despawn"` in `NewManager`).

### 2. SetVFXManager setter (AC2)

**File:** `internal/engine/combat/projectile/manager.go`

```go
func (m *Manager) SetVFXManager(v contractsvfx.Manager)
```

Stores the provided VFX manager on the `Manager` struct. This is a setter-based injection to keep VFX optional and maintain backward compatibility.

### 3. Spawn forwards VFX references to projectile (AC3)

**File:** `internal/engine/combat/projectile/manager.go`

`Spawn()` passes three VFX-related values when constructing each `projectile`:

```go
p := &projectile{
    // ... existing fields ...
    vfxManager:    m.vfxManager,
    impactEffect:  m.impactEffect,
    despawnEffect: m.despawnEffect,
}
```

The manager's default effect name strings are forwarded to every spawned projectile.

**Known gap (out of scope):** `ProjectileConfig.ImpactEffect` and `DespawnEffect` fields are not read during `Spawn()`. The manager always uses its hardcoded defaults. A follow-up story will wire config overrides into spawn logic.

### 4. No signature change to SpawnProjectile (AC6)

The public `SpawnProjectile(projectileType string, x16, y16, vx16, vy16 int, owner interface{})` signature remains unchanged. It delegates to `Spawn()` internally with a default `ProjectileConfig`.

### 5. Projectile struct -- VFX fields

**File:** `internal/engine/combat/projectile/projectile.go`

The `projectile` struct holds:

- `vfxManager contractsvfx.Manager`
- `impactEffect string`
- `despawnEffect string`

### 6. VFX trigger points via spawnVFX helper

**File:** `internal/engine/combat/projectile/projectile.go`

A private `spawnVFX(typeKey string)` method handles all VFX emission. It:

1. Guards against nil `vfxManager` and empty `typeKey` (AC4).
2. Reads the projectile's fp16 position via `body.GetPosition16()`.
3. Converts to float64 with `float64(x16) / 16.0` and `float64(y16) / 16.0`.
4. Calls `vfxManager.SpawnPuff(typeKey, x, y, 1, 0.0)`.

Trigger sites:

| Callback   | Effect key used   | Trigger condition                          |
|------------|-------------------|--------------------------------------------|
| `OnTouch`  | `p.impactEffect`  | Collision with a non-owner collidable body |
| `OnBlock`  | `p.impactEffect`  | Collision with a wall/static body          |
| `Update`   | `p.despawnEffect`  | Projectile exits tilemap bounds            |

### 7. ProjectileConfig VFX fields (AC8)

**File:** `internal/engine/combat/projectile/config.go`

```go
type ProjectileConfig struct {
    Width         int
    Height        int
    Damage        int
    ImpactEffect  string `json:"impact_effect,omitempty"`
    DespawnEffect string `json:"despawn_effect,omitempty"`
}
```

These fields exist on the config struct and are JSON-serializable. They support per-projectile-type VFX overrides. As noted in section 3, these fields are not yet consumed during `Spawn()` -- that is deferred to a follow-up story.

### 8. Nil safety / backward compatibility (AC4)

When `SetVFXManager` is never called (VFX manager remains nil):

- `spawnVFX` returns immediately without side effects.
- No panics occur during `OnTouch`, `OnBlock`, or `Update` out-of-bounds removal.
- All projectile lifecycle behavior (spawn, move, collide, despawn) works identically to pre-VFX behavior.

## Pre-conditions

- `Manager` struct and `NewManager(space)` constructor exist in `internal/engine/combat/projectile/manager.go`.
- `projectile` struct exists in `internal/engine/combat/projectile/projectile.go`.
- `vfx.Manager` contract exists at `internal/engine/contracts/vfx/vfx.go` and includes the `SpawnPuff(typeKey string, x float64, y float64, count int, randRange float64)` method.
- `body.BodiesSpace` contract provides `GetTilemapDimensionsProvider()` for bounds checking.

## Post-conditions

- `Manager` has `vfxManager`, `impactEffect`, and `despawnEffect` fields.
- `SetVFXManager` setter stores the injected VFX manager.
- `NewManager` initializes default effect names (`"bullet_impact"`, `"bullet_despawn"`).
- Every spawned projectile receives the VFX manager reference and both effect name strings.
- `spawnVFX` is nil-safe and empty-key-safe.
- `ProjectileConfig` includes `ImpactEffect` and `DespawnEffect` string fields with JSON tags.
- `SpawnProjectile` signature is unchanged.
- All existing tests continue to pass.

## Integration Points

- **Contract:** `internal/engine/contracts/vfx/vfx.go` -- the `Manager` interface, specifically `SpawnPuff`.
- **Contract:** `internal/engine/contracts/body/` -- `BodiesSpace`, `Collidable`, `Movable` interfaces.
- **Injection site:** Game-level code calls `projectileManager.SetVFXManager(vfxMgr)` after constructing both managers.
- **Follow-up:** A future story will wire `ProjectileConfig.ImpactEffect` / `DespawnEffect` into `Spawn()` to allow per-type overrides of the manager defaults.

## Red Phase -- Failing Test Scenarios

### Test file: `internal/engine/combat/projectile/manager_test.go`

All tests are table-driven per constitution standards. A `mockVFXManager` must be added to `mocks_test.go` that records `SpawnPuff` calls.

#### Scenario 1: SetVFXManager stores the manager (AC2, AC5a)

```go
func TestManager_SetVFXManager(t *testing.T) {
    tests := []struct {
        name       string
        setManager bool
    }{
        {name: "stores non-nil VFX manager", setManager: true},
        {name: "manager remains nil when setter not called", setManager: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSpace := &mockBodiesSpace{}
            mgr := NewManager(mockSpace)
            vfxMgr := &mockVFXManager{}

            if tt.setManager {
                mgr.SetVFXManager(vfxMgr)
            }

            if tt.setManager && mgr.vfxManager != vfxMgr {
                t.Error("expected vfxManager to be set")
            }
            if !tt.setManager && mgr.vfxManager != nil {
                t.Error("expected vfxManager to remain nil")
            }
        })
    }
}
```

**Expected failure:** Test accesses `mgr.vfxManager` -- passes only if field and setter exist.

#### Scenario 2: Spawned projectile receives VFX manager and effect strings (AC3, AC5b, AC7)

```go
func TestManager_Spawn_ForwardsVFX(t *testing.T) {
    tests := []struct {
        name              string
        setVFX            bool
        wantImpactEffect  string
        wantDespawnEffect string
    }{
        {
            name:              "projectile receives VFX manager and default effects",
            setVFX:            true,
            wantImpactEffect:  "bullet_impact",
            wantDespawnEffect: "bullet_despawn",
        },
        {
            name:              "projectile receives nil VFX manager when not set",
            setVFX:            false,
            wantImpactEffect:  "bullet_impact",
            wantDespawnEffect: "bullet_despawn",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSpace := &mockBodiesSpace{}
            mgr := NewManager(mockSpace)
            vfxMgr := &mockVFXManager{}

            if tt.setVFX {
                mgr.SetVFXManager(vfxMgr)
            }

            cfg := ProjectileConfig{Width: 2, Height: 1}
            mgr.Spawn(cfg, 100<<4, 50<<4, 5<<4, 0, nil)

            p := mgr.projectiles[0]
            if tt.setVFX && p.vfxManager != vfxMgr {
                t.Error("expected projectile to have VFX manager")
            }
            if !tt.setVFX && p.vfxManager != nil {
                t.Error("expected projectile VFX manager to be nil")
            }
            if p.impactEffect != tt.wantImpactEffect {
                t.Errorf("impactEffect = %q, want %q", p.impactEffect, tt.wantImpactEffect)
            }
            if p.despawnEffect != tt.wantDespawnEffect {
                t.Errorf("despawnEffect = %q, want %q", p.despawnEffect, tt.wantDespawnEffect)
            }
        })
    }
}
```

**Expected failure:** Test accesses `p.vfxManager`, `p.impactEffect`, `p.despawnEffect` -- passes only if projectile struct and spawn logic include these fields.

#### Scenario 3: Nil VFX manager does not panic (AC4, AC5c)

```go
func TestProjectile_NilVFXManager_NoPanic(t *testing.T) {
    tests := []struct {
        name    string
        trigger string // "impact", "block", or "despawn"
    }{
        {name: "OnTouch with nil VFX manager", trigger: "impact"},
        {name: "OnBlock with nil VFX manager", trigger: "block"},
        {name: "Update out-of-bounds with nil VFX manager", trigger: "despawn"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSpace := &mockBodiesSpace{
                tilemapProvider: &mockTilemapDimensionsProvider{
                    width: 100, height: 100,
                },
            }
            mgr := NewManager(mockSpace)
            // Do NOT call SetVFXManager -- vfxManager stays nil
            cfg := ProjectileConfig{Width: 2, Height: 1}
            mgr.Spawn(cfg, 95<<4, 50<<4, 10<<4, 0, nil)

            // Should not panic regardless of trigger
            switch tt.trigger {
            case "impact":
                mgr.projectiles[0].OnTouch(nil)
            case "block":
                mgr.projectiles[0].OnBlock(nil)
            case "despawn":
                mgr.Update() // projectile will go out of bounds
            }
        })
    }
}
```

**Expected failure:** Panics if `spawnVFX` does not guard against nil manager.

#### Scenario 4: VFX manager SpawnPuff is called on impact and despawn (AC3, AC7)

```go
func TestProjectile_VFX_SpawnPuffCalled(t *testing.T) {
    tests := []struct {
        name          string
        trigger       string
        wantTypeKey   string
        wantCallCount int
    }{
        {name: "SpawnPuff on OnTouch", trigger: "impact", wantTypeKey: "bullet_impact", wantCallCount: 1},
        {name: "SpawnPuff on OnBlock", trigger: "block", wantTypeKey: "bullet_impact", wantCallCount: 1},
        {name: "SpawnPuff on out-of-bounds", trigger: "despawn", wantTypeKey: "bullet_despawn", wantCallCount: 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            vfxMgr := &mockVFXManager{}
            mockSpace := &mockBodiesSpace{
                tilemapProvider: &mockTilemapDimensionsProvider{
                    width: 100, height: 100,
                },
            }
            mgr := NewManager(mockSpace)
            mgr.SetVFXManager(vfxMgr)
            cfg := ProjectileConfig{Width: 2, Height: 1}
            mgr.Spawn(cfg, 95<<4, 50<<4, 10<<4, 0, "owner")

            switch tt.trigger {
            case "impact":
                // Create a mock collidable that is not the owner
                other := &mockCollidable{id: "enemy"}
                mgr.projectiles[0].OnTouch(other)
            case "block":
                other := &mockCollidable{id: "wall"}
                mgr.projectiles[0].OnBlock(other)
            case "despawn":
                mgr.Update()
            }

            if vfxMgr.spawnPuffCallCount != tt.wantCallCount {
                t.Errorf("SpawnPuff call count = %d, want %d", vfxMgr.spawnPuffCallCount, tt.wantCallCount)
            }
            if vfxMgr.lastTypeKey != tt.wantTypeKey {
                t.Errorf("SpawnPuff typeKey = %q, want %q", vfxMgr.lastTypeKey, tt.wantTypeKey)
            }
        })
    }
}
```

**Expected failure:** `mockVFXManager.spawnPuffCallCount` is 0 if VFX calls are not wired.

#### Scenario 5: ProjectileConfig VFX fields (AC8)

Covered by existing `config_test.go` -- `TestProjectileConfig_VFXFields` and `TestProjectileConfig_JSONRoundTrip` validate that `ImpactEffect` and `DespawnEffect` fields exist on `ProjectileConfig` and round-trip through JSON correctly.

### Mock additions required in `mocks_test.go`

```go
type mockVFXManager struct {
    spawnPuffCallCount int
    lastTypeKey        string
    lastX, lastY       float64
}

func (m *mockVFXManager) SpawnPuff(typeKey string, x, y float64, count int, randRange float64) {
    m.spawnPuffCallCount++
    m.lastTypeKey = typeKey
    m.lastX = x
    m.lastY = y
}

// Remaining vfx.Manager interface methods as no-ops...
```

### Test file: `internal/engine/combat/projectile/config_test.go`

Already contains table-driven tests for AC8 (`TestProjectileConfig_VFXFields`, `TestProjectileConfig_JSONRoundTrip`). No additional tests needed for config fields.
