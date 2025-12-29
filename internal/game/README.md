# Game Module

This module contains the specific implementation and logic for the _Firefly_ game. It is built upon the reusable components and contracts provided by the `internal/engine` module.

## Game-Specific Logic

- `app/`: Contains the game-specific setup and initialization code, configuring the engine to run _Firefly_.
- `entity/`: Defines the concrete game entities.
  - `actors/`: Implements the `Player` and specific `Enemies` for the game.
  - `items/`: Implements collectible items like `Coin`.
  - `obstacles/`: Defines game-specific obstacles.
- `scenes/`: Implements the actual game scenes, such as the `IntroScene`, `MenuScene`, and gameplay levels. It orchestrates the actors, items, and UI for each part of the game.

## Customization and Implementation

- `render/`: Contains game-specific rendering logic, such as custom camera behaviors tailored for _Firefly_.
- `ui/`: Implements the game's specific user interface, including the main `HUD` and `SpeechBubble` styles, built upon the engine's UI components.
