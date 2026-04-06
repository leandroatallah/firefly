# 027 — Skill Factory (JSON-to-Skills Instantiation)

**Branch:** `027-skill-factory`
**Bounded Context:** Engine (`internal/engine/physics/skill/`, `internal/engine/entity/actors/builder/`)

## Story

As a game developer using this boilerplate,
I want skills to be instantiated automatically from the entity JSON config via a factory,
so that adding a new character requires only a JSON file, not Go code changes.

## Context

US-024 introduced `SkillsConfig` and JSON parsing. This story wires the factory into the entity builder pipeline so that `builder.go` calls `skill.CreateSkillsFromConfig()` instead of manually constructing skills. It also ensures `SkillDeps` (inventory, projectile manager, input profile) flows correctly from `AppContext` through the builder.

## Acceptance Criteria

- **AC1** — `skill.CreateSkillsFromConfig(config SkillsConfig, deps SkillDeps) []Skill` implemented in `factory.go`.
- **AC2** — `SkillDeps` carries `*inventory.Inventory`, `*projectile.Manager`, `*input.Profile`.
- **AC3** — `builder.ApplySkills(character, config, deps)` calls the factory and adds all returned skills.
- **AC4** — `internal/game/scenes/phases/player.go` no longer manually calls `AddSkill()`.
- **AC5** — Unknown or unsupported skill config fields are silently ignored (no panic, no error).
- **AC6** — Unit tests cover: all four skill types instantiated from config, disabled skills omitted, deps correctly threaded.

## Notes

- Depends on US-024 (JSON schema), US-025 (inventory), US-026 (shooting skill).
- Factory is the final integration point; all prior stories must be complete first.
