# Technical Specification — 050-move-weapon-to-kit

## Branch

`050-move-weapon-to-kit`

## Summary

Pure relocation refactor (no behaviour change), splitting the current `internal/engine/physics/skill/` package into two genre-aware homes:

1. **Engine layer — infrastructure.** Promote the genre-agnostic skill primitives (`Skill`, `ActiveSkill` interfaces, `SkillBase` struct + accessors, `SkillState` enum + constants, and the registry/Set type with `Add`/`Get`/`Update`/`ActiveCount`) from `internal/engine/physics/skill/` into a new top-level engine package `internal/engine/skill/`. This package sits alongside `internal/engine/entity/` and `internal/engine/scene/` — it is **not** under `internal/engine/contracts/` because it carries timing/cooldown logic and a registry, not just thin interfaces.
2. **Kit layer — concretes.** Move all concrete skill implementations and the factory from `internal/engine/physics/skill/` into a new kit package `internal/kit/skills/` (package name `kitskills`).
3. **Rewire** all importers in `internal/engine/entity/actors/{character,platformer,builder}` (engine layer) to use `internal/engine/skill`, and importers in `internal/game/` to use `internal/kit/skills` (for concretes/factory) and `internal/engine/skill` (for the interface).
4. **Delete** `internal/engine/physics/skill/` once empty and add a directory-absent guard test.

## Bounded Context

- **Engine / Physics → Engine / Skill** (infrastructure promotion to a sibling top-level package).
- **Engine / Physics → Kit / Skills** (concrete relocation).
- Layering rule reaffirmed: `internal/engine/**` MUST NOT import `internal/kit/**` or `internal/game/**`.
- Layering rule for kit: `internal/kit/**` MAY import `internal/engine/**` and MUST NOT import `internal/game/**`.

## Why `internal/engine/skill/` and not `internal/engine/contracts/skill/`

`internal/engine/contracts/` is reserved for thin interface declarations with no logic and no state (e.g. `body.MovableCollidable`, `combat.Inventory`). `SkillBase` carries:

- Mutable state (`state`, `timer`, `duration`, `cooldown`, `speed`)
- Tick logic (`Update`, `IsActive`)
- Accessor/mutator behaviour

The registry/Set type carries collection logic (`Add`, `Get`, `Update` over all members, `ActiveCount`).

The right structural analogy is `internal/engine/entity/` (holds the `Body` and `Actor` structs) or `internal/engine/scene/` (holds the scene lifecycle infrastructure), not `internal/engine/contracts/`. The engine destination is therefore `internal/engine/skill/` as a top-level engine sub-package.

## Package Layout (Target)

### New: `internal/engine/skill/`

Package name: `skill`.

Files:

- `skill.go` — `SkillState` enum + constants; `Skill`, `ActiveSkill` interfaces; `SkillBase` struct + accessors.
- `set.go` — registry/Set type (`Set` or `Registry` — see decision below) with `Add`, `Get`, `Update`, `ActiveCount`.
- `doc.go` — package documentation describing layering and intent.
- `skill_test.go` — unit tests for `SkillBase` lifecycle and `Set` operations (migrated from `engine/physics/skill/skill_test.go` and `skill_lifecycle_test.go` *minus* concrete-skill scenarios, which move to kit).

Public surface:

```go
package skill

import (
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
    physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
    "github.com/hajimehoshi/ebiten/v2"
)

// SkillState enumerates lifecycle states.
type SkillState string

const (
    StateReady    SkillState = "ready"
    StateActive   SkillState = "active"
    StateCooldown SkillState = "cooldown"
)

// Skill is the genre-agnostic skill contract.
type Skill interface {
    Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel)
    IsActive() bool
}

// ActiveSkill adds player input + key binding semantics on top of Skill.
type ActiveSkill interface {
    Skill
    HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace)
    ActivationKey() ebiten.Key
}

// SkillBase is the shared, embeddable infrastructure for concrete skills.
// Fields are unexported; concrete skills (in kit) interact via accessors.
type SkillBase struct {
    state    SkillState
    duration int
    cooldown int
    speed    int
    timer    int
}

func (s *SkillBase) State() SkillState           { /* ... */ }
func (s *SkillBase) SetState(st SkillState)      { /* ... */ }
func (s *SkillBase) Duration() int               { /* ... */ }
func (s *SkillBase) SetDuration(d int)           { /* ... */ }
func (s *SkillBase) Cooldown() int               { /* ... */ }
func (s *SkillBase) SetCooldown(c int)           { /* ... */ }
func (s *SkillBase) Speed() int                  { /* ... */ }
func (s *SkillBase) SetSpeed(sp int)             { /* ... */ }
func (s *SkillBase) Timer() int                  { /* ... */ }
func (s *SkillBase) SetTimer(t int)              { /* ... */ }
func (s *SkillBase) IncTimer()                   { /* ... */ }

func (s *SkillBase) Update(body.MovableCollidable, *physicsmovement.PlatformMovementModel) {}
func (s *SkillBase) IsActive() bool { return s.state == StateActive }
```

