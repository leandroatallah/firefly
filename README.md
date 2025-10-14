# Firefly

A 2D game built with Ebitengine, featuring a modular architecture.

## Architecture Overview

The project is structured into two main packages: `engine` and `game`.

-   **`internal/engine`**: The core game engine, providing reusable components for scenes, physics, actors, and other systems.
-   **`internal/game`**: The specific implementation of the game, including scenes, characters, and items.

This separation allows the engine to be developed independently from the game's content.

## Folder Structure

```
.
├── assets/              # Game assets (images, sounds, etc.)
├── cmd/game/            # Application entry point
├── internal/
│   ├── config/          # Game configuration
│   ├── engine/          # Core game engine
│   │   ├── actors/      # Base actor components
│   │   ├── core/        # Scene management, game loop
│   │   └── systems/     # Physics, audio, input
│   └── game/            # Game-specific implementation
│       ├── scenes/      # Game scenes
│       ├── actors/      # Game characters
│       └── items/       # Game items
├── go.mod               # Go module definition
└── README.md
```

## Dependencies

-   **Ebitengine**: A dead simple 2D game engine for Go.
-   **Go**: The programming language.
