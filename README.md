# Growbel

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
│   ├── audio/           # Audio files
│   ├── fonts/           # Font files
│   ├── images/          # Image files
│   ├── particles/       # Particle effect configurations
│   ├── sequences/       # Scripted sequences (JSON)
│   └── tilemap/         # Tilemap related assets (TMJ, TSX, PNG)
├── main.go              # Application entry point
├── internal/
│   ├── engine/          # Core game engine components
│   │   ├── app/         # Main engine loop, context, and initialization
│   │   ├── assets/      # Asset loading and management (images, fonts)
│   │   ├── audio/       # Audio playback functionality
│   │   ├── contracts/   # Interfaces for engine components (animation, body, config, context, navigation, sequences, tilemaplayer)
│   │   ├── data/        # Data loading, management, and configuration schemas
│   │   │   └── config/  # Engine configuration
│   │   ├── entity/      # Foundational structures for in-game objects (actors, items)
│   │   │   ├── actors/  # Actor management and movement (e.g., characters, enemies)
│   │   │   └── items/   # Item management
│   │   ├── event/       # Event handling system
│   │   ├── input/       # User input handling
│   │   ├── mocks/       # Test mocks for engine components
│   │   ├── physics/     # Physics simulation (body, movement, skill, space)
│   │   ├── render/      # Rendering tasks (camera, particles, screenutil, sprites, tilemap, vfx)
│   │   │   └── camera/  # Camera control and rendering
│   │   ├── scene/       # Game scene management and transitions
│   │   ├── sequences/   # Game sequences and command processing
│   │   ├── ui/          # Building blocks for user interface elements (hud, speech)
│   │   └── utils/       # Utility functions (fixed-point arithmetic, timing, triggers)
│   └── game/            # Game-specific implementation
│       ├── app/         # Game-specific setup and initialization (config, phases list, setup)
│       ├── entity/      # Concrete game entities (actors, items, obstacles, types)
│       │   ├── actors/  # Game-specific characters (player, npcs, enemies) and state logic
│       │   ├── items/   # Game-specific items (coins, falling platforms)
│       │   └── obstacles/ # Game-specific obstacles (walls, hazards)
│       ├── physics/     # Game-specific physics behaviors and skills
│       ├── render/      # Game-specific rendering logic (vfx)
│       ├── scenes/      # Game scenes and phases (intro, menu, story, summary)
│       └── ui/          # Game's specific user interface (hud, speech)
├── go.mod               # Go module definition
└── README.md
```

## Dependencies

- **Ebitengine**: A dead simple 2D game engine for Go.
- **EbitenUI**: A UI library for Ebitengine.
- **Kamera/v2**: A camera library for Ebitengine.
- **Go**: The programming language.

## Code Style Guidelines

### Avoid `_ = variable` Pattern

Do **not** use `_ = variable` to silence unused variable warnings in production code. This pattern clutters code and hides potential issues.

**Instead, use one of these approaches:**

1. **Use blank identifier in parameter list** (for unused params):

   ```go
   func (t *Transition) Draw(_ *ebiten.Image) {}
   ```

2. **Remove unused variables entirely** if not needed

3. **Actually use the variable** if it should be used

**Acceptable uses of `_`:**

- Ignoring return values: `_, err := someFunc()`
- Blank identifier in assignments: `val, _ = map[key]`
