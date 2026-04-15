# SPEC — US-032: Projectile Lifetime and Despawn VFX

**Branch:** `032-projectile-lifetime-despawn`
**Bounded Context:** `internal/engine/combat/projectile/`
**Depends on:** US-035 (VFX manager in `Manager`), US-036 (VFX config fields)

## 1. Goal

Add a frame-counted lifetime to projectiles. When a projectile's lifetime reaches zero, it queues for removal and optionally spawns a despawn VFX at its last fp16 position. The existing out-of-bounds check remains a silent safety fallback (no VFX).

## 2. Technical Requirements

### 2.1 `ProjectileConfig` (`config.go`)

Add one field:

```go
type ProjectileConfig struct {
    Width         int
    Height        int
    Damage        int
    ImpactEffect  string `json:"impact_effect,omitempty"`
    DespawnEffect string `json:"despawn_effect,omitempty"`
    LifetimeFrames int   `json:"lifetime_frames,omitempty"` // 0 = infinite (backward compat)
}
```

Semantics:
- `LifetimeFrames == 0` → lifetime disabled; projectile lives until collision or out-of-bounds (current behavior preserved).
- `LifetimeFrames > 0` → projectile despawns after N frames.
- Negative values are treated as 0 (infinite) and MUST NOT cause immediate removal.

### 2.2 `projectile` struct (`projectile.go`)

Add two fields (keep existing field order / visibility — all unexported):

```go
type projectile struct {
    movable         contractsbody.Movable
    body            contractsbody.Collidable
    space           contractsbody.BodiesSpace
    speedX16        int
    speedY16        int
    vfxManager      contractsvfx.Manager
    impactEffect    string
    despawnEffect   string
    lifetimeFrames  int // configured total lifetime (0 = infinite)
    currentLifetime int // frames remaining; only meaningful when lifetimeFrames > 0
}
```

Initialization contract (in `Manager.Spawn`):
- `lifetimeFrames = config.LifetimeFrames` (clamped to `>= 0`).
- `currentLifetime = lifetimeFrames` at spawn time.

### 2.3 `Manager.Spawn` changes (`manager.go`)

When constructing `p`, populate the two new fields from `config.LifetimeFrames`. No changes to public method signatures. `SpawnProjectile` (interface method) continues to use the minimal default config with `LifetimeFrames: 0` (infinite).

### 2.4 `projectile.Update()` behavior

New logic, in this order each frame:

1. Advance position (`x += speedX16`, `y += speedY16`, `SetPosition16`).
2. `space.ResolveCollisions(p.body)`.
3. **Lifetime tick (new)**: if `p.lifetimeFrames > 0`:
   - `p.currentLifetime--`
   - If `p.currentLifetime <= 0`:
     - Call `p.spawnVFX(p.despawnEffect)` at current fp16 position.
     - Call `p.space.QueueForRemoval(p.body)`.
     - `return` (skip out-of-bounds check — already removed).
4. Out-of-bounds check (unchanged). Note: current code calls `spawnVFX(p.despawnEffect)` on OOB; per this story, OOB is a **silent safety fallback** with **no VFX**. Update the OOB branch to call `space.QueueForRemoval(p.body)` only (remove the `spawnVFX` call from the OOB path).

### 2.5 Contracts

- `contractsvfx.Manager` (already exists) — `SpawnPuff(typeKey string, x, y float64, count int, randRange float64)` is the boundary.
- `contractsbody.BodiesSpace.QueueForRemoval(b Collidable)` — used for removal.
- No new contracts are introduced.

## 3. Pre-/Post-Conditions

### `Manager.Spawn(cfg, x16, y16, vx16, vy16, owner)`
- **Pre:** `cfg` is a `ProjectileConfig`. `space != nil`. `cfg.LifetimeFrames >= 0` (negatives clamped to 0).
- **Post:** A new `projectile` is registered in `space` with `lifetimeFrames == max(cfg.LifetimeFrames, 0)` and `currentLifetime == lifetimeFrames`.

### `projectile.Update()`
- **Pre:** projectile has been spawned; `body` is in `space`.
- **Post (lifetime path, `lifetimeFrames > 0`):** After `lifetimeFrames` invocations of `Update()` total, `currentLifetime == 0`, body is queued for removal exactly once, and — if `vfxManager != nil` and `despawnEffect != ""` — exactly one call to `SpawnPuff(despawnEffect, float64(x16)/16.0, float64(y16)/16.0, 1, 0.0)` has been emitted at the projectile's last position.
- **Post (infinite path, `lifetimeFrames == 0`):** `currentLifetime` is never decremented; lifecycle is identical to the pre-change implementation (collision + OOB fallback only).
- **Post (OOB path):** Body is queued for removal; **no** VFX is emitted.
- **Invariant:** `QueueForRemoval` for the same body is issued at most once per frame.

### `projectile.spawnVFX(typeKey)` (existing helper, reused)
- **Pre:** called with any string.
- **Post:** no-op when `vfxManager == nil` or `typeKey == ""`. Otherwise emits exactly one `SpawnPuff` call at the projectile's current fp16 position converted to float64 pixels.

