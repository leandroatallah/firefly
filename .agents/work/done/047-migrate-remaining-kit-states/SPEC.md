# SPEC — 047-migrate-remaining-kit-states

## Branch Name

`047-migrate-remaining-kit-states`

## Summary

Relocate every remaining genre-reusable Actor State from `internal/game/entity/actors/states/` (package `gamestates`) into `internal/kit/states/` (package `kitstates`). After this story:

- `internal/kit/states/` is the single source of truth for all platformer states (`GroundedState`, `DashState`, `MeleeAttackState` adapter, sub-states, `GroundedInput`/`GroundedDeps`, and `SubState*`/`StateGrounded`/`StateDashing`/`StateMeleeAttack` enums).
- `internal/game/entity/actors/states/` contains only project-specific shims (`actor_state.go`, `actor_state_concrete.go`) and any non-state utilities not in scope (`offset_toggler.go` stays per USER_STORY).
- All importers in `internal/game/` reference the moved symbols via `kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"`.
- The three-layer dependency rule still holds (engine clean of kit/game; kit clean of game).

This is a **refactor with no behaviour change**. All existing tests must pass unchanged in semantics; only their package/import lines move with the code.

## Scope (Files to Move)

From `internal/game/entity/actors/states/` -> `internal/kit/states/`:

| Source file | Destination | Notes |
|---|---|---|
| `walking_sub_state.go` | `internal/kit/states/walking_sub_state.go` | Package `kitstates`. Type stays unexported. |
| `ducking_sub_state.go` | `internal/kit/states/ducking_sub_state.go` | Package `kitstates`. Type stays unexported. |
| `aim_lock_sub_state.go` | `internal/kit/states/aim_lock_sub_state.go` | Package `kitstates`. Type stays unexported. |
| `grounded_sub_state.go` | `internal/kit/states/grounded_sub_state.go` | Package `kitstates`. `GroundedSubStateEnum`, `groundedSubState` (interface stays unexported), `SubStateIdle/Walking/Ducking/AimLock` exported. |
| `grounded_input.go` | `internal/kit/states/grounded_input.go` | Package `kitstates`. `GroundedInput`, `GroundedDeps` exported. |
| `grounded_state.go` | `internal/kit/states/grounded_state.go` | Package `kitstates`. `GroundedState`, `NewGroundedState`, `StateGrounded` (var + `init()` registration). |
| `grounded_state_test.go` | `internal/kit/states/grounded_state_test.go` | Package `kitstates_test`. Imports `kitstates`. |
| `dash_state.go` | `internal/kit/states/dash_state.go` | Package `kitstates`. `DashState`, `NewDashState`, `DashConfig`, `StateDashing` (var + `init()` registration). |
| `dash_state_test.go` | `internal/kit/states/dash_state_test.go` | Package `kitstates_test`. Imports `kitstates`. |
| `melee_state.go` | `internal/kit/states/melee_state.go` | Package `kitstates`. `MeleeAttackState` (alias), `NewMeleeAttackState`, `InstallMeleeAttackState`, `TryMeleeFromFalling`, `ResetComboOnInterrupt`, `MeleeAttackStepStates`, `StateMeleeAttack` (var + `init()` registration). |
| `melee_state_test.go` | `internal/kit/states/melee_state_test.go` | Package `kitstates_test`. Imports `kitstates`. |
| `melee_step_states_test.go` | `internal/kit/states/melee_step_states_test.go` | Package `kitstates_test`. |
| `mocks_test.go` | `internal/kit/states/mocks_test.go` | Package `kitstates_test`. Provides `MockInputSource` and any other shared in-package test doubles already present. |

Files that **stay** in `internal/game/entity/actors/states/`:

- `actor_state.go` (currently empty `package gamestates`)
- `actor_state_concrete.go` (re-exports `Dying`, `Dead`, `Exiting` from `internal/engine/entity/actors`)
- `offset_toggler.go` and `offset_toggler_test.go` (utility, explicitly out of scope per USER_STORY)

## Package Renames

- All files moved into `internal/kit/states/` MUST declare `package kitstates` (matching the existing `idle_sub_state.go`).
- All test files moved MUST declare `package kitstates_test` (current `gamestates_test`).
- Import path: `kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"`.

## Internal Symbol Adjustments

Inside `kitstates`:

1. `grounded_state.go` — replace the qualified reference `kitstates.IdleSubState[...]` with the unqualified local type `IdleSubState[...]` (now in the same package). The instantiation:
   ```go
   return &IdleSubState[GroundedSubStateEnum, GroundedInput]{
       Idle:    SubStateIdle,
       Walking: SubStateWalking,
       Ducking: SubStateDucking,
       AimLock: SubStateAimLock,
   }
   ```
