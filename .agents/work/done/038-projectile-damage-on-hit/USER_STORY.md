# US-038 — Projectiles Deal Damage on Hit

**Branch:** `038-projectile-damage-on-hit`
**Bounded Context:** Engine (`internal/engine/combat/projectile/`) + Game (`internal/game/`)

## Story

As a game designer,
I want projectiles to deal damage to Actors and destructible objects they collide with,
so that combat has real consequence — enemies are hurt by player shots, the player is hurt by enemy shots, and breakable objects respond to projectile impacts.

## Context

US-031 and US-036 established the projectile lifecycle (spawn, move, lifetime despawn, VFX). Projectiles already call `OnTouch` / `OnBlock` when they collide with any `Collidable` body (see `projectile.go`). The `Character` already has a `Hurt(damage int)` method that deducts health and transitions to the `Hurted` state with a 2-second invulnerability window. `ProjectileConfig` already carries a `Damage int` field.

What is missing is the bridge: when a projectile's `OnTouch` fires, it receives the colliding `Collidable`, but it has no way to ask that body "can you take damage, and if so, how?". This story introduces that bridge as a new contract and wires it through the projectile, weapon, and game-layer actors.

The owner stored on a projectile body (`SetOwner`) identifies the shooter. On collision, the projectile must:
1. Ignore its own owner (already done).
2. Resolve whether the hit body (or its owner) implements a `Damageable` contract.
3. If it does, call `TakeDamage(amount int)`.
4. Despawn itself (already done via `QueueForRemoval`).

Friendly-fire is **not** required for this story. A simple team/faction tag on the projectile and on `Damageable` is sufficient to distinguish player-owned from enemy-owned projectiles. If the projectile faction matches the target faction, the hit is ignored.

## Acceptance Criteria

- **AC1** — A new contract `contracts/combat.Damageable` is defined with a single method `TakeDamage(amount int)`. It lives in `internal/engine/contracts/combat/`.
- **AC2** — `projectile.OnTouch` and `projectile.OnBlock` attempt to resolve `Damageable` from the hit `Collidable`. Resolution order:
  1. Does `other` itself implement `Damageable`? Use it.
  2. Does `other.Owner()` implement `Damageable`? Use it.
  3. Neither → skip damage (projectile still despawns on solid hits).
- **AC3** — Damage is taken only when the projectile's faction differs from the target's faction, OR when either faction is `FactionNeutral` (neutral projectiles hurt everyone). A `Faction` type (`int`) and the values `FactionNeutral`, `FactionPlayer`, `FactionEnemy` are defined in `internal/engine/combat/`.
- **AC4** — `ProjectileConfig` gains a `Faction combat.Faction` field. `Manager.Spawn` stores the faction on the `projectile` struct.
- **AC5** — `Character.Hurt(damage int)` already exists and satisfies `Damageable`. No changes to `Character` are required unless a wrapper is needed to expose `TakeDamage`.
- **AC6** — `Character` (or a thin adapter registered at the game layer) is registered as `Damageable` so that both player and enemy `Character` instances respond to `TakeDamage`.
- **AC7** — A `Destructible` interface is defined in `internal/engine/contracts/combat/` as: `Damageable` + `IsDestroyed() bool`. Destructible objects (e.g. breakable walls) implement this interface. The projectile collision handler uses the same resolution path — no special-casing needed.
- **AC8** — `ProjectileConfig.Damage` of `0` is treated as a no-op: `TakeDamage(0)` is never called (guard in the projectile hit handler).
- **AC9** — The player's weapon config (`climber.json` or equivalent) sets a non-zero `damage` value. Enemy weapon configs also set a non-zero `damage`.
- **AC10** — Unit tests (table-driven, no GPU calls) cover:
  - Hit on a `Damageable` owner → `TakeDamage` called with correct amount.
  - Hit on a body that directly implements `Damageable` → `TakeDamage` called.
  - Hit on a non-damageable body → no panic, projectile still despawns.
  - Same-faction hit → `TakeDamage` not called.
  - `Damage == 0` → `TakeDamage` not called.
  - Invulnerable character ignores damage (existing `Hurt` guard, tested via `Character.Hurt`).

## Proposed Contract

```go
// internal/engine/contracts/combat/damageable.go
package combat

// Damageable is implemented by any entity that can receive damage.
type Damageable interface {
    TakeDamage(amount int)
}

// Destructible is a Damageable that can also report whether it has been destroyed.
type Destructible interface {
    Damageable
    IsDestroyed() bool
}
```

```go
// internal/engine/combat/faction.go
package combat

// Faction identifies which side a projectile or actor belongs to.
type Faction int

const (
    FactionNeutral Faction = iota // hurts everyone
    FactionPlayer
    FactionEnemy
)
```

## Design Notes

- **No global mutable state.** `Faction` constants are package-level `const`, not `var`, so they are safe.
- **Projectile struct change:** add `damage int` and `faction combat.Faction` fields. Both are set from `ProjectileConfig` in `Manager.Spawn`.
- **Faction on `Character`:** Add a `Faction() combat.Faction` method to `Character` (backed by a field set at construction or via a setter). The game layer sets `FactionPlayer` on the player and `FactionEnemy` on enemies. The `Damageable` resolution in the projectile checks `faction` before calling `TakeDamage`.
- **`TakeDamage` adapter:** `Character.Hurt` already has the correct semantics. Add `func (c *Character) TakeDamage(amount int) { c.Hurt(amount) }` to satisfy the interface without duplicating logic.
- **Destructible objects:** Out of scope for a concrete game-layer implementation in this story, but the contract and resolution path must support them. A test stub implementing `Destructible` is sufficient.
- **Pass-through projectiles:** Not in scope for this story. All projectiles despawn on hit (existing behavior preserved).
- **Faction check location:** Inside `projectile.applyDamage(other Collidable)` — a new private helper extracted from `OnTouch`/`OnBlock`.
- **Package layout:** All new engine types go in `internal/engine/combat/` (faction) or `internal/engine/contracts/combat/` (interfaces). No changes to physics packages.

## Dependencies

- US-031 (projectile lifetime / despawn VFX) — merged.
- US-037 (per-state spawn offset) — merged.
- No dependency on any unmerged story.
