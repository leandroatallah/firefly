# US-039 — Enemies Shoot Projectiles

**Branch:** `039-enemies-shoot-projectiles`
**Bounded Context:** Game (`internal/game/entity/actors/enemies/`) + Engine (`internal/engine/combat/`)

## Story
As a game designer,
I want enemies to be able to fire projectiles at the player,
so that the gameplay includes ranged threats and requires more tactical movement.

## Context
US-031 and US-036 established the projectile lifecycle.
US-038 (dependency) establishes how projectiles deal damage based on Factions.
Currently, only the Player has an inventory and weapons. Enemies like `BatEnemy` only deal damage on touch.
We need to bridge the gap so enemies can also own and use Weapons.

## Acceptance Criteria
- **AC1** — `internal/engine/entity/actors/Character` (or a specialized sub-struct) gains an optional `Weapon` or `Inventory` that can be configured via JSON.
- **AC2** — A new enemy type (e.g., `ShooterBat`) or an update to `BatEnemy` implements shooting logic: it fires when the player is within `Range` and a `Cooldown` has passed.
- **AC3** — Enemy weapons spawn projectiles with `FactionEnemy`.
- **AC4** — Enemy projectiles deal damage to `FactionPlayer` (and `FactionNeutral`) actors as defined in US-038.
- **AC5** — The enemy JSON configuration (`bat.json` or similar) includes a `weapon` section with: `projectile_type`, `speed`, `cooldown`, `damage`, and `range`.
- **AC6** — `internal/engine/entity/actors/builder.ConfigureCharacter` is updated to initialize the enemy's weapon from the JSON config.
- **AC7** — Unit tests verify:
  - Enemy fires when player is in range.
  - Enemy does not fire during cooldown.
  - Projectile correctly identifies the enemy as owner and has `FactionEnemy`.

## Proposed Changes
- Update `internal/game/entity/actors/enemies/bat.go` (or create a new one) to handle shooting.
- Update `internal/engine/entity/actors/builder/` to parse weapon configs for enemies.
- Ensure `ProjectileWeapon` can be easily attached to any `Character`.

## Dependencies
- US-038 — Projectile Damage on Hit (Backlog)
