# Project Constitution

Every agent in this pipeline MUST read this file before producing any output. These constraints are non-negotiable and apply to all specs, tests, and production code.

## Project Identity

- **Engine**: Ebitengine (Go)
- **Language**: Go 1.25+
- **Architecture**: three-layer — `engine`, `kit`, `game`; dependency rules: engine must not import kit or game; kit may import engine, must not import game; game may import both engine and kit
- **Core module**: `internal/engine/`
- **Kit module**: `internal/kit/` — genre-reusable concrete implementations; may import engine, must not import game
- **Game module**: `internal/game/`

## Ubiquitous Language (DDD)

Use these terms consistently across all stories, specs, and code:

| Term | Meaning |
|---|---|
| `Actor` | An entity with a state machine (player, enemy) |
| `Item` | A collectible or interactable entity |
| `Body` | A physics body with position and velocity |
| `Space` / `BodiesSpace` | The physics world that holds and updates bodies |
| `Scene` | A self-contained game state (menu, phase, cutscene) |
| `Phase` | A playable game level scene |
| `Sequence` | A scripted, frame-by-frame event chain |
| `Contract` | A Go interface in `internal/engine/contracts/` |
| `State` | A named node in an Actor's state machine |

## Bounded Contexts

| Context | Package |
|---|---|
| Physics | `internal/engine/physics/` |
| Entity | `internal/engine/entity/` |
| Scene | `internal/engine/scene/` |
| Sequences | `internal/engine/sequences/` |
| Input | `internal/engine/input/` |
| i18n | `internal/engine/data/i18n/` |
| Game Logic | `internal/game/` |
| Kit | `internal/kit/` |

## Non-Negotiable Standards

### Production Code
- No `_ = variable` pattern. Use blank identifier in parameter lists: `func (t *T) Method(_ Type) {}`.
- No global mutable state. Inject dependencies via interfaces (`internal/engine/contracts/`).
- Fixed-point positions use `x16`/`y16` — always validate with `fp16.From16()` / `fp16.To16()`.

### Tests
- Table-driven tests for all logic with multiple input/output scenarios.
- No `ebiten.RunGame` or GPU-dependent calls in unit tests. Use `ebiten.NewImage(w, h)` headlessly.
- No `time.Sleep`. Use frame counters or virtual time.
- Tests must be deterministic and non-flaky.
- Mock only at system boundaries using interfaces from `internal/engine/contracts/`.
- Shared mocks (used across packages) → `internal/engine/mocks/`.
- Package-local mocks → `mocks_test.go` in the same package.

### Git
- Every story maps to a feature branch: `[ID]-[short-kebab-description]` (e.g., `42-player-dash-state`).
- Branch is created before any code is written.
- No commits are made by agents. The developer owns all version control operations.
- No `Co-Authored-By` or similar AI attribution trailers in commit messages.

## Coverage Goal

80%+ across `internal/engine/` and `internal/game/`. The Gatekeeper must confirm a positive coverage delta before closing any story.

## Extension Hooks

Before writing a spec, check if `.agents/hooks/before_spec.md` exists. If it does, read and follow its instructions before proceeding.
