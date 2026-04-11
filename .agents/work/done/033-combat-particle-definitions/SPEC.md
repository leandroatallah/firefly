# SPEC — US-033 — Combat Particle Definitions

**Branch:** `033-combat-particle-definitions`

## Overview

Foundation story. Adds three combat particle types (`muzzle_flash`, `bullet_impact`, `bullet_despawn`) to `assets/particles/vfx.json`. Because the current loader supports only image-based particles, this SPEC also extends the schema and `vfx.Manager` to consume **pixel-based** particle definitions from JSON.

Velocity/spread (`randRange`) and particle `count` remain caller-controlled via `SpawnPuff` — they are NOT part of the JSON definition (kept that way to match existing US-030/031/032 SPECs which pass them inline).

## Technical Requirements

### 1. Schema extension

**File:** `internal/engine/data/schemas/json.go`

Add a `Pixel` field to `ParticleData` and a new `PixelParticleData` struct.

```go
// ParticleData defines the configuration for a particle effect.
// A particle entry is either image-based (Image set) or pixel-based (Pixel set).
type ParticleData struct {
    Image       string             `json:"image,omitempty"`
    FrameWidth  int                `json:"frame_width,omitempty"`
    FrameHeight int                `json:"frame_height,omitempty"`
    FrameRate   int                `json:"frame_rate,omitempty"`
    Scale       float64            `json:"scale,omitempty"`
    Pixel       *PixelParticleData `json:"pixel,omitempty"`
}

// PixelParticleData defines a pixel-based (no image asset) particle.
// Color is a hex string limited to "#000000" or "#FFFFFF" (1-bit palette).
type PixelParticleData struct {
    Size           int    `json:"size"`
    Color          string `json:"color"`
    LifetimeFrames int    `json:"lifetime_frames"`
}
```

`omitempty` on the legacy fields keeps existing `jump`/`landing` JSON unchanged (they keep their current keys, no rewrites required).

### 2. `particles.Config` extension

**File:** `internal/engine/render/particles/particle.go`

Add optional pixel-mode fields. Image-based path keeps current behavior.

```go
type Config struct {
    Image       *ebiten.Image
    FrameWidth  int
    FrameHeight int
    FrameCount  int
    FrameRate   int

    // Pixel mode (when set, SpawnPuff uses Lifetime instead of FrameCount*FrameRate
    // and applies Color to each spawned particle).
    Lifetime int
    Color    color.Color // nil → no override (white)
}
```

No changes to `Particle.Update`/`Draw`.

### 3. Loader branch in `vfx.NewManager`

**File:** `internal/engine/render/particles/vfx/vfx.go`

In the `for _, vfx := range vfxList` loop, branch on `vfx.Pixel != nil`:

```go
for _, vfx := range vfxList {
    if vfx.Pixel != nil {
        configs[vfx.Type] = createConfigFromPixel(vfx.Pixel)
        continue
    }
    // ... existing image-based path unchanged ...
}
```

Add `createConfigFromPixel`:

```go
func createConfigFromPixel(pd *schemas.PixelParticleData) *particles.Config {
    size := pd.Size
    if size <= 0 {
        size = 1
    }
    img := ebiten.NewImage(size, size)
    img.Fill(color.White)

    c, err := parseHexColor(pd.Color)
    if err != nil {
        log.Printf("invalid pixel particle color %q: %v", pd.Color, err)
        c = color.White
    }

    return &particles.Config{
        Image:       img,
        FrameWidth:  size,
        FrameHeight: size,
        FrameCount:  1,
        FrameRate:   1,
        Lifetime:    pd.LifetimeFrames,
        Color:       c,
    }
}

// parseHexColor accepts only "#000000" or "#FFFFFF" (1-bit palette enforcement).
func parseHexColor(hex string) (color.Color, error) {
    switch hex {
    case "#FFFFFF":
        return color.White, nil
    case "#000000":
        return color.Black, nil
    default:
        return nil, fmt.Errorf("color must be #000000 or #FFFFFF, got %q", hex)
    }
}
```

`fmt` import added.

