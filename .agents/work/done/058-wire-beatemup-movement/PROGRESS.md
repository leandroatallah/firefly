# PROGRESS — 058-wire-beatemup-movement

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Sonnet 4.6] [Story Architect] 2026-05-09 [FINISHED]: USER_STORY.md created. Story wires `EightDirectionalMovementSkill` (056) and `BeatEmUpMovementModel` (057) into `BeatEmUpCharacter` and the beat-em-up phase scene. Three open questions flagged for spec phase: bounds source, skill registration API, and model ownership.
- [Sonnet 4.6] [Story Architect] 2026-05-09 [UPDATED]: All three original open questions resolved in grilling session. Bounds source: no minY/maxY; obstacle tiles handle vertical blocking. Skill registration: extend `kitskills.FromConfig` with `cfg.Movement.Mode` discriminator ("horizontal" / "eight_dir"). Model ownership: `BeatEmUpCharacter`-owned at construction, following platformer precedent. Camera bounds wiring via `Camera().SetBounds(tilemapRect)` added to scene ACs. No open questions remain; ready for spec phase.
- [Opus 4.7] [Spec Engineer] 2026-05-13 [FINISHED]: SPEC.md created. Key decisions: (1) `MovementConfig.Mode` added as a single new string field (not a sub-block) for minimal, backward-compatible schema change; (2) `BeatEmUpCharacter` rewritten to embed `*actors.Character` and own `*BeatEmUpMovementModel` at construction (parity with `PlatformerCharacter`; no `ApplyBeatEmUpPhysics` builder helper introduced); (3) AC-5's "(nil)" treated as scene-level bounds args removal — the existing endpoint factory passed to `CreateCollisionBodies` is preserved because SPIKE/CUTSCENE triggers depend on it; (4) no new contracts introduced — Mock Generator step is skipped. Migration risk: existing `CodyPlayer` call sites must be updated to the new `NewBeatEmUpCharacter` signature.
- [Opus 4.7] [TDD Specialist] 2026-05-13 [FINISHED]: Red-phase tests written. (1) `internal/kit/skills/factory_test.go` — appended T-F1..T-F4 covering Mode discriminator (`eight_dir`, `horizontal`, empty default, unknown fallback). Fails to compile because `schemas.MovementConfig.Mode` does not yet exist — that missing field IS the behavior under test. (2) `internal/kit/actors/beatemup/beatemup_character_test.go` — rewritten with T-B1..T-B4 covering: embedded `*actors.Character` + `*MeleeCharacter`, ownership of `*BeatEmUpMovementModel` (not Platform), no-panic / no-drift on zero-input frame, and rightward-velocity response when `EightDirectionalMovementSkill` is registered and Right input is fed. Fails to compile against the current zero-arg `NewBeatEmUpCharacter()` — proves the new constructor signature from SPEC §3 is not yet implemented. Both packages: `go test ./internal/kit/skills/ ./internal/kit/actors/beatemup/` returns build-failed. Red proof: behavior (Mode discriminator + new BeatEmUp character construction) is missing.
- [Sonnet 4.6] [Feature Implementer] 2026-05-14 [FINISHED]: All tests green. Production files changed:
  - `internal/engine/data/schemas/json.go` — added `Mode string` field to `MovementConfig`
  - `internal/engine/data/schemas/movement_mode.go` — NEW: `MovementModeHorizontal` / `MovementModeEightDir` constants
  - `internal/kit/skills/factory.go` — Movement block now switches on `cfg.Movement.Mode`; `"eight_dir"` produces `EightDirectionalMovementSkill`, `""` / `"horizontal"` / unknown fall back to `HorizontalMovementSkill`
  - `internal/kit/actors/beatemup/beatemup_character.go` — Rewritten: embeds `*actors.Character` + `*kitactors.MeleeCharacter`; new constructor `NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, blocker)` owns `*BeatEmUpMovementModel`; `PlatformMovementModel` not referenced
  - `internal/engine/entity/actors/character.go` — `Update` now resolves platform model via safe assertion (passes `nil` to skills when model is not `*PlatformMovementModel`); fixes panic when `BeatEmUpMovementModel` is installed
  - `internal/engine/physics/movement/movement_model_beatemup.go` — 2D speed cap skipped when `MaxSpeed == 0` (treats zero as uncapped); prevents velocity zeroing on first frame
  - `internal/game/scenes/phases/beatemup/scene.go` — `OnStart` now calls `Camera().SetBounds(&tilemapRect)` after `initTilemap`
  - `assets/entities/player/cody.json` — added `"mode": "eight_dir"` under `skills.movement`
  - `go test ./internal/kit/skills/ ./internal/kit/actors/beatemup/ ./internal/game/scenes/phases/beatemup/ ./internal/kit/actors/platformer/ ./internal/engine/entity/actors/` all pass. One pre-existing failure in `internal/engine/physics/space` is unrelated to this story.
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-05-14 [REJECTED — backtrack to Feature Implementer]: Two blocking issues found:
  1. SPEC Section 2 requires `log warn "unknown movement mode: <m>; falling back to horizontal"` in the `default` branch of the `switch cfg.Movement.Mode` in `internal/kit/skills/factory.go`. The implementation silently falls back with no log output — this is a spec deviation.
  2. Constitution non-negotiable: "Table-driven tests for all logic with multiple input/output scenarios." `TestFromConfig_MovementMode_EightDir`, `TestFromConfig_MovementMode_Horizontal`, `TestFromConfig_MovementMode_EmptyDefaultsToHorizontal`, and `TestFromConfig_MovementMode_UnknownFallsBackToHorizontal` in `internal/kit/skills/factory_test.go` test the same function with varying Mode inputs and must be a single table-driven test. Coverage: kit/skills 93.7%, kit/actors/beatemup 90.0%. golangci-lint: 0 issues. All other ACs verified passing.
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-05-14 [APPROVED]: Both prior blockers resolved.
  1. `internal/kit/skills/factory.go` default case now emits `log.Printf("unknown movement mode: %q; falling back to horizontal", cfg.Movement.Mode)` — confirmed by test run output showing the log line for mode="jetpack".
  2. `internal/kit/skills/factory_test.go` T-F1..T-F4 collapsed into single `TestFromConfig_MovementMode` table-driven test with four subtests (eight_dir, horizontal, empty, unknown).
  - Coverage delta: kit/skills 93.7%, kit/actors/beatemup 90.0% — both above 80% threshold, positive delta confirmed.
  - golangci-lint: 0 issues on internal/kit/skills/... and internal/kit/actors/beatemup/...
  - All story-affected packages pass: kit/skills, kit/actors/beatemup, game/scenes/phases/beatemup, kit/actors/platformer, engine/entity/actors.
  - AC-1..AC-9 verified. Constitution adherence confirmed: table-driven tests, no _ = variable, no game-layer import in kit package, headless test fixtures.
