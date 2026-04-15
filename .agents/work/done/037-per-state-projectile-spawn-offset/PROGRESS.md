# PROGRESS — 037-per-state-projectile-spawn-offset

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Dependencies

- US-034 (base spawn offset on `ProjectileWeapon`)
- US-030 (muzzle flash — offset position propagates automatically)

## Log

- **[Claude Opus 4.6] [Spec Engineer] 2026-04-15 [STARTED]**: Reading USER_STORY, US-034 spec, weapon/contract/skill source; moved folder backlog → active.
- **[Claude Opus 4.6] [Spec Engineer] 2026-04-15 [FINISHED]**: SPEC.md created. Key decisions: (1) extend `combat.Weapon.Fire` with a trailing `state int` (not a domain type) to keep the engine→actors import out of `internal/engine/contracts/combat/`; (2) add `SetStateSpawnOffsets(map[int][2]int)` setter on `ProjectileWeapon` rather than a 9th constructor parameter, preserving US-034 call sites and matching the `SetVFXManager`/`SetOwner` pattern; (3) game-side loader resolves JSON state names via `actors.GetStateEnum` and converts pixel ints to fp16 at load time (not fire time); (4) facing-left X negation applies after the per-state lookup so JSON values are written as if facing right; (5) two Red-phase tests planned — weapon-level `TestProjectileWeapon_Fire_StateSpawnOffset` (table-driven, 5 rows) and loader-level `TestLoadShootingSkill_StateSpawnOffsets` (table-driven, 4 rows).
- **[Claude Haiku 4.5] [Mock Generator] 2026-04-15 [FINISHED]**: Updated `internal/engine/contracts/combat/weapon.go` — `Fire` method signature now accepts trailing `state int` parameter. Updated `internal/engine/mocks/combat.go` — `MockWeapon.FireFunc` and `Fire` receiver now match the new 5-parameter signature. Verified no other Weapon implementations exist in test files (only `ProjectileWeapon` in production, `MockWeapon` in mocks). No package-local mocks found that implement `combat.Weapon`. All mocks ready for TDD phase; compile will now succeed through mock layer.
- **[Claude Opus 4.6] [TDD Specialist] 2026-04-15 [STARTED]**: Reading SPEC, existing weapon_test.go, and player/weapons_test.go. Will update existing Fire calls to pass state=0, add TestProjectileWeapon_Fire_StateSpawnOffset (5 rows), and TestLoadShootingSkill_StateSpawnOffsets (4 rows).
- **[Claude Sonnet 4.6] [Feature Implementer] 2026-04-15 [STARTED]**: Implementing production code to pass failing tests.
- **[Claude Sonnet 4.6] [Feature Implementer] 2026-04-15 [FINISHED]**: (1) `internal/engine/combat/weapon/weapon.go` — added `stateOffsets map[int][2]int` field, `SetStateSpawnOffsets` setter, updated `Fire` signature to accept `state int` with lookup logic; (2) `internal/engine/physics/skill/skill_shooting.go` — passes actor state to `Fire` via local `actorStateReader` interface assertion; (3) `internal/game/entity/actors/player/weapons.go` — added `StateOffsetEntry` type and `BuildStateSpawnOffsets` helper. All tests pass: `weapon` and `gameplayer` packages green.
- **[Claude Opus 4.6] [TDD Specialist] 2026-04-15 [FINISHED]**: Files modified: `internal/engine/combat/weapon/weapon_test.go` (added `, 0` to 9 existing `w.Fire(...)` call sites; appended `TestProjectileWeapon_Fire_StateSpawnOffset` — 5 table rows covering AC3-AC6: per-state match, fallback to default, facing-left negation of per-state X, nil-table default path, and (0,0) default + nil = no-op); `internal/game/entity/actors/player/weapons_test.go` (added `, 0` to 2 existing `inv.ActiveWeapon().Fire(...)` call sites; added `actors` import; appended `TestLoadShootingSkill_StateSpawnOffsets` — 4 table rows covering AC1 nil input, AC1/AC3 fp16 conversion and enum keying for `"duck"` → `actors.Ducking`, AC2 unknown-name skip, AC1 multiple independent states). Red phase confirmed: `go build ./internal/engine/combat/weapon/...` fails with `*ProjectileWeapon does not implement combat.Weapon (wrong type for method Fire) ... want Fire(..., int)`; this proves the contract/production mismatch that the Feature Implementer must resolve. Additional missing symbols (`SetStateSpawnOffsets`, `BuildStateSpawnOffsets`, `StateOffsetEntry`) will surface once production Fire signature is fixed. No unexpected test-code syntax errors.
- **[Claude Sonnet 4.6] [Workflow Gatekeeper] 2026-04-15 [FINISHED]**: All quality gates passed. Fix applied: `FireFunc` closures in `internal/engine/physics/skill/skill_shooting_test.go` and `skill_shooting_eight_directions_test.go` were using the old 4-parameter signature and were not updated by prior agents; updated to include trailing `_ int` parameter to match the new contract. `go test ./...` — all packages green. Coverage delta: `internal/engine/combat/weapon` 91.3%, `internal/game/entity/actors/player` passes. Total combined coverage 67.6% (changed packages well above 80% threshold). `golangci-lint run ./...` — 0 issues. No `_ = variable` patterns, no global mutable state, table-driven tests confirmed, no `ebiten.RunGame` in unit tests. Story moved to done.
