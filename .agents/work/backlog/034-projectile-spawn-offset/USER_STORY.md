# US-034 — Projectile Spawn Offset Configuration

**Branch:** `034-projectile-spawn-offset`
**Bounded Context:** Engine (`internal/engine/combat/weapon/`)

## Story

As a game developer using this boilerplate,
I want to configure spawn offset for projectiles,
so that projectiles spawn aligned with the weapon sprite position instead of the entity's origin.

## Context

The `ShootingSkill` in `internal/engine/physics/skill/skill_shooting.go` has a TODO comment about needing an offset to align shooting with the sprite. Currently, projectiles spawn at the entity's position with only a basic width adjustment for facing direction.

## Acceptance Criteria

- **AC1** — `ProjectileWeapon` accepts `spawnOffsetX16 int` and `spawnOffsetY16 int` in constructor (fp16 units).
- **AC2** — `Fire()` method applies offset to spawn position: `spawnX16 = x16 + offsetX16`, `spawnY16 = y16 + offsetY16`.
- **AC3** — Offset respects facing direction: when facing left, X offset is negated.
- **AC4** — Offset of (0, 0) maintains current behavior (backward compatible).
- **AC5** — Unit tests verify offset applied correctly for both facing directions.
- **AC6** — Muzzle flash VFX (if configured) spawns at offset position, not entity origin.

## Notes

- Package: `internal/engine/combat/weapon/`
- Offsets in fp16 units (scale factor 16: 1 pixel = 16 units)
- Resolves TODO in `skill_shooting.go`
- Should be implemented after US-030 (muzzle flash) to ensure VFX spawns at correct position
