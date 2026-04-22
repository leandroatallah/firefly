# Engine Module

This module contains the core, reusable game engine components. It is designed to be game-agnostic and provides the fundamental building blocks for creating a 2D game.

## Core Components

- `app/`: Manages the main engine loop, context, and initialization (`engine.go`, `context.go`).
- `combat/`: Handles weapon management, inventory, projectile lifecycles, and faction-gated damage.
  - `inventory/`: Manages weapon collections and ammo.
  - `weapon/`: Logic for firing and cooldowns; includes `EnemyShooting` for automatic enemy fire and `ProjectileWeapon` with per-state spawn offsets and muzzle-flash VFX.
  - `projectile/`: High-performance projectile management with lifetime, damage, and impact/despawn VFX hooks.
  - `faction.go`: `FactionNeutral | FactionPlayer | FactionEnemy` — side identification used to prevent self-damage.
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
  - `HorizontalAxis`: Last-pressed-wins directional input — when both left and right are held, the most recently pressed direction wins.
- `mocks/`: Contains mock implementations of engine components for testing purposes, facilitating unit and integration tests for the game module.
- `sequences/`: Manages scripted event sequences, commands, and cutscenes. See [`sequences/README.md`](sequences/README.md).
  - `player.go`: Executes sequences of commands.
  - `commands_*.go`: Scriptable actions for actors, camera, music, and visual effects.
- `utils/`: Contains various utility functions (e.g., fixed-point arithmetic `fp16/`, timing `timing/`, and `delay_trigger.go`).

## Game Object Management

- `entity/`: Provides the foundational structures for all in-game objects.
  - `actors/`: Base structures and logic for character-like entities.
    - `StateContributor`: Optional hook polled by `Character.handleState` before default movement transitions. Lets adapters (e.g., dash, shooting) override the target state without subclassing `Character`. See [ADR-008](../../docs/adr/ADR-008-state-contributor-pattern.md).
  - `items/`: Base structures and logic for collectible or interactive items.
  - `animation_utils.go`: Helper functions for animation logic.
- `physics/`: Implements the physics simulation.
  - `body/`: Defines physical body interfaces and implementations.
  - `movement/`: Provides movement models (e.g., platformer physics). Includes one-way platform drop-through logic.
  - `skill/`: Manages physics-related skills or abilities. See [`physics/skill/README.md`](physics/skill/README.md).
    - `JumpSkill` (variable jump height), `DashSkill` (tween-based deceleration), `HorizontalMovementSkill`, `ShootingSkill`.
    - `factory.go` (`FromConfig`) — builds the skill set from a JSON `SkillsConfig`.
  - `space/`: Handles collision detection and spatial partitioning.
  - `tween/`: Interpolation utilities.
    - `InOutSineTween`: Smooth `InOutSine` tween used by the dash deceleration.
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
  - `freeze.go`: `FreezeController` — pauses all Actor and Body updates for a given number of frames (hit-stop effect). Exposed via the `Freezable` contract in `contracts/scene/`.

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

## Architecture Decision Records

Key non-obvious design choices are documented in [`docs/adr/`](../../docs/adr/):

- [ADR-001](../../docs/adr/ADR-001-fp16-fixed-point-arithmetic.md) — Why positions use FP16 fixed-point instead of `float64`
- [ADR-002](../../docs/adr/ADR-002-registry-based-state-pattern.md) — Why actor states use a global registry with `init()` registration
- [ADR-003](../../docs/adr/ADR-003-goroutine-audio-looping.md) — Why audio looping uses goroutines instead of Ebitengine's built-in loop
- [ADR-004](../../docs/adr/ADR-004-space-body-model-physics.md) — Why physics is split into Space / Body / MovementModel layers
- [ADR-005](../../docs/adr/ADR-005-composite-grounded-sub-state.md) — Why the grounded state uses a sub-state machine instead of flat states
- [ADR-006](../../docs/adr/ADR-006-engine-game-layer-separation.md) — Engine/Game two-layer architecture
- [ADR-007](../../docs/adr/ADR-007-fp16-scale-factor.md) — FP16 scale factor is 16, not 65536
- [ADR-008](../../docs/adr/ADR-008-state-contributor-pattern.md) — StateContributor hook for extensible state transitions
