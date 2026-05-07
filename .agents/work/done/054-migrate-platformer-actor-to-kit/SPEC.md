# SPEC — 054 Migrate Platformer Actor to Kit

**Branch:** `054-migrate-platformer-actor-to-kit`
**Bounded Context:** Kit (`internal/kit/`)
**Type:** Pure relocation, no behaviour changes.

---

## 1. Goal

Relocate the concrete `PlatformerCharacter` actor and its `PlatformerActorEntity` interface from the engine layer to the kit layer:

```
internal/engine/entity/actors/platformer/  →  internal/kit/actors/platformer/
```

After the move:
- `package platformer` declaration is unchanged.
- All public symbols (`PlatformerCharacter`, `PlatformerActorEntity`, `AlivePlayer`, `NewPlatformerCharacter`) are unchanged in name, signature, and behaviour.
- The engine layer no longer references the platformer package, directly or transitively (verified by `internal/engine/layering_test.go`).
- All callers compile and all existing tests pass.

---

## 2. Resolution of the `builder.go` engine→kit violation (decided first)

### 2.1 Problem

`internal/engine/entity/actors/builder/builder.go` currently does:

```go
import "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"

func PreparePlatformer(ctx *app.AppContext, jsonPath string) (*platformer.PlatformerCharacter, schemas.SpriteData, actors.StatData, map[string]animation.SpriteState, error) {
    ...
    character, err := platformer.NewPlatformerCharacter(ctx.Assets, stateMap, spriteData, rect)
    ...
}
```

After moving `platformer` to `internal/kit/actors/platformer`, the engine builder would import kit, which is forbidden. The existing layering tests at `internal/engine/layering_test.go` and `internal/engine/entity/actors/builder/builder_layering_test.go` would both fail.

### 2.2 Options considered

| Option | Verdict |
|---|---|
| (a) Abstract the dependency behind an engine-side interface | Rejected. `PreparePlatformer` constructs and returns a concrete `*PlatformerCharacter`. An interface would only hide the construction; the call to `platformer.NewPlatformerCharacter(...)` would still need to live somewhere, and that somewhere cannot be the engine layer. |
| (b) Move the entire `builder` package to kit | Rejected. Most of `builder.go` (`BuildStateMap`, `BodyRectFromSpriteData`, `SetCharacterBodies`, `SetCharacterStats`, `ConfigureCharacter`, `ApplySkills`, the `collisionRectSetter` private interface) is genre-agnostic and operates on `actors.ActorEntity`. Other engine-side helpers (e.g. `configure_enemy_weapon.go` referenced by `internal/engine/layering_test.go` line 26) live in the same builder package and would be wrongly relocated. |
| (c) Remove the direct platformer reference | Same problem as (a): the call to `NewPlatformerCharacter` must exist somewhere. |
| **(d) Split — keep generic helpers in engine builder, move platformer-specific factory to kit** | **Chosen.** Surgically minimal: only `PreparePlatformer` (the single function that actually mentions a platformer type) moves. Existing engine builder helpers are untouched and their unit tests in `builder_test.go` keep working unchanged. |

### 2.3 Decision

**Move only `PreparePlatformer` out of the engine builder.** Place it in the new kit platformer package as `internal/kit/actors/platformer/builder.go` (same package `platformer`). The four call sites in the game layer update their import for this one symbol.

Rationale for placing it inside `internal/kit/actors/platformer/` rather than a sibling `internal/kit/actors/builder/`:
- Co-locates the constructor with the type it constructs.
- Avoids creating a new kit subpackage for a single function.
- Mirrors how `NewPlatformerCharacter` already lives next to `PlatformerCharacter`; `PreparePlatformer` is a thin orchestration over `NewPlatformerCharacter`.

Result:
- Engine builder retains all generic helpers and remains pure (no kit imports).
- Game callers replace `builder.PreparePlatformer(...)` with `platformer.PreparePlatformer(...)`, picking up the new kit path.
- Other engine builder symbols (`ApplyPlatformerPhysics`, `ConfigureCharacter`, `ApplySkills`, `BuildStateMap`, `BodyRectFromSpriteData`, `SetCharacterBodies`, `SetCharacterStats`) keep their import path: `internal/engine/entity/actors/builder`.

Note on `ApplyPlatformerPhysics`: despite its name it only references `actors.ActorEntity` and `physicsmovement` types — no `platformer` import. It stays in engine builder. The naming is left as-is (renaming is out of scope).

---

## 3. File operations

### 3.1 Create

| Path | Source | Notes |
|---|---|---|
| `internal/kit/actors/platformer/platformer.go` | Verbatim copy of `internal/engine/entity/actors/platformer/platformer.go` | `package platformer` unchanged. Imports unchanged (all engine-layer imports remain valid; kit may import engine). |
| `internal/kit/actors/platformer/platformer_test.go` | Verbatim copy of `internal/engine/entity/actors/platformer/platformer_test.go` | `package platformer` unchanged. Test imports unchanged. |
| `internal/kit/actors/platformer/builder.go` | New file containing only the extracted `PreparePlatformer` function | `package platformer`. Imports the same set as in engine `builder.go`, except `engine/entity/actors/platformer` is dropped (same package now). Calls `NewPlatformerCharacter` directly (unqualified). |

