# SPEC — 039-enemies-shoot-projectiles

**Branch:** `039-enemies-shoot-projectiles`
**Bounded Contexts:** Engine (`internal/engine/combat/*`, `internal/engine/entity/actors/builder/`, `internal/engine/data/schemas/`) + Game (`internal/game/entity/actors/enemies/`).
**Depends on:** US-031, US-036, US-038 (projectile damage + faction gating already merged in commit `a15ee47`).

## 1. Goal

Allow enemies to own a `Weapon` and fire projectiles. Shooting behaviour is data-driven per enemy archetype and may be:

- **Unconditional** ("always" mode) — fire on cooldown regardless of whether a target exists (used by **Simple Immobile Shooter** archetypes like `BatEnemy`). No line-of-sight or range check.
- **Conditional** ("on_sight" mode) — fire only when a target is set AND within `Range` (used by **Smart Patrol Shooter** archetypes like `WolfEnemy`). When a patroller transitions into a configured shoot state it becomes immobile and fires.

In addition, every enemy weapon is configured with:
- A fixed **shoot direction axis** (`horizontal` or `vertical`) — no aim-at-target logic in this story.
- An optional **shoot state** gate — firing is only attempted while the owner's `State()` matches the configured state. If unset, firing is allowed in any state.

Enemy-owned projectiles are tagged `FactionEnemy` so the existing damage-gating logic (US-038) forwards damage to `FactionPlayer` / `FactionNeutral` actors while ignoring other enemies.

## 2. Domain Model Additions

### 2.1 JSON schema — `EnemyWeaponConfig`

Add to `internal/engine/data/schemas/json.go`:

```go
// EnemyWeaponConfig describes a simple ranged weapon on an enemy actor.
// Range is expressed in pixels and Cooldown in frames.
type EnemyWeaponConfig struct {
    ProjectileType string `json:"projectile_type"`
    Speed          int    `json:"speed"`                     // projectile speed (pixels/frame, converted to fp16 at use)
    Cooldown       int    `json:"cooldown"`                  // frames between shots
    Damage         int    `json:"damage"`                    // projectile damage
    Range          int    `json:"range"`                     // activation distance in pixels (only consulted when ShootMode == "on_sight")

    // ShootMode controls target / range gating. Allowed values:
    //   "always"   — fire on cooldown regardless of target (Simple Immobile Shooter).
    //   "on_sight" — fire only when a target is set AND within Range (Smart Patrol Shooter).
    // An empty / missing value defaults to "on_sight" (backward compatible with the initial WolfEnemy design).
    ShootMode string `json:"shoot_mode,omitempty"`

    // ShootDirection is the fixed axis the projectile travels on. Allowed values:
    //   "horizontal" — maps to body.ShootDirectionStraight (sign follows owner.FaceDirection()).
    //   "vertical"   — maps to body.ShootDirectionDown (immobile shooters fire downward by default).
    // An empty / missing value defaults to "horizontal".
    ShootDirection string `json:"shoot_direction,omitempty"`

    // ShootState is the name of the actor state (as registered via actors.RegisterState)
    // in which firing is active. When empty, firing is permitted in any state.
    // Looked up via actors.GetStateEnum at builder time; unknown names produce a build error.
    ShootState string `json:"shoot_state,omitempty"`
}
```

Extend `SpriteData` with an optional `Weapon *EnemyWeaponConfig` field (`json:"weapon,omitempty"`). Nil means "no weapon", preserving backward compatibility for enemies without ranged attacks.

Validation performed at unmarshal / builder stage (see §3):

| Field | Allowed | On invalid |
|---|---|---|
| `ShootMode` | `""`, `"always"`, `"on_sight"` | error: `"invalid shoot_mode: %q"` |
| `ShootDirection` | `""`, `"horizontal"`, `"vertical"` | error: `"invalid shoot_direction: %q"` |
| `ShootState` | `""` or a name registered via `actors.RegisterState` | error: `"unknown shoot_state: %q"` |

### 2.2 Enemy shooting AI contract

Update `internal/engine/contracts/combat/enemy_shooter.go`:

