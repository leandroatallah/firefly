# SPEC — US-038 Projectiles Deal Damage on Hit

**Branch:** `038-projectile-damage-on-hit`
**Bounded Contexts:** `internal/engine/combat/` (projectile, faction), `internal/engine/contracts/combat/`, `internal/engine/entity/actors/` (Character adapter), `internal/game/entity/actors/` (wiring).

---

## 1. Goal

Bridge projectile collisions to the existing `Character.Hurt` path via a new `Damageable` contract, gated by a faction check. Preserve all existing projectile despawn / VFX behaviour (US-031, US-036, US-037).

---

## 2. New / Changed Packages

### 2.1 `internal/engine/combat/faction.go` (NEW)
```go
package combat

type Faction int

const (
    FactionNeutral Faction = iota
    FactionPlayer
    FactionEnemy
)
```

Rules:
- Package-level `const` only (no global `var`). Constitution compliant.
- Lives at `internal/engine/combat/` (sibling of existing `projectile/`, `weapon/`, `inventory/` sub-packages). No new sub-dir.

### 2.2 `internal/engine/contracts/combat/damageable.go` (NEW)
```go
package combat

type Damageable interface {
    TakeDamage(amount int)
}

type Destructible interface {
    Damageable
    IsDestroyed() bool
}
```

Notes:
- `contracts/combat` already hosts `weapon.go`, `inventory.go`, `projectile_manager.go`. Adding `damageable.go` alongside them.
- Interfaces are referenced by the projectile package via type assertion only — no cross-package cycle created.

### 2.3 `internal/engine/combat/projectile/config.go` (CHANGED)
Add import of engine combat package:
```go
import enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"

type ProjectileConfig struct {
    Width          int
    Height         int
    Damage         int
    Faction        enginecombat.Faction `json:"faction,omitempty"` // NEW (default 0 = Neutral)
    ImpactEffect   string `json:"impact_effect,omitempty"`
    DespawnEffect  string `json:"despawn_effect,omitempty"`
    LifetimeFrames int    `json:"lifetime_frames,omitempty"`
}
```

### 2.4 `internal/engine/combat/projectile/projectile.go` (CHANGED)
Add two fields to `projectile`:
```go
damage  int
faction enginecombat.Faction
```

Extract a helper:
```go
// applyDamage resolves a Damageable from the hit body and calls TakeDamage,
// honouring faction and zero-damage guards. Safe on nil / non-damageable others.
func (p *projectile) applyDamage(other contractsbody.Collidable) {
    if p.damage == 0 { return }
    if other == nil { return }

    target, tFaction, ok := p.resolveDamageable(other)
    if !ok { return }

    // Faction gate: skip only when both sides are non-neutral AND equal.
    if p.faction != combat.FactionNeutral &&
        tFaction != combat.FactionNeutral &&
        p.faction == tFaction {
        return
    }
    target.TakeDamage(p.damage)
}

// resolveDamageable tries (1) the body itself, (2) body.Owner().
// Returns the Damageable plus the target's faction (FactionNeutral when
// target does not implement a Faction() accessor).
func (p *projectile) resolveDamageable(other contractsbody.Collidable) (combat.Damageable, combat.Faction, bool)
```

The target faction is read via an ad-hoc interface (type assertion) inside `resolveDamageable`:
```go
type factioned interface { Faction() combat.Faction }
```
If the target does not implement `factioned`, its faction is `FactionNeutral`.

`OnTouch` / `OnBlock` changes:
```go
func (p *projectile) OnTouch(other contractsbody.Collidable) {
    if other == p.body.Owner() { return }     // unchanged (self-owner ignore)
    p.applyDamage(other)                      // NEW
    p.spawnVFX(p.impactEffect)                // unchanged
    p.space.QueueForRemoval(p.body)           // unchanged
}

func (p *projectile) OnBlock(other contractsbody.Collidable) {
    p.applyDamage(other)                      // NEW (blocking hits also damage)
    p.spawnVFX(p.impactEffect)
    p.space.QueueForRemoval(p.body)
}
```

