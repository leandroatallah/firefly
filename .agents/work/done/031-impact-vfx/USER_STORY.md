# US-031 — Impact VFX on Projectile Hit

**Branch:** `031-impact-vfx`
**Bounded Context:** Engine (`internal/engine/combat/projectile/`)

## Story

As a game developer using this boilerplate,
I want impact VFX to spawn when projectiles hit targets or blocks,
so that collisions have visual feedback at the impact position.

## Context

The `projectile` struct in `internal/engine/combat/projectile/projectile.go` handles collisions via `OnTouch()` and `OnBlock()` but provides no visual feedback. This story adds impact VFX spawning before projectile removal.

## Acceptance Criteria

- **AC1** — `projectile` struct stores `impactEffectType string` and `vfxManager vfx.Manager` fields.
- **AC2** — `OnTouch()` spawns impact VFX at projectile position before queuing removal (if VFX configured).
- **AC3** — `OnBlock()` spawns impact VFX at projectile position before queuing removal (if VFX configured).
- **AC4** — Position converted from fp16 to float64 using `float64(x16) / 16.0`.
- **AC5** — No error when VFX manager is nil or impact effect type is empty (backward compatible).
- **AC6** — Unit tests verify `SpawnPuff` called with correct position and effect type.

## Notes

- Package path: `internal/engine/combat/projectile/`
- Depends on US-035 (VFX manager in Manager) and US-036 (config fields)
- Impact effect type passed from `Manager.Spawn()` using config
