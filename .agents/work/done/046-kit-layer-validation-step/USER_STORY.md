# User Story — 046-kit-layer-validation-step

## Title

Introduce `internal/kit/` Layer: Validation Step (Move One Concrete State)

## As a...

Engine developer maintaining and extending this codebase across multiple games

## I want...

A validated proof-of-concept showing that a concrete Actor State (`IdleSubState`) can live in a new `internal/kit/states/` package — importing only `internal/engine/` — with all tests green and CI dependency-layer checks passing

## So that...

The three-layer architecture (`engine` ← `kit` ← `game`) is confirmed structurally sound at minimum viable scale before the full migration begins, and we can abort cheaply if the approach reveals unforeseen coupling.

## Background

ADR-006 currently describes a two-layer architecture (`engine` / `game`). As the codebase has grown to support multiple games, concrete state implementations (idle, walk, dash, duck, etc.) that are genre-reusable have accumulated in `internal/game/entity/actors/states/`. These belong neither to the pure-substrate `engine` nor to the project-specific `game` — they are the missing middle layer.

The plan is to introduce `internal/kit/` as a genre-level library: reusable concrete implementations that depend only on engine contracts, never on game logic. Before migrating all states (6+ files), a single validation step de-risks the entire initiative by moving one self-contained state and confirming the dependency rule holds end-to-end.

The chosen candidate for this validation step is `IdleSubState` (`idle_sub_state.go`), which has no game-specific dependencies, implements an engine Contract, and has clear test coverage.

This story delivers only the validation step. Remaining migration phases are tracked as follow-up stories (see "Out of Scope" below).

## Acceptance Criteria

- [ ] **New package exists**: `internal/kit/states/` is created with a `doc.go` declaring that it imports only `internal/engine/` packages and has zero knowledge of `internal/game/`.
- [ ] **`IdleSubState` relocated**: The type (and its tests) moves from `internal/game/entity/actors/states/` to `internal/kit/states/`. The original file is deleted; `internal/game/` imports the type from its new location.
- [ ] **Dependency rule — engine is clean**: `go list -deps ./internal/engine/...` produces no paths containing `internal/kit` or `internal/game`. CI check must pass.
- [ ] **Dependency rule — kit is clean**: `go list -deps ./internal/kit/...` produces no paths containing `internal/game`. CI check must pass.
- [ ] **All existing tests pass**: `go test ./...` is green with no regressions after the move.
- [ ] **Coverage does not drop**: `internal/kit/states/` reaches 80%+ coverage (tests travel with the type). Overall coverage delta across `internal/engine/` and `internal/game/` is non-negative.
- [ ] **No hardcoded state-name strings in kit**: `IdleSubState` (and any support code) must not reference state names as raw string literals. State wiring remains in `internal/game/`.
- [ ] **ADR-006 amended**: `docs/adr/ADR-006-engine-game-layer-separation.md` Status is updated to `Superseded by 046` and a note links to the forthcoming ADR rewrite. (Full rewrite is a follow-up story.)

## Out of Scope

The following are explicitly deferred to follow-up stories (each on its own branch):

- **047 (planned)**: Migrate remaining concrete states (`WalkingSubState`, `DashState`, `DuckingSubState`, `GroundedState`, `MeleeState`) to `internal/kit/states/`.
- **048 (planned)**: Character archetypes to `internal/kit/actors/` (platformer, shooter, melee).
- **049 (planned)**: Move `internal/engine/combat/` to `internal/kit/combat/`.
- **050 (planned)**: Move `internal/engine/weapon/` to `internal/kit/weapon/`.
- **051 (planned)**: Skills system to `internal/kit/skills/`.
- **052 (planned)**: UI split — primitive widgets stay in `engine/ui/`; menus/HUD patterns move to `kit/ui/`.
- Full ADR-006 rewrite and optional new ADR for kit rationale — deferred until the layer is populated enough to document with confidence.
- `kit/mocks/` package — not needed until kit defines its own interfaces (e.g., `Skill`, `WeaponBehavior`).
- New gameplay features during migration.
- Splitting `kit` into sub-modules.

## Domain Notes

- A **State** is a named node in an Actor's state machine, implementing a Contract from `internal/engine/contracts/`.
- A **Contract** is a Go interface in `internal/engine/contracts/`. Contracts stay in engine. `kit` implements them; `game` wires them.
- The **three-layer dependency rule**: `engine` imports nothing from `kit` or `game`; `kit` imports only from `engine`; `game` imports from both `engine` and `kit`.
- CI enforcement (to be added to the lint/check step):
  ```
  go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)' && exit 1
  go list -deps ./internal/kit/...    | grep    'internal/game'        && exit 1
  ```
- `.agents/constitution.md` Bounded Contexts table must be updated to add a `Kit` row. The Gatekeeper must confirm this is done before closing the story.
- Relevant source file to move: `internal/game/entity/actors/states/idle_sub_state.go` (and its test counterpart).

## Branch Name

`046-kit-layer-validation-step`
