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

- `internal/engine/README.md`: Details about the engine core.
- `internal/game/`: Example game implementation (menu + platformer phases).

## License

MIT
