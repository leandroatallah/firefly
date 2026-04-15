# SPEC — US-031: Impact VFX on Projectile Hit

**Branch:** `031-impact-vfx`
**Bounded Context:** `internal/engine/combat/projectile/`
**Depends on:** US-035 (VFX manager wiring on `projectile.Manager`), US-036 (config fields `ImpactEffect`, `DespawnEffect`)

## 1. Objective

Give every projectile collision visible feedback by spawning an impact VFX puff at the projectile's position immediately before the projectile is queued for removal. The feature must be opt-in: projectiles created without a VFX manager or without a configured effect type must behave exactly as they do today (no panic, no spawn).

## 2. Contracts (no new interfaces)

Reuse existing contracts — no new interface files required.

- `internal/engine/contracts/vfx/vfx.go` → `vfx.Manager`
  - Only method consumed in this story: `SpawnPuff(typeKey string, x float64, y float64, count int, randRange float64)`.
- `internal/engine/contracts/body/...` → `body.Collidable`, `body.BodiesSpace`, `body.Movable` (unchanged).

No global mutable state is introduced. The VFX manager is injected via `(*Manager).SetVFXManager(contractsvfx.Manager)` and propagated into each `projectile` at `Spawn` time.

## 3. Struct Changes

### 3.1 `projectile` (internal, `projectile.go`)

Add the following fields (already present in current code; spec locks them in):

| Field | Type | Purpose |
|---|---|---|
| `vfxManager`    | `contractsvfx.Manager` | Injected VFX manager; may be `nil`. |
| `impactEffect`  | `string`               | VFX type key used on `OnTouch` / `OnBlock`. Empty = disabled. |
| `despawnEffect` | `string`               | VFX type key used when leaving tilemap bounds (used by US-037; included for completeness). Empty = disabled. |

### 3.2 `Manager` (`manager.go`)

| Field | Type | Default |
|---|---|---|
| `vfxManager`    | `contractsvfx.Manager` | `nil` |
| `impactEffect`  | `string`               | `"bullet_impact"` |
| `despawnEffect` | `string`               | `"bullet_despawn"` |

- `SetVFXManager(v contractsvfx.Manager)` sets the manager.
- `Spawn(...)` copies `m.vfxManager`, `m.impactEffect`, `m.despawnEffect` into the newly created `projectile`.
- `ProjectileConfig.ImpactEffect` / `DespawnEffect` exist per US-036 but are not consumed by US-031 (manager-level defaults drive behavior today). This is noted to avoid regression when US-038 wires per-config overrides.

## 4. Behavior Specification

### 4.1 `OnTouch(other body.Collidable)`

Pre-conditions:
- `p.body` and `p.space` are non-nil (guaranteed by `Manager.Spawn`).

Logic:
```
if other == p.body.Owner():
    no-op                               # projectile does not hit its owner
else:
    p.spawnVFX(p.impactEffect)          # before removal
    p.space.QueueForRemoval(p.body)
```

Post-conditions (when `other != p.body.Owner()`):
- If `p.vfxManager != nil && p.impactEffect != ""`:
  - `SpawnPuff` was called **exactly once** with:
    - `typeKey  == p.impactEffect`
    - `x        == float64(x16) / 16.0`
    - `y        == float64(y16) / 16.0`
    - `count    == 1`
    - `randRange == 0.0`
- `QueueForRemoval(p.body)` was called exactly once.
- `SpawnPuff` is invoked **before** `QueueForRemoval` (ordering matters so position fetch happens on a still-valid body).

### 4.2 `OnBlock(other body.Collidable)`

Logic:
```
p.spawnVFX(p.impactEffect)
p.space.QueueForRemoval(p.body)
```

Post-conditions: identical to `OnTouch` (happy path), with no owner check — every block triggers the effect.

### 4.3 `spawnVFX(typeKey string)` helper

```
if p.vfxManager == nil || typeKey == "":
    return                              # AC5: backward compatible
x16, y16 := p.body.GetPosition16()
p.vfxManager.SpawnPuff(typeKey, float64(x16)/16.0, float64(y16)/16.0, 1, 0.0)
```

Guards (must all hold simultaneously to be safe):
- Nil-manager guard → no panic.
- Empty-key guard → no SpawnPuff call.
- Position fetched from `body.GetPosition16()` (fp16 units), converted via `float64(v) / 16.0` per constitution's fp16 rule (equivalent to `fp16.To16` semantics in this module).

