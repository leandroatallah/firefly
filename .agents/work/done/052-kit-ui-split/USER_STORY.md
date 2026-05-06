# User Story — 052-kit-ui-split

## Title

UI Layer Split — Dialogue Orchestrator Moves to `kit/ui/speech/`; All Primitives and Menu Stay in `engine/ui/`

## As a...

Engine developer building reusable genre-level tooling across multiple games

## I want...

The `internal/engine/ui/` package reorganised so that only the dialogue orchestrator — which composes `engine/audio`, `engine/data/config`, and `engine/input` into a multi-line typing flow — is relocated to a new `internal/kit/ui/speech/` package, while all primitive widgets (HUD, menu selection, raw speech base) remain in `engine/ui/`

## So that...

Reusable dialogue orchestration code lives at the correct layer (`kit`) and game code can import it from `internal/kit/ui/speech/` without mixing composition logic into the engine layer, keeping the three-layer dependency rule clean and seeding `kit/ui/` for future composite patterns.

## Background

ADR-006 (superseded by story 046) established a three-layer architecture: `engine` ← `kit` ← `game`. Story 046 introduced `internal/kit/` and validated the layer with the `IdleSubState` migration. Subsequent stories (047–051) extended kit with states, actors, combat, weapons, and skills.

The `internal/engine/ui/` package currently contains code at two distinct levels of abstraction:

**Primitive layer (belongs in `engine/ui/`)** — types that provide raw drawing, lifecycle contracts, or simple selection mechanics with no knowledge of engine subsystems beyond Ebiten and `engine/input`/`engine/assets/font` at most:
- `internal/engine/ui/hud/hud.go` — `HUD` interface, `BaseHUD` struct, `Manager` struct. Import: only `github.com/hajimehoshi/ebiten/v2`. Pure primitive; stays.
- `internal/engine/ui/menu/menu.go` — `Menu` and `MenuItem`. Imports `internal/engine/input` and `internal/engine/assets/font`. Implements a vertical-list keyboard-navigated selection primitive (index, arrow keys, callbacks). The imports are engine-internal and do not violate the dependency rule. Menu selection is a primitive mechanic; stays in `engine/ui/menu/`.
- `internal/engine/ui/speech/speech.go` — `Speech` interface and `SpeechBase` struct. Core spelling/visibility/color state machine with no dependency on `input`, `audio`, or `config`. Primitive; stays.
- `internal/engine/ui/speech/font.go` — `SpeechFont` wrapper around `engine/assets/font`. Thin rendering helper; stays.
- `internal/engine/ui/speech/speech_ext.go` — `SetPosition`/`GetPosition`/`SetSpeed`/`GetSpeed` accessors on `SpeechBase`. Primitive accessors; stay.

**Composite layer (belongs in `kit/ui/`)** — types that compose multiple engine subsystems into reusable orchestration patterns that no individual game should need to rewrite:
- `internal/engine/ui/speech/dialogue.go` — `Manager` (speech dialogue orchestrator). Imports `internal/engine/audio`, `internal/engine/data/config`, `internal/engine/input`. Manages active speech, line sequencing, typing-sound scheduling, spelling-skip logic. Carries `BubbleSpeechID` and `StorySpeechID` constants. Move to `internal/kit/ui/speech/dialogue.go`.

`internal/engine/scene/pause/pause.go` (`PauseScreen`) uses `engine/ui/menu` and is **not** affected by this story — menu stays in engine, so the pause scene import path is unchanged.

The game layer already has its own concrete speech types in `internal/game/ui/speech/` (`SpeechBubble`, `StorySpeech`, `baseSpeech`) which correctly extend the engine primitives and will update their import of `Manager` to `internal/kit/ui/speech/` after this split.

## Acceptance Criteria

- [ ] **New package created**: `internal/kit/ui/speech/` exists and contains a `doc.go` declaring that it imports only `internal/engine/` packages and has zero knowledge of `internal/game/`.

- [ ] **`speech.Manager` (dialogue orchestrator) relocated**: `internal/engine/ui/speech/dialogue.go` moves to `internal/kit/ui/speech/dialogue.go`. The original file is deleted. The `BubbleSpeechID` and `StorySpeechID` string constants travel with it to `kit/ui/speech/`. All callers inside `internal/game/` (`SpeechBubble`, `StorySpeech`) update their import path from `internal/engine/ui/speech` (for `Manager`, `BubbleSpeechID`, `StorySpeechID`) to `internal/kit/ui/speech`.

