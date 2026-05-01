# SPEC — 049-move-combat-to-kit

**Branch:** `049-move-combat-to-kit`
**Type:** Structural refactor (no behaviour change)
**Bounded Contexts touched:** `engine/combat` (deleted), `kit` (new sub-tree), `engine/entity/actors`, `engine/entity/actors/builder`, `engine/app`, `game/*`

---

## 1. Goal

Relocate every concrete combat implementation from `internal/engine/combat/` to `internal/kit/combat/`. Eliminate the `engine → engine/combat` couplings that currently violate the layering rule (engine must not depend on genre-specific code). The migration is purely structural — no public API surface, function signature, or test assertion changes — except where required to remove the engine-layer dependency on the moved package.

## 2. Pre-Conditions

- `internal/engine/combat/` exists with sub-packages `inventory`, `weapon`, `melee`, `projectile`, plus `faction.go` and `faction_test.go` at root.
- `internal/engine/contracts/combat/` already defines: `Faction`, `FactionNeutral`, `FactionPlayer`, `FactionEnemy`, `Factioned`, `Weapon`, `Inventory`, `ProjectileManager`, `Damageable`, `Destructible`, `EnemyShooter`, `ShootMode*` constants.
- `internal/engine/combat/faction.go` contains a duplicate definition of `Faction` and the three constants (already mirrored in `contracts/combat/weapon.go`).
- `internal/kit/combat/` does NOT yet exist.
- `go build ./...` and `go test ./...` are green at HEAD.

## 3. Post-Conditions

- `internal/engine/combat/` is removed entirely (directory absent).
- `internal/kit/combat/` exists with sub-packages `inventory`, `weapon`, `melee`, `projectile` and a root package containing `Faction` aliases.
- `internal/engine/contracts/combat/` is unchanged in API. The canonical `Faction` type lives only there.
- No file under `internal/engine/` imports `internal/kit/...` or the (now-deleted) `internal/engine/combat/...`.
- All callers in `internal/kit/` and `internal/game/` import `internal/kit/combat/...` instead of `internal/engine/combat/...`.
- `go build ./...` passes. `go test ./...` passes. Coverage delta ≥ 0 across `internal/engine/` + `internal/kit/`.

## 4. File Mapping (old path → new path)

All files keep their package name and contents. Only the import path of the package changes (because it moves to a new directory). Internal cross-references between sub-packages must be updated to the new module path.

### 4.1 Files that move verbatim (only update internal import paths)

