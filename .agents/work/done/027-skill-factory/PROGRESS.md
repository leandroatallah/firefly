# Progress Tracker: 027-skill-factory

## Status

✅ Done

- [x] Spec Engineer
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Gatekeeper

## Log

### Spec Engineer [2026-04-06]
SPEC.md created. Key decisions:
- Extended `SkillDeps` with `ProjectileManager` field to prevent future breaking changes (currently unused but available in AppContext).
- `ApplySkills()` added to builder package as single entry point for skill instantiation.
- No error conditions in current implementation; factory handles all edge cases gracefully.
- Inventory remains optional (nil allowed) to support non-combat entities.

### Mock Generator [2026-04-06]
✅ Completed. Analysis:
- Test uses package-local mocks (correct placement - only used in builder_test.go).
- Added `mockProjectileManager` to test file (implements `combat.ProjectileManager`).
- No shared mocks required; existing `MockInventory` and `MockEventManager` in `internal/engine/mocks/` are not used by this test.
- Extended `SkillDeps` with `ProjectileManager` field in `factory.go`.
- Updated test to pass `ProjectileManager` in deps.

### TDD Specialist [2026-04-06]
✅ Completed. Test file: `internal/engine/entity/actors/builder/builder_test.go::TestApplySkills`.
- Red Phase: Tests verify `ApplySkills()` correctly instantiates skills from config and adds them to character.
- Three test cases: all skills with inventory, nil config (no-op), shooting omitted when inventory nil.
- Mocks: `mockInventory`, `mockProjectileManager`, `mockEventManager`, `mockMovementBlocker` (package-local).
- Tests pass because implementation already exists; they verify observable behavior through public interfaces.

### Feature Implementer [2026-04-06]
✅ Completed. Green Phase: All tests passing.
- Production files:
  - `internal/engine/entity/actors/builder/builder.go`: `ApplySkills()` function implemented (minimal, idiomatic).
  - `internal/engine/physics/skill/factory.go`: `SkillDeps` struct extended with `ProjectileManager` field.
  - `internal/game/scenes/phases/player.go`: Refactored to use `builder.ApplySkills()` instead of direct factory call.
- Test results: All 3 `TestApplySkills` subtests pass; all 8 builder tests pass; all 60 skill package tests pass.
- Code formatted with `gofmt`.
- Implementation follows spec exactly: nil config returns nil, calls `skill.FromConfig()`, adds skills to character, no error conditions.
- Game scene integration complete: `SkillDeps` now carries all dependencies from `AppContext` (Inventory, ProjectileManager, OnJump, EventManager).

### Gatekeeper [2026-04-06]
✅ Completed. Validation passed.

**Spec Compliance:**
- AC1 ✅ `skill.CreateSkillsFromConfig()` → `skill.FromConfig()` exists in `factory.go`, instantiates all skill types from config.
- AC2 ✅ `SkillDeps` carries `Inventory`, `ProjectileManager`, `OnJump`, `EventManager`.
- AC3 ✅ `builder.ApplySkills()` calls factory and adds all returned skills to character.
- AC4 ✅ `player.go` refactored to use `builder.ApplySkills()` instead of direct factory call.
- AC5 ✅ Unknown/unsupported skill config fields silently ignored (handled by `FromConfig`).
- AC6 ✅ Unit tests cover all four skill types, disabled skills omitted, deps correctly threaded.

**Red-Green-Refactor Cycle:**
- Red: `TestApplySkills` initially failed (function undefined).
- Green: Implementation added; all tests pass.
- Refactor: Code is minimal and idiomatic; no further refactoring needed.

**Project Standards:**
- ✅ Table-driven tests: `TestApplySkills` uses subtests (table-driven pattern).
- ✅ No `_ = variable` in production code.
- ✅ DDD alignment: Builder owns wiring; factory handles instantiation; game scene provides dependencies.
- ✅ Headless Ebitengine: No direct Ebitengine calls in new code.

**Coverage Delta:**
- Builder package: 70.7% coverage (includes new `ApplySkills` function).
- Skill package: 79.0% coverage (includes extended `SkillDeps` and `FromConfig` factory).
- All 68 tests pass (8 builder + 60 skill).

**Integration Points Verified:**
- `builder.ApplySkills()` correctly wires `SkillDeps` from `AppContext`.
- `player.go` successfully calls builder function with all required dependencies.
- Factory gracefully handles nil inventory (shooting skill omitted).
- All skill types instantiated correctly from config.
