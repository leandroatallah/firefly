# US-035 — VFX Manager Integration for Projectile Manager

**Branch:** `035-vfx-projectile-manager`
**Bounded Context:** Engine / Combat (`internal/engine/combat/projectile/`)

## Story

As a game developer using this boilerplate,
I want the projectile manager to accept a VFX manager dependency,
so that projectiles can spawn visual effects for impacts and despawn events.

## Context

The `Manager` in `internal/engine/combat/projectile/manager.go` now integrates with the VFX system via the `contractsvfx.Manager` interface (`internal/engine/contracts/vfx/`). The VFX manager is injected via a setter to maintain backward compatibility. Projectiles use `SpawnPuff` on the VFX contract to emit effects at impact (collision via `OnTouch`/`OnBlock`) and despawn (out-of-bounds removal in `Update`).

## Acceptance Criteria

- **AC1** — `Manager` struct has a `vfxManager` field typed to the `contracts/vfx.Manager` interface.
- **AC2** — `SetVFXManager(v contractsvfx.Manager)` setter method added to `Manager`.
- **AC3** — `Manager.Spawn()` passes the VFX manager reference and effect name strings to each new `projectile` instance.
- **AC4** — All existing functionality works when VFX manager is nil (backward compatible). The `projectile.spawnVFX` method guards against nil manager and empty type key.
- **AC5** — Unit tests verify: (a) `SetVFXManager` stores the manager on `Manager`, (b) spawned projectiles receive the VFX manager, (c) nil VFX manager does not panic during impact/despawn.
- **AC6** — No changes to existing `SpawnProjectile` signature.
- **AC7** — `Manager` holds default effect name strings (`impactEffect`, `despawnEffect`) that are forwarded to each spawned projectile.
- **AC8** — `ProjectileConfig` includes `ImpactEffect` and `DespawnEffect` fields for per-projectile-type VFX overrides.

## Notes

- Package path: `internal/engine/combat/projectile/`
- VFX contract: `internal/engine/contracts/vfx/Manager`
- VFX manager is an optional dependency injected via setter
- Projectiles call `vfxManager.SpawnPuff(typeKey, x, y, count, randRange)` with fp16-to-float64 coordinate conversion
- Effect is triggered in `OnTouch` (collision with non-owner), `OnBlock` (wall collision), and `Update` (out-of-bounds despawn)
- Should be implemented after US-033 (particle definitions)
- **Known gap:** `ProjectileConfig.ImpactEffect` and `DespawnEffect` fields exist but are not read during `Spawn()` — the manager's hardcoded defaults are always used. A follow-up story should wire config overrides into spawn logic.
