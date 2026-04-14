# Combat VFX Stories — Implementation Order

## Overview

Stories for implementing combat visual effects system with proper dependency ordering.

## Recommended Implementation Order

### Phase 1: Foundation (No Dependencies)

**1. US-033 — Combat Particle Definitions**
- **Why first:** All other stories reference these particle types
- **What:** Add `muzzle_flash`, `bullet_impact`, `bullet_despawn` to `vfx.json`
- **Tech:** Pixel-based particles, 1-bit colors (black/white only)
- **Blocking:** US-030, US-031, US-032

**2. US-036 — VFX Configuration for Projectiles**
- **Why second:** Config fields needed before implementation
- **What:** Add `ImpactEffect` and `DespawnEffect` string fields to `ProjectileConfig`
- **Tech:** Struct fields only, no logic
- **Blocking:** US-031, US-032

**3. US-035 — VFX Manager Integration**
- **Why third:** Foundation for projectile VFX
- **What:** Add `vfxManager` field and `SetVFXManager()` to projectile `Manager`
- **Tech:** Dependency injection via setter
- **Blocking:** US-031, US-032

### Phase 2: Weapon VFX (Independent)

**4. US-030 — Muzzle Flash VFX**
- **Depends on:** US-033 (particles)
- **What:** Add muzzle flash VFX to `ProjectileWeapon.Fire()`
- **Tech:** Add `muzzleEffectType` field, `SetVFXManager()` method, spawn VFX before projectile
- **Note:** Constructor signature changes (adds 6th parameter)

**5. US-034 — Projectile Spawn Offset**
- **Depends on:** US-030 (so VFX spawns at correct position)
- **What:** Add spawn offset configuration to align projectiles with sprite
- **Tech:** Add `spawnOffsetX16`, `spawnOffsetY16` fields, apply in `Fire()`
- **Resolves:** TODO in `skill_shooting.go`

### Phase 3: Projectile VFX (Dependent)

**6. US-031 — Impact VFX**
- **Depends on:** US-033, US-035, US-036
- **What:** Add impact VFX to `projectile.OnTouch()` and `OnBlock()`
- **Tech:** Add `impactEffectType` and `vfxManager` fields, spawn VFX before removal

**7. US-032 — Projectile Lifetime and Despawn VFX**
- **Depends on:** US-033, US-035, US-036
- **What:** Add frame-based lifetime system with despawn VFX
- **Tech:** Add `LifetimeFrames` to config, `lifetimeFrames`/`currentLifetime` to projectile, spawn VFX on expiration
- **Note:** Includes config field (not in US-036)

## Story Changes Summary

### US-033 — Updated
- Specified pixel-based particles (no images)
- Limited to 1-bit colors (black/white)
- Marked as foundation story (implement first)

### US-034 — Replaced
- **Old:** Combat VFX integration tests (redundant)
- **New:** Projectile spawn offset configuration
- **Reason:** Addresses TODO in codebase, more valuable than redundant tests

### US-036 — Simplified
- Removed `LifetimeFrames` field (moved to US-032)
- Now only VFX config fields
- Breaks circular dependency with US-032

### US-032 — Expanded
- Added `LifetimeFrames` config field (AC1)
- Now owns complete lifetime feature including config
- No longer depends on US-036 for lifetime config

### US-035 — Clarified
- Added note about implementation order (after US-033)
- Clarified projectiles receive VFX manager for both collision and lifetime events

## Dependency Graph

```
US-033 (particles)
  ├─→ US-030 (muzzle flash)
  │     └─→ US-034 (spawn offset)
  ├─→ US-031 (impact VFX) ←─┐
  └─→ US-032 (lifetime VFX) ←┤
                             │
US-036 (config) ─────────────┤
US-035 (VFX manager) ────────┘
```

## Breaking Changes

### US-030
- `NewProjectileWeapon()` signature changes from 5 to 6 parameters
- Adds `muzzleEffectType string` parameter

### US-034
- `NewProjectileWeapon()` signature changes from 6 to 8 parameters
- Adds `spawnOffsetX16 int` and `spawnOffsetY16 int` parameters

## Backward Compatibility

All stories maintain backward compatibility through:
- Optional VFX (nil manager or empty effect type = no VFX)
- Lifetime of 0 = infinite (existing behavior)
- Spawn offset of (0, 0) = no offset (existing behavior)

## Testing Strategy

Each story includes unit tests (AC in each story):
- US-030: Muzzle flash spawning with correct position
- US-031: Impact VFX on collision
- US-032: Despawn VFX after lifetime
- US-033: JSON validation
- US-034: Offset applied correctly for both facing directions
- US-035: VFX manager injection
- US-036: Config struct with new fields

No separate integration test story needed (covered by unit tests).
