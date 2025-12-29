# Firefly

A 2D game built with Ebitengine, featuring a modular architecture.

## Architecture Overview

The project is structured into two main packages: `engine` and `game`.

- **`internal/engine`**: The core game engine, providing reusable components for scenes, physics, actors, and other systems.
- **`internal/game`**: The specific implementation of the game, including scenes, characters, and items.

This separation allows the engine to be developed independently from the game's content.

## Folder Structure

```
.
├── assets/              # Game assets (images, sounds, etc.)
├── main.go            # Application entry point
├── internal/
│   ├── config/          # Game configuration
│   ├── engine/          # Core game engine components
│   │   ├── actors/      # Actor management and movement (e.g., characters, enemies)
│   │   ├── assets/      # Asset loading and management
│   │   ├── camera/      # Camera control and rendering
│   │   ├── contracts/   # Interfaces and contracts for engine components
│   │   ├── core/        # Game loop, scene management, and core utilities
│   │   ├── items/       # Item management
│   │   ├── sequences/   # Game sequences and command processing
│   │   └── systems/     # Various game systems (audio, physics, input, speech, etc.)
│   └── game/            # Game-specific implementation
│       ├── actors/      # Game-specific characters and enemies
│       ├── camera/      # Game-specific camera settings
│       ├── hud/         # Game-specific HUD elements
│       ├── items/       # Game-specific items
│       ├── scenes/      # Game scenes and phases
│       ├── setup/       # Game setup and initialization
│       ├── speech/      # Game-specific speech bubbles and dialogues
│       └── state/       # Game state management
├── go.mod               # Go module definition
└── README.md
```

## Dependencies

- **Ebitengine**: A dead simple 2D game engine for Go.
- **EbitenUI**: A UI library for Ebitengine.
- **Kamera/v2**: A camera library for Ebitengine.
- **Go**: The programming language.
