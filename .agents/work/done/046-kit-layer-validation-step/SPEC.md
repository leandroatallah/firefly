# SPEC — 046-kit-layer-validation-step

**Branch:** `046-kit-layer-validation-step`
**Bounded Context (new):** `Kit` → `internal/kit/`
**Type:** Pure structural relocation. No behaviour change.

## 1. Goal

Stand up `internal/kit/` as the third architectural layer (`engine` ← `kit` ← `game`) and prove the dependency rule end-to-end by relocating exactly one concrete Actor sub-state — `IdleSubState` — from `internal/game/entity/actors/states/` to `internal/kit/states/`. Behaviour and state-machine semantics are identical before and after.

## 2. Source-Code Audit (current state)

`internal/game/entity/actors/states/` (package `gamestates`) currently holds:

| Symbol | Visibility | Used by | Notes |
|---|---|---|---|
| `idleSubState` (struct) | unexported | `newSubState()` factory in `grounded_state.go` | Move target |
| `OnStart(int)` / `OnFinish()` | exported on unexported type | sub-state contract | Method signatures unchanged |
| `transitionTo(GroundedInput) GroundedSubStateEnum` | unexported method | `groundedSubState` interface | Must become exported for kit→game structural typing |
| `groundedSubState` (interface) | unexported | `GroundedState.activeSub`, `newSubState` | Stays in game; widened to require exported method |
| `GroundedInput` (interface) | exported | already public | No change |
| `GroundedSubStateEnum` + `SubStateIdle/Walking/Ducking/AimLock` | exported | already public | No change |
| `GroundedDeps`, `GroundedState`, `NewGroundedState` | exported | already public | No change |

Key observation: `idleSubState` has zero imports outside its own package and only references the symbols listed above. Once the contract method is exported, the type can be relocated without any further surface widening.

There is **no** dedicated `idle_sub_state_test.go`. Idle behaviour is exercised today through `grounded_state_test.go::TestGroundedSubStateTransitions` (cases `"no input stays Idle"` and `"duck released with clearance transitions to Idle"`) and the package-local `MockInputSource` in `mocks_test.go`. Per the User Story ("tests travel with the type"), this story introduces a new dedicated unit test alongside the relocated type and **does not** delete or weaken the existing `gamestates` integration tests — they must continue to pass against the kit-provided implementation.

## 3. Contracts Touched

### 3.1 Existing engine contracts
None. `internal/engine/contracts/` is **not** modified by this story. The `groundedSubState` contract remains a *game-internal* interface (it is parameterised on `GroundedInput` and `GroundedSubStateEnum`, both of which live in the game `gamestates` package and will move to kit only in story 047). This story deliberately does not lift the sub-state contract into engine.

### 3.2 Game-package contract widening (in-place edit, not a move)
File: `internal/game/entity/actors/states/grounded_sub_state.go`

```go
// groundedSubState is the internal contract every sub-state must satisfy.
type groundedSubState interface {
    OnStart(currentCount int)
    OnFinish()
    TransitionTo(input GroundedInput) GroundedSubStateEnum   // was: transitionTo
}
```

Rationale: Go does not allow a type defined in package `kit/states` to satisfy an interface that names an unexported method belonging to package `gamestates`. The interface method must therefore be exported. All four existing sub-state implementations (`idleSubState`, `walkingSubState`, `duckingSubState`, `aimLockSubState`) get their method renamed `transitionTo` → `TransitionTo` in this story. Three of them (Walking, Ducking, AimLock) **stay in game** and are migrated in story 047; only Idle is *moved*.

### 3.3 New kit type
Package: `kitstates` (path `internal/kit/states/`)
File: `internal/kit/states/idle_sub_state.go`

```go
package kitstates

import gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
```

Wait — this would invert the dependency rule. See §4 for the resolution: kit imports `GroundedInput` and `GroundedSubStateEnum` from where they will live in story 047. For story 046, since those symbols still live in `gamestates`, we use **Go generics + a tiny local alias** so the kit package depends on no concrete game type. Concretely:

The kit package declares its own enum-and-input shape via type parameters supplied by the caller. To keep this story minimal and avoid a generics rewrite, we instead use a **structural-typing trampoline** in game: kit declares `IdleSubState` parameterised on two interfaces it defines locally, and the game package uses a thin adapter. See §4 for the exact form chosen.

## 4. Resolution: Avoiding the Engine→Game Inversion

