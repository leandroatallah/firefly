# SPEC-024 — JSON-Driven Skills Configuration

**Branch:** `024-json-skills-config`  
**Bounded Context:** Engine (Data Schemas, Physics/Skill)  
**Contracts:** None new; uses existing `body.Shooter`, `body.MovableCollidable`

---

## Overview

Move skill configuration (movement, jump, dash, shooting) from hardcoded values in `internal/game/scenes/phases/player.go` to JSON schema in `internal/engine/data/schemas/`. A factory in `internal/engine/physics/skill/` will instantiate skills from the parsed config.

---

## Technical Requirements

### 1. Schema Extension (`internal/engine/data/schemas/json.go`)

Add `SkillsConfig` struct to existing schema:

```go
type MovementConfig struct {
    Enabled              *bool   `json:"enabled,omitempty"`
    HorizontalSpeed      float64 `json:"horizontal_speed,omitempty"`
}

type JumpConfig struct {
    Enabled            *bool   `json:"enabled,omitempty"`
    JumpCutMultiplier  float64 `json:"jump_cut_multiplier,omitempty"`
    CoyoteTimeFrames   int     `json:"coyote_time_frames,omitempty"`
    JumpBufferFrames   int     `json:"jump_buffer_frames,omitempty"`
}

type DashConfig struct {
    Enabled     *bool `json:"enabled,omitempty"`
    DurationMs  int   `json:"duration_ms,omitempty"`
    CooldownMs  int   `json:"cooldown_ms,omitempty"`
    Speed       int   `json:"speed,omitempty"`
    CanAirDash  *bool `json:"can_air_dash,omitempty"`
}

type ShootingConfig struct {
    Enabled         *bool `json:"enabled,omitempty"`
    CooldownFrames  int   `json:"cooldown_frames,omitempty"`
    ProjectileSpeed int   `json:"projectile_speed,omitempty"` // fp16
    ProjectileRange int   `json:"projectile_range,omitempty"` // fp16
    Directions      int   `json:"directions,omitempty"`       // 4 or 8
}

type SkillsConfig struct {
    Movement *MovementConfig `json:"movement,omitempty"`
    Jump     *JumpConfig     `json:"jump,omitempty"`
    Dash     *DashConfig     `json:"dash,omitempty"`
    Shooting *ShootingConfig `json:"shooting,omitempty"`
}
```

Add `Skills` field to `SpriteData`:

```go
type SpriteData struct {
    // ... existing fields
    Skills *SkillsConfig `json:"skills,omitempty"`
}
```

### 2. Skill Factory (`internal/engine/physics/skill/factory.go`)

Create new file with:

```go
type SkillDeps struct {
    Shooter      body.Shooter
    EventManager interface{ Publish(interface{}) }
    OnJump       func(body.MovableCollidable)
}

func FromConfig(cfg *schemas.SkillsConfig, deps SkillDeps) []skill.Skill
```

**Logic:**
- Return empty slice if `cfg == nil`
- For each skill sub-config (Movement, Jump, Dash, Shooting):
  - If sub-config is `nil` or `Enabled` is `false`, skip
  - Instantiate corresponding skill with config values
  - Append to result slice
- Jump skill: attach `deps.OnJump` callback if provided
- Shooting skill: requires `deps.Shooter != nil`

### 3. Player Factory Update (`internal/game/scenes/phases/player.go`)

Replace `addJumpSkill`, `addShootingSkill`, and manual skill additions with:

