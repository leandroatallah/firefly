# ADR-006 — Three-Layer Architecture: Engine / Kit / Game

## Quick Reference

- **When to cite:** Deciding where to place new code (engine vs kit vs game).
- **Key constraint:** `game` → `kit` → `engine`. No reverse imports.
- **DO:** 80%+ genre reuse → `kit`; game-specific → `game`; abstract/interface → `engine`.
- **DON'T:** Add genre-specific code to `engine`; add game-specific code to `kit`.

## Status
Accepted (amended from two-layer; kit layer validated by stories 046–052)

## Context
The codebase separates reusable engine code from game-specific implementation. The original decision (two-layer: `engine` + `game`) placed both abstractions and genre-reusable concrete implementations in `engine`. As the concrete implementations grew (skills, combat, actor archetypes, UI dialogue), this created a tension: `engine` was accumulating code that wasn't truly engine-level but was reusable across games in the same genre.

Story `046-kit-layer-validation-step` introduced a `kit` layer as a proof of concept, and stories 047–052 completed the migration.

## Decision
Use a **three-layer architecture**:

- **`internal/engine/`**: Abstractions (contracts/interfaces) + minimal core systems (`app`, `entity`, `physics`, `scene`, `sequences`, `render`, `audio`, `input`, `data`). Nothing here should depend on genre-specific concepts.
- **`internal/kit/`**: Genre-reusable concrete implementations that depend on `engine` contracts but are not game-specific. Includes:
  - `combat/` — weapon inventory, projectile lifecycle, melee controller, faction system.
  - `skills/` — `JumpSkill`, `DashSkill`, `HorizontalMovementSkill`, `ShootingSkill`, and a JSON `FromConfig` factory.
  - `actors/` — reusable actor trait compositions (`MeleeCharacter`, `ShooterCharacter`, `DeathBehavior`).
  - `states/` — genre-reusable state machine states.
  - `ui/speech/` — `speech.Manager` dialogue implementation.
- **`internal/game/`**: Project-specific concrete implementations (player, enemies, scenes, custom states, level-specific logic).

## Placement Rules

| Question | Layer |
|---|---|
| Is this an abstraction (interface/contract)? | `engine` |
| Is this a core system needed by all game types? | `engine` |
| Would 80%+ of games in this genre use this as-is? | `kit` |
| Is this specific to *this* game's design? | `game` |

## Consequences
- Engine stays minimal and genre-agnostic.
- Kit provides a reusable genre library without polluting engine.
- Game layer stays thin: it wires kit components to project-specific art, levels, and rules.
- Dependency direction is strict: `game` → `kit` → `engine`. No reverse dependencies.
- Risk: Kit may accumulate "maybe reusable" code. Apply the 80% rule strictly — if a component is only likely used in this one game, it belongs in `game`.
