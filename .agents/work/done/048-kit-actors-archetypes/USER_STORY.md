# User Story — 048-kit-actors-archetypes

## Title

Introduce `internal/kit/actors/`: Reusable Character Trait Components

## As a...

Engine developer building platformer games on top of this framework

## I want...

Reusable, genre-level character trait components in `internal/kit/actors/` (shooter behavior, melee behavior, death behavior) that can be independently composed with `PlatformerCharacter`, so concrete game characters can embed only the traits they need instead of duplicating structural boilerplate

## So that...

The three-layer architecture's kit layer delivers its first actor-level reuse through composable building blocks: game packages wire game-specific logic while kit packages own the structural patterns that are genre-reusable across any platformer game built on this engine. Any combination (platformer + shooter, platformer + melee, platformer + shooter + melee) is possible without creating new kit types.

## Background

Story 047 (commit `0e282d4`) completed the migration of all concrete Actor States to `internal/kit/states/`. The kit layer now owns state logic but has no actor-level types.

The current game layer contains structural boilerplate that is not game-specific:

- `internal/game/entity/actors/methods/death_behavior.go` — `PlayerDeathBehavior` only imports `internal/engine/entity/actors/platformer`; it has zero game-specific dependencies and is already kit-eligible as-is.
- `internal/game/entity/actors/enemies/bat.go` (`BatEnemy`) and `internal/game/entity/actors/enemies/wolf.go` (`WolfEnemy`) share an identical structural pattern: hold a `combat.EnemyShooter` and duplicate the same `Update` logic (tick shooter, then character). A reusable `ShooterCharacter` trait in `kit/actors/` can own this pattern; game-specific touch reactions remain in the concrete enemy type.
- `internal/game/entity/actors/player/climber.go` (`ClimberPlayer`) holds a `*meleeengine.Controller` field. A reusable `MeleeCharacter` trait in `kit/actors/` can own this field and accessor; game-specific logic (`Hurt`, `SetInventory`, `SetMelee` wiring with `gamestates`) stays in the concrete player type.

The pre-scoped work from Story 046's "Out of Scope" section framed this as "character archetypes to `internal/kit/actors/` (platformer, shooter, melee)". This story delivers those three traits as **independent, composable components** rather than monolithic archetype bases, avoiding combinatorial explosion (e.g., a future brawler enemy that needs both shooter + melee can embed both traits without requiring a new `PlatformerShooterMeleeEnemy` type).

## Acceptance Criteria

- [ ] **New package exists**: `internal/kit/actors/` is created with a `doc.go` declaring that it imports only `internal/engine/` packages and has zero knowledge of `internal/game/`.
- [ ] **`PlayerDeathBehavior` relocated**: The type moves from `internal/game/entity/actors/methods/death_behavior.go` to `internal/kit/actors/`. The original file is deleted; `internal/game/` imports the type from its new location.
- [ ] **`ShooterCharacter` trait created**: A new `ShooterCharacter` struct in `internal/kit/actors/` holds a `combat.EnemyShooter`. It provides `Shooter() combat.EnemyShooter`, `SetShooter(combat.EnemyShooter)`, and `UpdateShooter()` (ticks `shooter.Update()` if non-nil). It does not contain any game-specific touch logic, faction checks, or platformer-specific methods.
- [ ] **`MeleeCharacter` trait created**: A new `MeleeCharacter` struct in `internal/kit/actors/` holds a `*meleeengine.Controller`. It provides `MeleeController() *meleeengine.Controller` and `SetMeleeController(*meleeengine.Controller)`. It does not contain game-specific input, inventory, or state-name logic.
- [ ] **`BatEnemy` and `WolfEnemy` updated**: Both types embed `*kitactors.ShooterCharacter` rather than duplicating the `shooter` field and `Update` boilerplate. `Update` calls `e.ShooterCharacter.UpdateShooter()` then `e.Character.Update(space)`. All existing tests remain green.
- [ ] **`ClimberPlayer` updated**: Embeds `*kitactors.MeleeCharacter` for the melee controller field. `MeleeController()` is promoted from the embedded trait. All existing tests remain green with no behavior change.
- [ ] **`PlayerDeathBehavior` import updated**: `ClimberPlayer` imports `PlayerDeathBehavior` from `internal/kit/actors/` instead of `internal/game/entity/actors/methods/`.
- [ ] **Dependency rule — engine is clean**: `go list -deps ./internal/engine/...` produces no paths containing `internal/kit` or `internal/game`. CI check must pass.
- [ ] **Dependency rule — kit is clean**: `go list -deps ./internal/kit/...` produces no paths containing `internal/game`. CI check must pass.
- [ ] **All existing tests pass**: `go test ./...` is green with no regressions after all moves and updates.
- [ ] **Coverage does not drop**: `internal/kit/actors/` reaches 80%+ coverage. Overall coverage delta across `internal/engine/`, `internal/kit/`, and `internal/game/` is non-negative.
- [ ] **No game-specific logic in kit**: No type in `internal/kit/actors/` references state names as raw string literals, game-package types, or faction-specific constants. Game wiring stays in `internal/game/`.

## Out of Scope

The following are explicitly deferred:

- **049 (planned)**: Move `internal/engine/combat/` to `internal/kit/combat/`.
- **050 (planned)**: Move `internal/engine/weapon/` to `internal/kit/weapon/`.
- **051 (planned)**: Skills system to `internal/kit/skills/`.
- **052 (planned)**: UI split — primitive widgets stay in `engine/ui/`; menus/HUD patterns move to `kit/ui/`.
- Moving `NPC` base types or AI movement patterns to kit — deferred until a second game's NPC needs are understood.
- Full ADR-006 rewrite — deferred until the layer is populated enough to document with confidence.
- `kit/mocks/` package — not needed until kit defines its own interfaces.
- New gameplay features during this migration.
- Splitting `kit` into sub-modules.
- Combining multiple combat traits in a single enemy (e.g., shooter + melee brawler) — possible with this design but not demonstrated in this story.

## Domain Notes

- An **Actor** is an entity with a state machine (player, enemy). A **trait** is a reusable structural component that holds a single behavioral capability (shooting, melee). Traits live in `kit`; concrete game characters live in `game` and compose traits as needed.
- A **Contract** is a Go interface in `internal/engine/contracts/`. Contracts stay in engine. `kit` implements and composes them; `game` wires them.
- The **three-layer dependency rule**: `engine` imports nothing from `kit` or `game`; `kit` imports only from `engine`; `game` imports from both `engine` and `kit`.
- CI enforcement (already specified in Story 046, must continue to pass):
  ```
  go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)' && exit 1
  go list -deps ./internal/kit/...    | grep    'internal/game'        && exit 1
  ```
- Relevant source files to move or refactor:
  - `internal/game/entity/actors/methods/death_behavior.go` → `internal/kit/actors/death_behavior.go`
  - `internal/game/entity/actors/enemies/bat.go` — update to embed `ShooterCharacter`
  - `internal/game/entity/actors/enemies/wolf.go` — update to embed `ShooterCharacter`
  - `internal/game/entity/actors/player/climber.go` — update to embed `MeleeCharacter`
- `ShooterCharacter` and `MeleeCharacter` are new types with no prior game counterpart; they are introduced in `internal/kit/actors/` and used by the game layer from day one.

## Branch Name

`048-kit-actors-archetypes`