| Old path | New path |
|---|---|
| `internal/engine/combat/inventory/inventory.go` | `internal/kit/combat/inventory/inventory.go` |
| `internal/engine/combat/inventory/inventory_test.go` | `internal/kit/combat/inventory/inventory_test.go` |
| `internal/engine/combat/inventory/README.md` | `internal/kit/combat/inventory/README.md` |
| `internal/engine/combat/weapon/weapon.go` | `internal/kit/combat/weapon/weapon.go` |
| `internal/engine/combat/weapon/weapon_test.go` | `internal/kit/combat/weapon/weapon_test.go` |
| `internal/engine/combat/weapon/factory.go` | `internal/kit/combat/weapon/factory.go` |
| `internal/engine/combat/weapon/factory_test.go` | `internal/kit/combat/weapon/factory_test.go` |
| `internal/engine/combat/weapon/melee.go` | `internal/kit/combat/weapon/melee.go` |
| `internal/engine/combat/weapon/melee_test.go` | `internal/kit/combat/weapon/melee_test.go` |
| `internal/engine/combat/weapon/melee_is_swinging_test.go` | `internal/kit/combat/weapon/melee_is_swinging_test.go` |
| `internal/engine/combat/weapon/melee_hitbox_rect_test.go` | `internal/kit/combat/weapon/melee_hitbox_rect_test.go` |
| `internal/engine/combat/weapon/enemy_shooting.go` | `internal/kit/combat/weapon/enemy_shooting.go` |
| `internal/engine/combat/weapon/enemy_shooting_test.go` | `internal/kit/combat/weapon/enemy_shooting_test.go` |
| `internal/engine/combat/weapon/mocks_test.go` | `internal/kit/combat/weapon/mocks_test.go` |
| `internal/engine/combat/weapon/README.md` | `internal/kit/combat/weapon/README.md` |
| `internal/engine/combat/melee/controller.go` | `internal/kit/combat/melee/controller.go` |
| `internal/engine/combat/melee/controller_test.go` | `internal/kit/combat/melee/controller_test.go` |
| `internal/engine/combat/melee/state.go` | `internal/kit/combat/melee/state.go` |
| `internal/engine/combat/melee/state_test.go` | `internal/kit/combat/melee/state_test.go` |
| `internal/engine/combat/projectile/manager.go` | `internal/kit/combat/projectile/manager.go` |
| `internal/engine/combat/projectile/manager_test.go` | `internal/kit/combat/projectile/manager_test.go` |
| `internal/engine/combat/projectile/manager_debug_test.go` | `internal/kit/combat/projectile/manager_debug_test.go` |
| `internal/engine/combat/projectile/projectile.go` | `internal/kit/combat/projectile/projectile.go` |
| `internal/engine/combat/projectile/projectile_test.go` | `internal/kit/combat/projectile/projectile_test.go` |
| `internal/engine/combat/projectile/config.go` | `internal/kit/combat/projectile/config.go` |
| `internal/engine/combat/projectile/config_test.go` | `internal/kit/combat/projectile/config_test.go` |
| `internal/engine/combat/projectile/damage_test.go` | `internal/kit/combat/projectile/damage_test.go` |
| `internal/engine/combat/projectile/faction_config_test.go` | `internal/kit/combat/projectile/faction_config_test.go` |
| `internal/engine/combat/projectile/friendly_fire_test.go` | `internal/kit/combat/projectile/friendly_fire_test.go` |
| `internal/engine/combat/projectile/friendly_fire_integration_test.go` | `internal/kit/combat/projectile/friendly_fire_integration_test.go` |
| `internal/engine/combat/projectile/interceptable_config_test.go` | `internal/kit/combat/projectile/interceptable_config_test.go` |
| `internal/engine/combat/projectile/mocks_test.go` | `internal/kit/combat/projectile/mocks_test.go` |
| `internal/engine/combat/projectile/README.md` | `internal/kit/combat/projectile/README.md` |
| `internal/engine/combat/README.md` | `internal/kit/combat/README.md` |

### 4.2 Files that move with content modification

| Old path | New path | Change |
|---|---|---|
| `internal/engine/combat/faction.go` | `internal/kit/combat/faction.go` | Replace duplicate `Faction` definition with type aliases to `contracts/combat`. See §5.1. |
| `internal/engine/combat/faction_test.go` | `internal/kit/combat/faction_test.go` | Update import path; assertions remain identical. |

### 4.3 Internal package import rewrites (within moved files)

Inside the relocated tree, any import of `github.com/boilerplate/ebiten-template/internal/engine/combat/<sub>` must be rewritten to `github.com/boilerplate/ebiten-template/internal/kit/combat/<sub>`. Known cases:

- `internal/kit/combat/melee/controller.go` — currently imports `engine/combat/weapon` → must become `kit/combat/weapon`.
- Any test file that imports a sibling sub-package — same rewrite rule.

Imports of `engine/contracts/...` inside the moved files are kept as-is (kit may import engine).

## 5. Engine-Layer Refactors (CRITICAL — these break the layering today)

Four engine files currently reference `engine/combat`. After the move, they CANNOT switch to `kit/combat` (engine must not import kit). They must be refactored to depend only on `engine/contracts/combat`.

### 5.1 `internal/kit/combat/faction.go` — replace duplicate type