### 4. `SpawnPuff` honors pixel-config fields

**File:** `internal/engine/render/particles/vfx/vfx.go`

Update `SpawnPuff` to use `config.Lifetime` when set, and apply `config.Color`:

```go
func (m *Manager) SpawnPuff(typeKey string, x, y float64, count int, randRange float64) {
    config, ok := m.configs[typeKey]
    if !ok {
        return
    }

    duration := config.FrameCount * config.FrameRate
    if config.Lifetime > 0 {
        duration = config.Lifetime
    }

    for i := 0; i < count; i++ {
        p := &particles.Particle{
            X:           x,
            Y:           y,
            VelX:        (rand.Float64() - 0.5) * randRange,
            VelY:        (rand.Float64() - 0.5) * randRange,
            Duration:    duration,
            MaxDuration: duration,
            Scale:       1.0,
            Config:      config,
        }
        if config.Color != nil {
            p.ColorScale.ScaleWithColor(config.Color)
        }
        m.system.Add(p)
    }
}
```

### 5. `vfx.json` additions

**File:** `assets/particles/vfx.json`

Append three new entries (existing `jump` and `landing` entries unchanged):

```json
[
  {
    "type": "jump",
    "image": "assets/images/jump-particles-24.png",
    "frame_width": 24,
    "frame_height": 24,
    "frame_rate": 5,
    "scale": 1.0
  },
  {
    "type": "landing",
    "image": "assets/images/land-particles-24.png",
    "frame_width": 24,
    "frame_height": 24,
    "frame_rate": 5,
    "scale": 1.0
  },
  {
    "type": "muzzle_flash",
    "pixel": {
      "size": 2,
      "color": "#FFFFFF",
      "lifetime_frames": 3
    }
  },
  {
    "type": "bullet_impact",
    "pixel": {
      "size": 1,
      "color": "#FFFFFF",
      "lifetime_frames": 6
    }
  },
  {
    "type": "bullet_despawn",
    "pixel": {
      "size": 1,
      "color": "#FFFFFF",
      "lifetime_frames": 8
    }
  }
]
```

## Constraints

- **Color palette:** loader rejects any hex other than `#000000` / `#FFFFFF`.
- **Particle system:** new entries must use `pixel`, never `image`.
- **Backward compatibility:** existing `jump` / `landing` definitions and `SpawnJumpPuff` / `SpawnLandingPuff` callers behave identically (image-based path untouched).
- **Pixel size:** `size <= 0` clamps to `1`.
- **Lifetime:** `lifetime_frames <= 0` falls through to `FrameCount*FrameRate` (=`1`) → 1-tick particle. Combat entries set lifetimes ≥ 3.

## Pre-conditions

- `assets/particles/vfx.json` exists with `jump` and `landing` entries.
- `vfx.NewManager` reads vfx.json and only handles image-based entries.
- `ParticleData` schema exposes only image-based fields.
- Pixel image rendering already works in `particles.Particle.Draw` (used by `SpawnFallingRocks`/`SpawnDeathExplosion`).

## Post-conditions

- `vfx.json` contains 5 particle types: `jump`, `landing`, `muzzle_flash`, `bullet_impact`, `bullet_despawn`.
- `vfx.NewManager` populates `configs["muzzle_flash"]`, `configs["bullet_impact"]`, `configs["bullet_despawn"]` from pixel definitions without loading any image asset.
- `SpawnPuff("muzzle_flash", ...)` spawns a 2×2 white particle with 3-tick lifetime.
- `SpawnPuff("bullet_impact", ...)` spawns 1×1 white particles with 6-tick lifetime.
- `SpawnPuff("bullet_despawn", ...)` spawns 1×1 white particles with 8-tick lifetime.
- US-030, US-031, US-032 can `SpawnPuff` these type keys with caller-supplied `count` and `randRange`.

## Integration Points