- [ ] **`engine/ui/menu/` is untouched**: `internal/engine/ui/menu/menu.go` is **not** moved. Its import path remains `internal/engine/ui/menu` everywhere it is used (`scene_menu.go`, `phases/scene.go`, `pause/pause.go`).

- [ ] **Primitives remain in `engine/ui/`**: `internal/engine/ui/hud/hud.go`, `internal/engine/ui/speech/speech.go`, `internal/engine/ui/speech/font.go`, and `internal/engine/ui/speech/speech_ext.go` are **not** moved. Their import paths are unchanged.

- [ ] **Dependency rule — engine is clean**: Running `go list -deps ./internal/engine/...` produces no paths containing `internal/kit` or `internal/game`. CI check must pass:
  ```
  go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)' && exit 1 || true
  ```

- [ ] **Dependency rule — kit is clean**: Running `go list -deps ./internal/kit/...` produces no paths containing `internal/game`. CI check must pass:
  ```
  go list -deps ./internal/kit/... | grep 'internal/game' && exit 1 || true
  ```

- [ ] **All existing tests pass**: `go test ./...` is green with no regressions after the move.

- [ ] **Tests travel with code**: Test files for `dialogue.go` that currently reside in `internal/engine/ui/speech/` move alongside the file to `internal/kit/ui/speech/`.

- [ ] **Coverage non-negative**: `internal/kit/ui/speech/` reaches 80%+ coverage. Overall coverage delta across `internal/engine/` and `internal/game/` is non-negative.

- [ ] **No hardcoded game-specific strings in `kit/ui/`**: `internal/kit/ui/speech/` must not reference any game-specific string literals (e.g., asset paths, character names). Audio paths and label strings remain wired in `internal/game/`.

- [ ] **`.agents/constitution.md` Bounded Contexts table is current**: Verify the `Kit` row already present covers `internal/kit/ui/`; no new row is required. The Gatekeeper must confirm this before closing the story.

## Out of Scope

- **`internal/engine/ui/menu/menu.go`** — menu selection is a primitive mechanic with engine-internal imports; it stays in `engine/ui/menu/`. Moving it is explicitly out of scope for this story.
- **`internal/engine/scene/pause/pause.go`** (`PauseScreen`) — unaffected because menu stays in engine.
- New UI features of any kind (new widget types, animations, layout algorithms).
- Theme system or visual redesign of any existing widget.
- Replacing Ebitengine UI primitives with a third-party UI framework.
- HUD redesign or adding new HUD elements.
- Full rewrite of ADR-006 (already deferred from story 046).
- Splitting `kit` into sub-modules.

## Domain Notes

- A **Contract** is a Go interface in `internal/engine/contracts/`. Contracts stay in `engine`. `kit` implements or orchestrates them; `game` wires them.
- The **three-layer dependency rule**: `engine` imports nothing from `kit` or `game`; `kit` imports only from `engine`; `game` imports from both.
- `menu.Menu` is a concrete primitive with engine-internal imports (`engine/input`, `engine/assets/font`). This does not violate the dependency rule. Its home remains `engine/ui/menu/`.
- `speech.Manager` (dialogue orchestrator) is a concrete orchestrator that composes `engine/audio`, `engine/data/config`, and `engine/input` — a multi-subsystem composition. Its home is `kit/ui/speech/`.
- `speech.SpeechBase`, `speech.Speech` (interface), `speech.SpeechFont` remain in `engine/ui/speech/` because they depend on nothing beyond Ebiten's drawing API and `engine/assets/font`.
- `hud.HUD`, `hud.BaseHUD`, `hud.Manager` remain in `engine/ui/hud/` because they depend only on Ebiten's drawing API.
- The game-layer concrete types `SpeechBubble` and `StorySpeech` (`internal/game/ui/speech/`) already correctly extend engine primitives; after this story they import `Manager`, `BubbleSpeechID`, and `StorySpeechID` from `kit/ui/speech/` instead of `engine/ui/speech/`.
- `kit/ui/` is intentionally seeded with only the dialogue manager in this story. A thin package is acceptable; it provides the correct home for future composite UI patterns.

## Branch Name

`052-kit-ui-split`