The accessor set (`State/SetState`, `Duration/SetDuration`, `Cooldown/SetCooldown`, `Speed/SetSpeed`, `Timer/SetTimer/IncTimer`) is the **minimum** required so that concrete skills in `internal/kit/skills/` can drive lifecycle transitions without touching unexported fields. Fields stay unexported to preserve encapsulation.

#### Registry / Set type

Public surface (preserves current behaviour from `engine/physics/skill/skill.go`):

```go
package skill

// Set is the per-actor registry of skills.
type Set struct { /* unexported fields */ }

func NewSet() *Set
func (s *Set) Add(sk Skill)
func (s *Set) Get(key ebiten.Key) (ActiveSkill, bool) // looks up by ActivationKey
func (s *Set) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel)
func (s *Set) ActiveCount() int
func (s *Set) All() []Skill // read-only access used by player state contributors
```

> Decision — name `Set` (not `Registry`). Matches the conventional Go vocabulary used elsewhere in the codebase (`Set` operations are minimal: add/iterate/lookup) and is shorter at call sites (`*skill.Set`).

If the existing type in `engine/physics/skill/` is currently named differently (e.g. `Skills` or unnamed/inlined inside `Actor`), **preserve** the existing exported name during the move and only rename if it is unexported or inconsistent. Implementer must read the current code first and update this spec via an addendum if a rename is required.

### New: `internal/kit/skills/`

Package name: `kitskills` (matches `internal/kit/states/` → `package kitstates` convention; avoids stutter at call sites: `kitskills.DashSkill`).

Files (one-to-one move from `engine/physics/skill/`, with updated imports and embedding):

| Source (engine) | Destination (kit) |
|---|---|
| `skill_shooting.go` | `internal/kit/skills/shooting.go` |
| `skill_dash.go` | `internal/kit/skills/dash.go` |
| `skill_platform_jump.go` | `internal/kit/skills/platform_jump.go` |
| `skill_platform_move.go` | `internal/kit/skills/platform_move.go` |
| `offset_toggler.go` | `internal/kit/skills/offset_toggler.go` |
| `factory.go` | `internal/kit/skills/factory.go` |
| `skill_shooting_test.go` | `internal/kit/skills/shooting_test.go` |
| `skill_shooting_eight_directions_test.go` | `internal/kit/skills/shooting_eight_directions_test.go` |
| `skill_lifecycle_test.go` (concrete-skill scenarios only) | `internal/kit/skills/lifecycle_test.go` |
| `skill_platform_jump_test.go` | `internal/kit/skills/platform_jump_test.go` |
| `factory_test.go` | `internal/kit/skills/factory_test.go` |
| `README.md` | `internal/kit/skills/README.md` |

After move:

- Each file's `package skill` declaration becomes `package kitskills`.
- Each concrete embeds `skill.SkillBase` (alias-imported as `engineskill` if needed for clarity).
- Each direct unexported field access (`s.state = …`, `s.timer++`, `s.duration`) is replaced with the accessor methods on `SkillBase`.

`doc.go`:

```go
// Package kitskills contains genre-reusable concrete Skill implementations
// (shooting, dash, platformer jump, platformer horizontal move) plus the
// FromConfig factory that wires them from a SkillsConfig.
//
// Dependency rule (enforced by layering tests):
//   - kitskills MAY import internal/engine/...
//   - kitskills MUST NOT import internal/game/...
package kitskills
```

`FromConfig` and `SkillDeps` keep their current shape, only the return type's import path changes:

