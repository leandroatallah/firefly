# User Story — 047-migrate-remaining-kit-states

## Title

Migrate Remaining Concrete States to `internal/kit/states/`

## As a...

Engine developer maintaining and extending this codebase across multiple games

## I want...

All remaining genre-reusable concrete Actor States (`walkingSubState`, `duckingSubState`, `aimLockSubState`, `GroundedSubStateEnum`/`groundedSubState` contract + `GroundedInput`, `GroundedState`, and `DashState`) moved from `internal/game/entity/actors/states/` into `internal/kit/states/`, with `MeleeAttackState` also migrated as its adapter imports only engine packages

## So that...

`internal/kit/states/` is the authoritative home for all genre-level platformer states, `internal/game/` contains only project-specific wiring, and the three-layer dependency rule is fully demonstrated across the complete state set.

## Background

Story 046 validated the three-layer architecture (`engine` ← `kit` ← `game`) by relocating `IdleSubState` to `internal/kit/states/` and confirming that the dependency rules held under CI. The states that remain in `internal/game/entity/actors/states/` are:

| File | Type(s) exported / used | Has test file |
|---|---|---|
| `walking_sub_state.go` | `walkingSubState` (unexported) | No |
| `ducking_sub_state.go` | `duckingSubState` (unexported) | No |
| `aim_lock_sub_state.go` | `aimLockSubState` (unexported) | No |
| `grounded_sub_state.go` | `GroundedSubStateEnum`, `groundedSubState` interface, `SubState*` consts | No |
| `grounded_input.go` | `GroundedInput`, `GroundedDeps` | No |
| `grounded_state.go` | `GroundedState`, `StateGrounded` | Yes (`grounded_state_test.go`) |
| `dash_state.go` | `DashState`, `StateDashing`, `DashConfig` | Yes (`dash_state_test.go`) |
| `melee_state.go` | `MeleeAttackState` (type alias for `internal/engine/combat/melee.State`) | Yes (`melee_state_test.go`) |

`melee_state.go` is a thin adapter that re-exports the engine-side melee state. It references `StateGrounded` and `StateMeleeAttack` — both of which will have moved to `kit` by the time this story completes — and imports only `internal/engine/` packages. It is in scope for this migration.

The three sub-states (`walkingSubState`, `duckingSubState`, `aimLockSubState`) and the `groundedSubState` contract are currently unexported and tightly coupled to `GroundedInput`. Moving them alongside `GroundedState` is the natural unit of work.

`GroundedInput` and `GroundedDeps` are engine-contract-only interfaces; they travel with `GroundedState`.

## Acceptance Criteria