The new `builder.go` content (verbatim spec for the Feature Implementer):

```go
package platformer

import (
    "fmt"

    "github.com/boilerplate/ebiten-template/internal/engine/app"
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
    "github.com/boilerplate/ebiten-template/internal/engine/data/jsonutil"
    "github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
    "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
    "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
)

// PreparePlatformer loads sprite and stat data, builds the state map, and initializes a PlatformerCharacter.
func PreparePlatformer(
    ctx *app.AppContext,
    jsonPath string,
) (*PlatformerCharacter, schemas.SpriteData, actors.StatData, map[string]animation.SpriteState, error) {
    spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData](ctx.Assets, jsonPath)
    if err != nil {
        return nil, schemas.SpriteData{}, actors.StatData{}, nil, err
    }

    stateMap, err := builder.BuildStateMap(spriteData)
    if err != nil {
        return nil, schemas.SpriteData{}, actors.StatData{}, nil, err
    }

    rect := builder.BodyRectFromSpriteData(spriteData)
    character, err := NewPlatformerCharacter(ctx.Assets, stateMap, spriteData, rect)
    if err != nil {
        return nil, schemas.SpriteData{}, actors.StatData{}, nil, fmt.Errorf("failed to create platformer character: %w", err)
    }
    character.SetAppContext(ctx)

    return character, spriteData, statData, stateMap, nil
}
```

### 3.2 Modify (engine builder)

**`internal/engine/entity/actors/builder/builder.go`** — delete:
- The `import` of `internal/engine/entity/actors/platformer`.
- The `PreparePlatformer` function (lines 23–46 in the current file).

Everything else stays identical: `collisionRectSetter`, `ApplyPlatformerPhysics`, `SetCharacterBodies`, `SetCharacterStats`, `BuildStateMap`, `BodyRectFromSpriteData`, `ConfigureCharacter`, `ApplySkills`.

**`internal/engine/entity/actors/builder/builder_test.go`** — no changes. It does not exercise `PreparePlatformer` and does not import the platformer package.

**`internal/engine/entity/actors/builder/builder_layering_test.go`** — no changes. Continues to assert the engine builder doesn't import kit. Will now pass on its current assertion because the platformer reference is gone.

### 3.3 Update consumer imports

For all 10 consumers, update the import path. For the four that call `builder.PreparePlatformer`, also rename the call site to `platformer.PreparePlatformer`.

| # | File | Change |
|---|---|---|
| 1 | `internal/game/scenes/phases/scene.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. No call-site change. |
| 2 | `internal/game/scenes/phases/player.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. No call-site change. |
| 3 | `internal/game/entity/actors/npcs/init_npcs.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. No call-site change. |
| 4 | `internal/game/entity/actors/enemies/bat.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. Call: `builder.PreparePlatformer(...)` → `platformer.PreparePlatformer(...)`. The `builder` import remains (still used for `ConfigureCharacter`, `ApplyPlatformerPhysics`, etc., if applicable; verify). |
| 5 | `internal/game/entity/actors/enemies/wolf.go` | Same as bat.go. |
| 6 | `internal/game/entity/actors/enemies/init_enemies.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. No call-site change (init file). |
| 7 | `internal/game/entity/actors/player/cody.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. Call: `builder.PreparePlatformer(...)` → `platformer.PreparePlatformer(...)`. |
| 8 | `internal/game/entity/actors/player/climber.go` | Same as cody.go. |
| 9 | `internal/engine/entity/actors/builder/builder.go` | Import deleted (see §3.2). `PreparePlatformer` deleted. |
| 10 | `internal/kit/actors/death_behavior.go` | Import `engine/entity/actors/platformer` → `kit/actors/platformer`. No call-site change. |

Old import string (must not appear anywhere after the migration):
```
github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer
```

New import string:
```
github.com/boilerplate/ebiten-template/internal/kit/actors/platformer
```

### 3.4 Delete

| Path | Notes |
|---|---|
| `internal/engine/entity/actors/platformer/platformer.go` | After successful copy. |
| `internal/engine/entity/actors/platformer/platformer_test.go` | After successful copy. |
| `internal/engine/entity/actors/platformer/` (directory) | Must not exist after the migration (AC-2). |

---

## 4. Pre-conditions

- Story 053 merged (already merged per `git log`).
- `go test ./...` is green at HEAD before changes start.
- Working tree clean on branch `054-migrate-platformer-actor-to-kit`.

## 5. Post-conditions