2. Remove the now-unused `kitstates` import alias from `grounded_state.go`.
3. `melee_state.go` references `StateGrounded` and `StateMeleeAttack` — both now siblings in the same `kitstates` package; no alias needed.

State-enum globals stay annotated `//nolint:gochecknoglobals` exactly as in the source; the `init()` registration via `actors.RegisterState("grounded" | "dash" | "melee_attack", ...)` is preserved verbatim. State **name strings** (`"grounded"`, `"dash"`, `"melee_attack"`, `"melee_attack_step_<i>"`) move with their `init()` calls and remain the single authoritative source for those names (per AC: no-duplicate-strings).

## Importer Updates (Game Layer)

The following call-sites must be updated to import from `kitstates` instead of `gamestates` for the moved symbols. The `gamestates` import remains where the file still references game-only symbols (`Dying`, `Dead`, `Exiting`).

| File | Change |
|---|---|
| `internal/game/entity/actors/player/state_contributors.go` | `gamestates.StateDashing` -> `kitstates.StateDashing`. Update import accordingly. |
| `internal/game/entity/actors/player/state_contributors_test.go` | Mirror import update. |
| `internal/game/entity/actors/player/climber.go` | `gamestates.MeleeAttackStepStates`, `gamestates.StateMeleeAttack`, `gamestates.StateGrounded` -> `kitstates.*`. `gamestates.Dying`/`Dead` stay. |
| `internal/game/entity/actors/player/climber_test.go` | Mirror updates. |
| `internal/game/entity/actors/player/player_test.go` | Update any `gamestates.<moved-symbol>` references to `kitstates`. |
| `internal/game/entity/actors/enemies/wolf.go` | Only references `gamestates.Dying` — no change required. |
| `internal/game/entity/actors/enemies/bat.go` | Only references `gamestates.Dying` — no change required. |
| `internal/game/entity/actors/enemies/bat_test.go` | Mirror — verify, no change unless it uses a moved symbol. |
| `internal/game/scenes/phases/scene.go` | Only references `gamestates.Dying`/`Dead` — no change required. |

(The implementer must run `grep -rn "gamestates\." internal/` after the moves and rewrite any remaining hits that point to a moved symbol.)

## Interface Contracts

No new contracts in `internal/engine/contracts/` are introduced or modified. All interfaces involved (`contractsbody.MovableCollidable`, `contractsbody.BodiesSpace`, `combat.Weapon`, `animation.FacingDirectionEnum`, etc.) are already in engine and unchanged.

The unexported `groundedSubState` interface (in `grounded_sub_state.go`) stays unexported in its new home — it is an internal kit contract.

`GroundedInputLike` (already in `idle_sub_state.go`) is structurally satisfied by the new sibling `GroundedInput` interface; the generic instantiation works because of structural typing on Go interfaces.

## State Machine / Transition Behaviour

No behavioural changes. The state graph remains:

```
Grounded (composite)
  |- Idle <-> Walking <-> Ducking <-> AimLock      (sub-state transitions, see walking/ducking/aim_lock TransitionTo)
  |
  +-- JumpPressed --> actors.Falling
  +-- DashPressed --> StateDashing
  +-- MeleePressed --> StateMeleeAttack

Dashing
  +-- tween done & grounded --> actors.Idle
  +-- tween done & airborne --> actors.Falling
  +-- wall blocked --> nextIdleOrFalling()

MeleeAttack (alias to engine combat/melee.State)
  +-- behaviour delegated to engine; no changes
```

## Pre-conditions

- Story 046 is merged (`IdleSubState` already in `kitstates`); confirmed by `internal/kit/states/idle_sub_state.go` existing.
- `go test ./...` is currently green on `main`.
- Working tree is on the `047-migrate-remaining-kit-states` branch (created by the developer before TDD/implementation begins).

## Post-conditions

1. `internal/kit/states/` contains all 13 files listed in Scope, plus the pre-existing `doc.go`, `idle_sub_state.go`, `idle_sub_state_test.go`.
2. `internal/game/entity/actors/states/` contains exactly: `actor_state.go`, `actor_state_concrete.go`, `offset_toggler.go`, `offset_toggler_test.go`.
3. `go build ./...` succeeds.
4. `go test ./...` is green with the same number of (or more) passing test cases.
5. CI dependency commands pass:
   ```bash
   ! go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)'
   ! go list -deps ./internal/kit/...    | grep    'internal/game'
   ```