```go
package combat

import (
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
    "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// ShootMode describes whether a shooter requires a target / range check.
type ShootMode int

const (
    // ShootModeOnSight fires only when a target is set and within Range.
    ShootModeOnSight ShootMode = iota
    // ShootModeAlways fires on cooldown regardless of target.
    ShootModeAlways
)

// EnemyShooter is an Actor-facing behaviour that attempts to fire each frame
// according to its configured ShootMode.
type EnemyShooter interface {
    // SetTarget sets the actor this shooter aims at. May be nil.
    // Required only by ShootModeOnSight; ignored by ShootModeAlways.
    SetTarget(target body.MovableCollidable)
    // Target returns the current target (nil if unset).
    Target() body.MovableCollidable
    // Range returns the activation distance in pixels (only used by ShootModeOnSight).
    Range() int
    // Mode returns the configured firing mode.
    Mode() ShootMode
    // Direction returns the fixed projectile direction (body.ShootDirectionStraight or body.ShootDirectionDown).
    Direction() body.ShootDirection
    // ShootState returns the state in which firing is active, and true if gating is enabled.
    // When the second return value is false, firing is allowed in any state.
    ShootState() (actors.ActorStateEnum, bool)
    // TryFire attempts to fire. Returns true if a projectile was spawned.
    TryFire() bool
    // Update advances weapon cooldown and, when possible, fires.
    Update()
}
```

Note: the contract may avoid importing `actors` by returning `(int, bool)` for `ShootState()`. Choice deferred to TDD Specialist; either form is acceptable provided the builder stores the resolved enum.

### 2.3 Implementation — `combat/weapon.EnemyShooting`

File `internal/engine/combat/weapon/enemy_shooting.go`.

Constructor signature:

```go
func NewEnemyShooting(
    owner body.MovableCollidable,
    weapon combat.Weapon,
    rangePx int,
    mode combat.ShootMode,
    direction body.ShootDirection,
    shootState actors.ActorStateEnum,
    shootStateActive bool,
) *EnemyShooting
```

Internal state:

```go
type EnemyShooting struct {
    owner            body.MovableCollidable
    weapon           combat.Weapon
    target           body.MovableCollidable
    rangePx          int
    mode             combat.ShootMode
    direction        body.ShootDirection // ShootDirectionStraight or ShootDirectionDown
    shootState       actors.ActorStateEnum
    shootStateActive bool
}
```

Algorithm per `Update()` (which calls `TryFire()` internally and returns nothing):

1. `weapon.Update()` (decrement cooldown). This runs every frame, independent of all gates below.
2. **State gate.** If `shootStateActive` and `owner.State() != shootState` → return (no fire, cooldown already ticked).
3. **Target gate.** Evaluated per mode:
   - `ShootModeAlways` → no target required; skip to step 4.
   - `ShootModeOnSight`:
     - If `target == nil` → return.
     - Compute `dx = targetCenterX − ownerCenterX`, `dy = targetCenterY − ownerCenterY` in pixel space.
     - If `dx*dx + dy*dy > rangePx*rangePx` → return. (Squared comparison, no sqrt.)
4. **Cooldown gate.** If `!weapon.CanFire()` → return.
5. **Face direction.** Only for `ShootModeOnSight` with a target: set `owner.SetFaceDirection(...)` to `FaceDirectionLeft` if `dx < 0`, else `FaceDirectionRight`. For `ShootModeAlways` the owner's current `FaceDirection()` is left untouched (immobile shooters keep their configured facing).
6. **Fire.** Read `x16, y16 := owner.GetPosition16()`; call `weapon.Fire(x16, y16, owner.FaceDirection(), s.direction, int(owner.State()))`.

`TryFire()` performs steps 2–6 and returns `true` only if step 6 executes. `Update()` is `w.Update(); _ = s.TryFire()` — note weapon cooldown decrements every frame regardless of gating.

### 2.4 Character factions at spawn

- `BatEnemy`, `WolfEnemy`, and every new enemy type with a weapon MUST call `SetFaction(FactionEnemy)` after builder configuration.
- `ClimberPlayer` MUST call `SetFaction(FactionPlayer)` (verify existing; document and keep).

## 3. Builder Integration

Extend `internal/engine/entity/actors/builder/builder.go`:

