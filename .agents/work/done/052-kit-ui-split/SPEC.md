# SPEC — 052-kit-ui-split

**Branch:** `052-kit-ui-split`
**Status:** Active
**Source story:** `USER_STORY.md` (same folder)

## 1. Summary

Relocate the dialogue orchestrator (`speech.Manager` and its `BubbleSpeechID` / `StorySpeechID` constants) from `internal/engine/ui/speech/` to a new `internal/kit/ui/speech/` package. Speech primitives (`Speech`, `SpeechBase`, `SpeechFont`, `speech_ext.go`) and all of `engine/ui/menu/` and `engine/ui/hud/` stay where they are. The result must satisfy the three-layer dependency rule: `engine` ← `kit` ← `game`.

## 2. Critical Pre-Spec Finding (Architecture)

A direct file move is **not sufficient** to satisfy AC-5 (engine has no `internal/kit` deps). Two engine-layer files currently depend on the concrete `speech.Manager` type:

1. `internal/engine/app/context.go` — `AppContext.DialogueManager *speech.Manager` (field uses concrete type).
2. `internal/engine/sequences/commands.go` — `DialogueCommand.dialogueManager *speech.Manager` plus the constant `speech.BubbleSpeechID`.

Naively moving `Manager` to `kit/ui/speech/` and updating these imports would make `engine/app` and `engine/sequences` import `internal/kit/...`, violating the constitution.

**Resolution (mandatory):** Introduce a new contract package `internal/engine/contracts/dialogue/` that exposes the engine-facing surface of the dialogue orchestrator as a Go interface. Engine code (`AppContext`, `DialogueCommand`) depends on this interface; `kit/ui/speech/Manager` implements it; `game/app/setup.go` wires the concrete `kit` Manager into the `AppContext` field typed as the contract.

This adheres to the constitution rule: "A Contract is a Go interface in `internal/engine/contracts/`. Contracts stay in engine. kit implements or orchestrates them; game wires them."

## 3. Target Layout

```
internal/engine/contracts/dialogue/
    dialogue.go                  ← NEW (interface)

internal/engine/ui/speech/
    speech.go                    ← unchanged (Speech interface + SpeechBase)
    speech_test.go               ← unchanged
    speech_ext.go                ← unchanged
    font.go                      ← unchanged
    font_test.go                 ← unchanged
    (dialogue.go REMOVED)
    (dialogue_test.go REMOVED)
    (mocks_test.go MOVED — see below)

internal/kit/ui/speech/
    doc.go                       ← NEW
    dialogue.go                  ← MOVED from engine/ui/speech/dialogue.go
    dialogue_test.go             ← MOVED from engine/ui/speech/dialogue_test.go
    mocks_test.go                ← MOVED from engine/ui/speech/mocks_test.go
                                   (only used by dialogue_test.go;
                                    speech_test.go does not reference mockSpeech)
```

If, after the move, `speech_test.go` or `font_test.go` still requires the `mockSpeech` helper, a copy stays in `engine/ui/speech/mocks_test.go`. **Pre-condition check (verify during implementation):** `grep -n "mockSpeech" internal/engine/ui/speech/speech_test.go internal/engine/ui/speech/font_test.go`. If empty, full move; otherwise duplicate or trim.

## 4. New Contract — `internal/engine/contracts/dialogue/dialogue.go`

Package name: `dialogue`. Defines the surface that `engine/sequences/DialogueCommand` and any other engine consumer needs. Derived directly from current call sites:

```go
package dialogue

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/boilerplate/ebiten-template/internal/engine/audio"
    "github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
)

// Manager is the engine-facing contract for a dialogue orchestrator.
// Concrete implementation lives in internal/kit/ui/speech.
type Manager interface {
    // Lifecycle / per-frame
    Update() error
    Draw(screen *ebiten.Image)
    IsSpeaking() bool

    // Speech registry
    AddSpeech(s speech.Speech)
    SetSpeech(id string)
    SetActiveSpeech(id string)
    GetActiveSpeech() speech.Speech

    // Audio wiring
    SetAudioManager(m *audio.AudioManager)
    SetTypingSound(path string)
    SetTypingSounds(paths []string)
    SetSpeechAudioQueue(paths []string)
    ClearSpeechAudioQueue()
    SetDefaultSpeechAudio(paths []string)
    ApplyDefaultSpeechAudio(lineCount int)

    // Behaviour flags
    SetSpeechSkipEnabled(enabled bool)

    // Display
    ShowMessages(lines []string, position string, speed int)
}
```