```go
func createPlayer(ctx *app.AppContext, playerType gameentitytypes.PlayerType) (platformer.PlatformerActorEntity, error) {
    // ... existing player creation
    
    spriteData := p.GetSpriteData() // Assume this method exists or extract from player
    
    deps := skill.SkillDeps{
        EventManager: ctx.EventManager,
        OnJump: func(body body.MovableCollidable) {
            pos := body.Position()
            jumpPos := image.Point{X: pos.Min.X + pos.Dx()/2, Y: pos.Max.Y}
            ctx.EventManager.Publish(&events.ActorJumpedEvent{
                X: float64(jumpPos.X),
                Y: float64(jumpPos.Y),
            })
        },
    }
    
    if currentScene := ctx.SceneManager.CurrentScene(); currentScene != nil {
        if shooter, ok := currentScene.(body.Shooter); ok {
            deps.Shooter = shooter
        }
    }
    
    skills := skill.FromConfig(spriteData.Skills, deps)
    for _, s := range skills {
        p.GetCharacter().AddSkill(s)
    }
    
    return p, nil
}
```

Remove `addJumpSkill` and `addShootingSkill` functions.

### 4. JSON Update (`assets/data/entities/player/climber.json`)

Add `"skills"` block:

```json
{
  "skills": {
    "movement": {
      "enabled": true
    },
    "jump": {
      "enabled": true,
      "jump_cut_multiplier": 0.4
    },
    "dash": {
      "enabled": true
    },
    "shooting": {
      "enabled": true,
      "cooldown_frames": 15,
      "projectile_speed": 6,
      "projectile_range": 4,
      "directions": 4
    }
  }
}
```

---

## Pre-conditions

- `internal/engine/data/schemas/json.go` exists with `SpriteData` struct
- `internal/engine/physics/skill/` contains skill implementations (`skill_platform_move.go`, `skill_platform_jump.go`, `skill_dash.go`, `skill_shooting.go`)
- `internal/game/scenes/phases/player.go` currently hardcodes skill instantiation
- `assets/data/entities/player/climber.json` exists

---

## Post-conditions

- All skill parameters are read from JSON
- No magic numbers in `player.go`
- Skills with `enabled: false` or missing config are not instantiated
- Factory returns empty slice for `nil` config (no panic)

---

## Integration Points

**Within Engine Bounded Context:**
- `internal/engine/data/schemas/` ← schema definition
- `internal/engine/physics/skill/` ← factory implementation
- Uses existing contracts: `body.Shooter`, `body.MovableCollidable`

**Game Layer:**
- `internal/game/scenes/phases/player.go` ← consumes factory
- `assets/data/entities/player/climber.json` ← data source

---

## Red Phase (Failing Test Scenario)

**File:** `internal/engine/physics/skill/factory_test.go`

**Test:** `TestFromConfig_AllSkillsEnabled`

**Scenario:**
```go
cfg := &schemas.SkillsConfig{
    Movement: &schemas.MovementConfig{Enabled: ptrBool(true)},
    Jump:     &schemas.JumpConfig{Enabled: ptrBool(true), JumpCutMultiplier: 0.4},
    Dash:     &schemas.DashConfig{Enabled: ptrBool(true)},
    Shooting: &schemas.ShootingConfig{Enabled: ptrBool(true), CooldownFrames: 15},
}
deps := skill.SkillDeps{Shooter: &mockShooter{}}
skills := skill.FromConfig(cfg, deps)

// Expected: 4 skills returned
// Expected: Jump skill has cut multiplier 0.4
// Expected: Shooting skill has cooldown 15
```

**Initial state:** `skill.FromConfig` does not exist → compilation fails.

**Additional tests:**
- `TestFromConfig_NilConfig` → returns empty slice
- `TestFromConfig_DisabledSkills` → `enabled: false` omits skill
- `TestFromConfig_MissingShooter` → shooting skill skipped if `deps.Shooter == nil`

---

## Design Decisions

1. **Pointer fields for `Enabled`:** Distinguishes between "not set" (nil) and "explicitly false".
2. **Factory returns slice, not map:** Simpler iteration; skill order doesn't matter for `AddSkill`.
3. **`SkillDeps` struct:** Encapsulates all injected dependencies (shooter, event manager, callbacks) in one parameter.
4. **No new contracts:** Reuses existing `body.Shooter` and `body.MovableCollidable` interfaces.
