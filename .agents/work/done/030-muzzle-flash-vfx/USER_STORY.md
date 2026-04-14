# US-030 — Muzzle Flash VFX on Weapon Fire

**Branch:** `030-muzzle-flash-vfx`
**Bounded Context:** Engine (`internal/engine/combat/weapon/`)

## Story

As a game developer using this boilerplate,
I want muzzle flash VFX to spawn when a weapon fires,
so that weapon firing has immediate visual feedback at the firing position.

## Context

The `ProjectileWeapon` in `internal/engine/combat/weapon/weapon.go` currently spawns projectiles without visual effects. This story adds muzzle flash VFX spawning when `Fire()` is called.

## Acceptance Criteria

- **AC1** — `ProjectileWeapon` accepts `muzzleEffectType string` parameter in constructor.
- **AC2** — `ProjectileWeapon` has `SetVFXManager(manager vfx.Manager)` method for dependency injection.
- **AC3** — `Fire()` spawns muzzle flash VFX before spawning projectile (if VFX configured).
- **AC4** — Position converted from fp16 to float64 using `float64(x16) / 16.0`.
- **AC5** — No error when `vfxManager` is nil or `muzzleEffectType` is empty (backward compatible).
- **AC6** — Unit tests verify `SpawnPuff` called with correct position, effect type, count=1, randRange=0.0.

## Notes

- Package: `internal/engine/combat/weapon/`
- VFX manager injected via setter to avoid circular dependencies
- Muzzle effect type comes from weapon configuration