Constants (well-known speech IDs) also become engine-level so `DialogueCommand` can reference them without importing `kit`:

```go
const (
    BubbleSpeechID = "bubble"
    StorySpeechID  = "story"
)
```

Rationale: these are protocol identifiers shared between engine sequence commands and game-layer speech setup. Their stable home is the contract package. The `kit/ui/speech` package will re-export them as aliases to preserve story AC wording ("constants travel with `Manager`"), so existing callers in `internal/game/` that import them through `kit/ui/speech` keep working.

In `kit/ui/speech/dialogue.go`:
```go
const (
    BubbleSpeechID = dialogue.BubbleSpeechID
    StorySpeechID  = dialogue.StorySpeechID
)
```

(Compile-time assertion in `kit/ui/speech` test or init: `var _ dialogue.Manager = (*Manager)(nil)`.)

## 5. File-Level Changes

### 5.1 Engine — modifications

| File | Change |
|---|---|
| `internal/engine/contracts/dialogue/dialogue.go` | NEW — interface + ID constants per §4. |
| `internal/engine/app/context.go` | Replace import `engine/ui/speech` → `engine/contracts/dialogue`. Change field `DialogueManager *speech.Manager` → `DialogueManager dialogue.Manager`. |
| `internal/engine/app/app_test.go` | Update import; constructor must build a test double or use `&fakeDialogueManager{}` implementing the contract. (Currently calls `speech.NewManager()`; replace with a minimal test stub or move test to `app_test.go` constructing through an injected mock.) |
| `internal/engine/sequences/commands.go` | Replace import `engine/ui/speech` → `engine/contracts/dialogue`. Change `dialogueManager *speech.Manager` → `dialogueManager dialogue.Manager`. Replace `speech.BubbleSpeechID` → `dialogue.BubbleSpeechID`. |
| `internal/engine/sequences/commands_test.go` | Update import to `engine/contracts/dialogue`; tests construct a fake `dialogue.Manager` (or use a generated mock from `internal/engine/mocks/`). |
| `internal/engine/ui/speech/dialogue.go` | DELETED. |
| `internal/engine/ui/speech/dialogue_test.go` | DELETED (moves to kit). |
| `internal/engine/ui/speech/mocks_test.go` | If only `dialogue_test.go` referenced `mockSpeech`, MOVE to kit; otherwise duplicate. |

### 5.2 Kit — new files

| File | Change |
|---|---|
| `internal/kit/ui/speech/doc.go` | Package comment per §6. |
| `internal/kit/ui/speech/dialogue.go` | Body of old `dialogue.go`, package renamed to `speech`. Imports unchanged except `dialogue` contract import for ID constant aliases. Add compile-time assertion `var _ dialogue.Manager = (*Manager)(nil)`. |
| `internal/kit/ui/speech/dialogue_test.go` | Moved test, package `speech`. |
| `internal/kit/ui/speech/mocks_test.go` | Moved (or copied) `mockSpeech` from engine. Imports `engine/ui/speech` to satisfy `speech.Speech` interface methods that take `*ebiten.Image`. |

### 5.3 Game — modifications

| File | Change |
|---|---|
| `internal/game/app/setup.go` | Add import `kitspeech "internal/kit/ui/speech"`. Keep existing import of `engine/ui/speech` for `SpeechFont`. Replace `speech.NewManager(...)` → `kitspeech.NewManager(...)` and `speech.BubbleSpeechID` → `kitspeech.BubbleSpeechID`. |
| `internal/game/ui/speech/bubble.go` | Add import `kitspeech "internal/kit/ui/speech"`. Replace `speech.BubbleSpeechID` → `kitspeech.BubbleSpeechID`. |
| `internal/game/ui/speech/story.go` | Same pattern for `StorySpeechID`. |
| `internal/game/ui/speech/common.go` | No change needed (only references `SpeechBase` / `SpeechFont`, which stay in engine). |
| `internal/game/ui/speech/speech_test.go` | Update if it references moved symbols (review during implementation). |