### 2.5 `internal/engine/combat/projectile/manager.go` (CHANGED)
In `Spawn`, propagate new fields:
```go
p := &projectile{
    ...,
    damage:  config.Damage,
    faction: config.Faction,
}
```

### 2.6 `internal/engine/entity/actors/character.go` (CHANGED)
Add:
```go
// New field
faction enginecombat.Faction

// Getter / setter
func (c *Character) Faction() enginecombat.Faction        { return c.faction }
func (c *Character) SetFaction(f enginecombat.Faction)    { c.faction = f }

// Damageable adapter (delegates to existing Hurt; invulnerability preserved).
func (c *Character) TakeDamage(amount int) { c.Hurt(amount) }
```

The default faction is `FactionNeutral` (zero value). Game-layer wiring sets `FactionPlayer` / `FactionEnemy` explicitly.

### 2.7 Game-layer wiring (CHANGED)
- `internal/game/entity/actors/player/climber.go`: after constructing the player `Character`, call `SetFaction(combat.FactionPlayer)`.
- `internal/game/entity/actors/enemies/bat.go` & `wolf.go`: call `SetFaction(combat.FactionEnemy)`.

### 2.8 JSON configs (CHANGED)
- `assets/entities/player/climber.json` weapon block: add `"damage": <non-zero>` (e.g. `10`). If the shoot/weapon JSON schema does not currently expose damage, thread it through the existing weapon loader path the same way `impact_effect` is threaded. *(Spec Engineer note: loader changes, if any, are limited to passing the existing `Damage` field forward; no new loader types required.)*
- Any enemy weapon JSON: same, non-zero damage.
- Faction on the `ProjectileConfig` is set by the weapon / firing site based on the shooter's faction, NOT by JSON (per AC3 — JSON damage only). Recommended: `Fire` call-site reads shooter faction and injects into the `ProjectileConfig` before delegating to `Manager.Spawn`. If this is out of scope for RED-phase tests, it may be set as a constant at the weapon wiring layer.

---

## 3. Pre- / Post-conditions per Acceptance Criterion

| AC | Pre-condition | Post-condition |
|----|---------------|----------------|
| AC1 | `contracts/combat` exists | `Damageable` interface compiles with exactly one method `TakeDamage(int)` |
| AC2 | `other` is non-nil `Collidable` | Resolution tries body → owner → skip; never panics |
| AC3 | Projectile and target have factions set | `TakeDamage` called iff factions differ OR at least one is `FactionNeutral` |
| AC4 | `ProjectileConfig.Faction` field exists | `Spawn` copies `config.Faction` into `projectile.faction` |
| AC5 | `Character.Hurt` unchanged | `Character` satisfies `Damageable` via `TakeDamage` adapter |
| AC6 | Player & enemy constructed via game layer | Both have their faction set; both respond to `TakeDamage` |
| AC7 | `Destructible` compiles | Projectile hit path works for an object implementing `Destructible` without extra code |
| AC8 | `config.Damage == 0` | No call to `TakeDamage`; projectile still despawns |
| AC9 | Weapon configs updated | Runtime projectiles carry non-zero damage |
| AC10 | Mocks for `Damageable` + `Collidable` available | All table-driven scenarios pass, no GPU calls |

---

## 4. Red-Phase Test Plan (feeds TDD Specialist)

All tests live in `internal/engine/combat/projectile/projectile_test.go` (new cases) unless noted.

### 4.1 Table: `TestProjectile_AppliesDamageOnHit`
Columns: `name`, `projectileFaction`, `projectileDamage`, `target` (fake Collidable factory), `expectTakeDamageCalls`, `expectAmount`.