```go
package kitskills

import (
    "github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
    "github.com/boilerplate/ebiten-template/internal/engine/skill"
)

type SkillDeps struct {
    Inventory         /* unchanged */
    ProjectileManager /* unchanged */
    OnJump            /* unchanged */
    EventManager      /* unchanged */
}

func FromConfig(cfg *schemas.SkillsConfig, deps SkillDeps) []skill.Skill
```

`SkillDeps` field set is unchanged from the current implementation. (Story 051 may extend it; not this story.)

### Removed: `internal/engine/physics/skill/`

The directory is deleted in its entirety. If `README.md` contains prose still relevant after the split, fold the engine-infrastructure portions into a brief `internal/engine/skill/README.md` and the genre-specific portions into `internal/kit/skills/README.md`. New documentation is not a hard requirement of this story.

## Interface Contracts

The `Skill` and `ActiveSkill` interface signatures are **unchanged**. Only the import path moves:

- Before: `import "…/internal/engine/physics/skill"` → `skill.Skill`
- After:  `import "…/internal/engine/skill"` → `skill.Skill` (same identifier, different package path)

`SkillBase`'s field set is unchanged; the only addition is the accessor methods table above.

## Call-Site Rewrites

| File | Current import | New import |
|---|---|---|
| `internal/engine/entity/actors/character.go` | `…/engine/physics/skill` | `…/engine/skill` |
| `internal/engine/entity/actors/platformer/platformer.go` | `…/engine/physics/skill` | `…/engine/skill` |
| `internal/engine/entity/actors/builder/builder.go` | `…/engine/physics/skill` | `…/engine/skill` (interface only — no kit import) |
| `internal/engine/entity/actors/builder/builder_test.go` | `…/engine/physics/skill` | `…/engine/skill` (+ test helpers may stub `skill.Skill` directly without kit) |
| `internal/game/scenes/phases/player.go` | `engineskill "…/engine/physics/skill"` | `kitskills "…/internal/kit/skills"` (for `FromConfig`/`SkillDeps`); `"…/internal/engine/skill"` if the `skill.Skill` interface is referenced explicitly |
| `internal/game/entity/actors/player/state_contributors.go` | `engineskill "…/engine/physics/skill"` | `"…/internal/engine/skill"` (for the `Skill` interface in receiver types) AND `kitskills "…/internal/kit/skills"` (for `*kitskills.DashSkill`, `*kitskills.ShootingSkill` type switches) |

### Builder — dependency inversion (load-bearing change)

`builder.go` currently declares parameters typed `skill.SkillDeps` and calls `skill.FromConfig`. After the split:

- `SkillDeps` and `FromConfig` live in **kit** (`kitskills`).
- The builder is in the **engine** layer and **MUST NOT** import kit.

Resolution: invert the dependency. The builder accepts a pre-built `[]skill.Skill` (engine type) instead of `SkillDeps`. The factory call moves up to the game-layer caller.

Concretely:

- `builder.ApplySkills(deps skill.SkillDeps, …)` → `builder.ApplySkills(skills []skill.Skill)` (where `skill` is `internal/engine/skill`).
- `internal/game/scenes/phases/player.go` constructs `kitskills.SkillDeps`, calls `kitskills.FromConfig(cfg, deps)`, and passes the resulting `[]skill.Skill` to the builder.
- `builder_test.go` adapts to the new signature; assertions about "factory called with X" become "builder attaches the skills it is given".

This is the only behavioural-shape change in the story (test seams move). No runtime behaviour changes.

## Struct & Transition Shapes (Unchanged)

The `SkillState` transition table is preserved verbatim:

| From | Trigger | To |
|---|---|---|
| `StateReady` | `HandleInput` activation | `StateActive` |
| `StateActive` | `timer >= duration` (in `Update`) | `StateCooldown` |
| `StateCooldown` | `timer >= cooldown` (in `Update`) | `StateReady` |

`Set.Update` iterates all skills and calls `Update` on each; `Set.Get(key)` returns the registered `ActiveSkill` whose `ActivationKey()` matches; `Set.ActiveCount()` returns the count of skills currently in `StateActive`.

## Pre-Conditions

