# US-022 — Weapon System

**Branch:** `022-weapon-system`
**Bounded Context:** Engine (`internal/engine/combat/`)

## Story

As a game developer using this boilerplate,
I want a flexible weapon system decoupled from the shooting skill and scene,
so that I can define multiple weapon types (projectile, melee, charge) without modifying engine code.

## Context

The current `ShootingSkill` is tightly coupled to `body.Shooter` (a scene dependency) and supports only one bullet type. This prevents weapon switching, charge shots, and reuse across scenes. This story introduces `internal/engine/combat/weapon/` as the foundation for all combat.

## Acceptance Criteria

- **AC1** — `Weapon` interface defined with `ID()`, `Fire()`, `CanFire()`, `Update()`, `Cooldown()`, `SetCooldown()`.
- **AC2** — `ProjectileWeapon` struct implements `Weapon`, managing its own cooldown in frames.
- **AC3** — `ProjectileWeapon` delegates projectile spawning to a `projectile.Manager` (injected dependency).
- **AC4** — Weapon JSON schema defined: `id`, `type`, `cooldown_frames`, `projectile` sub-object (`speed`, `damage`, `behavior`).
- **AC5** — `weapon.Factory` creates a `Weapon` from JSON config.
- **AC6** — No scene types are imported by any package under `internal/engine/combat/`.
- **AC7** — Unit tests cover: fire, cooldown decrement, `CanFire` gating, factory instantiation from config.

## Notes

- Package path: `internal/engine/combat/weapon/`
- `ProjectileWeapon` holds `projectileType string` (resolved by manager at spawn time).
- Melee and charge weapon types are out of scope for this story.
