# SPEC-016 — PhasesScene Test Coverage

**Branch:** `016-phases-scene-tests`
**Bounded Context:** Game Logic
**Package under test:** `internal/game/scenes/phases`

---

## Revision Note

Revision 2 expands scope to allow `Update()` path testing via mocks. The original SPEC excluded `Update` entirely due to GPU concerns, but the non-GPU paths inside `Update` (goal completion, death/completion triggers, bullet cleanup, sequence player update) are fully exercisable with package-local mocks. `Draw`, `OnStart`, `initTilemap`, and `NewPhasesScene` remain excluded (GPU/AppContext-dependent).

---

## Technical Requirements

### No new production code
All changes are test-only. No interfaces or production files are added or modified.

### Test file layout
| File | Purpose |
|---|---|
| `scene_test.go` | All tests (existing + new `Update` path tests) |
| `mocks_test.go` | Package-local mocks |

### Contracts used (read-only)
| Contract | Import path |
|---|---|
| `body.BodiesSpace` | `internal/engine/contracts/body` |
| `sequences.Player` | `internal/engine/contracts/sequences` |
| `phases.Goal` | `internal/engine/scene/phases` |
| `scene.FreezeController` | `internal/engine/scene` |

---

## Pre-conditions

- `PhasesScene` is constructed directly as a struct literal (zero-value fields + only the fields needed per test), bypassing `NewPhasesScene` entirely.
- `Update()` is called on a minimal `PhasesScene` with:
  - `pauseScreen = nil` (skips pause branch)
  - `screenFlipper = nil` (skips flip branch)
  - `gameCamera = nil` (skips camera branch)
  - `hasPlayer = false` (skips player-fall-death and collision branches)
  - `AppContext` set to a minimal `*app.AppContext{Space: mockBodiesSpace, VFX: nil}`
  - `BaseScene` initialised via `scene.NewTilemapScene(ctx)` so `BaseScene.Update()` does not panic
- GPU-dependent functions (`Draw`, `OnStart`, `initTilemap`, `NewPhasesScene`) remain excluded.

## Post-conditions

- Coverage for `internal/game/scenes/phases` ≥ 20%.
- Zero GPU or `ebiten.RunGame` calls in any test.
- All tests pass with `go test ./internal/game/scenes/phases/...`.

---

## Integration Points

- `phases.Goal` interface — tested via all three concrete implementations (existing).
- `scene.FreezeController` — tested in isolation (existing).
- `BodyCounter.setBodyCounter` — tested with `mockBodiesSpace` (existing).
- Bullet cleanup loop — tested via struct-literal harness (existing).
- `Update()` non-GPU paths — new tests targeting:
  - Goal completion branch: `goal.IsCompleted() == true && !completionTrigger.IsEnabled()` → calls `goal.OnCompletion()`.
  - Goal partial: `goal.IsCompleted() == false` → `OnCompletion` not called.
  - `completionTrigger` and `deathTrigger` update paths (trigger fires → `SceneManager.NavigateTo` called).
  - `sequencePlayer.Update()` called each frame.

---

## Package-local Mocks (`mocks_test.go`)

### Existing mocks (unchanged)
- `mockBodiesSpace` — implements `body.BodiesSpace`.
- `mockSequencePlayer` — implements `sequences.Player`.
- `mockCollidable` — minimal `body.Collidable` stub.

### New mocks
#### `mockGoal`
Implements `phases.Goal`. Exposes `completed bool` and `onCompletionCalled bool`.

#### `mockSceneManager`
Implements `navigation.SceneManager`. Exposes `navigateToCalled bool` and `navigateBackCalled bool`. Required to satisfy `AppContext.SceneManager` without a real scene graph.

---

## Red Phase — Failing Test Scenarios

**New failing tests (red until implemented):**

| Test name | Red condition |
|---|---|
| `TestUpdate_GoalCompletion_CallsOnCompletion` | no test yet |
| `TestUpdate_GoalPartial_DoesNotCallOnCompletion` | no test yet |
| `TestUpdate_SequencePlayerUpdated` | no test yet |
| `TestUpdate_DeathTrigger_NavigatesToReboot` | no test yet |

These tests fail (compile error / assertion failure) until `scene_test.go` is extended and `mockGoal` / `mockSceneManager` are added to `mocks_test.go`.

---

## Design Decisions

- `Update()` is called with `hasPlayer = false` and `pauseScreen = nil` to avoid all GPU-dependent branches. The goal/trigger/sequence paths execute unconditionally of player presence.
- `mockGoal` is package-local (single-package use) — not added to `internal/engine/mocks/`.
- `mockSceneManager` is package-local for the same reason.
- `config.Get().CamDebug` is `false` by default in tests (no config file present), so the `CamDebug` branch is skipped without any setup.
- `completionTrigger.Trigger()` fires after `completionTrigger.Enable(n)` + `n` calls to `Update()`. Tests drive this by calling `s.Update()` in a loop or by pre-enabling the trigger with a 1-frame delay.