6. Coverage of `internal/kit/states/` is >= 80%; overall coverage delta on `internal/engine/...` and `internal/game/...` is non-negative (reported by Coverage Analyzer).
7. `golangci-lint run ./...` passes (Gatekeeper final gate).
8. No raw string literal duplicating any of `"grounded"`, `"dash"`, `"melee_attack"`, `"melee_attack_step_<i>"` exists outside the original `actors.RegisterState(...)` call sites.

## Integration Points

- Bounded Context: **Entity / Actor State Machine** (`internal/engine/entity/actors`).
- All moved code interacts with the actor state machine through the existing `actors.ActorStateEnum`, `actors.RegisterState`, and `actors.GetStateEnum` API. These are unaffected.
- Combat integration (`internal/engine/combat/melee.State`, `melee.InstallState`, `melee.New`) is consumed but unchanged.
- Skills system (`internal/engine/physics/skill.DashSkill`, `ShootingSkill`) — `state_contributors.go` updates its import path only.

## Red Phase Scenarios (TDD Specialist Brief)

This story is a **structural refactor**. The "Red Phase" therefore consists of:

### R1. Compile-failure red

Before the move, copy each `*_test.go` file to its new location under `internal/kit/states/` with package `kitstates_test` and import `kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"`. The tests will fail to compile because the production symbols (`GroundedState`, `DashState`, `MeleeAttackState`, `StateGrounded`, `StateDashing`, `StateMeleeAttack`, `GroundedInput`, `GroundedDeps`, `GroundedSubStateEnum`, `SubState*`, `NewGroundedState`, `NewDashState`, `NewMeleeAttackState`, `MeleeAttackStepStates`, `TryMeleeFromFalling`, `ResetComboOnInterrupt`, `DashConfig`) do not yet exist in `kitstates`. This is the failing test signal.

### R2. Dependency-rule red

Add (or reuse from story 046) a CI-mirroring shell test that runs:

```bash
go list -deps ./internal/kit/... | grep 'internal/game'
```

After the imports inside the moved files are flipped (production code is moved but its `import` block still names `gamestates`), this guard fails. Green is reached when only `kitstates` references the moved symbols and `gamestates` no longer holds them.

### R3. Behavioural-equivalence red (no new tests required)

The migrated tests `TestGroundedSubStateTransitions`, `TestGroundedSubStateTransitions_JumpExitsGrounded`, `TestGroundedStateOnFinishCallsSubOnFinish`, `TestGroundedStateReEntryResetsToIdle`, `TestDashStateUpdate` (+ subtest "second OnStart is no-op while dashing"), and the full `melee_state_test.go` / `melee_step_states_test.go` suites must pass unchanged in semantics from their new location. Any divergence is a regression.

The TDD Specialist does **not** need to author new behavioural tests — the existing tests, after relocation, are the executable specification. Their failure to compile (R1) is the Red Phase; their passing after the production move (Green) closes the cycle.

## Implementation Order (Feature Implementer Brief)

1. Move all 8 production `.go` files from `gamestates` to `kitstates`; rewrite the package clause and update intra-file references (drop `kitstates.` prefix on `IdleSubState`).
2. Move all 4 test files (`grounded_state_test.go`, `dash_state_test.go`, `melee_state_test.go`, `melee_step_states_test.go`) and `mocks_test.go`; rewrite their package clause to `kitstates_test` and update import paths.
3. Delete the now-empty source files in `internal/game/entity/actors/states/`.
4. Update `state_contributors.go`, `climber.go`, and any other game-layer file whose `gamestates.<symbol>` refers to a moved symbol — switch the import to `kitstates` (keep the original `gamestates` import only where `Dying`/`Dead`/`Exiting` are still used).
5. Run `go build ./...`, `go test ./...`, and the two `go list -deps` guards. Iterate until all are green.

## Out of Scope

Per the user story: archetypes (048), combat move (049), weapon move (050), skills move (051), UI split (052), ADR rewrite, `kit/mocks/`, new behaviour, sub-module split, and `offset_toggler` migration are all deferred.

## Risks & Mitigations

| Risk | Mitigation |
|---|---|
| `init()` ordering shifts because state-enum registration moves package | Names registered via `actors.RegisterState` are deduplicated by the engine; `init()` runs once per package import. The set of package imports is preserved (any package that previously imported `gamestates` for these states now imports `kitstates`), so ordering equivalence holds. |
| Generic instantiation `IdleSubState[GroundedSubStateEnum, GroundedInput]` fails because `GroundedInput` no longer satisfies `GroundedInputLike` after rename | `GroundedInputLike` matches structurally on three method names; `GroundedInput` retains those methods unchanged. No risk in practice; verified during compile. |
| Hidden game-package test relies on the symbol still living in `gamestates` | The grep step in the Implementation Order phase catches every remaining `gamestates.<moved-symbol>` reference. |
