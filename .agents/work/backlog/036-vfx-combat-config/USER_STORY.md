# US-036 — VFX Configuration for Projectiles

**Branch:** `036-vfx-projectile-config`
**Bounded Context:** Engine (`internal/engine/combat/projectile/`)

## Story

As a game developer using this boilerplate,
I want to configure VFX effects in projectile configs,
so that different projectile types can have unique visual feedback for impacts and despawn events.

## Context

The current `ProjectileConfig` in `internal/engine/combat/projectile/config.go` only defines Width, Height, and Damage. This story extends it with optional VFX effect type fields.

## Acceptance Criteria

- **AC1** — `ProjectileConfig` struct has `ImpactEffect string` field for hit impact VFX type.
- **AC2** — `ProjectileConfig` struct has `DespawnEffect string` field for lifetime expiration VFX type.
- **AC3** — All VFX fields are optional (empty string means no effect).
- **AC4** — JSON marshaling/unmarshaling works with new fields.
- **AC5** — Unit tests verify config struct with VFX fields.
- **AC6** — Existing code using `ProjectileConfig` continues to work (backward compatible).

## Notes

- Package path: `internal/engine/combat/projectile/`
- VFX effect types are string identifiers referencing particle definitions in `assets/particles/vfx.json`
- Muzzle flash is configured in weapon config (US-030), not projectile config
- Lifetime configuration added in US-032, not here
