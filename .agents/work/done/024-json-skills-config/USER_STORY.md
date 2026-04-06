# 024 — JSON-Driven Skills Configuration

**Branch:** `024-json-skills-config`
**Bounded Context:** Engine (`internal/engine/data/schemas/`, `internal/engine/physics/skill/`)

## Story

As a game developer using this boilerplate,
I want to configure skills (movement, jump, dash, shooting) in the entity JSON file,
so that I can tune character behavior without changing Go code.

## Context

Skills are currently hardcoded in `internal/game/scenes/phases/player.go` with magic numbers. Different characters need different parameters (e.g., jump cut multiplier, dash speed). Moving this to JSON enables per-character customization and removes magic numbers from code.

## Acceptance Criteria

- **AC1** — `SkillsConfig` struct added to `internal/engine/data/schemas/` covering `movement` (horizontal, jump, dash) and `combat` (shooting) sub-configs.
- **AC2** — `SkillsConfig` is parsed as part of the entity JSON (added to the existing entity schema struct).
- **AC3** — `skill.Factory` in `internal/engine/physics/skill/factory.go` instantiates skills from `SkillsConfig`.
- **AC4** — `player/climber.json` updated with a `"skills"` block matching the schema.
- **AC5** — `internal/game/scenes/phases/player.go` replaced with a call to `skill.Factory`; no magic numbers remain.
- **AC6** — Missing or `false` `enabled` fields disable the corresponding skill (no panic).
- **AC7** — Unit tests cover: factory produces correct skill types, disabled skills are omitted, jump cut multiplier is applied.

## Notes

- `SkillDeps` struct carries `ProjectileManager` and any other injected dependencies needed by skills.
- Jump config fields: `jump_cut_multiplier`, `coyote_time_frames`, `jump_buffer_frames`.
- Dash config fields: `duration_ms`, `cooldown_ms`, `speed`, `can_air_dash`.
