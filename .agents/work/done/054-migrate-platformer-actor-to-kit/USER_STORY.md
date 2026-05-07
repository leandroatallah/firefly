# User Story 054 — Migrate Platformer Actor to Kit

**Branch:** `054-migrate-platformer-actor-to-kit`
**Bounded Context:** Kit (`internal/kit/`)

---

## Story

As an engine developer, I want `internal/engine/entity/actors/platformer/` relocated to `internal/kit/actors/platformer/` with all import references updated, so that genre-reusable concrete actors live in the `kit` layer alongside `beatemup/`, `melee_character.go`, and `shooter_character.go` rather than inside the `engine` layer.

---

## Background

The three-layer architecture (engine → kit → game) requires that genre-reusable concrete actor implementations live in `internal/kit/`, not in `internal/engine/entity/`. The `PlatformerCharacter` struct is a concrete, game-genre-specific actor — it is not an engine primitive. It currently sits in `internal/engine/entity/actors/platformer/`, which violates the layering rule and was flagged at the bottom of story 053 as a follow-up relocation task.

This story is a **pure package relocation**. No behaviour changes, no new logic, no interface modifications. The only deliverables are:

1. Files copied to the new path with the `package` declaration updated.
2. All import paths updated throughout the codebase.
3. Old path deleted.
4. All existing tests pass at the new location.

This story depends on story 053 being stable but is not blocked on its completion.

---

## Acceptance Criteria

### AC-1: Package exists at new path and compiles

Given the source files moved to `internal/kit/actors/platformer/`, when `go build ./internal/kit/actors/platformer/...` is run, then the package compiles without errors.

### AC-2: Old path removed

Given the relocation is complete, when `ls internal/engine/entity/actors/platformer/` is run, then no such directory exists. The old import path `github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer` must not appear in any `.go` file in the repository.

### AC-3: All import references updated

Given the files that previously imported `internal/engine/entity/actors/platformer`, when each file is inspected after the migration, then every import path has been updated to `github.com/boilerplate/ebiten-template/internal/kit/actors/platformer`.

The known consumers at story-write time are:

| File | Old import |
|---|---|
| `internal/game/scenes/phases/scene.go` | `engine/entity/actors/platformer` |
| `internal/game/scenes/phases/player.go` | `engine/entity/actors/platformer` |
| `internal/game/entity/actors/npcs/init_npcs.go` | `engine/entity/actors/platformer` |
| `internal/game/entity/actors/enemies/bat.go` | `engine/entity/actors/platformer` |
| `internal/game/entity/actors/enemies/wolf.go` | `engine/entity/actors/platformer` |
| `internal/game/entity/actors/enemies/init_enemies.go` | `engine/entity/actors/platformer` |
| `internal/game/entity/actors/player/cody.go` | `engine/entity/actors/platformer` |
| `internal/game/entity/actors/player/climber.go` | `engine/entity/actors/platformer` |
| `internal/engine/entity/actors/builder/builder.go` | `engine/entity/actors/platformer` |
| `internal/kit/actors/death_behavior.go` | `engine/entity/actors/platformer` |

A `grep -r "entity/actors/platformer" . --include="*.go"` must return no results after the migration.

### AC-4: Existing platformer tests pass at new location

Given the test file `platformer_test.go` relocated alongside the production code to `internal/kit/actors/platformer/`, when `go test ./internal/kit/actors/platformer/...` is run, then all tests pass with no failures.

The tests to preserve are:

- `TestPlatformerCharacter_SetOnJump`
- `TestPlatformerCharacter_SetOnLand`
- `TestPlatformerCharacter_SetOnFall`
- `TestPlatformerCharacter_OnJump` (with and without handler sub-cases)
- `TestPlatformerCharacter_OnLand` (with and without handler sub-cases)
- `TestPlatformerCharacter_OnFall` (with and without handler sub-cases)
- `TestPlatformerCharacter_CoinCount`
- `TestPlatformerCharacter_MovementBlockers`
- `TestPlatformerCharacter_Fields`

### AC-5: Full test suite passes with no regressions

Given the complete migration, when `go test ./...` is run, then there are no failures and no regressions in any package. Coverage delta across `internal/kit/` must be non-negative.

---

## Behavioural Edge Cases

| Scenario | Expected Behaviour |
|---|---|
| `death_behavior.go` in `internal/kit/actors/` imports the platformer type | Import path updated; the `kit` layer importing `kit` is valid by architecture rules |
| `builder.go` in `internal/engine/entity/actors/builder/` imports platformer | Import path updated to `kit/actors/platformer`; engine importing kit violates layering — if this import is load-bearing it must be resolved (see note below) |
| Any test file importing old path | Test file import path updated to new location |
| `package platformer` declaration | Remains `package platformer`; only the filesystem path and import path change |

> **Note on `builder.go`:** If `internal/engine/entity/actors/builder/builder.go` directly imports the platformer type from the engine layer, moving platformer to kit will create a dependency from engine to kit, which is forbidden by the architecture. The Feature Implementer must inspect `builder.go` and either (a) abstract the dependency behind an interface, (b) move the builder to the `kit` layer, or (c) remove the direct platformer type reference. This is the only non-trivial aspect of the relocation and must be resolved before the story closes.

---

## Out of Scope

- Changes to `PlatformerCharacter` behaviour, fields, or methods.
- Changes to any skill, physics, or state machine logic.
- Adding new tests beyond relocating the existing test file.
- Introducing a `doc.go` (optional, not required for this story).
