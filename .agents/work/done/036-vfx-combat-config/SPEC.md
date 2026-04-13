# SPEC — US-036 — VFX Configuration for Projectiles

**Branch:** `036-vfx-combat-config`

## Technical Requirements

### Struct Modification
**File:** `internal/engine/combat/projectile/config.go`

```go
type ProjectileConfig struct {
    Width         int
    Height        int
    Damage        int
    ImpactEffect  string `json:"impact_effect,omitempty"` // VFX type for collision impacts
    DespawnEffect string `json:"despawn_effect,omitempty"` // VFX type for lifetime expiration
}
```

### Field Semantics
- `ImpactEffect`: Particle type key from `vfx.json` (e.g., `"bullet_impact"`)
- `DespawnEffect`: Particle type key from `vfx.json` (e.g., `"bullet_despawn"`)
- Empty string = no VFX (backward compatible)
- JSON tags use snake_case with `omitempty` so existing JSON configs without these fields continue to unmarshal cleanly

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

### Failing Test Scenarios

```go
package projectile

import (
    "encoding/json"
    "testing"
)

func TestProjectileConfig_VFXFields(t *testing.T) {
    tests := []struct {
        name          string
        cfg           ProjectileConfig
        wantImpact    string
        wantDespawn   string
    }{
        {
            name: "both VFX fields set",
            cfg: ProjectileConfig{
                Width:         2,
                Height:        1,
                Damage:        10,
                ImpactEffect:  "bullet_impact",
                DespawnEffect: "bullet_despawn",
            },
            wantImpact:  "bullet_impact",
            wantDespawn: "bullet_despawn",
        },
        {
            name: "only impact effect set",
            cfg: ProjectileConfig{
                Width:        2,
                Height:       1,
                Damage:       10,
                ImpactEffect: "bullet_impact",
            },
            wantImpact:  "bullet_impact",
            wantDespawn: "",
        },
        {
            name: "no VFX fields (backward compat)",
            cfg: ProjectileConfig{
                Width:  2,
                Height: 1,
                Damage: 10,
            },
            wantImpact:  "",
            wantDespawn: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.cfg.ImpactEffect != tt.wantImpact {
                t.Errorf("ImpactEffect = %q, want %q", tt.cfg.ImpactEffect, tt.wantImpact)
            }
            if tt.cfg.DespawnEffect != tt.wantDespawn {
                t.Errorf("DespawnEffect = %q, want %q", tt.cfg.DespawnEffect, tt.wantDespawn)
            }
        })
    }
}

func TestProjectileConfig_JSONRoundTrip(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    ProjectileConfig
    }{
        {
            name:  "full config with VFX fields",
            input: `{"Width":2,"Height":1,"Damage":10,"impact_effect":"bullet_impact","despawn_effect":"bullet_despawn"}`,
            want: ProjectileConfig{
                Width:         2,
                Height:        1,
                Damage:        10,
                ImpactEffect:  "bullet_impact",
                DespawnEffect: "bullet_despawn",
            },
        },
        {
            name:  "legacy config without VFX fields",
            input: `{"Width":2,"Height":1,"Damage":10}`,
            want: ProjectileConfig{
                Width:  2,
                Height: 1,
                Damage: 10,
            },
        },
        {
            name:  "marshal omits empty VFX fields",
            input: "",
            want: ProjectileConfig{
                Width:  2,
                Height: 1,
                Damage: 10,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.name == "marshal omits empty VFX fields" {
                data, err := json.Marshal(tt.want)
                if err != nil {
                    t.Fatalf("json.Marshal error: %v", err)
                }
                var got ProjectileConfig
                if err := json.Unmarshal(data, &got); err != nil {
                    t.Fatalf("json.Unmarshal error: %v", err)
                }
                if got != tt.want {
                    t.Errorf("got %+v, want %+v", got, tt.want)
                }
                return
            }

            var got ProjectileConfig
            if err := json.Unmarshal([]byte(tt.input), &got); err != nil {
                t.Fatalf("json.Unmarshal error: %v", err)
            }
            if got != tt.want {
                t.Errorf("got %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

**Expected failure:** `ImpactEffect` and `DespawnEffect` fields do not exist on `ProjectileConfig` struct.