### 5.4 Engine — untouched (must verify)

- `internal/engine/ui/menu/menu.go` and all its callers (`scene_menu.go`, `phases/scene.go`, `pause/pause.go`).
- `internal/engine/ui/hud/hud.go`.
- `internal/engine/ui/speech/speech.go`, `font.go`, `speech_ext.go`.

## 6. `kit/ui/speech/doc.go` — Required Content

```go
// Package speech provides the dialogue orchestrator for genre-reusable UI.
//
// The orchestrator (Manager) composes engine subsystems — audio, config,
// input — into a multi-line typing flow with typing-sound scheduling and
// spelling-skip behaviour. Speech primitives (Speech, SpeechBase,
// SpeechFont) live in internal/engine/ui/speech and are not duplicated here.
//
// Dependency rule (enforced by CI):
//   - kit/ui/speech MAY import internal/engine/...
//   - kit/ui/speech MUST NOT import internal/game/...
//   - It implements internal/engine/contracts/dialogue.Manager.
package speech
```

## 7. Pre-Conditions

- Working tree on branch `052-kit-ui-split`.
- `go test ./...` is green at HEAD.
- No file under `internal/engine/contracts/dialogue/` exists yet.
- No directory `internal/kit/ui/` exists yet.

## 8. Post-Conditions

1. Directory `internal/kit/ui/speech/` exists with `doc.go`, `dialogue.go`, `dialogue_test.go`, `mocks_test.go`.
2. Directory `internal/engine/contracts/dialogue/` exists with `dialogue.go`.
3. `internal/engine/ui/speech/dialogue.go` and `internal/engine/ui/speech/dialogue_test.go` no longer exist.
4. `go build ./...` passes.
5. `go test ./...` passes with no regressions.
6. `go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)'` produces no output.
7. `go list -deps ./internal/kit/... | grep 'internal/game'` produces no output.
8. `internal/kit/ui/speech/` reaches ≥ 80 % coverage (essentially preserved by moved tests).
9. Coverage delta across `internal/engine/` + `internal/game/` is non-negative.
10. No game-specific string literals appear in `internal/kit/ui/speech/` (audio paths and labels remain in `internal/game/`).

## 9. State Machine / Behavioural Changes

**None.** This story is a pure relocation + interface extraction. `Manager`'s logic — line sequencing, typing-sound cooldown, spelling-skip, speech audio queue, accumulative mode — is preserved verbatim. No new states; no transition changes.

## 10. Integration Points

Within the Bounded Context **Kit ↔ Engine**:

- `kit/ui/speech.Manager` implements `engine/contracts/dialogue.Manager`.
- `engine/sequences.DialogueCommand` consumes `dialogue.Manager` via `AppContext.DialogueManager`.
- `engine/app.AppContext` exposes `DialogueManager dialogue.Manager` (interface field).
- `game/app.setup` constructs `kitspeech.NewManager(...)` and assigns it to `AppContext.DialogueManager`.
- `game/ui/speech.SpeechBubble` / `StorySpeech` continue extending `engine/ui/speech.SpeechBase` and now reference IDs via `kitspeech.BubbleSpeechID` / `kitspeech.StorySpeechID`.

## 11. Red Phase — Failing Test Scenario (for TDD Specialist)

Add a new file `internal/engine/contracts/dialogue/dialogue_test.go` (or a layer-rule test under `internal/engine/app/`) that fails until the relocation is performed. Two complementary failing tests:

### 11.1 Layer-rule test (primary Red)

`internal/engine/dependency_test.go` (package `engine_test`, build-tagged or plain) executing:

