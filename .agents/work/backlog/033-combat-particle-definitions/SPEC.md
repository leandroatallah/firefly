# SPEC — US-033 — Combat Particle Definitions

**Branch:** `033-combat-particle-definitions`

## Technical Requirements

### File Modification
- **File:** `assets/particles/vfx.json`
- **Action:** Add three new particle type definitions

### Particle Definitions

Add pixel-based particle configurations (no image assets):

```json
{
  "type": "muzzle_flash",
  "pixel": {
    "size": 2,
    "color": "#FFFFFF",
    "lifetime_frames": 3,
    "velocity_range": 1.0
  }
}
```

```json
{
  "type": "bullet_impact",
  "pixel": {
    "size": 1,
    "color": "#FFFFFF",
    "lifetime_frames": 6,
    "velocity_range": 2.0
  }
}
```

```json
{
  "type": "bullet_despawn",
  "pixel": {
    "size": 1,
    "color": "#FFFFFF",
    "lifetime_frames": 8,
    "velocity_range": 0.5
  }
}
```

### Constraints
- **Color palette:** Only `#000000` (black) or `#FFFFFF` (white)
- **Particle system:** Pixel-based only (no `image`, `frame_width`, `frame_height` fields)
- **Backward compatibility:** Existing `jump` and `landing` definitions remain unchanged

## Pre-conditions
- `assets/particles/vfx.json` exists with `jump` and `landing` definitions
- Particle system supports pixel-based particles

## Post-conditions
- `vfx.json` contains 5 particle types: `jump`, `landing`, `muzzle_flash`, `bullet_impact`, `bullet_despawn`
- JSON is valid and parseable
- All combat VFX stories can reference these particle types

## Integration Points
- **Contract:** `vfx.Manager.SpawnPuff(typeKey, x, y, count, randRange)`
- **Consumers:** US-030 (muzzle flash), US-031 (impact), US-032 (despawn)

## Red Phase

### Test File
`assets/particles/vfx_test.go` (new file)

### Failing Test Scenario
```go
func TestCombatParticleDefinitions(t *testing.T) {
    // Load vfx.json
    // Assert "muzzle_flash" type exists
    // Assert "bullet_impact" type exists
    // Assert "bullet_despawn" type exists
    // Assert all use pixel-based config (no image field)
    // Assert all colors are #000000 or #FFFFFF
}
```

**Expected failure:** Particle types not found in JSON.

## Notes
- Foundation story — implement first
- No code changes, only asset configuration
- Pixel particle system must already exist in engine
