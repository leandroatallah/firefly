# Engine Module

This module contains the core, reusable game engine components. It is designed to be game-agnostic and provides the fundamental building blocks for creating a 2D game.

## Core Components

- `app/`: Manages the main engine loop, context, and initialization (`engine.go`, `context.go`).
- `contracts/`: Defines the Go interfaces (contracts) for key engine components like animations, bodies, configuration, context, navigation, sequences, tilemap layers, and visual effects (vfx). This promotes a decoupled architecture.
- `data/`: Handles data loading, management, and configuration schemas (e.g., from JSON files).
  - `config/`: Engine-specific configuration structures.
  - `datamanager/`: Centralized manager for loading and accessing game data.
  - `i18n/`: Internationalization system for loading and managing translations.
    - `i18n.go`: `I18nManager` that loads translations from `assets/lang/{langCode}.json` and provides `T(key, args...)` for translated strings with `fmt.Sprintf`-style formatting.
  - `jsonutil/`: Helpers for JSON parsing and schema validation.
  - `schemas/`: Definitions for data structures used in asset files.
- `event/`: Provides a basic event handling system for inter-component communication.
- `input/`: Manages user input from keyboard, mouse, or gamepads.
- `mocks/`: Contains mock implementations of engine components for testing purposes, facilitating unit and integration tests for the game module.
- `sequences/`: Manages scripted event sequences, commands, and cutscenes.
  - `player.go`: Executes sequences of commands.
  - `commands_*.go`: Scriptable actions for actors, camera, music, and visual effects.
- `utils/`: Contains various utility functions (e.g., fixed-point arithmetic `fp16/`, timing `timing/`, and `delay_trigger.go`).

## Game Object Management

- `entity/`: Provides the foundational structures for all in-game objects.
  - `actors/`: Base structures and logic for character-like entities.
  - `items/`: Base structures and logic for collectible or interactive items.
  - `animation_utils.go`: Helper functions for animation logic.
- `physics/`: Implements the physics simulation.
  - `body/`: Defines physical body interfaces and implementations.
  - `movement/`: Provides movement models (e.g., platformer physics).
  - `skill/`: Manages physics-related skills or abilities.
  - `space/`: Handles collision detection and spatial partitioning.
- `scene/`: Manages game scenes, scene transitions, and the overall scene lifecycle.
  - `scene_manager.go`: Orchestrates scene loading, updating, and drawing.
  - `scene_base.go`: Provides a common base for all scenes.
  - `scene_factory.go`: Responsible for creating new scene instances.
  - `transition/`: Handles scene transitions (e.g., fades).
  - `pause/`: Implements pause menu functionality.
  - `phases/`: Manages different states or phases within a single scene.
  - `camera_config.go`: Defines camera behavior for scenes.
  - `screen_flipper.go`: Manages screen flipping effects.
  - `scene_tilemap.go`: Handles tilemap-based scene elements.

## Presentation

- `assets/`: Handles the loading and management of game assets.
  - `imagemanager/`: Manages loading and caching of images.
  - `font/`: Handles font loading and text rendering.
- `audio/`: Provides the core audio playback functionality.
  - `loader.go`: Facilitates the loading of audio files into the engine.
- `render/`: Responsible for all rendering tasks.
  - `camera/`: Controls the game camera's position and zoom.
  - `particles/`: Manages particle systems and emitters.
    - `vfx/`: Particle-based visual effects.
  - `sprites/`: Handles sprite rendering, layering, and animations.
  - `tilemap/`: Renders tilemaps and handles tile-based collisions.
  - `vfx/`: Provides visual effects and screen-wide overlays.
    - `text/`: Text-based visual effects (e.g., damage numbers, popups).
  - `screenutil/`: Utility functions for screen coordinates, rendering, and screen-wide effects like flashes.
- `ui/`: Provides building blocks for user interface elements.
  - `hud/`: Base components for Heads-Up Displays.
  - `menu/`: Components for creating interactive menus.
  - `speech/`: Components for speech bubbles and dialogue systems.
