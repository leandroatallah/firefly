# Technical Specification: Skill Factory

**Branch:** `027-skill-factory`  
**Bounded Context:** Engine (`internal/engine/physics/skill/`, `internal/engine/entity/actors/builder/`)

## Overview

Integrate the existing `skill.FromConfig()` factory into the entity builder pipeline, eliminating manual skill instantiation in game scenes. The factory already exists and handles JSON-to-skill conversion; this story wires it into `builder.go` and extends `SkillDeps` to carry all required dependencies from `AppContext`.

## Technical Requirements

### 1. Extend `SkillDeps` Structure

**File:** `internal/engine/physics/skill/factory.go`

Current `SkillDeps`:
```go
type SkillDeps struct {
    Inventory    combat.Inventory
    OnJump       func(interface{})
    EventManager interface{ Publish(interface{}) }
}
```

**Change:** Add `ProjectileManager` field (currently unused but required for future weapon system integration):
```go
type SkillDeps struct {
    Inventory         combat.Inventory
    ProjectileManager projectile.Manager
    OnJump            func(interface{})
    EventManager      interface{ Publish(interface{}) }
}
```

**Rationale:** `ProjectileManager` is available in `AppContext` and will be needed when weapons are fully integrated with the shooting skill. Including it now prevents a future breaking change.

### 2. Add Builder Function

**File:** `internal/engine/entity/actors/builder/builder.go`

**New Function:**
```go
func ApplySkills(
    character actors.ActorEntity,
    spriteData schemas.SpriteData,
    deps skill.SkillDeps,
) error
```

**Behavior:**
- If `spriteData.Skills == nil`, return `nil` (no-op).
- Call `skill.FromConfig(spriteData.Skills, deps)`.
- For each returned skill, call `character.GetCharacter().AddSkill(skill)`.
- Return `nil` (no error conditions).

**Integration Point:** This function is called by game-layer code (e.g., `internal/game/scenes/phases/player.go`) after `PreparePlatformer` and `ConfigureCharacter`.

### 3. Update Game Scene

**File:** `internal/game/scenes/phases/player.go`

**Current Implementation:**
```go
deps := skill.SkillDeps{
    OnJump: func(b interface{}) { /* ... */ },
}
skills := skill.FromConfig(spriteData.Skills, deps)
for _, s := range skills {
    p.GetCharacter().AddSkill(s)
}
```

**Change:** Replace direct factory call with `builder.ApplySkills()`:
```go
deps := skill.SkillDeps{
    Inventory:         nil, // TODO: wire inventory when available
    ProjectileManager: ctx.ProjectileManager,
    OnJump:            func(b interface{}) { /* ... */ },
    EventManager:      ctx.EventManager,
}
if err := builder.ApplySkills(p, *spriteData, deps); err != nil {
    return nil, err
}
```

**Note:** `Inventory` remains `nil` until the inventory system is wired into the player entity. The factory already handles `nil` inventory by skipping the shooting skill.

## Pre-Conditions

- US-024 complete: `schemas.SkillsConfig` exists and is parsed from JSON.
- US-025 complete: `combat.Inventory` interface exists.
- US-026 complete: `ShootingSkill` constructor accepts `combat.Inventory`.
- `skill.FromConfig()` already implemented in `factory.go`.

## Post-Conditions

- `builder.ApplySkills()` is the single entry point for skill instantiation.
- Game scenes no longer call `skill.FromConfig()` directly.
- `SkillDeps` carries all dependencies from `AppContext`.
- Unknown skill config fields are silently ignored (already handled by `FromConfig`).

## Integration Points

- **Builder Package:** New `ApplySkills()` function added to `builder.go`.
- **Skill Factory:** `SkillDeps` extended with `ProjectileManager`.
- **Game Scene:** `player.go` refactored to use builder function.

## Red Phase: Failing Test Scenario

**Test File:** `internal/engine/entity/actors/builder/builder_test.go`

**Test Name:** `TestApplySkills`

**Scenario:**
1. Create a mock `ActorEntity` with a character that tracks added skills.
2. Create a `schemas.SpriteData` with a non-nil `Skills` field containing:
   - Movement skill enabled.
   - Jump skill enabled with custom `JumpCutMultiplier`.
   - Dash skill enabled with custom duration/cooldown.
   - Shooting skill enabled.
3. Create `SkillDeps` with mock inventory, projectile manager, and callbacks.
4. Call `builder.ApplySkills(character, spriteData, deps)`.
5. Assert:
   - No error returned.
   - Character has exactly 4 skills added.
   - Each skill type is present (movement, jump, dash, shooting).

**Expected Failure:** `undefined: ApplySkills` (function does not exist).

**Additional Test Cases:**
- `spriteData.Skills == nil` → no skills added, no error.
- Disabled skills in config → omitted from result.
- `deps.Inventory == nil` → shooting skill omitted.

## Design Decisions

1. **No Error Handling:** `ApplySkills` returns `error` for consistency with other builder functions, but currently always returns `nil`. The factory handles all edge cases (nil config, disabled skills) gracefully.

2. **ProjectileManager in Deps:** Added now to avoid breaking changes later. Not yet used by any skill but will be required when weapon firing is fully integrated.

3. **Builder Owns Wiring:** The builder package is responsible for connecting JSON config to entity state. Game scenes only provide dependencies and call builder functions.

4. **Inventory Remains Optional:** The shooting skill is skipped if `Inventory` is `nil`, allowing characters without combat capabilities to use the same pipeline.