```go
// ConfigureEnemyWeapon builds an EnemyShooter from SpriteData.Weapon and attaches it
// to the given character. If cfg is nil, returns (nil, nil).
func ConfigureEnemyWeapon(
    character actors.ActorEntity,
    cfg *schemas.EnemyWeaponConfig,
    manager combat.ProjectileManager,
) (combat.EnemyShooter, error)
```

Behaviour:

- Returns `(nil, nil)` if `cfg == nil`.
- Returns `(nil, error)` if `manager == nil` when `cfg != nil`.
- Parses `cfg.ShootMode`:
  - `""` or `"on_sight"` → `combat.ShootModeOnSight`.
  - `"always"` → `combat.ShootModeAlways`.
  - Anything else → `(nil, error)`.
- Parses `cfg.ShootDirection`:
  - `""` or `"horizontal"` → `body.ShootDirectionStraight`.
  - `"vertical"` → `body.ShootDirectionDown`.
  - Anything else → `(nil, error)`.
- Parses `cfg.ShootState`:
  - `""` → `shootStateActive = false`, enum zero.
  - Non-empty → `enum, ok := actors.GetStateEnum(cfg.ShootState)`; on `!ok` return `(nil, error)`.
- Constructs a `ProjectileWeapon` via `weapon.NewProjectileWeapon(id, cfg.Cooldown, cfg.ProjectileType, fp16.To16(cfg.Speed), manager, "", 0, 0)` with `id = fmt.Sprintf("%s_weapon", character.ID())`.
- Calls `w.SetDamage(cfg.Damage)` and `w.SetOwner(character)`.
- Wraps via `weapon.NewEnemyShooting(character.GetCharacter(), w, cfg.Range, mode, direction, shootStateEnum, shootStateActive)`.

`ConfigureCharacter` is NOT changed; enemies invoke `ConfigureEnemyWeapon` explicitly. Rationale: the projectile manager is only available at the game layer and `ConfigureCharacter` does not currently receive it; adding an optional call keeps the current API stable.

## 4. Game-Layer Wiring

**Archetype note — "immobile" vs. "smart patrol":** "Immobile" in this story means the game-layer wiring does NOT call `SetTarget()` on the enemy in the `InitEnemyMap` factory — the enemy spawns once and fires on its configured axis without ever acquiring a target. A "smart patrol" enemy patrols via its movement state and dynamically calls `SetTarget(player)` so the on_sight gate can evaluate range. Both archetypes still run the same `EnemyShooting.Update()` algorithm; the difference is purely in mode / direction configuration and whether the game loop ever forwards a target.

### 4.1 `BatEnemy` — Simple Immobile Shooter archetype

Update `internal/game/entity/actors/enemies/bat.go`:

1. Struct gains `shooter combat.EnemyShooter`.
2. `NewBatEnemy` reads `spriteData.Weapon` and, if non-nil, calls `builder.ConfigureEnemyWeapon(enemy, spriteData.Weapon, ctx.ProjectileManager)`, storing the result.
3. `enemy.SetFaction(enginecombat.FactionEnemy)` called once after building.
4. `BatEnemy` is immobile (no `SetMovementState` that produces motion beyond the existing hover, plus `SetGravityEnabled(false)`). It does NOT require `SetTarget()` to fire; the `InitEnemyMap` factory for `BatEnemyType` MUST NOT call `enemy.SetTarget(player)`. If `SetTarget` is retained on the struct for API parity, it becomes a no-op for weapon purposes (it may still forward to the movement state, but the shooter does not consume it).
5. `Update(space)` calls `e.shooter.Update()` after `e.Character.Update(space)` when `shooter != nil`.
6. `bat.json` declares `"shoot_mode": "always"` and `"shoot_direction": "vertical"`. No `shoot_state` is required — the bat fires continuously on cooldown. `range` is `0` since it is never consulted.

### 4.2 `WolfEnemy` — Smart Patrol Shooter archetype

Update `internal/game/entity/actors/enemies/wolf.go`:

