# US-032 — Projectile Lifetime and Despawn VFX

**Branch:** `032-projectile-lifetime-despawn`
**Bounded Context:** Engine (`internal/engine/combat/projectile/`)

## Story

As a game developer using this boilerplate,
I want projectiles to despawn after a configurable lifetime,
so that weapons have limited range and despawn events can show optional visual feedback.

## Context

The `projectile` struct in `internal/engine/combat/projectile/projectile.go` currently lives until it hits something or exits bounds. This story adds a frame-based lifetime system with optional despawn VFX at the projectile's last position.

## Acceptance Criteria

- **AC1** — `ProjectileConfig` struct has `LifetimeFrames int` field (0 = infinite, backward compatible).
- **AC2** — `projectile` struct has `lifetimeFrames int` and `currentLifetime int` fields.
- **AC3** — `Update()` decrements `currentLifetime` each frame.
- **AC4** — When `currentLifetime` reaches 0, projectile queues for removal.
- **AC5** — Optional despawn VFX spawned at projectile position before removal (if configured).
- **AC6** — `despawnEffectType string` field determines VFX type to spawn.
- **AC7** — Position converted from fp16 to float64 using `float64(x16) / 16.0`.
- **AC8** — No error when VFX manager is nil or despawn effect type is empty.
- **AC9** — Unit tests verify projectile despawns after lifetime expires with correct VFX.

## Notes

- Package path: `internal/engine/combat/projectile/`
- Depends on US-035 (VFX manager in Manager) and US-036 (VFX config fields)
- Lifetime and despawn effect type passed from `Manager.Spawn()` using config
- Out-of-bounds check remains as safety fallback (no VFX)
- Includes `LifetimeFrames` config field (not in US-036)