```go
func TestEngineDoesNotDependOnKitOrGame(t *testing.T) {
    out, err := exec.Command("go", "list", "-deps", "./internal/engine/...").Output()
    if err != nil { t.Fatalf("go list: %v", err) }
    re := regexp.MustCompile(`internal/(kit|game)`)
    if loc := re.FindString(string(out)); loc != "" {
        t.Fatalf("engine has forbidden dep: %s", loc)
    }
}
```

This passes today (engine doesn't import kit yet because nothing is in kit). It must **continue** to pass after the move — and would fail loudly if the implementer naively rewrites engine imports to point at `internal/kit/ui/speech`. It is the safety net.

### 11.2 Contract-implementation test (the actual Red)

`internal/kit/ui/speech/dialogue_test.go` adds:

```go
func TestManager_ImplementsDialogueContract(t *testing.T) {
    var _ dialogue.Manager = (*Manager)(nil)
}
```

Before any work: this test does not compile because:
- The package `internal/kit/ui/speech` does not exist.
- The package `internal/engine/contracts/dialogue` does not exist.
- `Manager` is still in `internal/engine/ui/speech`.

This is the canonical failing test that proves the architecture is in place once the implementer has finished.

### 11.3 Behavioural regression test (carried over)

The moved `dialogue_test.go` continues asserting:
- `NewManager(s1, s2)` registers two speeches.
- `SetActiveSpeech("speech1")` then `SetSpeech("speech2")` switches active id.
- `SetSpeech("non-existent")` is a no-op.
- `ShowMessages` initialises state and triggers `Show()` on the active speech.

These tests must pass post-move under the new package path with no logic change.

## 12. Risks & Mitigations

| Risk | Mitigation |
|---|---|
| Hidden engine consumer of `speech.Manager` not listed above | `grep -RIn "speech\\.Manager\\|speech\\.NewManager\\|speech\\.BubbleSpeechID\\|speech\\.StorySpeechID" internal/engine` before implementation; add to plan. |
| `mocks_test.go` shared by primitive tests | If `speech_test.go` or `font_test.go` use `mockSpeech`, keep a copy in `engine/ui/speech/mocks_test.go`. |
| `app_test.go` depends on concrete `*speech.Manager` zero value | Replace with a small in-test stub satisfying `dialogue.Manager`, or with a generated mock under `internal/engine/mocks/mock_dialogue.go` (Mock Generator stage handles this). |
| Coverage drop in `internal/engine/ui/speech/` after dialogue tests leave | Acceptable — coverage migrates to `internal/kit/ui/speech/`; aggregate delta is non-negative per AC. |

## 13. Out of Scope (re-stated from story)

- `internal/engine/ui/menu/menu.go` and `pause/pause.go` are not touched.
- `internal/engine/ui/hud/` is not touched.
- No new UI features, themes, or widget types.
- No splitting of `kit` into sub-modules.

## 14. Constitution Compliance Checklist

- [x] No `_ = variable` patterns introduced.
- [x] No global mutable state added (interface field on `AppContext` is dependency injection).
- [x] No `time.Sleep`; no `ebiten.RunGame` in tests.
- [x] New contract added under `internal/engine/contracts/` per the standard.
- [x] Bounded Contexts table already lists `Kit` (`internal/kit/`); no addition required.
- [x] Branch name matches story slug.

## 15. Pipeline Next Steps

This SPEC introduces a **new contract** (`internal/engine/contracts/dialogue/`). Therefore the next agents in order are:

1. **Mock Generator** — generate `internal/engine/mocks/mock_dialogue.go` for `dialogue.Manager` (used by `engine/sequences/commands_test.go`, `engine/app/app_test.go`, and any future consumer).
2. **TDD Specialist** — author the failing tests of §11.1, §11.2, and the relocated §11.3.
3. **Feature Implementer** — perform the file moves, create `dialogue` contract package, update imports across engine and game, add ID-constant aliases, and the `var _ dialogue.Manager = (*Manager)(nil)` assertion.
4. **Workflow Gatekeeper** — verify ACs, layer-rule script, coverage delta, run `golangci-lint run ./...`.