1. Struct gains `shooter combat.EnemyShooter`.
2. `NewWolfEnemy` reads `spriteData.Weapon` and, if non-nil, calls `builder.ConfigureEnemyWeapon(enemy, spriteData.Weapon, ctx.ProjectileManager)`, storing the result.
3. `enemy.SetFaction(enginecombat.FactionEnemy)` called once after building.
4. `SetTarget(target)` forwards the target both to the existing movement state (patrol / chase logic) AND to `e.shooter` when non-nil, so the `ShootModeOnSight` range gate can evaluate distance to the player.
5. `Update(space)` calls `e.shooter.Update()` after `e.Character.Update(space)` when `shooter != nil`.
6. `wolf.json` declares `"shoot_mode": "on_sight"` and `"shoot_direction": "horizontal"`. A `shoot_state` gate MAY be supplied (e.g. `"StateChase"` or another patrol/chase state registered via `actors.RegisterState`) so the wolf only fires once the movement state transitions out of idle. When omitted, the wolf fires in any state whenever the player is within range and cooldown has elapsed.
7. The `InitEnemyMap` entry for `WolfEnemyType` continues to call `enemy.SetTarget(player)` (existing behaviour); this now also primes the shooter's range gate.

## 5. JSON Examples

### 5.1 Updated `assets/entities/enemies/bat.json`

Add a top-level `weapon` block (beside `stats` and `sprites`):

```json
{
  "sprites": { ... },
  "stats": { "health": 1, "speed": 0, "max_speed": 0 },
  "weapon": {
    "projectile_type": "bullet_small",
    "speed": 5,
    "cooldown": 60,
    "damage": 1,
    "range": 0,
    "shoot_mode": "always",
    "shoot_direction": "vertical"
  }
}
```

### 5.2 Updated `assets/entities/enemies/wolf.json`

Add a top-level `weapon` block (beside `stats` and `sprites`):

```json
{
  "sprites": { ... },
  "stats": { "health": 2, "speed": 6, "max_speed": 6 },
  "weapon": {
    "projectile_type": "bullet_small",
    "speed": 6,
    "cooldown": 90,
    "damage": 1,
    "range": 160,
    "shoot_mode": "on_sight",
    "shoot_direction": "horizontal"
  }
}
```

Existing enemies without a `weapon` block remain melee-only.

## 6. Projectile Faction

Unchanged from prior iteration:

- **Preferred (explicit):** `ProjectileWeapon.SetOwner(owner)` already tags ownership. The projectile manager resolves the projectile's faction from the owner actor's `Faction()` when spawning (existing behaviour after US-038). Setting the enemy's faction to `FactionEnemy` (§2.4) is sufficient.
- If the manager does not yet derive faction from owner, fall back to an explicit `SetFaction(FactionEnemy)` hook on the weapon. Verify during TDD; adjust only if tests fail.

## 7. Pre/Post-Conditions

| # | Pre-condition | Post-condition |
|---|---|---|
| P1 | `cfg.Weapon` present in JSON and valid | `enemy.shooter != nil`; `Mode()`, `Direction()`, and (when non-empty) `ShootState()` reflect the parsed config |
| P2 | `Mode == OnSight` AND target within `Range` AND `weapon.CanFire()` AND (no state gate OR `owner.State() == shootState`) | Exactly one projectile spawned at owner's fp16 position with velocity matching `Direction` × owner `FaceDirection()` |
| P3 | Same as P2 but `!weapon.CanFire()` | No projectile spawned; cooldown continues to decrement |
| P4 | `Mode == OnSight` AND target farther than `Range` | No projectile; `weapon.CanFire()` preserved (no cooldown consumed) |
| P5 | `Mode == OnSight` AND `target == nil` | No projectile; no panic |
| P6 | `Mode == Always` AND `weapon.CanFire()` AND (no state gate OR state match) — target may be nil | Exactly one projectile spawned on the configured `Direction` axis |
| P7 | `Mode == Always` AND `!weapon.CanFire()` | No projectile; cooldown continues to decrement |
| P8 | State gate active AND `owner.State() != shootState` (any mode) | No projectile; cooldown still decrements |
| P9 | Projectile spawned by enemy | `projectile.Owner() == enemy` AND projectile faction == `FactionEnemy` |
| P10 | Enemy projectile collides with `FactionPlayer` actor | Target's `TakeDamage(cfg.Damage)` invoked |
| P11 | Enemy projectile collides with another `FactionEnemy` actor | No damage applied (existing US-038 gating) |
| P12 | `cfg.ShootDirection == "horizontal"` | `shooter.Direction() == body.ShootDirectionStraight` and spawned projectile velocity has `vy16 == 0` |
| P13 | `cfg.ShootDirection == "vertical"` | `shooter.Direction() == body.ShootDirectionDown` and spawned projectile velocity has `vx16 == 0`, `vy16 > 0` |

