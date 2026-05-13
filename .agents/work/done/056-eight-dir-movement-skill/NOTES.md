# NOTES — 056-eight-dir-movement-skill

Human-only context. Not consumed by downstream agents directly.

## Investigation Findings

- `internal/engine/skill/skill.go` defines `ActiveSkill.HandleInput(b, model *physicsmovement.PlatformMovementModel, space)`. The model parameter is a **concrete type**, not an interface. Every existing kit skill (`HorizontalMovementSkill`, `JumpSkill`, `DashSkill`, `ShootingSkill`) imports `physicsmovement` for this reason. The story requirement "no import of any movement model package" is therefore impossible at the strict letter; the practical reading is "no import of a *genre-specific* movement model (beat-em-up, top-down, platformer-specific helpers)". This skill imports `physicsmovement` purely for the type signature and only ever calls `model.IsInputBlocked()` — which is shared across all `*MovementModel` implementations conceptually. Loosening this signature to a `PlayerInputBlocker` interface is a separate refactor and out of scope here.

- `HorizontalMovementSkill` uses an `input.HorizontalAxis` for press/release smoothing and a `HorizontalInertia` branch. After re-reading the story acceptance criteria, neither smoothing nor inertia is required for this skill: AC-2 says directly "calls `OnMove*` on active directions". Beat-em-up / top-down feel is governed by the movement model (story 057). Keeping the skill dead simple aligns with the "passive model + active skill" architecture choice in the grilling session.

- `MovementConfig` schema (`internal/engine/data/schemas/json.go`) currently has only `Enabled` and `HorizontalSpeed`. There is no `Mode` field. The user story's "selected via `cfg.Movement.Mode == 'eight_dir'`" is **story 058's** job. We document it here but do not modify the schema or `FromConfig` in this story.

## Design Rationale

- **Guard order = IsInputBlocked first.** Matches `HorizontalMovementSkill` (`platform_move.go` L43–L52) and `JumpSkill` (`platform_jump.go` L57). A scripted/cutscene block (`IsInputBlocked`) is a higher-precedence external override than a per-actor `Immobile()` flag. Reversing it would let immobile-but-not-blocked actors short-circuit before the cutscene check (irrelevant here, but inconsistent with sibling skills, which trips reviewer expectations).

- **Immobile zeroes both axes.** `HorizontalMovementSkill` only zeroes `vx16` (preserves `vy16` for gravity). For 8-dir, Y is ground-plane motion, not gravity — there is no falling state to preserve. So zeroing both X and Y is correct.

- **No new contracts.** Mock Generator can be skipped. Existing `newMockMovableCollidable` covers the body surface; if Up/Down recording fields are missing the TDD Specialist will add them.

## Risks

- **Risk:** `newMockMovableCollidable` may not track `OnMoveUp`/`OnMoveDown` invocations (the platformer mock historically only needed Left/Right). **Mitigation:** TDD Specialist extends the local mock; documented in SPEC §8.

- **Risk:** Future refactor to decouple skill signature from `*PlatformMovementModel` will touch this file. **Mitigation:** Acceptable — every kit skill carries the same coupling and will migrate together. Tracked separately.

## Out of Scope

- Schema `MovementConfig.Mode` field, JSON loader, validation → 058.
- `FromConfig` dispatch on `Mode == "eight_dir"` → 058.
- `BeatEmUpMovementModel` (passive model that consumes velocity from the skill) → 057.
- Diagonal speed normalization, depth-z sorting integration, altitude/jump for beat-em-up → later stories.
