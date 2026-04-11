# SPEC — US-036 — VFX Configuration for Projectiles

**Branch:** `036-vfx-projectile-config`

## Technical Requirements

### Struct Modification
**File:** `internal/engine/combat/projectile/config.go`

```go
type ProjectileConfig struct {
    Width         int
    Height        int
    Damage        int
    ImpactEffect  string // VFX type for collision impacts
    DespawnEffect string // VFX type for lifetime expiration
}
```

### Field Semantics
- `ImpactEffect`: Particle type key from `vfx.json` (e.g., `"bullet_impact"`)
- `DespawnEffect`: Particle type key from `vfx.json` (e.g., `"bullet_despawn"`)
- Empty string = no VFX (backward compatible)

## Pre-conditions
- `ProjectileConfig` exists with `Width`, `Height`, `Damage` fields
- No existing VFX-related fields

## Post-conditions
- `ProjectileConfig` has two new string fields
- JSON marshaling includes new fields
- Zero values (empty strings) maintain backward compatibility
- Existing code compiles without changes

## Integration Points
- **Package:** `internal/engine/combat/projectile`
- **Consumers:** US-031 (impact VFX), US-032 (despawn VFX)
- **Contract:** Fields reference `vfx.Manager.SpawnPuff(typeKey, ...)`

## Red Phase

### Test File
`internal/engine/combat/projectile/config_test.go`

### Failing Test Scenario
```go
func TestProjectileConfig_VFXFields(t *testing.T) {
    cfg := ProjectileConfig{
        Width:         2,
        Height:        1,
        Damage:        10,
        ImpactEffect:  "bullet_impact",
        DespawnEffect: "bullet_despawn",
    }
    
    // Assert fields are accessible
    assert.Equal(t, "bullet_impact", cfg.ImpactEffect)
    assert.Equal(t, "bullet_despawn", cfg.DespawnEffect)
}

func TestProjectileConfig_JSONMarshaling(t *testing.T) {
    cfg := ProjectileConfig{
        Width:         2,
        Height:        1,
        Damage:        10,
        ImpactEffect:  "bullet_impact",
        DespawnEffect: "bullet_despawn",
    }
    
    // Marshal to JSON
    // Unmarshal back
    // Assert all fields match
}
```

**Expected failure:** Fields do not exist on struct.
