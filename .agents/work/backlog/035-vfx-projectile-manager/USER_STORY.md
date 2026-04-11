# US-035 — VFX Manager Integration for Projectile Manager

**Branch:** `035-vfx-projectile-manager`
**Bounded Context:** Engine (`internal/engine/combat/projectile/`)

## Story

As a game developer using this boilerplate,
I want the projectile manager to accept a VFX manager dependency,
so that projectiles can spawn visual effects for impacts and despawn events.

## Context

The current `Manager` in `internal/engine/combat/projectile/manager.go` has no VFX integration. This story adds VFX manager injection to enable projectiles to spawn visual effects while maintaining backward compatibility.

## Acceptance Criteria

- **AC1** — `Manager` struct has `vfxManager vfx.Manager` field.
- **AC2** — `SetVFXManager(manager vfx.Manager)` method added to `Manager`.
- **AC3** — `Manager` passes VFX manager to projectiles during spawn in `Spawn()` method.
- **AC4** — All existing functionality works when VFX manager is nil (backward compatible).
- **AC5** — Unit tests verify setter stores manager and handles nil gracefully.
- **AC6** — No changes to existing `SpawnProjectile` signature.

## Notes

- Package path: `internal/engine/combat/projectile/`
- VFX manager is optional dependency injected via setter
- Projectiles receive VFX manager reference for use in collision callbacks and lifetime despawn
- Should be implemented after US-033 (particle definitions)