- [ ] **`walkingSubState` relocated**: `walking_sub_state.go` (and any new test) moves to `internal/kit/states/`. The original file is deleted; if `internal/game/` needs the type it imports from `kit`.
- [ ] **`duckingSubState` relocated**: `ducking_sub_state.go` moves to `internal/kit/states/`. Original deleted.
- [ ] **`aimLockSubState` relocated**: `aim_lock_sub_state.go` moves to `internal/kit/states/`. Original deleted.
- [ ] **`GroundedSubStateEnum`, `groundedSubState` contract, and `SubState*` constants relocated**: `grounded_sub_state.go` moves to `internal/kit/states/`. Original deleted.
- [ ] **`GroundedInput` and `GroundedDeps` relocated**: `grounded_input.go` moves to `internal/kit/states/`. Original deleted.
- [ ] **`GroundedState` relocated**: `grounded_state.go` and `grounded_state_test.go` move to `internal/kit/states/`. Original files deleted; `internal/game/` imports from `kit`.
- [ ] **`DashState` relocated**: `dash_state.go` and `dash_state_test.go` move to `internal/kit/states/`. Original files deleted; `internal/game/` imports from `kit`.
- [ ] **`MeleeAttackState` adapter relocated**: `melee_state.go` and `melee_state_test.go` move to `internal/kit/states/`. Original files deleted; `internal/game/` imports from `kit`.
- [ ] **`StateGrounded`, `StateDashing`, `StateMeleeAttack` enums accessible from kit**: State-enum `var` declarations travel with their respective state types so `internal/game/` wiring code can import them from `kit`.
- [ ] **Dependency rule — engine is clean**: `go list -deps ./internal/engine/...` produces no paths containing `internal/kit` or `internal/game`. CI check must pass.
- [ ] **Dependency rule — kit is clean**: `go list -deps ./internal/kit/...` produces no paths containing `internal/game`. CI check must pass.
- [ ] **All existing tests pass**: `go test ./...` is green with no regressions after the move.
- [ ] **Coverage does not drop**: `internal/kit/states/` maintains 80%+ coverage (tests travel with types). Overall coverage delta across `internal/engine/` and `internal/game/` is non-negative.
- [ ] **No hardcoded state-name strings in kit**: No duplicate string literals for the same state name may exist across `kit` and `game`. State registration calls (`actors.RegisterState(...)`) that pass a string literal are the single authoritative source for that name — that pattern is acceptable. Any additional raw string references to the same state name elsewhere are not.
- [ ] **`internal/game/entity/actors/states/` retains only game-specific wiring**: After the migration, the only remaining files in that package are those which are genuinely game-specific (e.g., `actor_state_concrete.go`, `actor_state.go`, any future game-unique states). No genre-reusable state logic remains in `game`.

## Out of Scope

The following are explicitly deferred to follow-up stories:

- **048 (planned)**: Character archetypes to `internal/kit/actors/` (platformer, shooter, melee).
- **049 (planned)**: Move `internal/engine/combat/` to `internal/kit/combat/`.
- **050 (planned)**: Move `internal/engine/weapon/` to `internal/kit/weapon/`.
- **051 (planned)**: Skills system to `internal/kit/skills/`.
- **052 (planned)**: UI split — primitive widgets stay in `engine/ui/`; menus/HUD patterns move to `kit/ui/`.
- Full ADR-006 rewrite and new ADR for kit rationale — deferred until the layer is fully populated.
- `kit/mocks/` package — not needed until kit defines its own interfaces (e.g., `Skill`, `WeaponBehavior`).
- New gameplay features during migration.
- Splitting `kit` into sub-modules.
- Moving `offset_toggler.go` (utility, not a state) — assess separately.

## Domain Notes

- A **State** is a named node in an Actor's state machine, implementing a Contract from `internal/engine/contracts/`.
- A **Contract** is a Go interface in `internal/engine/contracts/`. Contracts stay in `engine`. `kit` implements them; `game` wires them.
- The **three-layer dependency rule**: `engine` imports nothing from `kit` or `game`; `kit` imports only from `engine`; `game` imports from both `engine` and `kit`.
- `GroundedInput` is a pure input contract consumed by `GroundedState` and its sub-states. It references only engine contracts (`contractsbody.MovableCollidable`, `physicsmovement.PlatformMovementModel`), so it belongs in `kit` alongside the states it serves.
- `MeleeAttackState` is a type alias (`type MeleeAttackState = meleeengine.State`). Moving the file to `kit` is a package rename of the alias declaration only; the underlying type remains in `internal/engine/combat/melee`.
- Sub-states (`walkingSubState`, `duckingSubState`, `aimLockSubState`) are currently unexported. They may remain unexported inside `internal/kit/states/` — export only if needed by tests or `game` wiring.
- CI enforcement (already in place from story 046; must continue to pass):
  ```
  go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)' && exit 1
  go list -deps ./internal/kit/...    | grep    'internal/game'        && exit 1
  ```
- Relevant source files to move: all `.go` files in `internal/game/entity/actors/states/` that contain genre-reusable state logic, i.e. all files except `actor_state_concrete.go` and `actor_state.go`.

## Branch Name

`047-migrate-remaining-kit-states`
