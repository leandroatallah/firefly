# PROGRESS-023 — Projectile Manager

**Status:** ✅ Done

| Agent | Status |
|---|---|
| Spec Engineer | ✅ |
| Mock Generator | ✅ |
| TDD Specialist | ✅ |
| Developer | ✅ |
| Gatekeeper | ✅ |

## Log

Spec Engineer 2026-04-05: SPEC.md created. Key decisions: introduced `internal/engine/combat/projectile/` as the new bounded-context package; added a thin `contracts/projectile` interface to avoid circular imports; `PhasesScene.SpawnBullet` is kept as a delegating shim so `body.Shooter` callers require no changes; `gamestates.Bullet` is deleted and replaced by the manager's internal `projectile` type; `ProjectileManager` is initialised in `setup.go` after `Space` is wired.

Mock Generator 2026-04-05: No shared mocks required — `BodiesSpace` is only needed by `manager_test.go`. Created package-local `internal/engine/combat/projectile/mocks_test.go` with:
- `mockBodiesSpace` — implements `body.BodiesSpace`; exposes `AddBodyFunc`, `QueueForRemovalFunc`, `RemoveBodyFunc` for call tracking in tests.
- `mockTilemapDimensionsProvider` — implements `tilemaplayer.TilemapDimensionsProvider` with configurable `width`/`height` for out-of-bounds scenario.

TDD Specialist 2026-04-05: Created `internal/engine/combat/projectile/manager_test.go` with three table-driven test scenarios: (1) Spawn increases projectile count and calls AddBody, (2) Update removes out-of-bounds projectiles and calls QueueForRemoval, (3) Clear removes all projectiles and calls RemoveBody. Tests fail with "undefined: ProjectileConfig", "undefined: NewManager" — proving missing behavior, not just missing symbols.

Developer 2026-04-05: Implemented `projectile.Manager` and `projectile` internal type. Wired `ProjectileManager` into `AppContext` and `setup.go`. Refactored `PhasesScene` to delegate projectile lifecycle to the manager. Updated `PhasesScene` tests to work with the new architecture. Deleted `gamestates.Bullet`. Verified all tests pass.

## Log

Gatekeeper 2026-04-06: ✅ Approved.
- All 4 tests in `internal/engine/combat/projectile` pass (Spawn, Update_OutOfBounds, Clear, Draw).
- All `internal/game/scenes/phases` tests pass.
- Coverage delta: new package `internal/engine/combat/projectile` at **80.0%** (net new — previously 0%).
- `gamestates.Bullet` deleted; no stale bullet fields or loops remain in `PhasesScene`.
- `SpawnBullet` retained as a delegating shim — `body.Shooter` contract satisfied.
- `ProjectileManager` initialised in `setup.go` after `Space` — pre-condition met.
- No `_ = variable` in production code. Table-driven tests confirmed. Headless Ebitengine setup intact.
