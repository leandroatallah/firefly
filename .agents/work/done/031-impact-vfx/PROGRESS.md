# PROGRESS — US-031

## Status: Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper ✅

## Notes

## Log

- 2026-04-08: Spec Engineer — SPEC.md created. Adds impactEffectType and vfxManager fields, spawns VFX in OnTouch() and OnBlock().
- 2026-04-15: Spec Engineer — SPEC.md rewritten against current code. Key decisions: reuse existing `contracts/vfx.Manager` (no new interface), reuse existing `internal/engine/mocks/MockVFXManager` plus package-local `mockVFXManager` in `mocks_test.go`, centralize nil/empty guards in private `spawnVFX` helper, VFX spawn must precede `QueueForRemoval` to preserve valid body position, per-config effect override deferred to US-038.
- 2026-04-15: TDD Specialist — Authored table-driven `TestProjectile_ImpactVFX` in `internal/engine/combat/projectile/projectile_test.go` covering SPEC §7's six cases (OnTouch hit, OnTouch owner no-op, OnBlock hit, nil manager on OnTouch, empty effect key, nil manager on OnBlock). Extended `mocks_test.go` `mockCollidable` with configurable fp16 position (`x16`/`y16` via `SetPosition16`/`GetPosition16`) and `mockVFXManager` with `lastCount`/`lastRandRange` capture so the suite asserts SpawnPuff is called exactly once with `(typeKey, x16/16.0, y16/16.0, 1, 0.0)` and that `QueueForRemoval` runs regardless of VFX availability. Red proof: the suite fails for any implementation missing the `spawnVFX` helper, the `vfxManager`/`impactEffect` fields, the owner-skip branch in `OnTouch`, or the nil/empty guards — i.e. it fails on behavior, not just missing symbols.
- 2026-04-15: Feature Implementer — Verified that the current implementation of `projectile.go` and `manager.go` exactly matches SPEC US-031. `projectile` struct has the required fields, `OnTouch` and `OnBlock` trigger `spawnVFX` before `QueueForRemoval`, and `spawnVFX` correctly converts fp16 positions to world-space floats with proper nil/empty guards. All tests in `internal/engine/combat/projectile/` passed successfully.
- 2026-04-15: Workflow Gatekeeper ✅ — All checks passed. Coverage: 95.8% for `internal/engine/combat/projectile/` (positive delta; `OnTouch`, `OnBlock`, `spawnVFX` all at 100%). All 6 `TestProjectile_ImpactVFX` table-driven cases pass. Full `go test ./...` clean. `golangci-lint run` reports 0 issues. Implementation matches SPEC §3–§6 exactly. Story moved to done.