The old `engine/combat/faction.go` redefined `Faction` and its constants as concrete types in package `combat`. The canonical definitions live in `engine/contracts/combat`. The new file at `internal/kit/combat/faction.go` MUST NOT redefine the type — it must alias the contract types so existing kit/game code that uses `kitcombat.Faction`, `kitcombat.FactionPlayer`, etc., keeps compiling:

```go
package combat

import contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"

// Faction aliases the canonical contracts/combat.Faction type so that kit
// callers can refer to `kitcombat.Faction` without importing the contracts
// package directly.
type Faction = contractscombat.Faction

const (
    FactionNeutral = contractscombat.FactionNeutral
    FactionPlayer  = contractscombat.FactionPlayer
    FactionEnemy   = contractscombat.FactionEnemy
)
```

`faction_test.go` continues to test these symbols via the kit import path; behaviour is unchanged because the underlying type is identical.

### 5.2 `internal/engine/entity/actors/character.go`

Current:
```go
enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"
...
faction enginecombat.Faction
func (c *Character) Faction() enginecombat.Faction { ... }
func (c *Character) SetFaction(f enginecombat.Faction) { ... }
```

Replace with import of `contracts/combat`:
```go
contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
...
faction contractscombat.Faction
func (c *Character) Faction() contractscombat.Faction { ... }
func (c *Character) SetFaction(f contractscombat.Faction) { ... }
```

Because `kit/combat.Faction` is a type alias for `contracts/combat.Faction` (§5.1), all callers passing `kitcombat.FactionPlayer` continue to compile.

### 5.3 `internal/engine/entity/actors/character_damage_test.go`

Same rewrite as §5.2 — switch the alias `enginecombat` to `contractscombat` and update each `enginecombat.Faction*` reference accordingly. No assertion changes.

### 5.4 `internal/engine/app/context.go`

Current:
```go
combatprojectile "github.com/boilerplate/ebiten-template/internal/engine/combat/projectile"
...
ProjectileManager *combatprojectile.Manager
```

Replace the concrete type with the existing interface:
```go
contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
...
ProjectileManager contractscombat.ProjectileManager
```

Justification: `contracts/combat.ProjectileManager` already declares `SpawnProjectile(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{})`, which `*projectile.Manager` satisfies. Game-layer code that previously assigned `&projectile.Manager{...}` to this field will continue to work because the kit type still satisfies the interface; call sites that relied on concrete fields/methods of `*projectile.Manager` (if any) must be audited and updated to use the new kit import. See §6.4.

### 5.5 `internal/engine/entity/actors/builder/configure_enemy_weapon.go`

This file currently calls `weapon.NewProjectileWeapon(...)` and `weapon.NewEnemyShooting(...)` from `engine/combat/weapon`. It is in the engine layer and instantiates concrete combat — this is the deepest layering violation.

**Resolution: relocate the file to the kit layer.**

| Old path | New path |
|---|---|
| `internal/engine/entity/actors/builder/configure_enemy_weapon.go` | `internal/kit/combat/weapon/configure_enemy.go` |

Rationale: this is genre-specific wiring (it parses `EnemyWeaponConfig` from data schemas and produces an `EnemyShooter`) — it belongs in the kit. The function signature and behaviour are unchanged; only the package and import path move. Callers that referenced `builder.ConfigureEnemyWeapon` must be updated to the new package. (Search confirms callers live in `internal/game/`; verify during implementation and update accordingly.)

If a sibling test file `configure_enemy_weapon_test.go` exists, it moves with the production file.

> Implementation note: if relocation proves disruptive, an acceptable alternative is to introduce a `WeaponFactory` contract in `engine/contracts/combat/` and inject the kit-layer factory at app wiring time. The TDD Specialist should prefer the simple relocation approach above; only fall back to the factory contract if a caller is provably outside the kit/game scope.

## 6. Caller Updates (kit + game layers)

Every file below currently imports a path under `internal/engine/combat/...` and must switch to the equivalent path under `internal/kit/combat/...`. No symbol renames are required.