Because `GroundedInput` and `GroundedSubStateEnum` still live in `gamestates` for story 046 (they migrate in 047), the kit package cannot import them without violating the dependency rule. Two acceptable shapes:

**Chosen — Option K1: Local input/enum interfaces in kit, parameterised by generics.**

`internal/kit/states/idle_sub_state.go`:

```go
package kitstates

// GroundedInputLike is the minimum input surface IdleSubState reads.
// Game-side GroundedInput satisfies this interface structurally.
type GroundedInputLike interface {
    AimLockHeld() bool
    DuckHeld() bool
    HorizontalInput() int
}

// IdleSubState is a no-op grounded sub-state that selects the next
// sub-state based on input. The concrete enum E is supplied by the
// caller so kit stays free of game-specific identifiers.
type IdleSubState[E comparable, I GroundedInputLike] struct {
    Idle    E
    Walking E
    Ducking E
    AimLock E
}

func (s *IdleSubState[E, I]) OnStart(_ int) {}
func (s *IdleSubState[E, I]) OnFinish()     {}

func (s *IdleSubState[E, I]) TransitionTo(input I) E {
    switch {
    case input.AimLockHeld():
        return s.AimLock
    case input.DuckHeld():
        return s.Ducking
    case input.HorizontalInput() != 0:
        return s.Walking
    default:
        return s.Idle
    }
}
```

This:
- Imports nothing from `internal/game/...` (dependency rule holds).
- Imports nothing from `internal/engine/...` either — it only uses Go builtins. That is acceptable (kit is *allowed* to import engine, not *required* to).
- Preserves the exact transition truth table from §2 of the original `idle_sub_state.go`.

Game-side wiring (`grounded_state.go::newSubState`):

```go
case SubStateIdle:
    return &kitstates.IdleSubState[GroundedSubStateEnum, GroundedInput]{
        Idle:    SubStateIdle,
        Walking: SubStateWalking,
        Ducking: SubStateDucking,
        AimLock: SubStateAimLock,
    }
```

**Rejected — Option K2:** Move `GroundedInput` and `GroundedSubStateEnum` to kit alongside `IdleSubState`. Rejected because it expands story scope: the other three sub-states (which stay in game in story 046) would then have to import kit just to reference the enum, doubling the import-graph churn before story 047 is ready.

## 5. Package Layout

```
internal/kit/
└── states/
    ├── doc.go                  # package documentation + dependency-rule contract
    ├── idle_sub_state.go       # IdleSubState[E, I]
    └── idle_sub_state_test.go  # unit tests (table-driven)
```

### 5.1 `internal/kit/states/doc.go`

```go
// Package kitstates contains genre-reusable concrete Actor sub-state
// implementations.
//
// Dependency rule (enforced by CI):
//   - kitstates MAY import internal/engine/...
//   - kitstates MUST NOT import internal/game/...
//   - kitstates MUST NOT import any other internal/kit/... package that
//     transitively imports internal/game/...
//
// Types here are parameterised on the caller's enum and input contract
// to avoid coupling to a specific game's state-machine vocabulary.
package kitstates
```

Package import name: `kitstates` (matches Go convention `<dir>states` collision-free; callers may rename via `import kitstates "..."`).

## 6. Migration Steps (deterministic, in order)

The Feature Implementer will execute these steps. The TDD Specialist writes the new tests in step 5 first (Red) before step 4 produces the Green code.

1. **Create kit package skeleton**
   - `mkdir -p internal/kit/states`
   - Add `internal/kit/states/doc.go` (content per §5.1).

2. **Export the sub-state contract in game**
   - Edit `internal/game/entity/actors/states/grounded_sub_state.go`: rename interface method `transitionTo` → `TransitionTo`.
   - Edit each of `idle_sub_state.go`, `walking_sub_state.go`, `ducking_sub_state.go`, `aim_lock_sub_state.go`: rename method receiver `transitionTo` → `TransitionTo`.
   - Edit `grounded_state.go::Update`: change call site `g.activeSub.transitionTo(input)` → `g.activeSub.TransitionTo(input)`.

3. **Add the kit type**
   - Create `internal/kit/states/idle_sub_state.go` with `IdleSubState[E, I]` per §4.

4. **Rewire game to consume kit**
   - Edit `internal/game/entity/actors/states/grounded_state.go::newSubState` so `case SubStateIdle` returns `&kitstates.IdleSubState[GroundedSubStateEnum, GroundedInput]{...}` (full literal per §4).
   - Add import `kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"`.
   - **Delete** `internal/game/entity/actors/states/idle_sub_state.go` (the old type is no longer referenced).