- **Schema:** `schemas.ParticleData` (extended), `schemas.PixelParticleData` (new).
- **Engine:** `particles.Config` (extended), `vfx.Manager.NewManager`, `vfx.Manager.SpawnPuff`.
- **Contract:** `vfx.Manager.SpawnPuff(typeKey, x, y, count, randRange)` — signature unchanged.
- **Consumers:**
  - US-030 muzzle flash → `SpawnPuff("muzzle_flash", x, y, 1, 0.0)`
  - US-031 impact → `SpawnPuff("bullet_impact", x, y, 3, 1.0)`
  - US-032 despawn → `SpawnPuff("bullet_despawn", x, y, 5, 1.5)`

## Red Phase

### Test files

1. **Schema test** — `internal/engine/data/schemas/json_test.go` (extend)
2. **Loader/SpawnPuff test** — `internal/engine/render/particles/vfx/vfx_combat_test.go` (new)

The original SPEC's `assets/particles/vfx_test.go` location is invalid — `assets/` has no Go package. Tests must live in the engine packages.

### Failing test scenarios

**`internal/engine/data/schemas/json_test.go`** — extend with:

```go
func TestParticleData_PixelMode(t *testing.T) {
    raw := []byte(`{
        "type": "muzzle_flash",
        "pixel": {"size": 2, "color": "#FFFFFF", "lifetime_frames": 3}
    }`)

    var pd ParticleData
    if err := json.Unmarshal(raw, &pd); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }
    if pd.Pixel == nil {
        t.Fatal("Pixel field nil; expected populated PixelParticleData")
    }
    if pd.Pixel.Size != 2 || pd.Pixel.Color != "#FFFFFF" || pd.Pixel.LifetimeFrames != 3 {
        t.Errorf("pixel fields wrong: %+v", pd.Pixel)
    }
}
```

**`internal/engine/render/particles/vfx/vfx_combat_test.go`** — new file:

```go
func TestVFXManager_LoadsCombatPixelTypes(t *testing.T) {
    fsys := os.DirFS("../../../../../assets/particles")
    m := NewManager(fsys, "vfx.json")

    for _, typeKey := range []string{"muzzle_flash", "bullet_impact", "bullet_despawn"} {
        if _, ok := m.configs[typeKey]; !ok {
            t.Errorf("config %q not loaded from vfx.json", typeKey)
        }
    }
}

func TestVFXManager_PixelConfigLifetimeAndColor(t *testing.T) {
    // Build a manager from an in-memory fs containing only the pixel entries.
    // Assert configs["muzzle_flash"].Lifetime == 3
    // Assert configs["bullet_impact"].Lifetime == 6
    // Assert configs["bullet_despawn"].Lifetime == 8
    // Assert each Color is non-nil and equal to color.White
}

func TestVFXManager_SpawnPuffUsesPixelLifetime(t *testing.T) {
    // Build manager with a single "test_pixel" entry, lifetime_frames=10.
    // m.SpawnPuff("test_pixel", 0, 0, 1, 0)
    // Inspect m.system.Particles(); assert Duration == 10 and MaxDuration == 10.
}

func TestVFXManager_RejectsInvalidPixelColor(t *testing.T) {
    // JSON with "color": "#123456" must produce a config (fallback white)
    // and log an error. Loader must not panic.
}

func TestVFXJSON_AllPixelEntriesUse1BitPalette(t *testing.T) {
    // Read vfx.json directly, unmarshal into []VFXConfig.
    // For every entry where Pixel != nil, assert Color is "#000000" or "#FFFFFF".
}
```

**Expected failures (red phase):**
- `ParticleData.Pixel` field does not exist → schema test fails to compile.
- `configs["muzzle_flash"]` etc. absent → loader test fails.
- `Config.Lifetime` does not exist → SpawnPuff test fails to compile.

## Notes

- Foundation story — **must merge before** US-030, US-031, US-032 (their consumers reference these type keys).
- The 1-bit palette restriction is enforced at load time (`parseHexColor`), not just convention.
- `velocity_range` was deliberately dropped from the JSON schema vs. the previous SPEC draft: callers (US-030/031/032) already pass `randRange` directly into `SpawnPuff`, so a JSON value would be either redundant or in conflict. Keep it caller-side.
- `count` likewise stays caller-side.