### 6.1 Kit layer

| File | Old import(s) | New import(s) |
|---|---|---|
| `internal/kit/actors/melee_character.go` | `engine/combat/melee` | `kit/combat/melee` |
| `internal/kit/actors/melee_character_test.go` | `engine/combat/melee` (if any) | `kit/combat/melee` |
| `internal/kit/states/melee_state.go` | `engine/combat/melee` | `kit/combat/melee` |
| `internal/kit/states/melee_state_test.go` | `engine/combat/melee` | `kit/combat/melee` |

### 6.2 Game layer (production)

| File | Old import path roots → New |
|---|---|
| `internal/game/app/setup.go` | `engine/combat/...` → `kit/combat/...` |
| `internal/game/entity/actors/player/climber.go` | same |
| `internal/game/entity/actors/player/weapons.go` | same |
| `internal/game/entity/actors/enemies/wolf.go` | same |
| `internal/game/entity/actors/enemies/bat.go` | same |
| `internal/game/entity/items/item_weapon_cannon.go` | same |

### 6.3 Game layer (tests)

| File | Old → New |
|---|---|
| `internal/game/scenes/phases/scene_collision_debug_test.go` | `engine/combat/...` → `kit/combat/...` |
| `internal/game/entity/actors/player/climber_test.go` | same |
| `internal/game/entity/actors/enemies/wolf_test.go` | same |
| `internal/game/entity/actors/enemies/bat_test.go` | same |

### 6.4 Builder caller fix-up

Any file importing `internal/engine/entity/actors/builder` and calling `ConfigureEnemyWeapon` must be updated to import `internal/kit/combat/weapon` and call `weapon.ConfigureEnemy` (or whatever the new exported name is). Verify the exact callers during implementation.

## 7. Integration Points

- **`engine/contracts/combat/`**: untouched. Public API stable.
- **`engine/entity/actors/Character`**: API surface (`Faction()`, `SetFaction(f)`) preserved by the type alias trick (§5.1). Existing kit/game callers compile without change.
- **`engine/app.AppContext.ProjectileManager`**: type changes from `*combatprojectile.Manager` to `contractscombat.ProjectileManager` (interface). All current production assignments use `kit.combat.projectile.Manager`, which satisfies the interface.
- **Builder for enemy weapons**: relocates from engine to kit. Callers in game must update their import path.

## 8. Migration Order (recommended sequence for the Feature Implementer)

1. Create `internal/kit/combat/` directory tree.
2. Copy every file from §4.1 to its new location, replacing the in-file `engine/combat/<sub>` self-imports with `kit/combat/<sub>`.
3. Add `internal/kit/combat/faction.go` per §5.1 and copy `faction_test.go` with import path update.
4. Refactor the four engine files in §5.2–5.5:
   - `engine/entity/actors/character.go` and its test — switch to `contracts/combat`.
   - `engine/app/context.go` — switch the field type to the contract interface.
   - Move `engine/entity/actors/builder/configure_enemy_weapon.go` to `kit/combat/weapon/configure_enemy.go` and update its package declaration.
5. Update every caller listed in §6 to the new import paths.
6. Delete `internal/engine/combat/` (entire subtree).
7. Run `go build ./...` then `go test ./...`.
8. Run `go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)'` — expect empty output.

## 9. Red Phase Scenario

Because this is a structural refactor with zero new behaviour, the Red Phase is expressed as **structural / build assertions**, not new logic tests. The TDD Specialist must add the following failing checks before any code is moved.

### 9.1 Layering test (new)

Create `internal/engine/layering_test.go` (package `engine_test` or similar) that uses `golang.org/x/tools/go/packages` (or equivalent) to assert:

> For every package P under `github.com/boilerplate/ebiten-template/internal/engine/...`, the transitive import set of P MUST NOT contain any path beginning with `github.com/boilerplate/ebiten-template/internal/kit` or `github.com/boilerplate/ebiten-template/internal/game`.