5. **Tests travel with the type**
   - Create `internal/kit/states/idle_sub_state_test.go` (see §8 Red Phase).
   - Leave `internal/game/entity/actors/states/grounded_state_test.go` untouched — it remains the integration-level guarantee that the kit type is wired correctly and that all four sub-state cases continue to work end-to-end.

6. **Run dependency-rule checks locally**
   - See §7. Both checks must exit 0.

7. **ADR-006 amendment** (see §11)

8. **Constitution update** (see §12)

## 7. Dependency-Rule CI Commands

These two commands are the gate. The Workflow Gatekeeper MUST run both; either non-zero exit fails the story. They mirror the snippet in `USER_STORY.md`:

```bash
# 7.1 — Engine must be clean: zero imports of kit or game.
go list -deps ./internal/engine/... \
  | grep -E 'github.com/boilerplate/ebiten-template/internal/(kit|game)' \
  && { echo "FAIL: engine depends on kit or game"; exit 1; } \
  || echo "OK: engine clean"

# 7.2 — Kit must be clean: zero imports of game.
go list -deps ./internal/kit/... \
  | grep    'github.com/boilerplate/ebiten-template/internal/game' \
  && { echo "FAIL: kit depends on game"; exit 1; } \
  || echo "OK: kit clean"
```

Note: the User Story renders the checks with `&& exit 1` only; the bash idiom above adds `|| echo OK` so a clean run exits 0 (grep with no matches returns 1, which would otherwise propagate). Functionally equivalent — the *pass condition* is "no matching lines."

These commands are intended to be added to the project's lint/check pipeline in a follow-up CI wiring task; for story 046 they are run manually by the Gatekeeper.

## 8. Red Phase — Failing Tests (TDD Specialist input)

New file: `internal/kit/states/idle_sub_state_test.go` (package `kitstates_test`).

The test must be table-driven, deterministic, and use a local fake input (no engine mocks needed — `GroundedInputLike` is a 3-method interface; declare a struct with func fields in the same test file). The table covers every branch of `TransitionTo`:

| # | Input setup | Expected return |
|---|---|---|
| 1 | all zero | `Idle` |
| 2 | `HorizontalInput() == 1` | `Walking` |
| 3 | `HorizontalInput() == -1` | `Walking` |
| 4 | `DuckHeld() == true` | `Ducking` |
| 5 | `AimLockHeld() == true` | `AimLock` |
| 6 | `AimLockHeld() == true` AND `DuckHeld() == true` | `AimLock` (precedence) |
| 7 | `DuckHeld() == true` AND `HorizontalInput() == 1` | `Ducking` (precedence) |

Plus two no-op lifecycle tests:

| # | Action | Expected |
|---|---|---|
| 8 | `OnStart(7)` does not panic and is observable as a no-op | no state mutation |
| 9 | `OnFinish()` does not panic | — |

Use a local enum type for `E`:

```go
type subStateEnum int
const ( idle subStateEnum = iota; walking; ducking; aimLock )
```

The test instantiates `IdleSubState[subStateEnum, *fakeInput]` with the four enum constants and calls `TransitionTo` per row.

**Acceptance:** all rows pass; `go test ./internal/kit/states/...` is green.

## 9. Pre-conditions

- Branch `046-kit-layer-validation-step` is checked out (already true).
- `internal/kit/` does **not** exist on disk before this story starts.
- `internal/game/entity/actors/states/idle_sub_state.go` exists with the content audited in §2.
- `go test ./...` is green at HEAD.
- No story-047 work has begun (the other three sub-states are still in `gamestates`).

## 10. Post-conditions

After Feature Implementer + Gatekeeper finish:

- `internal/kit/states/{doc.go,idle_sub_state.go,idle_sub_state_test.go}` exist.
- `internal/game/entity/actors/states/idle_sub_state.go` does **not** exist.
- `internal/game/entity/actors/states/grounded_state.go` imports `kitstates` and constructs `&kitstates.IdleSubState[...]{...}` in `newSubState` for `SubStateIdle`.
- `groundedSubState` interface and the four sub-state types use exported `TransitionTo`.
- `go build ./...` succeeds.
- `go test ./...` is green; no test was deleted or skipped.
- `go test -cover ./internal/kit/states/...` reports ≥ 80% (target ≥ 90% — the new file is small enough that all branches are reachable).
- The two dependency-rule commands in §7 both report OK.
- `docs/adr/ADR-006-engine-game-layer-separation.md` Status reads `Superseded by 046` with the note in §11.
- `.agents/constitution.md` Bounded Contexts table contains a `Kit` row (see §12).
- No `internal/kit/...` source file contains a state-name string literal (verified by `grep -RnE '"(idle|walking|ducking|aim[-_]?lock|grounded)"' internal/kit/` returning zero matches).

