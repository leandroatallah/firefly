# Game Module

This module contains the specific implementation and logic for the _Growbel_ game. It is built upon the reusable components and contracts provided by the `internal/engine` module.

## Game-Specific Logic

- `app/`: Contains the game-specific setup and initialization code, configuring the engine to run _Growbel_.
  - `config.go`: Game configuration settings.
  - `phases_list.go`: Defines the list and order of game phases.
  - `setup.go`: Handles game initialization and setup routines.
- `entity/`: Defines the concrete game entities.
  - `actors/`: Implements the `Player`, `NPCs`, and specific `Enemies`. It includes state machines and behavior logic for these characters.
  - `items/`: Implements collectible items like `Coin` and interactive environmental items like `FallingPlatform`.
  - `obstacles/`: Defines game-specific obstacles like walls and movement-restricting boundaries.
  - `types/`: Custom types and interfaces related to game entities.
- `scenes/`: Implements the actual game scenes, such as the `IntroScene`, `MenuScene`, and gameplay phases. It orchestrates actors, items, and UI.
  - `init_scenes.go`: Initializes the registry of all game scenes.
  - `scene_intro.go`: The introductory cinematic/tutorial scene.
  - `scene_menu.go`: The main menu and options scene.
  - `scene_phase_title.go`: Interstitial scene for displaying phase names.
  - `scene_phase_reboot.go`: Scene for handling phase restarts.
  - `scene_story.go`: Scenes dedicated to cinematic storytelling.
  - `scene_summary.go`: Scene for displaying results and scores.
  - `sounds.go`: Management of game-specific audio triggers.
  - `phases/`: Phase-specific logic and entity layouts for gameplay levels.
  - `types/`: Custom types for game scenes and state management.

## Customization and Implementation

- `render/`: Contains game-specific rendering logic.
  - `vfx/`: Game-specific visual effects, such as overhead and screen text.
- `ui/`: Implements the game's specific user interface.
  - `hud/`: Game's main Heads-Up Display elements.
  - `speech/`: Game-specific speech bubbles and dialogue styles.
