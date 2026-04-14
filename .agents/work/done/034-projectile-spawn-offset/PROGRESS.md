# PROGRESS — US-034

## Status: Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Notes

New story created to address projectile spawn offset configuration. Currently, projectiles spawn at the entity origin without a configurable offset.

### Validation Findings (2026-04-14)

1. **TODO reference inaccurate**: USER_STORY.md Context section claims a TODO exists in `internal/engine/physics/skill/skill_shooting.go`. No such TODO exists in the file or in git history. The underlying need is still valid (spawn position lacks configurable offset), but the story's justification text is factually incorrect. Recommend removing the TODO reference and rephrasing context to describe the current behavior directly.

2. **Spec tests not table-driven**: The Red Phase tests in SPEC.md use separate test functions for facing-right, facing-left, and zero-offset scenarios. The constitution requires "Table-driven tests for all logic with multiple input/output scenarios." The TDD Specialist should combine these into a single table-driven test when implementing.

3. **Bounded context not in constitution**: `internal/engine/combat/` is not listed in the constitution's Bounded Contexts table, though the package exists. This is a constitution gap, not a story issue.

4. **All other aspects validated**: Constructor signature matches codebase (6 params), struct fields match, Fire() method signature matches, fp16 terminology used correctly, acceptance criteria are clear and testable, referenced files exist.

### Mock Generator Findings (2026-04-14)

No new mocks or contract changes required. All existing mocks are sufficient:

- **Local mock** `mockProjectileManager` in `weapon/mocks_test.go` -- unchanged, captures `SpawnProjectile` calls with position args.
- **Shared mock** `mocks.MockVFXManager` in `internal/engine/mocks/mock_vfx_manager.go` -- unchanged, captures `SpawnPuff` calls.
- **Contract** `combat.ProjectileManager` in `internal/engine/contracts/combat/projectile_manager.go` -- unchanged, receives absolute coordinates.

### TDD Specialist (2026-04-14)

Red Phase test written:

- **Test file:** `internal/engine/combat/weapon/weapon_test.go`
- **New test:** `TestProjectileWeapon_Fire_SpawnOffset` -- table-driven with 5 sub-cases covering AC1-AC6.
- **Existing tests updated:** All `NewProjectileWeapon` calls changed from 6 args to 8 args (appending `0, 0` for backward compatibility).
- **Expected failure:** Will not compile until constructor signature is updated to accept 8 params. Once compiled, offset assertions will fail until `Fire()` applies offsets.

## Test Files

- `internal/engine/combat/weapon/weapon_test.go`
- `internal/engine/combat/weapon/mocks_test.go`

## Log

- 2026-04-08: Story created to resolve spawn position alignment issue. Replaces previous integration test story.
- 2026-04-08: Spec Engineer — SPEC.md created. Adds spawnOffsetX16/Y16 parameters, applies offset in Fire() with facing direction support.
- 2026-04-14: Story Architect — Validation completed. Story and spec are valid and aligned with codebase. Two minor findings documented: (1) non-existent TODO reference in story context, (2) spec Red Phase tests should use table-driven format per constitution. Neither finding blocks implementation.
- 2026-04-14: Spec Engineer — SPEC.md revised. Removed inaccurate TODO reference from context and integration points. Rewrote Red Phase to use a single table-driven test with five sub-cases (facing-right, facing-left, zero-offset, no-VFX, negative-Y-offset). Added call-site migration notes for existing 6-arg constructor calls.
- 2026-04-14: Mock Generator — No new mocks needed. All contracts and mocks are already in place for the Red Phase test.
- 2026-04-14: TDD Specialist — Red Phase test `TestProjectileWeapon_Fire_SpawnOffset` added to weapon_test.go. All existing constructor calls updated to 8-arg signature.
- 2026-04-14: Feature Implementer — Production code implemented: added `spawnOffsetX16`/`spawnOffsetY16` fields to struct, updated constructor to 8 params, applied offset logic in `Fire()` with left-facing negation, updated factory.go and weapons.go with `0, 0` defaults. All 16 tests pass, build succeeds.
- 2026-04-14: Workflow Gatekeeper — All quality gates passed. Coverage delta: `internal/engine/combat/weapon` at 90.0% (above 80% goal). All 16 tests pass across the full suite. golangci-lint reports 0 issues. TDD Red-Green-Refactor cycle confirmed. Table-driven tests present, no `_ = variable` in production code, fp16 conventions followed. Story closed.