## 4. Integration Points

- **Physics boundary:** `contractsbody.BodiesSpace.QueueForRemoval` — unchanged usage.
- **VFX boundary:** `contractsvfx.Manager.SpawnPuff` — only call site affected is `projectile.spawnVFX`.
- **Config boundary:** `ProjectileConfig` gains `LifetimeFrames`. JSON tag `lifetime_frames,omitempty` keeps on-disk configs backward compatible.
- **Other callers:** `SpawnProjectile(projectileType, ...)` still passes the minimal default config (infinite lifetime) — preserves current behavior for callers that do not yet opt in.

## 5. Acceptance Criteria Mapping

| AC | Satisfied by |
|---|---|
| AC1 | §2.1 — `LifetimeFrames int` added with `omitempty`; 0 = infinite |
| AC2 | §2.2 — `lifetimeFrames` and `currentLifetime` added to `projectile` |
| AC3 | §2.4 step 3 — `currentLifetime--` each frame when active |
| AC4 | §2.4 step 3 — `QueueForRemoval` when `currentLifetime <= 0` |
| AC5 | §2.4 step 3 — `spawnVFX(despawnEffect)` before `QueueForRemoval` |
| AC6 | §2.2 — `despawnEffect` field already exists and is reused as the VFX key |
| AC7 | §2.5 — `spawnVFX` uses `float64(x16)/16.0`, `float64(y16)/16.0` |
| AC8 | §2.5 — early return in `spawnVFX` when `vfxManager == nil` or `typeKey == ""` |
| AC9 | §6 Red Phase — table-driven tests |

## 6. Red Phase (Failing Test Scenarios)

All tests live in `internal/engine/combat/projectile/projectile_test.go` (new) and `config_test.go` (extend). Tests are table-driven and use the existing `mocks_test.go` patterns plus a local VFX mock (`mockVFXManager`) recording `SpawnPuff` calls — keep it in `mocks_test.go` (package-local).

### 6.1 `TestProjectile_LifetimeDespawn`

Table with rows:

| name | lifetimeFrames | despawnEffect | vfxManagerNil | updates | wantQueued | wantVFXCalls | wantVFXKey |
|---|---|---|---|---|---|---|---|
| infinite lifetime never despawns | 0 | "bullet_despawn" | false | 100 | false | 0 | "" |
| lifetime expires and queues removal | 3 | "bullet_despawn" | false | 3 | true | 1 | "bullet_despawn" |
| lifetime expires without vfx manager | 2 | "bullet_despawn" | true | 2 | true | 0 | "" |
| lifetime expires with empty effect | 2 | "" | false | 2 | true | 0 | "" |
| lifetime not yet expired | 5 | "bullet_despawn" | false | 4 | false | 0 | "" |
| negative lifetime treated as infinite | -1 | "bullet_despawn" | false | 10 | false | 0 | "" |

Assertions per row:
- `mockSpace.queuedForRemoval` contains `p.body` iff `wantQueued`.
- `mockVFX.spawnPuffCalls` length equals `wantVFXCalls`.
- When `wantVFXCalls == 1`, the recorded call has `typeKey == wantVFXKey`, `x == float64(x16)/16.0`, `y == float64(y16)/16.0`, `count == 1`, `randRange == 0.0`.
- `QueueForRemoval` is invoked at most once across all updates.

### 6.2 `TestProjectile_OOBHasNoVFX`

Given a projectile stepped out of bounds (via fake tilemap dimensions provider returning small width/height), after one `Update()`:
- Body is queued for removal.
- `mockVFX.spawnPuffCalls` length is `0`.

### 6.3 `TestProjectileConfig_LifetimeFrames_Default`

Zero-value `ProjectileConfig{}` has `LifetimeFrames == 0` and, when passed through `Manager.Spawn`, the resulting projectile behaves as infinite (never despawns by lifetime in 1000 updates).

### 6.4 `TestManager_Spawn_PropagatesLifetime`

Spawning with `ProjectileConfig{Width:2, Height:1, LifetimeFrames: 7}` produces a projectile whose `lifetimeFrames == 7` and `currentLifetime == 7` at spawn. (Use a package-level helper or expose via `len(m.projectiles)` and read through an unexported test accessor — prefer observing behavior: after exactly 7 `Update()` calls the body is queued for removal.)

## 7. Non-Goals

- No per-state overrides (that is US-037).
- No impact-VFX changes (that is US-031).
- No new contracts, no changes to `contractsvfx.Manager` or `contractsbody.BodiesSpace`.
- No change to `Manager.SpawnProjectile` defaults beyond what US-035/036 already introduced.

## 8. Risks / Notes

- **Behavioral change:** The current OOB branch emits `despawnEffect` VFX. This spec explicitly removes that call (OOB becomes silent). Any test that relied on OOB producing VFX must be updated — none currently do (verified in `manager_test.go`).
- **Order of operations:** Lifetime check runs **after** position update and collision resolution so the last-position VFX reflects the frame's final fp16 coordinates.
- **Determinism:** Frame counter is the only time source; no `time.Sleep`, no wall clock. Complies with constitution test rules.
