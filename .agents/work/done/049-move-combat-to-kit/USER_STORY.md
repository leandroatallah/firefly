# User Story — 049-move-combat-to-kit

## Story

**As a** game developer using the engine,
**I want** the combat system to live in `internal/kit/` rather than `internal/engine/combat/`,
**so that** the `engine` layer remains a genre-agnostic foundation and beat-em-up combat concerns are correctly encapsulated at the `kit` layer.

---

## Background

The `engine` layer must not contain genre-specific logic — it should only provide primitives that any game genre could use (physics, entities, scenes, inputs, contracts). Combat (melee swings, projectile weapons, inventory, faction logic, enemy shooting) is a beat-em-up concern, not a generic engine concern.

Currently, `internal/engine/combat/` contains:

| Sub-package | Responsibility |
|---|---|
| `combat/` (root) | `Faction`, `FactionNeutral/Player/Enemy` constants |
| `combat/inventory/` | `Inventory` — weapon collection and ammo tracking |
| `combat/weapon/` | `ProjectileWeapon`, `MeleeWeapon`, `EnemyShooting`, `WeaponFactory` |
| `combat/melee/` | `Controller` and `State` — per-Actor melee swing state machine |
| `combat/projectile/` | `Manager` — projectile lifecycle and physics registration |

The `kit` layer already imports from `internal/engine/combat/` in:
- `internal/kit/actors/melee_character.go` — holds a `*melee.Controller`
- `internal/kit/states/melee_state.go` — aliases `melee.State` and wraps `melee.InstallState`

These are architectural violations: `engine/combat` is a concrete implementation, not a contract, and should not be a dependency of anything other than `game` or `kit`.

The interfaces that belong in `engine` are already correctly placed in `internal/engine/contracts/combat/` (`Weapon`, `Inventory`, `ProjectileManager`, `Damageable`, `Destructible`, `Factioned`, `EnemyShooter`). Those contracts must stay in `engine` — only the concrete implementations move.

---

## Acceptance Criteria

### AC-1: Combat implementations relocate to `internal/kit/combat/`

All concrete implementations from `internal/engine/combat/` are moved to `internal/kit/combat/`, preserving the same sub-package structure:

- `internal/kit/combat/` — `Faction`, `FactionNeutral/Player/Enemy` constants
- `internal/kit/combat/inventory/` — `Inventory` struct
- `internal/kit/combat/weapon/` — `ProjectileWeapon`, `MeleeWeapon`, `EnemyShooting`, weapon factory
- `internal/kit/combat/melee/` — `Controller`, `State`, `InstallState`
- `internal/kit/combat/projectile/` — `Manager`, `ProjectileConfig`, projectile lifecycle

### AC-2: `internal/engine/combat/` is deleted

The package `internal/engine/combat/` (and all sub-packages) no longer exists in the repository. No production file in `internal/engine/` imports from `internal/kit/`.

### AC-3: Contracts remain in `internal/engine/contracts/combat/`

The interfaces `Weapon`, `Inventory`, `ProjectileManager`, `Damageable`, `Destructible`, `Factioned`, and `EnemyShooter` remain unchanged in `internal/engine/contracts/combat/`. The `Faction` type and constants are duplicated/moved to `internal/kit/combat/` only — the contract file in `engine/contracts/combat/weapon.go` may define `Faction` there or import it from kit (engine must not import kit; if `Faction` is needed by contracts, it stays in contracts).

> Note: If `Faction` is referenced by the `contracts/combat/weapon.go` file, it must remain defined in `internal/engine/contracts/combat/` (since engine cannot import kit). The duplicate definition in `internal/engine/combat/faction.go` is what gets removed.

### AC-4: All callers updated

Every file in `internal/kit/` and `internal/game/` that previously imported `internal/engine/combat/...` is updated to import `internal/kit/combat/...` instead. No file outside `internal/engine/` still imports from `internal/engine/combat/`.

Affected callers include (non-exhaustive):
- `internal/kit/actors/melee_character.go`
- `internal/kit/states/melee_state.go`
- `internal/game/app/setup.go`
- `internal/game/entity/actors/player/climber.go`
- `internal/game/entity/actors/player/weapons.go`
- `internal/game/entity/actors/enemies/wolf.go`
- `internal/game/entity/actors/enemies/bat.go`
- `internal/game/entity/items/item_weapon_cannon.go`
- `internal/game/scenes/phases/scene_collision_debug_test.go`

### AC-5: All existing tests pass without modification

All tests that previously lived in `internal/engine/combat/...` are relocated alongside their production code to `internal/kit/combat/...`. Test behaviour is not changed — only import paths are updated. No test logic is rewritten.

### AC-6: Build compiles with zero errors

`go build ./...` completes with no errors after the migration.

### AC-7: Coverage delta is non-negative

The Workflow Gatekeeper confirms that overall coverage across `internal/engine/` and `internal/kit/` does not decrease as a result of this migration.

### AC-8: Dependency rule is enforced

`go list -f '{{.Imports}}' github.com/boilerplate/ebiten-template/internal/engine/...` returns no path that contains `internal/kit` or `internal/game`. The engine layer remains free of kit or game imports.

---

## Out of Scope

- No behaviour changes to the combat system. This is a pure structural relocation.
- No new combat features.
- No changes to `internal/engine/contracts/combat/` interfaces (unless required to resolve the `Faction` type dependency — see AC-3 note).
- No changes to game-layer combat logic (actor update loops, damage values, etc.).

---

## Definition of Done

- [ ] `internal/engine/combat/` and all sub-packages are deleted.
- [ ] `internal/kit/combat/` and all sub-packages exist with the full implementation.
- [ ] All `internal/kit/` and `internal/game/` imports of `engine/combat` are replaced with `kit/combat`.
- [ ] `go build ./...` passes.
- [ ] All tests pass (`go test ./...`).
- [ ] Coverage delta is non-negative (confirmed by Gatekeeper).
- [ ] `internal/engine/` does not import `internal/kit/` (confirmed by Gatekeeper).