Rows:
1. `"hit on owner that is Damageable"` — target body's `Owner()` returns a `*fakeDamageable` with faction `FactionEnemy`; projectile faction `FactionPlayer`, damage `10` → 1 call, amount `10`.
2. `"hit on body that directly implements Damageable"` — body itself is a `*fakeDamageableBody` implementing both `Collidable` + `Damageable`; different faction → 1 call.
3. `"hit on non-damageable body"` — plain `Collidable` with nil owner → 0 calls, no panic, projectile queued for removal.
4. `"same-faction hit ignored"` — projectile `FactionPlayer`, target owner `FactionPlayer` → 0 calls.
5. `"neutral projectile hurts player"` — projectile `FactionNeutral`, target `FactionPlayer` → 1 call.
6. `"neutral target hurt by player projectile"` — projectile `FactionPlayer`, target `FactionNeutral` (target has no Faction() method) → 1 call.
7. `"zero damage no-op"` — projectile damage `0`, damageable target → 0 calls, projectile still queued for removal.
8. `"hit on self owner"` — `other == body.Owner()` → 0 calls, projectile NOT removed (existing OnTouch short-circuit).

Every row also asserts: `space.QueueForRemoval` was called exactly once (except row 8: zero).

### 4.2 `TestProjectile_ResolvesDestructible`
- Target body's `Owner()` returns a struct implementing `Destructible`.
- Assert `TakeDamage` is called; no special-casing required (same path as Damageable).

### 4.3 `TestCharacter_TakeDamageDelegatesToHurt`
Location: `internal/engine/entity/actors/character_test.go`.
- Table: `(initialHealth, invulnerable, damage) → (expectedHealth, expectedState)`.
- Covers: damage applies + Hurted transition; invulnerable → no-op (AC10 last bullet).

### 4.4 `TestCharacter_FactionAccessors`
- `SetFaction(FactionEnemy)` then `Faction() == FactionEnemy`.
- Default (unset) is `FactionNeutral`.

### 4.5 `TestProjectileConfig_FactionField`
In `config_test.go`: JSON round-trip preserves `Faction` value.

---

## 5. Required Mocks (feeds Mock Generator)

Package-local, in `internal/engine/combat/projectile/mocks_test.go` (augment existing file):

- `fakeDamageable` — records `TakeDamage` calls (count + last amount), has `Faction() combat.Faction`.
- `fakeDamageableBody` — implements `contracts/body.Collidable` AND `contracts/combat.Damageable` AND `factioned`. Built on top of the existing fake body harness.
- `fakeDestructible` — `TakeDamage` + `IsDestroyed() bool`; optional `Faction()`.
- `fakeCollidableWithOwner(owner interface{})` — helper to build a Collidable returning a specific `Owner()`.

No new shared mocks under `internal/engine/mocks/` are required for this story.

---

## 6. Integration Points

- **Physics:** no changes. Projectile continues to call `ResolveCollisions` and react to `OnTouch` / `OnBlock`.
- **VFX:** no changes. Impact VFX still spawned before removal.
- **Weapon:** weapon `Fire` site must pass `ProjectileConfig.Faction`. Out of scope for unit tests in this story — wiring only.
- **Scene / Actor Manager:** unchanged; player/enemy construction simply calls `SetFaction`.

---

## 7. Non-goals (explicit)

- Pass-through projectiles (multiple hits per projectile).
- Per-target resistance / armour.
- Knockback / hitstun beyond existing `Hurted` state.
- Concrete breakable-wall game entity — only the interface and its hit path are spec'd.
- Friendly-fire toggling beyond faction equality.

---

## 8. File Checklist

Created:
- `internal/engine/combat/faction.go`
- `internal/engine/contracts/combat/damageable.go`

Modified:
- `internal/engine/combat/projectile/config.go`
- `internal/engine/combat/projectile/projectile.go`
- `internal/engine/combat/projectile/manager.go`
- `internal/engine/combat/projectile/projectile_test.go`
- `internal/engine/combat/projectile/mocks_test.go`
- `internal/engine/combat/projectile/config_test.go`
- `internal/engine/entity/actors/character.go`
- `internal/engine/entity/actors/character_test.go`
- `internal/game/entity/actors/player/climber.go`
- `internal/game/entity/actors/enemies/bat.go`
- `internal/game/entity/actors/enemies/wolf.go`
- `assets/entities/player/climber.json` (+ any enemy weapon JSON)
