# Ebitengine Boilerplate

A minimal, clean boilerplate for starting new game projects with Ebitengine, based on the Firefly engine core.

## Features

- **Core Engine:** Robust entity management, scene management, and physics.
- **Platformer Foundation:** Ready-to-use platformer logic in `internal/game/scenes/phases/`.
- **I18n:** Built-in internationalization support.
- **Input:** Unified input handling with last-pressed-wins directional priority.
- **VFX:** Particle and floating text systems.
- **Boilerplate Scene:** A minimal main menu to get you started.
- **Advanced Player Mechanics:** Duck state, variable jump height, tween-based dash, one-way platform drop-through, and composite grounded sub-state machine.
- **Scene Freeze Frame:** Hit-stop effect for impactful gameplay moments.

## Getting Started

1.  Clone this repository.
2.  Install Go 1.25 or later.
3.  Run `go run main.go` to see the boilerplate in action.
4.  Modify `internal/game/` to implement your own game logic.

## Documentation

- **[`.agents/WORKFLOW.md`](.agents/WORKFLOW.md)**: Spec-Driven Development (SDD) pipeline for implementing features with formal specs and TDD.
- **[`.agents/work/`](.agents/work/)**: Active, backlog, and done stories — source of truth for in-progress work.
- **[`AGENTS.md`](AGENTS.md)**: Testing strategy, patterns, coverage guidelines, and code style rules for agent work.
- `internal/engine/README.md`: Details about the engine core.
- `internal/engine/combat/README.md`: Combat system overview (weapons, projectiles, inventory).
- `internal/game/`: Example game implementation (menu + platformer phases).
- `docs/adr/`: Architecture Decision Records explaining key non-obvious design choices:
  - [ADR-001](docs/adr/ADR-001-fp16-fixed-point-arithmetic.md) — FP16 fixed-point arithmetic for positions
  - [ADR-002](docs/adr/ADR-002-registry-based-state-pattern.md) — Registry-based actor state pattern with `init()` registration
  - [ADR-003](docs/adr/ADR-003-goroutine-audio-looping.md) — Goroutine-based audio looping
  - [ADR-004](docs/adr/ADR-004-space-body-model-physics.md) — Space / Body / MovementModel physics layers
  - [ADR-005](docs/adr/ADR-005-composite-grounded-sub-state.md) — Composite grounded sub-state machine
  - [ADR-006](docs/adr/ADR-006-engine-game-layer-separation.md) — Engine/Game two-layer architecture
  - [ADR-007](docs/adr/ADR-007-fp16-scale-factor.md) — FP16 scale factor is 16, not 65536

## License

MIT
