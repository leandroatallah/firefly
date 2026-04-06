# User Story — 028 Dual Weapon Setup

**As a** player,  
**I want** the climber character to start with two different weapons in their inventory,  
**So that** I can switch between a fast light blaster and a slow heavy cannon during gameplay.

## Acceptance Criteria

- The player starts with two weapons already loaded in their inventory.
- Weapon 1 (`light_blaster`): small bullet (`bullet_small`), fast rate of fire (cooldown: 8 frames).
- Weapon 2 (`heavy_cannon`): large bullet (`bullet_large`), slow rate of fire (cooldown: 30 frames), higher projectile speed.
- The player can switch between weapons using the existing `WeaponNext` / `WeaponPrev` input commands.
- Both weapons fire projectiles through the existing `ProjectileManager`.
- Weapon definitions live in the game-specific layer (`internal/game/`), not in the engine.
- No new engine contracts or packages are needed — this is a game-layer wiring change only.