## 8. Red Phase — Failing Tests

Tests to author in the TDD Specialist stage.

### 8.1 `internal/engine/combat/weapon/enemy_shooting_test.go` — firing gates

Table-driven tests using a `mockProjectileManager` (already present in package) and a minimal `fakeOwner` that satisfies the subset of `body.MovableCollidable` required (`GetPosition16`, `FaceDirection`, `SetFaceDirection`, `GetShape`, `State`). Each row configures `mode`, `direction`, and optional `shootState`/`ownerState`. Rows exercise the two archetypes: `ShootModeOnSight` + `ShootDirectionStraight` (the WolfEnemy case) and `ShootModeAlways` + `ShootDirectionDown` (the BatEnemy case).

| Name | Mode | Owner | Target offset | Cooldown | State gate | Expect `TryFire` | Expect spawn | Expect faceDir |
|---|---|---|---|---|---|---|---|---|
| on_sight fires when target in range, left | OnSight | (100,100) | (−50,0) | ready | none | true | 1 | Left |
| on_sight fires when target in range, right | OnSight | (100,100) | (+40,0) | ready | none | true | 1 | Right |
| on_sight skips when out of range | OnSight | (100,100) | (+500,0), range=160 | ready | none | false | 0 | unchanged |
| on_sight skips during cooldown | OnSight | (100,100) | (+40,0) | cooldown=5 | none | false | 0 | unchanged |
| on_sight skips when target is nil | OnSight | (100,100) | — | ready | none | false | 0 | unchanged |
| on_sight fires after cooldown elapses | OnSight | (100,100) | (+40,0) | cooldown=1 → Update → 0 | none | second call true | 1 additional | — |
| always fires with no target | Always | (50,50) | — | ready | none | true | 1 | unchanged |
| always fires on cooldown end even without target | Always | (50,50) | — | cooldown=1 → Update | none | second call true | 1 additional | unchanged |
| always skips during cooldown | Always | (50,50) | — | cooldown=5 | none | false | 0 | unchanged |
| state gate blocks fire when state mismatches | Always | (50,50) | — | ready | active, mismatched | false | 0 | unchanged |
| state gate permits fire when state matches | Always | (50,50) | — | ready | active, matched | true | 1 | unchanged |
| state gate also applies to on_sight | OnSight | (100,100) | (+40,0) | ready | active, mismatched | false | 0 | unchanged |

### 8.2 Direction axis mapping

Cases asserting the velocity of the recorded spawn call:

| Name | Direction | FaceDir | Expected `vx16` sign | Expected `vy16` |
|---|---|---|---|---|
| horizontal right (WolfEnemy archetype) | ShootDirectionStraight | Right | > 0 | 0 |
| horizontal left (WolfEnemy archetype) | ShootDirectionStraight | Left | < 0 | 0 |
| vertical down (BatEnemy archetype) | ShootDirectionDown | Right | 0 | > 0 |
| vertical down regardless of facing (BatEnemy archetype) | ShootDirectionDown | Left | 0 | > 0 |

### 8.3 Faction assertion

One test: given a `fakeOwner` whose `Faction()` returns `FactionEnemy`, `SpawnProjectile` is called with that owner identity preserved on the recorded call.

### 8.4 `internal/engine/entity/actors/builder/builder_test.go`

