# SPEC — 045-projectile-friendly-fire-prevention

## Summary

Prevent projectile Bodies from triggering damage and despawn against other projectile Bodies in the physics Space. Filtering happens at the projectile's collision callbacks (`OnTouch` / `OnBlock`) and is **opt-in per projectile** through a new `body.Projectile` contract that exposes an `Interceptable()` flag. Default projectiles (current bullets) are non-interceptable, so projectile-vs-projectile contact is a no-op for both sides. Future projectile types (e.g. rockets) can opt in by returning `true` from `Interceptable()`.

## Branch

`045-projectile-friendly-fire-prevention` (already created and checked out).

## Bounded Context

- **Physics**: `internal/engine/physics/space/` (collision dispatch — unchanged in this story).
- **Combat**: `internal/engine/combat/projectile/` (filtering applied in `projectile.OnTouch` and `projectile.OnBlock`).
- **Contracts**: `internal/engine/contracts/body/` (new `Projectile` interface).

## Design Decisions

1. **Where to filter**: in the projectile's own `OnTouch` / `OnBlock`, mirroring the existing `isPassthrough` and `isOwner` short-circuits. The Space remains an oblivious dispatcher; this keeps filtering deterministic and isolated to the Combat context. Adding filtering to `space.HasCollision` would force Space to know about projectiles, violating the bounded contexts.
2. **Contract shape**: a tiny, single-method interface `body.Projectile { Interceptable() bool }` colocated with `body.Passthrough`. The collidable body itself satisfies it (delegated through the projectile struct's own marker). The same `Interceptable() bool` method is also implemented on the runtime `*projectile` struct, exposed via the body's `Owner()` (the projectile sets itself as the body's `Touchable` already; ownership-style discovery works the same way as `isPassthrough`). This mirrors the existing pattern where the body or its owner is interrogated for traits.
3. **Identification**: a Body is "a projectile" iff it (or its owner) implements `body.Projectile`. No new global registry, no collision-layer enums. The `Interceptable()` value is read on the **other** body during `OnTouch`/`OnBlock` to decide whether the contact is allowed.
4. **Symmetry**: when `A` (projectile) touches `B` (projectile), Space invokes both `A.OnTouch(B)` and `B.OnTouch(A)`. Each side independently runs its filter. With both default (non-interceptable), both early-return — neither despawns and neither calls `applyDamage`.
5. **Opt-in interceptability**: when a projectile body's `Interceptable()` returns `true`, it behaves like an Actor for incoming projectiles — i.e. the *other* projectile's filter does NOT short-circuit and proceeds to `applyDamage` + despawn. The interceptable projectile itself, on its own callback, may still skip damage if it has no `Damageable`, but it should still despawn (it was hit). See "Behavior Matrix" below.
6. **No regressions on melee**: melee hitboxes are not Bodies registered as projectiles (they do not implement `body.Projectile`). The new branch in `OnTouch`/`OnBlock` only short-circuits when the **other** body is a projectile; melee paths are untouched.

## New Contract

File: `internal/engine/contracts/body/projectile.go`

```go
package body

// Projectile marks a Body (or its owner) as a projectile spawned by a combat
// system. It is queried by other projectiles during collision callbacks to
// implement projectile-vs-projectile filtering.
//
// Interceptable reports whether this projectile may be hit (and destroyed) by
// other projectiles. The default for engine-spawned bullets is false: two
// non-interceptable projectiles ignore each other on contact. Specialised
// projectile types (e.g. a rocket) may return true to opt in to being shot
// down mid-flight.
type Projectile interface {
    Interceptable() bool
}
```

The interface lives next to `body.Passthrough` and follows the same lookup convention (check the body itself, then `body.Owner()`).

## Implementation Touchpoints

### `internal/engine/combat/projectile/projectile.go`

- Add a file-local helper, in line with `isPassthrough` / `isOwner`:

  ```go
  // isProjectile reports whether other (or its owner) is itself a projectile,
  // and whether that projectile opted in to being intercepted.
  // Returns (isProj, interceptable).
  func isProjectile(other contractsbody.Collidable) (bool, bool) { ... }
  ```

  Lookup order: `other.(body.Projectile)`, then `other.Owner().(body.Projectile)`.

- The runtime `projectile` struct implements `body.Projectile` with `Interceptable() bool` returning a new `interceptable` field (default `false`). The field is wired from `ProjectileConfig.Interceptable` in `manager.Spawn`. Because `projectile` is set as the body's `Touchable` and also as its `Owner` is *not* the projectile (owner is the firing actor), exposing `Projectile` on the **body** directly is required. We extend the collidable body wrapper used by the manager to carry a small adapter, OR — simpler and preferred — the manager wraps the body so that its `Touchable` is the projectile, and the projectile is reachable via `body.Owner()` only when the firing actor is nil. To avoid breaking owner semantics, we instead make the body itself implement `Projectile` by attaching a typed wrapper.

  **Concrete approach (chosen)**: introduce an unexported `projectileBody` type in `internal/engine/combat/projectile/` that embeds `contractsbody.Collidable` and adds `Interceptable() bool`. `manager.Spawn` constructs and registers `projectileBody` instead of the bare `collidableBody`. The `Projectile` trait is then directly discoverable on the body (no owner traversal needed), keeping the contract symmetric with `Passthrough`.

### `internal/engine/combat/projectile/projectile.go` — callbacks

`OnTouch` becomes:

```go
func (p *projectile) OnTouch(other contractsbody.Collidable) {
    if p.isOwner(other) { return }
    if isPassthrough(other) { return }
    if isProj, interceptable := isProjectile(other); isProj && !interceptable {
        return // projectile-vs-non-interceptable-projectile: no interaction
    }
    p.applyDamage(other)
    p.spawnVFX(p.impactEffect)
    p.space.QueueForRemoval(p.body)
}
```

`OnBlock` mirrors the same guard.

### `internal/engine/combat/projectile/config.go`

Add `Interceptable bool` to `ProjectileConfig`. Default zero-value is `false`, preserving current behaviour.

### `internal/engine/combat/projectile/manager.go`

`Spawn` propagates `config.Interceptable` into the new `projectileBody` wrapper and into the `projectile` struct.

## Pre-conditions

- Two projectile Bodies, each with unique IDs, registered in the same `BodiesSpace`.
- Each projectile body implements `body.Projectile` via the new `projectileBody` wrapper.
- Both projectiles are non-interceptable (default).

## Post-conditions

- After `space.ResolveCollisions(A)` is called and `A.Position()` overlaps `B.Position()`:
  - `A` and `B` remain in the Space; neither is queued for removal.
  - No `applyDamage` is invoked on either side.
  - No impact VFX is spawned.
- For projectile-vs-Actor (existing path): unchanged. Damage is applied, projectile despawns.

## Behavior Matrix

| `A` (calling OnTouch) | `B` (other) | Outcome on `A` |
|---|---|---|
| Projectile, default | Projectile, default | early-return, no damage, no despawn |
| Projectile, default | Projectile, interceptable | proceed: damage attempt + despawn `A` |
| Projectile, interceptable | Projectile, default | early-return on `A`'s callback (B is default) — `A` survives |
| Projectile, default | Actor (Damageable) | proceed: damage + despawn (existing) |
| Projectile, default | Passthrough | early-return (existing) |
| Projectile, default | Owner | early-return (existing) |

> Note: the third row preserves symmetry — interceptability is a property of the *target*, not the shooter. An interceptable rocket only dies when the *other* projectile's `Interceptable()` check on the rocket returns `true`. The rocket's own filter sees the bullet as non-interceptable and ignores it. Net effect: rocket is destroyed, bullet survives — exactly the "shoot down a rocket" semantic.

## Acceptance Criteria → Test Mapping (Red Phase)

These tests will be authored by the TDD Specialist. All live in `internal/engine/combat/projectile/`.

### AC1: Projectile vs. Projectile — no interaction

**Test**: `TestProjectile_OnTouch_IgnoresOtherDefaultProjectile`
- Arrange: two `projectileBody` instances, both `Interceptable()==false`. One projectile struct `p` set as Touchable on body `A`. Mock `BodiesSpace`: assert `QueueForRemoval` is **not** called for `A`. Mock VFX manager: assert no `SpawnPuff`. Both bodies have a `Damageable` adapter that records calls (must record zero).
- Act: `p.OnTouch(B)` and `p.OnBlock(B)` separately.
- Assert: no removal, no VFX, no damage call.

**Test**: `TestSpace_TwoDefaultProjectilesOverlap_BothSurvive` (integration-level, using real `space.Space`)
- Arrange: real `Space`; spawn two projectiles via `Manager.Spawn` with overlapping positions and opposite velocities; one frame of `Manager.Update`.
- Assert: both bodies still discoverable via `space.Find(id)` after the frame.

### AC2: Projectile vs. Actor — interaction preserved

**Test**: `TestProjectile_OnTouch_HitsActorBody_AppliesDamageAndDespawns`
- Arrange: `p` with `damage=5`, faction `FactionPlayer`. `other` is an Actor body whose owner implements `Damageable` and `Faction()==FactionEnemy`.
- Act: `p.OnTouch(other)`.
- Assert: `Damageable.TakeDamage(5)` called once; `space.QueueForRemoval(p.body)` called once; impact VFX spawned.

### AC3: No regression on melee hitboxes

**Test**: `TestProjectile_OnTouch_DoesNotShortCircuitOnNonProjectileBody`
- Arrange: `other` is a melee-style Collidable that does **not** implement `body.Projectile`, but does implement `Damageable` via owner.
- Act: `p.OnTouch(other)`.
- Assert: damage applied, projectile despawned (i.e. existing behaviour unchanged when `other` is not a projectile).

### AC4: Deterministic behavior

**Test**: `TestProjectile_PvP_DeterministicAcrossOrdering`
- Table-driven: invoke `p1.OnTouch(p2body)` and `p2.OnTouch(p1body)` in both orders, and also independently. In every order: zero removals, zero damage.

### AC5: Opt-in interceptability

**Test**: `TestProjectile_OnTouch_InterceptableTargetIsHit`
- Arrange: `p` is default; `other` is a `projectileBody` whose `Interceptable()==true`, with a `Damageable` exposed via owner.
- Act: `p.OnTouch(other)`.
- Assert: `applyDamage` reaches `TakeDamage`, `p.body` queued for removal.

**Test**: `TestProjectile_OnTouch_InterceptableShooterStillIgnoresDefaultBullet`
- Arrange: `p` is interceptable; `other` is a default projectile.
- Act: `p.OnTouch(other)`.
- Assert: no damage, no despawn (target-property semantic).

**Test**: `TestProjectileConfig_DefaultInterceptableIsFalse`
- Assert: zero-value `ProjectileConfig{}.Interceptable == false`; `manager.Spawn` produces a body whose `Interceptable()==false`.

## Red Phase Scenario (canonical)

```
Given a BodiesSpace containing two non-interceptable projectile bodies A and B
And A and B are positioned so their collision rectangles overlap
When the projectile manager updates one frame
Then neither A nor B is queued for removal
And no damage callback is invoked on either side
And no impact VFX is emitted
```

This scenario fails today because the current `projectile.OnTouch` proceeds straight to `applyDamage` + `QueueForRemoval` whenever `other` is not the owner and not passthrough — including when `other` is another bullet.

## Out of Scope (re-affirming USER_STORY.md)

- Friendly-fire between a player projectile and the player Actor (handled by existing faction gate).
- Projectile-vs-environment / wall behavior.
- Implementing any specific interceptable projectile type. This spec only delivers the hook.

## Files to Modify / Create

- **Create**: `internal/engine/contracts/body/projectile.go`
- **Modify**: `internal/engine/combat/projectile/projectile.go` (add `isProjectile` helper, guard in `OnTouch`/`OnBlock`, expose `Interceptable()` on a new `projectileBody` wrapper).
- **Modify**: `internal/engine/combat/projectile/config.go` (add `Interceptable bool`).
- **Modify**: `internal/engine/combat/projectile/manager.go` (wire `config.Interceptable` through Spawn; wrap registered body as `projectileBody`).
- **Create / Modify (TDD Specialist)**: tests in `internal/engine/combat/projectile/projectile_test.go`, `manager_test.go`, `config_test.go` per the AC table above.

## Mocks Required

The new contract `body.Projectile` is trivially satisfied by the unexported `projectileBody` wrapper and by ad-hoc test doubles in `internal/engine/combat/projectile/mocks_test.go`. No shared mock generation is needed in `internal/engine/mocks/` because the interface has a single one-line method and is consumed only inside the `combat/projectile` package via type-assertion.

> Recommendation: **skip Mock Generator** for this story. The TDD Specialist can add a tiny `fakeProjectileBody` test-double in `mocks_test.go` directly.
