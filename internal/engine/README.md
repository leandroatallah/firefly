# Engine Module

This module contains the core, reusable game engine components for the Firefly project. It is designed to be game-agnostic and provides the fundamental building blocks for creating a 2D game.

## Core Components

- `app/`: Manages the main engine loop, context, and initialization.
- `contracts/`: Defines the Go interfaces (contracts) for key engine components like bodies, animations, and scenes. This promotes a decoupled architecture.
- `data/`: Handles data loading, management, and configuration schemas (e.g., from JSON files).
- `sequences/`: Manages scripted event sequences, commands, and cutscenes.

## Game Object Management

- `entity/`: Provides the foundational structures for all in-game objects, primarily `actors` (like characters) and `items`.
- `physics/`: Implements the physics simulation, including movement models (platformer, top-down), collision detection, and physical body representations.
- `scene/`: Manages game scenes, scene transitions, and the overall scene lifecycle. It includes a phase manager to handle different states within a single scene.

## Presentation

- `assets/`: Handles the loading and management of game assets, including images, fonts, and audio files.
- `audio/`: Provides the core audio playback functionality.
- `input/`: Manages user input from keyboard, mouse, or gamepads.
- `render/`: Responsible for all rendering tasks, including the game camera, sprites, and tilemaps.
- `ui/`: Provides building blocks for user interface elements like HUDs and dialogue systems.