## 11. ADR-006 Amendment

Edit `docs/adr/ADR-006-engine-game-layer-separation.md`:

1. Change line 4 from `Accepted` to `Superseded by 046`.
2. Append (immediately after the Status line):

   ```markdown
   > **Note:** Superseded by story `046-kit-layer-validation-step`, which validates a three-layer architecture (`engine` ← `kit` ← `game`) by relocating `IdleSubState` to `internal/kit/states/`. A full ADR rewrite is deferred until follow-up stories 047–052 populate the `kit` layer; see `.agents/work/active/046-kit-layer-validation-step/SPEC.md` for the validation rationale.
   ```

The Decision and Consequences sections are left intact for historical record.

## 12. Constitution Bounded Contexts Update

Edit `.agents/constitution.md` → `## Bounded Contexts` table. Insert a new row (alphabetical or grouped placement is fine; recommended placement immediately after `Game Logic`):

```markdown
| Kit | `internal/kit/` |
```

No other changes to the constitution. The Non-Negotiable Standards already cover kit transitively (kit code is production code under the same rules).

## 13. State Machine / Behaviour Invariance (proof obligation)

The transition table for `IdleSubState` before and after must be byte-identical. Mapping:

| Pre-condition (input) | Pre-relocation result (`SubState…`) | Post-relocation result (`E`) |
|---|---|---|
| `AimLockHeld()` true | `SubStateAimLock` | `s.AimLock` (= `SubStateAimLock` at call site) |
| `DuckHeld()` true | `SubStateDucking` | `s.Ducking` (= `SubStateDucking`) |
| `HorizontalInput() != 0` | `SubStateWalking` | `s.Walking` (= `SubStateWalking`) |
| otherwise | `SubStateIdle` | `s.Idle` (= `SubStateIdle`) |

Branch order is preserved (AimLock → Duck → Walk → Idle). `OnStart` and `OnFinish` remain pure no-ops. `GroundedState.Update` is unchanged except for the method-name rename. Therefore `TestGroundedSubStateTransitions` continues to pass with no edits.

## 14. Test Coverage Plan

| Package | Target | How achieved |
|---|---|---|
| `internal/kit/states/` (new) | ≥ 80% (expected ≥ 90%) | `idle_sub_state_test.go` table per §8 covers all branches of `TransitionTo` plus `OnStart` / `OnFinish` no-op assertions |
| `internal/game/entity/actors/states/` | unchanged | Existing `grounded_state_test.go` still exercises Idle indirectly; method rename touches all four sub-states but not their semantics |
| Overall `internal/engine/` + `internal/game/` delta | non-negative | No production code paths are added or removed; only one type relocates and one method renames |

The Coverage Analyzer must report a non-negative delta. Because the only code change in `gamestates` is a method rename and one factory branch, percentage coverage there is unchanged. The new `kit/states` package adds covered code → overall percentage rises slightly or stays flat.

## 15. Out of Scope (re-confirmed)

- Moving `GroundedInput`, `GroundedSubStateEnum`, or any other sub-state to kit (story 047).
- Adding any contract to `internal/engine/contracts/`.
- Generating mocks: kit defines no new interface that crosses a system boundary; `GroundedInputLike` is satisfied directly by the existing public `GroundedInput`. **No Mock Generator work is required.**
- Wiring the dependency-rule checks into actual CI YAML (deferred — manual run for now).
- Full ADR-006 rewrite.

## 16. Pipeline Hand-off

Because no new contracts are introduced, the next agents are:

1. **TDD Specialist** — author `internal/kit/states/idle_sub_state_test.go` per §8 (must fail to compile / fail to run because the type does not yet exist → Red).
2. **Feature Implementer** — execute migration steps §6 in order; verify §10 post-conditions.
3. **Workflow Gatekeeper** — run §7 dependency-rule checks, full test suite, coverage delta, then `golangci-lint run ./...`; on success move folder to `done/`.

Mock Generator is **skipped** for this story.
