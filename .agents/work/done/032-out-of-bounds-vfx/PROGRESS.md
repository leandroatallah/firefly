# PROGRESS — US-032

## Status: Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Notes

Story rewritten from out-of-bounds VFX to projectile lifetime system with despawn VFX.

## Log

- 2026-04-08: Story rewritten. Changed from out-of-bounds VFX (off-screen, not visible) to frame-based lifetime system with optional despawn VFX at last visible position.
- Spec Engineer 2026-04-15: SPEC.md created. Key decisions: (1) OOB path becomes silent (no VFX) — behavioral change from current code which emits despawnEffect on OOB; lifetime expiration is the only VFX-emitting despawn path. (2) Negative `LifetimeFrames` clamped to 0 (infinite) rather than triggering immediate despawn. (3) Lifetime tick runs after position update + collision resolution so VFX position reflects the frame's final fp16 coords. (4) No new contracts introduced — reuses existing `contractsvfx.Manager.SpawnPuff` and `BodiesSpace.QueueForRemoval`. (5) VFX mock is package-local in `mocks_test.go` (not shared).
- TDD Specialist 2026-04-15: Red Phase committed. New failing tests in `internal/engine/combat/projectile/projectile_test.go` (TestProjectile_LifetimeDespawn, TestProjectile_OOBHasNoVFX, TestProjectileConfig_LifetimeFrames_Default, TestManager_Spawn_PropagatesLifetime) and extended `internal/engine/combat/projectile/config_test.go` (TestProjectileConfig_LifetimeFramesField, TestProjectileConfig_LifetimeFramesJSON). Added `spawnPuffCalls []spawnPuffCall` recorder to `mockVFXManager` in `mocks_test.go`. Red proof: `go test ./internal/engine/combat/projectile/` fails to build because `ProjectileConfig.LifetimeFrames` does not yet exist — proving both the missing field (AC1) and, once the field lands, the downstream behavioral expectations (lifetime tick, silent OOB, propagation from Spawn) remain to be implemented.
- Feature Implementer 2026-04-15: Green Phase completed. All 28 tests in `internal/engine/combat/projectile/` now PASS.
  - `internal/engine/combat/projectile/config.go`: Added `LifetimeFrames int` field with `json:"lifetime_frames,omitempty"` tag.
  - `internal/engine/combat/projectile/projectile.go`: Added `lifetimeFrames` and `currentLifetime` fields to `projectile` struct. Inserted lifetime tick logic in `Update()` after `ResolveCollisions` and before OOB check. Removed `spawnVFX` call from OOB path (now silent).
  - `internal/engine/combat/projectile/manager.go`: Updated `Spawn` to propagate `LifetimeFrames` (clamped via `max(val, 0)`) and conditionally apply VFX field fallbacks (only when `LifetimeFrames == 0`, preserving explicit empty values for lifetime configs).
- Workflow Gatekeeper 2026-04-15: All quality gates passed. ✅
  - Spec compliance: all 9 ACs verified against implementation (AC1–AC9 satisfied).
  - TDD cycle: Red phase documented (build failure on missing field); Green phase completed with 28 passing tests.
  - `go test ./...`: all packages pass.
  - Coverage delta: `internal/engine/combat/projectile/` at 96.4% (well above 80% threshold). Positive delta confirmed.
  - `golangci-lint run ./...`: 0 issues.
  - Story folder moved to `done/`.