- `TestConfigureEnemyWeapon_NilConfig` → returns `(nil, nil)`.
- `TestConfigureEnemyWeapon_MissingManager` → returns `(nil, error)`.
- `TestConfigureEnemyWeapon_Builds_OnSightHorizontal` → returns non-nil `EnemyShooter` whose `Range()`, `Mode()`, and `Direction()` match cfg (Wolf archetype) and whose internal weapon spawns a projectile with the configured damage on forced `TryFire`.
- `TestConfigureEnemyWeapon_Builds_AlwaysVertical` → verifies `Mode() == ShootModeAlways`, `Direction() == ShootDirectionDown` (Bat archetype).
- `TestConfigureEnemyWeapon_InvalidShootMode` → returns `(nil, error)` with `shoot_mode = "sometimes"`.
- `TestConfigureEnemyWeapon_InvalidShootDirection` → returns `(nil, error)` with `shoot_direction = "diagonal"`.
- `TestConfigureEnemyWeapon_UnknownShootState` → returns `(nil, error)` for a state name not registered.
- `TestConfigureEnemyWeapon_ValidShootState` → returns shooter whose `ShootState()` returns `(enum, true)` matching `actors.GetStateEnum`.

### 8.5 `internal/engine/data/schemas/json_test.go`

- Unmarshals a minimal JSON with a `weapon` block including `shoot_mode`, `shoot_direction`, and `shoot_state`; verifies all fields populate on `SpriteData.Weapon`.
- Verifies `SpriteData.Weapon` is nil when the block is absent (backward compatibility).
- Verifies default empty-string values when optional keys are omitted inside the block.

### 8.6 `internal/game/entity/actors/enemies/bat_test.go` (new) — always integration (Simple Immobile Shooter)

Integration-level unit test using existing test helpers (mock `AppContext` with a fake projectile manager). Assertions:

- After `NewBatEnemy`, `enemy.Faction() == FactionEnemy` and `enemy.shooter != nil`.
- Without ever calling `SetTarget`, after enough `Update` calls to elapse one cooldown, exactly one projectile is spawned on the fake manager.
- Spawned projectile velocity has `vx16 == 0` and `vy16 > 0` (vertical down).
- Repeated `Update` calls within cooldown do NOT spawn additional projectiles.

### 8.7 `internal/game/entity/actors/enemies/wolf_test.go` (new) — on_sight integration (Smart Patrol Shooter)

Integration-level unit test using existing test helpers (mock `AppContext` with a fake projectile manager). Assertions:

- After `NewWolfEnemy`, `enemy.Faction() == FactionEnemy` and `enemy.shooter != nil`.
- When `SetTarget(player)` is called with the player within `range`, a subsequent `Update` spawns exactly one projectile on the fake manager and records owner == enemy.
- When the target is beyond `range`, no projectile is spawned.
- When `SetTarget` has never been called, no projectile is spawned.

(If `app.AppContext` wiring proves heavy in a unit test, this test may be reduced to asserting faction only, with §8.1/§8.4 covering firing logic.)

## 9. Out of Scope

- **Aim-at-target direction logic.** Enemies only shoot on a fixed axis configured in JSON. A future story will add aim vectors / diagonal shots for smart enemies.
- Line-of-sight checks through walls.
- Multi-weapon inventories for enemies (a single `Weapon` is enough).
- New projectile types; enemies reuse existing `bullet_small`.
- Transitions from patrol state to a dedicated "shoot" state in `WolfEnemy`. The state gate plumbing is present, but wiring a patrol-to-shoot state machine change for wolves is not part of US-039 (use a future story; the BatEnemy always-mode archetype validates the non-gated path and the state gate is covered by unit tests in §8.1).
- Vertical aim for patrol shooters on arbitrary direction (only horizontal or downward-fixed per config here).

## 10. Definition of Done

- All tests in §8 pass; coverage in affected packages does not regress.
- `EnemyWeaponConfig` schema extensions (`ShootMode`, `ShootDirection`, `ShootState`) implemented and validated.
- `bat.json` updated with a `weapon` block using `shoot_mode: "always"` and `shoot_direction: "vertical"`; in-game `BatEnemy` fires continuously downward on its cooldown without ever acquiring a target.
- `wolf.json` updated with a `weapon` block using `shoot_mode: "on_sight"` and `shoot_direction: "horizontal"`; in-game `WolfEnemy` patrols and fires at the player when within range.
- Enemy projectiles carry `FactionEnemy` and deal damage to `FactionPlayer` only.
- `golangci-lint run ./...` clean.
- Acceptance Criteria AC1–AC7 verified by Gatekeeper, extended to cover the two archetypes and the `shoot_mode` / `shoot_direction` / `shoot_state` config keys.