- `internal/kit/combat/weapon/` exists (stories 040–049, complete).
- Layering test `TestEngineLayerHasNoKitOrGameDependencies` exists in `internal/engine/layering_test.go` and is green.
- Combat-absent guard `TestEngineCombatDirectoryAbsent` exists in `internal/engine/combat_absent_test.go` (template for the skill-absent test).
- Red Phase tests authored under the previous spec (RP-1, RP-2, RP-3, RP-6) currently target `internal/engine/contracts/skill/`; they MUST be retargeted to `internal/engine/skill/` before Feature Implementer begins (see "Red Phase" below).

## Post-Conditions

- `internal/engine/skill/` exists, package `skill`, exporting the surface above.
- `internal/kit/skills/` exists, package `kitskills`, with all concretes + factory + tests.
- `internal/engine/physics/skill/` does not exist on disk.
- `internal/engine/contracts/skill/` does NOT exist (the previous Spec's destination is explicitly rejected).
- `go build ./...` succeeds.
- `go test ./internal/...` is green.
- Coverage on `internal/engine/skill/` ≥ 80 %.
- Coverage on `internal/kit/skills/` ≥ 80 %.
- No engine package imports `internal/kit/...` (existing layering test passes).
- No call site imports `internal/engine/physics/skill`.
- `golangci-lint run ./...` reports no new warnings.

## Red Phase (Failing Tests to Author / Re-target)

The TDD Specialist already authored RP-1, RP-2, RP-3, RP-6 against `internal/engine/contracts/skill/`. They must be **retargeted** to the new layout. The list below is authoritative.

### RP-1: Engine `physics/skill` directory absent

- **File:** `internal/engine/skill_absent_test.go` (already exists — keep)
- **Test:** `TestEnginePhysicsSkillDirectoryAbsent`
- **Body:** `os.Stat("physics/skill")` returns an `os.IsNotExist` error.

### RP-2: Engine `skill` package surface present

- **File:** `internal/engine/skill/skill_surface_test.go` (NEW location — replaces the previous `internal/engine/contracts/skill/skill_contract_test.go`).
- **Test:** `TestEngineSkillPackageSurface`
- **Body:** Compile-time interface and enum assertions:
  ```go
  var _ skill.Skill = (*skill.SkillBase)(nil)
  var _ skill.SkillState = skill.StateReady
  var _ skill.SkillState = skill.StateActive
  var _ skill.SkillState = skill.StateCooldown
  ```
  Plus runtime checks:
  - Fresh `&SkillBase{}` has `IsActive() == false`.
  - After `b.SetState(skill.StateActive)`, `b.IsActive() == true`.
  - `b.IncTimer()` increments `b.Timer()` by 1.
  - All accessor pairs are coherent (`SetX(v); X() == v`) for `Duration`, `Cooldown`, `Speed`, `Timer`.

### RP-2b: Engine `skill.Set` registry surface present

- **File:** `internal/engine/skill/set_test.go`
- **Test:** `TestEngineSkillSetSurface`
- **Body:**
  - `s := skill.NewSet()` returns a non-nil `*skill.Set`.
  - `s.ActiveCount()` returns 0 on an empty set.
  - `s.Add(stub)` then `s.ActiveCount()` reflects the count of stubs in `StateActive`.
  - `s.Update(actor, model)` invokes `Update` on each registered stub (use a local stub `Skill` that records calls).
  - `s.Get(key)` round-trips an `ActiveSkill` previously added with that activation key.
  - `s.All()` returns all registered skills in insertion order.

### RP-3: Kit `skills` package present and exposes concretes

- **File:** `internal/kit/skills/package_surface_test.go` (already exists — update import path to `internal/engine/skill`).
- **Test:** `TestKitSkillsPackageSurface`
- **Body:** Compile-time assertions:
  ```go
  import (
      "…/internal/engine/skill"
      kitskills "…/internal/kit/skills"
  )
  var _ skill.Skill        = (*kitskills.HorizontalMovementSkill)(nil)
  var _ skill.ActiveSkill  = (*kitskills.JumpSkill)(nil)
  var _ skill.ActiveSkill  = (*kitskills.DashSkill)(nil)
  var _ skill.ActiveSkill  = (*kitskills.ShootingSkill)(nil)
  var _ kitskills.OffsetToggler // type exists
  ```
  Plus a smoke test that `kitskills.FromConfig(nil, kitskills.SkillDeps{})` returns an empty (or nil-but-zero-len) `[]skill.Skill`.

### RP-4: Engine layer does not import kit

Coverage by existing `TestEngineLayerHasNoKitOrGameDependencies`. No new test required, but it MUST remain green after the move.

### RP-5: Migrated behavioural tests

The existing tests (`shooting_test.go`, `factory_test.go`, `lifecycle_test.go`, `shooting_eight_directions_test.go`, `platform_jump_test.go`) move verbatim into `internal/kit/skills/` and continue to assert the same behaviour, with updated package + import paths only. Any test that exercises pure `SkillBase` lifecycle (no concrete skill) moves into `internal/engine/skill/skill_test.go` instead.

### RP-6: Builder no longer imports kit

- **File:** `internal/engine/entity/actors/builder/builder_layering_test.go` (already exists — keep).
- **Test:** `TestBuilderDoesNotImportKit`
- **Body:** Read each `.go` file in the builder package and assert none contain the substring `"internal/kit/"`. Currently passes; guards against regression during the dependency inversion.

## Integration Points

- **Schema layer (`internal/engine/data/schemas`)**: unchanged. `SkillsConfig` stays where it is.
- **Event layer (`internal/engine/event`)**: unchanged. `SkillDeps.EventManager` keeps its anonymous-interface shape.
- **Combat layer (`internal/engine/contracts/combat`, `internal/kit/combat`)**: unchanged.
- **Phases (`internal/game/scenes/phases`)**: takes responsibility for calling `kitskills.FromConfig` and passing `[]skill.Skill` to the builder.
- **Player state contributors (`internal/game/entity/actors/player/state_contributors.go`)**: continues to type-switch on `*kitskills.DashSkill` and `*kitskills.ShootingSkill`; the receiver interface listing `Skills() []skill.Skill` uses the engine-skill type.

## Out of Scope

- Adding new skill types or extending `SkillsConfig` (story 051).
- Changing field names on `SkillBase` (only adding accessor methods is allowed).
- Removing `SkillBase` embedding in favour of composition.
- Re-evaluating the `EventManager` anonymous interface in `SkillDeps`.
- Renaming the registry/Set type beyond what the existing source already calls it.

## Risks & Decisions

1. **Decision — engine destination is `internal/engine/skill/`, not `internal/engine/contracts/skill/`.** `SkillBase` and the registry carry state and logic; `internal/engine/contracts/` is for thin interface declarations only. Sibling precedents: `internal/engine/entity/`, `internal/engine/scene/`.
2. **Decision — package name `kitskills` over `skills`.** Aligns with the existing `kitstates` precedent and avoids stutter at call sites (`kitskills.DashSkill` reads cleanly).
3. **Decision — invert builder dependency rather than parameterising with a factory function.** Lifting the `FromConfig` call up to the game layer is the cleanest way to keep the engine free of kit imports while preserving the ergonomic `FromConfig` API.
4. **Decision — registry type name `Set`.** Pending verification of the current name in `engine/physics/skill/`. If a different exported name already exists, preserve it; do not rename in this story.
5. **Risk — unexported field access.** Concrete skills currently mutate `SkillBase.state`, `timer`, etc. directly. Mitigation: minimal accessor methods on `SkillBase` (in `internal/engine/skill/`). Mechanical change, no semantic surface added beyond what is strictly necessary.
6. **Risk — test regressions.** All existing test files move verbatim with import-path edits. Coverage must not drop below 80 % on either new package; Gatekeeper verifies.
7. **Risk — Red Phase tests authored under previous spec.** RP-2 and RP-3 currently reference `internal/engine/contracts/skill/`. TDD Specialist must retarget them to `internal/engine/skill/` before Feature Implementer begins.

## Pipeline Sequence (Next Agents)

This story introduces a new engine package (`internal/engine/skill/`) but does **not** introduce any new mockable interface boundaries beyond what the existing `body.MovableCollidable`, `body.BodiesSpace`, `combat.Inventory`, `combat.ProjectileManager` mocks already cover. Skip Mock Generator.

Recommended order:

1. **TDD Specialist (re-entry)** — retarget RP-2 to `internal/engine/skill/skill_surface_test.go`, add RP-2b `set_test.go`, update RP-3 import paths to `internal/engine/skill`. RP-1 and RP-6 stay as-is.
2. **Feature Implementer** — perform the relocation, accessor additions, registry move, builder dependency inversion, and call-site rewrites until all tests pass.
3. **Workflow Gatekeeper** — verify layering tests, coverage delta on both new packages, and `golangci-lint`.
