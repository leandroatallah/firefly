# ADR-006 — Engine/Game Two-Layer Architecture

## Status
Superseded by 046
> **Note:** Superseded by story `046-kit-layer-validation-step`, which validates a three-layer architecture (`engine` ← `kit` ← `game`) by relocating `IdleSubState` to `internal/kit/states/`. A full ADR rewrite is deferred until follow-up stories 047–052 populate the `kit` layer; see `.agents/work/active/046-kit-layer-validation-step/SPEC.md` for the validation rationale.

## Context
The codebase separates reusable engine code from game-specific implementation. This raises the question: where should concrete implementations (states, skills, scenes) live? Three options exist:

1. **Two-layer:** Engine contains both abstractions and common concrete implementations; game contains project-specific code.
2. **Strict two-layer:** Engine contains only abstractions; all concrete code lives in game.
3. **Three-layer:** Engine (abstractions) + Library (reusable concrete) + Game (project-specific).

## Decision
Use a **two-layer architecture** with pragmatic placement rules:

- **`internal/engine/`**: Abstractions (contracts, core systems) + common concrete implementations used in 80%+ of games in the genre (e.g., `IdleState`, `WalkingState`, `JumpingState`, `JumpSkill` for platformers).
- **`internal/game/`**: Project-specific concrete implementations (e.g., `ClimbingState`, `SwimmingState`, custom skills, all level scenes).

## Consequences
- Clear separation without unnecessary folder overhead.
- Engine stays reusable for similar game genres without bloat.
- Developers extend the engine by adding game-specific states/skills/scenes to `internal/game/` and registering them with the engine's systems.
- If multiple games share the same concrete implementations in the future, refactor to a three-layer architecture (engine/library/game).
- Risk: Engine may accumulate "maybe reusable" code. Mitigation: Apply the 80% rule strictly—if a state isn't genre-defining, it belongs in game.
