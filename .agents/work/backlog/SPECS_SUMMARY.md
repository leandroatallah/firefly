# Combat VFX Specs — Summary

**Date:** 2026-04-08  
**Status:** All specs completed, stories remain in backlog

## Specs Created

### US-033 — Combat Particle Definitions
- **File:** `.agents/work/backlog/033-combat-particle-definitions/SPEC.md`
- **Changes:** Add 3 particle types to `vfx.json` (muzzle_flash, bullet_impact, bullet_despawn)
- **Tech:** Pixel-based particles, 1-bit colors (#000000, #FFFFFF)
- **Tests:** JSON validation test

### US-036 — VFX Configuration for Projectiles
- **File:** `.agents/work/backlog/036-vfx-combat-config/SPEC.md`
- **Changes:** Add `ImpactEffect` and `DespawnEffect` fields to `ProjectileConfig`
- **Tech:** String fields for particle type keys
- **Tests:** Struct field access, JSON marshaling

### US-035 — VFX Manager Integration
- **File:** `.agents/work/backlog/035-vfx-projectile-manager/SPEC.md`
- **Changes:** Add `vfxManager` field and `SetVFXManager()` to projectile `Manager`
- **Tech:** Dependency injection via setter, pass to projectiles in `Spawn()`
- **Tests:** Setter functionality, nil safety

### US-030 — Muzzle Flash VFX
- **File:** `.agents/work/backlog/030-muzzle-flash-vfx/SPEC.md`
- **Changes:** Add `muzzleEffectType` parameter to constructor, `SetVFXManager()` method, VFX spawn in `Fire()`
- **Tech:** Constructor signature: 5 → 6 parameters, fp16 to float64 conversion
- **Tests:** VFX spawning with correct position, nil manager safety

### US-034 — Projectile Spawn Offset
- **File:** `.agents/work/backlog/034-projectile-spawn-offset/SPEC.md`
- **Changes:** Add `spawnOffsetX16` and `spawnOffsetY16` parameters to constructor, apply in `Fire()`
- **Tech:** Constructor signature: 6 → 8 parameters, X offset negated when facing left
- **Tests:** Offset applied for both facing directions, zero offset backward compatible

### US-031 — Impact VFX
- **File:** `.agents/work/backlog/031-impact-vfx/SPEC.md`
- **Changes:** Add `impactEffectType` and `vfxManager` to projectile, spawn VFX in `OnTouch()` and `OnBlock()`
- **Tech:** VFX spawned before removal, count=3, randRange=1.0
- **Tests:** VFX on collision, nil manager safety

### US-032 — Projectile Lifetime and Despawn VFX
- **File:** `.agents/work/backlog/032-out-of-bounds-vfx/SPEC.md`
- **Changes:** Add `LifetimeFrames` to config, `lifetimeFrames`/`currentLifetime` to projectile, despawn VFX in `Update()`
- **Tech:** Frame-based countdown, VFX on expiration (count=5, randRange=1.5), 0 = infinite
- **Tests:** Lifetime expiration, infinite lifetime, out-of-bounds no VFX

## Implementation Order

1. **US-033** — Particle definitions (foundation)
2. **US-036** — Config fields (no logic)
3. **US-035** — VFX manager integration (infrastructure)
4. **US-030** — Muzzle flash (weapon VFX)
5. **US-034** — Spawn offset (weapon enhancement)
6. **US-031** — Impact VFX (projectile collision)
7. **US-032** — Lifetime + despawn (projectile lifecycle)

## Key Design Decisions

### Minimal Code Approach
- All specs focus on absolute minimal changes
- No verbose implementations
- Direct, focused modifications

### Backward Compatibility
- Empty effect types = no VFX
- Nil VFX manager = no VFX
- Lifetime 0 = infinite (existing behavior)
- Spawn offset (0, 0) = no offset

### Position Conversion Standard
- All specs use: `float64(x16) / 16.0` for fp16 to float64 conversion
- Consistent across all VFX spawning

### VFX Parameters
- Muzzle flash: count=1, randRange=0.0 (precise)
- Impact: count=3, randRange=1.0 (small burst)
- Despawn: count=5, randRange=1.5 (larger spread)

### Constructor Evolution
- US-030: 5 → 6 parameters (adds muzzleEffectType)
- US-034: 6 → 8 parameters (adds spawnOffsetX16, spawnOffsetY16)

## Red Phase Strategy

All specs define failing test scenarios:
- **US-033:** Particle types not found in JSON
- **US-036:** Struct fields don't exist
- **US-035:** Setter method doesn't exist
- **US-030:** Constructor signature mismatch
- **US-034:** Offset not applied
- **US-031:** VFX not spawned on collision
- **US-032:** Lifetime logic not implemented

## Next Steps

Ready for Mock Generator to create test mocks for:
- `mockVFXManager` (implements `vfx.Manager`)
- `mockProjectileManager` (implements `combat.ProjectileManager`)
- `mockBodiesSpace` (implements `body.BodiesSpace`)
- `mockCollidableBody` (implements `body.Collidable`)
