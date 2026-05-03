# User Story — 050-move-weapon-to-kit

## Story

**As a** kit-layer contributor,
**I want** the skill system split so that genre-agnostic infrastructure lives in `internal/engine/skill/` and concrete beat-em-up / platformer skills live in `internal/kit/skills/`,
**so that** the engine layer is self-contained and reusable across game genres, while the kit layer owns which skills actually exist.

---

## Background

All skill code currently lives under `internal/engine/physics/skill/`. The prior stories (040–049) completed migration of combat types (`ProjectileWeapon`, `MeleeWeapon`, etc.) into `internal/kit/combat/weapon/`. The skill layer is the remaining piece.

On closer analysis the skill package contains two distinct concerns:

| Concern | Files | Correct home |
|---|---|---|
| Genre-agnostic skill infrastructure | `skill.go` — `Skill`/`ActiveSkill` interfaces, `SkillBase` struct, `SkillState` enum, registry/set type (`Add`, `Get`, `Update`, `ActiveCount`) | `internal/engine/skill/` |
| Genre-specific concrete skills | `skill_shooting.go`, `skill_dash.go`, `skill_platform_jump.go`, `skill_platform_move.go`, `factory.go` | `internal/kit/skills/` |

`SkillBase` and the registry are as genre-agnostic as `Body` or `Actor` — they encode timing, cooldown, and state-machine mechanics that any game genre would need. Placing them in kit would force every future genre to re-implement or re-import the same primitives. Keeping them in the engine (outside the `physics` sub-package) respects the same reasoning used for `internal/engine/entity/` and `internal/engine/scene/`.

The dependency flow is: `game` → `kit` → `engine`. Kit concrete skills embed or reference `engine/skill.SkillBase`. Game code constructs `kit/skills.SkillDeps` and wires the resulting `[]engine/skill.Skill` slice into the engine-layer builder.

---

## Scope

This story covers:

1. Moving the genre-agnostic infrastructure (`Skill`, `ActiveSkill`, `SkillBase`, `SkillState`, registry/set type) out of `internal/engine/physics/skill/` into a new top-level engine package `internal/engine/skill/`.
2. Moving all concrete skill implementations and `FromConfig` / `SkillDeps` into `internal/kit/skills/`.
3. Updating all call sites in `internal/game/`, `internal/engine/entity/actors/`, and `internal/engine/entity/actors/builder/` to import from the new locations.
4. Deleting `internal/engine/physics/skill/` once it is empty.

> Story 051 (`kit-skills-system`) covers new skill types, config schema extensions, and further API evolution. This story is a pure relocation with no behaviour change.

---

## Acceptance Criteria

### AC-1: Genre-agnostic skill infrastructure in engine layer

- `internal/engine/skill/` exists and declares:
  - `Skill` interface
  - `ActiveSkill` interface
  - `SkillBase` struct with accessor methods (`State`/`SetState`, `Duration`/`SetDuration`, `Cooldown`/`SetCooldown`, `Speed`/`SetSpeed`, `Timer`/`SetTimer`/`IncTimer`) and tick/duration logic
  - `SkillState` enum + constants
  - Registry / set type with `Add`, `Get`, `Update`-all, and `ActiveCount` operations
- No concrete, genre-specific skill lives in `internal/engine/skill/`.
- `internal/engine/entity/actors/character.go`, `platformer.go`, and `builder.go` import `internal/engine/skill` (not `engine/physics/skill`).

### AC-2: Concrete skills relocated to kit layer

- `internal/kit/skills/` contains all concrete skill types previously in `engine/physics/skill/`:
  - `ShootingSkill`
  - `DashSkill`
  - `JumpSkill`
  - `HorizontalMovementSkill`
  - `OffsetToggler`
  - `FromConfig` factory with `SkillDeps`
- Package declaration is `package skills` (consistent with `internal/kit/states/` convention).
- `internal/kit/skills/` may import `internal/engine/skill` and other engine packages freely; it must not import `internal/game/`.

### AC-3: `engine/physics/skill/` removed

- The directory `internal/engine/physics/skill/` no longer exists.
- A guard test (e.g. `TestEnginePhysicsSkillDirectoryAbsent`) asserts the path is gone and passes in CI.

### AC-4: Import direction respected

- No package under `internal/engine/` imports `internal/kit/` (existing layering test `TestEngineLayerHasNoKitOrGameDependencies` continues to pass).
- `internal/kit/skills/` imports `internal/engine/skill` — this is the correct downward dependency.
- `internal/game/` imports `internal/kit/skills/` for `FromConfig` / `SkillDeps`.
- `internal/game/` imports `internal/engine/skill` for the `Skill` interface where needed.

### AC-5: No behaviour change

- All previously existing unit tests for `ShootingSkill`, `DashSkill`, `JumpSkill`, `HorizontalMovementSkill`, `SkillBase`, and `FromConfig` are preserved and pass at their new locations.
- No test is deleted; only package paths in import statements change.

### AC-6: All tests pass

- `go test ./internal/...` passes with zero failures.
- Coverage on `internal/engine/skill/` and `internal/kit/skills/` is >= 80% (matching the project-wide goal).
- `golangci-lint run ./...` reports no new warnings.

---

## Design Notes

### Why not `internal/engine/contracts/skill/`?

`internal/engine/contracts/` is reserved for thin interface contracts (no structs, no logic). `SkillBase` carries timer and cooldown logic — it is infrastructure, not a contract. The correct analogy is `internal/engine/entity/` (holds `Body`, `Actor` structs) rather than `internal/engine/contracts/` (holds `Updater`, `Renderer` interfaces). Therefore the engine-side destination is `internal/engine/skill/`, not `internal/engine/contracts/skill/`.

### Dependency inversion for builder

`internal/engine/entity/actors/builder/builder.go` must not call `kit/skills.FromConfig` directly (that would introduce an engine → kit import). Instead:

1. `game` code (e.g. `internal/game/scenes/phases/player.go`) calls `kit/skills.FromConfig(deps)` and receives `[]engine/skill.Skill`.
2. That slice is passed into the engine-layer builder via an existing or new setter/option.
3. The builder stores and applies the `[]engine/skill.Skill` slice without knowing about kit.

---

## Out of Scope

- Adding new skill types or changing `ShootingSkill` behaviour (story 051).
- Migrating documentation artefacts (`README.md`) — update or move as needed during implementation, but no new docs are required by this story.
- Changing `SkillDeps` struct fields beyond what is necessary for the relocation.
