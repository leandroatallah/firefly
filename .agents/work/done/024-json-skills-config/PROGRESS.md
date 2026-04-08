# PROGRESS-024 — JSON-Driven Skills Configuration

## Status

- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Gatekeeper

**Status: ✅ Done**

---

## Log

**Spec Engineer [2026-04-06]:** SPEC.md created. Key decisions: (1) Factory returns slice instead of map for simplicity. (2) `SkillDeps` struct encapsulates all injected dependencies (shooter, event manager, callbacks). (3) Pointer `*bool` for `Enabled` fields distinguishes nil (not set) from false (explicitly disabled). (4) No new contracts needed—reuses existing `body.Shooter` and `body.MovableCollidable`.

**Mock Generator [2026-04-06]:** Created `internal/engine/mocks/event_manager.go` with `MockEventManager` for shared use across test packages. Existing `MockShooter` in `internal/engine/mocks/shooter.go` already covers the `body.Shooter` interface. No package-local mocks required—all dependencies can be satisfied with shared mocks.

**TDD Specialist [2026-04-06]:** Created `internal/engine/physics/skill/factory_test.go` with 7 table-driven tests covering: (1) all skills enabled with correct config values, (2) nil config returns empty slice, (3) disabled skills omitted, (4) missing Shooter skips shooting skill, (5) partial config, (6) OnJump callback attachment, (7) nil sub-configs. Red Phase: tests fail with "undefined: SkillDeps" and "undefined: FromConfig"—proving factory function and dependency struct are missing.

**Feature Implementer [2026-04-06]:** Created `internal/engine/physics/skill/factory.go` with `SkillDeps` struct and `FromConfig` function. Implementation: (1) Returns empty slice for nil config. (2) Checks `Enabled` field (nil or true = enabled, false = disabled). (3) Instantiates Movement, Jump, Dash, Shooting skills with config values. (4) Attaches OnJump callback to JumpSkill. (5) Skips Shooting skill if Shooter is nil. Fixed test compilation error (unused variable) and verified all 45 tests pass.

**Gatekeeper [2026-04-06]:** ✅ **Verification Complete**
- **Red-Green-Refactor**: Confirmed. Tests created first (red), factory implemented (green), all 45 tests passing.
- **Spec Compliance**: All acceptance criteria met:
  - AC1: SkillsConfig struct added to schemas with movement, jump, dash, shooting sub-configs ✅
  - AC2: SkillsConfig parsed as part of entity JSON (added to SpriteData) ✅
  - AC3: skill.Factory (FromConfig) instantiates skills from SkillsConfig ✅
  - AC4: climber.json updated with skills block matching schema ✅
  - AC5: player.go replaced with factory call; no magic numbers remain ✅
  - AC6: Disabled/missing skills handled gracefully (no panic) ✅
  - AC7: Unit tests cover factory behavior, disabled skills, jump cut multiplier ✅
- **Code Quality**:
  - Table-driven tests: All 7 factory tests use table-driven pattern ✅
  - No unused variables: go vet passes ✅
  - DDD alignment: Factory in skill package, schema in data/schemas, integration in game layer ✅
  - Headless Ebitengine: Tests run without display initialization ✅
- **Coverage Delta**: Factory tests provide comprehensive coverage of FromConfig logic (nil config, enabled/disabled states, dependency injection, callback attachment). Skill package coverage: 77.9%.
- **Integration**: Player creation now uses factory instead of hardcoded skill instantiation. ClimberPlayer stores SpriteData for access to skills config.
- **Final Verification**: All 45 tests pass. go vet clean. Builds successfully. Ready for production.
