# Game Module

This module contains the game-specific implementation built on top of `internal/engine`. Modify this module to implement your own game.

## Structure

- `app/`: Game setup and initialization.
  - `config.go`: Game configuration constants.
  - `phases_list.go`: Defines the ordered list of game phases.
  - `setup.go`: Wires all engine systems together and starts the game.
  - `setup_audio.go`: Collects speech bleep audio files for the dialogue system.
- `entity/`: Concrete game entities.
  - `actors/`: Player, NPCs, and enemies.
    - `enemies/`: Enemy implementations (`bat`, `wolf`).
    - `events/`: Actor-specific event handling.
    - `methods/`: Reusable actor behaviors (e.g., death).
    - `npcs/`: NPC implementations. The `ClimberPlayer` is reused as an NPC type.
    - `player/`: Player character (`ClimberPlayer`).
    - `states/`: Custom actor state machine states (`Dying`, `Dead`, `Exiting`), plus the composite grounded sub-state machine:
      - `GroundedState`: Composite state owning `Idle`, `Walking`, `Ducking`, and `AimLock` sub-states. Plugs into the parent state machine as a single `StateGrounded` value.
      - `DuckingState`: Entered when duck input is held while grounded. Shrinks the hitbox via `ResizeFixedBottom`, zeroes horizontal velocity, and blocks jumping. Exits when duck is released and there is ceiling clearance.
      - `DashState`: Tween-based dash. Velocity follows an `InOutSine` curve from `DashSpeed` to `0`. Gravity is suspended and hitbox is duck-height for the full duration. One air dash per jump; cooldown prevents immediate re-trigger.
  - `items/`: Collectible and interactive items.
    - `fall_platform.go`: A platform that falls when touched.
    - `item_power_base.go`: Base struct for power-up items. See `POWERUPS.md` for how to add new ones.
  - `obstacles/`: Static collision obstacles.
  - `types/`: Shared interfaces for game entities (`EnemyActor`, `PlayerActor`).
  - `OffsetToggler`: Alternates bullet spawn Y-offset between `+N` and `-N` on each `Next()` call, producing a double-barrel visual effect for rapid fire.
- `render/`: Game-specific rendering.
  - `camera/`: Camera controller with screen-shake support.
  - `vfx/`: Game-specific VFX helpers (aura particles, overhead/screen text).
- `scenes/`: Game scenes.
  - `scene_menu.go`: Main menu with "Start" and "Exit" options.
  - `phases/`: Platformer phase scene — loads tilemaps, manages actors, items, and goals.
  - `types/`: Scene type constants (`SceneMenu`, `ScenePhases`).
- `ui/`: Game UI.
  - `hud/`: Heads-Up Display (health, etc.).
  - `speech/`: Speech bubble and story dialogue implementations.