## 5. Position Conversion

- Bodies store positions in fp16 (scale factor 16: 1 pixel = 16 units).
- Conversion to world-space pixel floats for VFX:
  - `worldX = float64(x16) / 16.0`
  - `worldY = float64(y16) / 16.0`
- No rounding, no truncation. Direct division preserves sub-pixel precision for particle spawn.

## 6. Nil / Empty-Config Safety (AC5)

Matrix:

| `vfxManager` | `impactEffect` | Expected |
|---|---|---|
| `nil`        | any            | no call, no panic |
| non-nil      | `""`           | no call, no panic |
| non-nil      | `"foo"`        | `SpawnPuff("foo", x, y, 1, 0.0)` |

`QueueForRemoval` runs regardless of VFX availability — removal logic must not be coupled to VFX presence.

## 7. Red Phase (failing test scenario)

Location: `internal/engine/combat/projectile/projectile_test.go` (new file) using the existing `mocks_test.go` doubles (`mockBodiesSpace`, `mockVFXManager`, `mockCollidable`). No new shared mocks are required; `internal/engine/mocks/mock_vfx_manager.go` already provides `MockVFXManager` for cross-package reuse.

Table-driven test `TestProjectile_ImpactVFX` cases:

| Case | vfxManager | impactEffect | Trigger | Expect SpawnPuff | Expect QueueForRemoval |
|---|---|---|---|---|---|
| `OnTouch spawns impact VFX`         | mock     | `"bullet_impact"` | `OnTouch(nonOwner)` at (32, 16) fp16 | `("bullet_impact", 2.0, 1.0, 1, 0.0)` | yes |
| `OnTouch with owner does nothing`   | mock     | `"bullet_impact"` | `OnTouch(owner)`                     | not called                            | no  |
| `OnBlock spawns impact VFX`         | mock     | `"bullet_impact"` | `OnBlock(other)` at (48, 24) fp16    | `("bullet_impact", 3.0, 1.5, 1, 0.0)` | yes |
| `nil manager is safe`               | `nil`    | `"bullet_impact"` | `OnTouch(nonOwner)`                  | not called                            | yes |
| `empty effect type is safe`         | mock     | `""`              | `OnTouch(nonOwner)`                  | not called                            | yes |
| `OnBlock nil manager is safe`       | `nil`    | `"bullet_impact"` | `OnBlock(other)`                     | not called                            | yes |

Additional deterministic assertions:
- `mockVFXManager.spawnPuffCallCount` increments exactly once in positive cases.
- `mockVFXManager.lastTypeKey`, `lastX`, `lastY` equal expected values.
- `mockBodiesSpace.queuedForRemoval` length matches the Expect column.
- No `time.Sleep`, no `ebiten.RunGame`, no GPU calls. The test constructs `projectile{}` directly with injected mocks.

Red step: writing these tests before `spawnVFX` helper / field wiring exists produces compile or assertion failures — the TDD Specialist authors them in that order.

## 8. Integration Points

- `projectile.Manager` is owned at the game-logic layer; a higher-level wiring step (outside this story) calls `SetVFXManager` once the `vfx.Manager` singleton is constructed.
- `physics/body.BodiesSpace.QueueForRemoval` ordering: VFX spawn happens first so the call site still has a valid fp16 position. `ProcessRemovals()` on the next frame finalizes the removal.
- Tilemap-bounds despawn in `Update()` uses `despawnEffect` (scope of US-037); this SPEC leaves that path untouched aside from the shared `spawnVFX` helper.

## 9. Out of Scope

- Per-`ProjectileConfig` override of effect types (US-038).
- Despawn/off-screen VFX behavior beyond reusing the `spawnVFX` helper (US-037).
- Particle authoring / art registration (`SpawnPuff` resolution by `typeKey` is the VFX package's responsibility).

## 10. Acceptance Mapping

| AC | Covered by |
|---|---|
| AC1 | §3.1 struct fields |
| AC2 | §4.1 `OnTouch` + §4.3 helper |
| AC3 | §4.2 `OnBlock` + §4.3 helper |
| AC4 | §5 position conversion |
| AC5 | §4.3 guards + §6 matrix |
| AC6 | §7 Red Phase table-driven tests |