- `internal/kit/actors/platformer/` exists with `platformer.go`, `platformer_test.go`, `builder.go`.
- `internal/engine/entity/actors/platformer/` does not exist.
- `grep -r "engine/entity/actors/platformer" . --include="*.go"` returns no results.
- All 9 listed test functions in `platformer_test.go` pass at the new path.
- `internal/engine/layering_test.go::TestEngineLayerHasNoKitOrGameDependencies/kit` passes.
- `internal/engine/entity/actors/builder/builder_layering_test.go::TestBuilderDoesNotImportKit` passes.
- Coverage delta across `internal/kit/` is non-negative (no tests were dropped or watered down).

---

## 6. Integration points (within the bounded contexts touched)

| Context | What changes |
|---|---|
| Engine — `entity/actors/builder` | Loses `PreparePlatformer` and the platformer import. All other exports unchanged. |
| Engine — `entity/actors/platformer` | Deleted. |
| Kit — `actors/platformer` (new) | Owns `PlatformerCharacter`, `PlatformerActorEntity`, `AlivePlayer`, `NewPlatformerCharacter`, `PreparePlatformer`. |
| Kit — `actors/death_behavior.go` | Import path updated. Compiles unchanged. |
| Game — phases, npcs, enemies, player | Import paths updated; four files also rename the `builder.PreparePlatformer` call to `platformer.PreparePlatformer`. |

No public API additions (`PreparePlatformer` keeps the same signature, just moves package). No new contracts in `internal/engine/contracts/`.

---

## 7. Red Phase (failing test scenario)

This story is a refactor. The classic Red→Green cycle is realised by **architectural assertions**, not new behavioural tests. The failing tests that must drive the change are the ones already in the codebase that will start failing the moment the platformer package is moved without resolving the builder violation:

### 7.1 Red — pre-implementation state demonstrating the gap

Before implementation, run:

```bash
go test ./internal/engine/...
```

…with a hypothetical naive move (just relocate the directory and update all import strings literally). The following will fail:

1. **`internal/engine/layering_test.go::TestEngineLayerHasNoKitOrGameDependencies/kit`** — fails because `internal/engine/entity/actors/builder/builder.go` would now import `internal/kit/actors/platformer`.
2. **`internal/engine/entity/actors/builder/builder_layering_test.go::TestBuilderDoesNotImportKit`** — fails for the same reason (string match on `"internal/kit/`).

These two pre-existing tests are the Red Phase: they encode the architectural rule the migration must satisfy.

### 7.2 Green — what makes them pass

Applying §2.3 (extract `PreparePlatformer` to the kit package) removes the offending import from the engine builder. Both tests then pass.

### 7.3 Behavioural regression guard

The 9 platformer test functions (listed in the user story under AC-4) move with the source. They are the behavioural safety net: any accidental change to `PlatformerCharacter` semantics during the move surfaces as a test failure.

```
TestPlatformerCharacter_SetOnJump
TestPlatformerCharacter_SetOnLand
TestPlatformerCharacter_SetOnFall
TestPlatformerCharacter_OnJump
TestPlatformerCharacter_OnLand
TestPlatformerCharacter_OnFall
TestPlatformerCharacter_CoinCount
TestPlatformerCharacter_MovementBlockers
TestPlatformerCharacter_Fields
```

No new tests are introduced. Adding tests is explicitly out of scope per the user story.

---

## 8. Verification commands (per AC)

| AC | Command | Expected |
|---|---|---|
| AC-1 | `go build ./internal/kit/actors/platformer/...` | exit 0 |
| AC-2 | `test ! -d internal/engine/entity/actors/platformer && grep -r "engine/entity/actors/platformer" . --include="*.go"` | first cmd exit 0; grep prints nothing, exit 1 |
| AC-3 | `grep -rn "kit/actors/platformer" . --include="*.go" \| wc -l` then spot-check each of the 10 consumer files | count ≥ 10 (the 10 consumer files plus the new package's own `builder.go`'s package-relative usage; `cody.go` and `climber.go` and `bat.go` and `wolf.go` will also reference `platformer.PreparePlatformer`) |
| AC-4 | `go test ./internal/kit/actors/platformer/...` | PASS, all 9 tests run |
| AC-5 | `go test ./...` and `go test ./internal/engine/... -run TestEngineLayerHasNoKitOrGameDependencies` and `go test ./internal/engine/entity/actors/builder/... -run TestBuilderDoesNotImportKit` | all PASS |

Coverage check (informational, for the Gatekeeper):
```bash
go test -cover ./internal/kit/...
```
The number from before vs after should be non-negative delta.

---

## 9. Out of scope (re-stated)

- No changes to `PlatformerCharacter` fields, methods, or behaviour.
- No new interfaces (the interface returned by `PreparePlatformer` is still `*PlatformerCharacter`; the type just lives in a new package).
- No renaming (`ApplyPlatformerPhysics` keeps its name even though it stays in engine builder).
- No `doc.go` for the new kit platformer package.
- No relocation of any other engine builder helper.
- No changes to `internal/kit/actors/death_behavior.go` beyond the import path.