This test fails today because — after the move — engine packages would (incorrectly) import kit, OR (before the move) `engine/combat` exists and creates the duplicate `Faction` problem. The test passes only when the engine layer is fully decoupled.

If `golang.org/x/tools` is not desirable, an equivalent test may shell out to `go list -deps -f '{{.ImportPath}}' ./internal/engine/...` and grep the output. Implementation should prefer whichever pattern is already used elsewhere in the repo.

### 9.2 Package-existence assertions (build-driven)

The TDD Specialist must add a single test file `internal/kit/combat/package_test.go` that compiles only if the kit-combat sub-tree exists with the expected exported symbols:

```go
package combat_test

import (
    _ "github.com/boilerplate/ebiten-template/internal/kit/combat"
    _ "github.com/boilerplate/ebiten-template/internal/kit/combat/inventory"
    _ "github.com/boilerplate/ebiten-template/internal/kit/combat/melee"
    _ "github.com/boilerplate/ebiten-template/internal/kit/combat/projectile"
    _ "github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

// TestKitCombatPackagesExist is a compile-time assertion that all kit-combat
// sub-packages exist and are importable. The test body is intentionally empty.
func TestKitCombatPackagesExist(_ *testing.T) {}
```

Initial state: the file is added → `go build ./...` fails because the imports do not resolve. After the implementation moves the files, the test passes.

### 9.3 Non-existence of old path

Add `internal/engine/combat_absent_test.go` (or fold into the layering test) asserting that the directory `internal/engine/combat` is absent on disk. Use `os.Stat` with `errors.Is(err, fs.ErrNotExist)`. Initial state: directory exists → test fails. After deletion: passes.

### 9.4 Behavioural regression coverage

No new behavioural tests are required. The existing 30+ tests inside the moved sub-packages (faction, inventory, weapon, melee controller/state, projectile manager / friendly fire / damage / config) provide the regression net. They MUST continue to pass after relocation with zero edits to assertions — only their import paths change.

## 10. Acceptance Mapping

| AC | Verified by |
|---|---|
| AC-1 (kit/combat exists with correct sub-packages) | §9.2 + file listing |
| AC-2 (engine/combat deleted) | §9.3 |
| AC-3 (contracts unchanged) | Diff review of `internal/engine/contracts/combat/` (no changes expected) |
| AC-4 (all callers updated) | §6 + `go build ./...` |
| AC-5 (existing tests pass unchanged) | `go test ./...` over relocated packages |
| AC-6 (build compiles) | `go build ./...` |
| AC-7 (coverage delta ≥ 0) | Gatekeeper coverage diff |
| AC-8 (engine has no kit/game imports) | §9.1 |

## 11. Risks & Mitigations

| Risk | Mitigation |
|---|---|
| Hidden type-identity assumption: code does `var f enginecombat.Faction = someContractsFaction` relying on identity. | Type alias in §5.1 preserves identity (`type Faction = ...`). Identity is exact, not just structural. |
| `*combatprojectile.Manager` was assigned to `AppContext.ProjectileManager` and later type-asserted/dereferenced as concrete in game code. | Audit all reads of `appCtx.ProjectileManager`. If a concrete-only method is required, add it to the contract or have callers import `kit/combat/projectile` directly. Default expectation: only `SpawnProjectile` is called via this field. |
| `ConfigureEnemyWeapon` may be referenced by tooling/docs at the old path. | Update README cross-links; ADRs are not affected. |
| Test files inside moved sub-packages may import sibling sub-packages by old path. | §4.3 mandates rewrite; CI's failing build will catch any miss. |

## 12. Definition of Done (Spec-side)

- SPEC.md present in active folder.
- PROGRESS.md updated with `[x] Spec Engineer` and a `[FINISHED]` log entry.
- Mock Generator is NOT required for this story (no new contracts introduced — the only contract reused, `ProjectileManager`, already exists; the `Faction` type is unchanged).
