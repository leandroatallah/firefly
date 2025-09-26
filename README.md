# Firefly

A 2D game built with Ebitengine (formerly known as Ebiten) featuring a modular architecture with scene management, physics system, and audio support.

## Architecture Overview

The project follows a layered architecture pattern with clear separation of concerns:

### Core Layers

1. **Application Layer** (`cmd/`) - Entry point and application setup
2. **Core Layer** (`internal/core/`) - Game engine, scene management, and transitions
3. **Domain Layer** (`internal/actors/`) - Game entities and business logic
4. **Systems Layer** (`internal/systems/`) - Cross-cutting concerns (physics, audio, input)
5. **Configuration Layer** (`internal/config/`) - Game constants and settings

## Package Structure

### Application Layer
- **`cmd/game/`** - Application entry point and initialization

### Core Layer
- **`internal/core/game/`** - Main game engine and state management
- **`internal/core/scene/`** - Scene management system with factory pattern
- **`internal/core/transition/`** - Scene transition effects
- **`internal/core/screenutil/`** - Screen utility functions

### Domain Layer
- **`internal/actors/`** - Game entities and character system
  - Character base class with physics integration
  - Player and enemy implementations
  - Sprite management and animation states
  - Enemy factory for entity creation

### Systems Layer
- **`internal/systems/physics/`** - Physics engine and collision detection
  - Body physics with movement and collision
  - Shape definitions (rectangles, etc.)
  - Obstacle system for level boundaries
- **`internal/systems/input/`** - Input handling utilities
- **`internal/systems/audiomanager/`** - Audio system with OGG/WAV support

### Configuration Layer
- **`internal/config/`** - Game constants and configuration

## Folder Structure

```text
.
├── README.md
├── assets
│   └── ...
├── cmd
│   └── game
│       └── main.go
├── go.mod
├── go.sum
└── internal/
    ├── actors/                      # Game entities and characters
    │   ├── character.go             # Base character class
    │   ├── enemy_factory.go         # Enemy creation factory
    │   ├── enemy.go                 # Enemy implementation
    │   ├── player.go                # Player implementation
    │   ├── sprite.go                # Sprite management
    │   └── state.go                 # Character state machine
    ├── config/
    │   └── constants.go             # Game configuration constants
    ├── core/
    │   ├── game/                    # Game engine core
    │   │   ├── engine.go            # Main game engine
    │   │   ├── setup.go             # Game initialization
    │   │   ├── state_mainmenu.go    # Main menu state
    │   │   └── state.go             # Game state interface
    │   ├── scene/                   # Scene management
    │   │   ├── base_scene.go        # Base scene implementation
    │   │   ├── factory.go           # Scene factory pattern
    │   │   ├── manager.go           # Scene manager
    │   │   ├── menu_scene.go        # Menu scene
    │   │   └── sandbox_scene.go     # Game play scene
    │   ├── screenutil/
    │   │   └── utils.go             # Screen utility functions
    │   └── transition/              # Scene transitions
    │       ├── fader_transition.go  # Fade transition effect
    │       └── transition.go        # Transition interface
    └── systems/                     # System services
        ├── audiomanager/
        │   └── audio.go             # Audio management system
        ├── input/
        │   └── input.go             # Input handling utilities
        └── physics/                 # Physics engine
            ├── body.go              # Physics body implementation
            ├── collision.go         # Collision detection
            ├── movement.go          # Movement physics
            ├── obstacle_factory.go  # Obstacle creation
            ├── obstacle.go          # Obstacle implementation
            ├── shape_rect.go        # Rectangle shape
            └── shape.go             # Shape interface
```

## Dependencies

- **Ebitengine v2.8.8** - 2D game engine
- **Audio libraries** - OGG/Vorbis and WAV support
- **Go 1.23.6** - Programming language
