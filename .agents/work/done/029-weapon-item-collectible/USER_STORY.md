# User Story — 029 Weapon Item Collectible

**As a** player,  
**I want** to collect a heavy cannon weapon from an item placed on the tilemap,  
**So that** I can discover and use a more powerful weapon during gameplay.

## Acceptance Criteria

- AC1: An item with `item_type: "ITEM_WEAPON_CANNON"` exists on the game layer.
- AC2: The item uses the same sprite as `item-power-grow.png`.
- AC3: When the player touches the item, it adds `heavy_cannon` to the player's inventory with 10 ammo.
- AC4: The item is removed from the scene after collection.
- AC5: The player starts with only `light_blaster` (infinite ammo); `heavy_cannon` must be collected.
- AC6: The item follows the same base behavior as `PowerUpItem` (collision detection, removal on touch).
- AC7: If the player collects the same weapon item while already owning `heavy_cannon`, it adds 10 ammo instead of adding a duplicate weapon.
- AC8: When `heavy_cannon` ammo reaches 0, the weapon is removed from the inventory.

## Technical Notes

- Item config: `internal/game/entity/items/item-weapon-cannon.json`
- Item implementation: `internal/game/entity/items/item_weapon_cannon.go`
- Register item type in `internal/game/entity/items/init_items.go`
- Use callback pattern to access player inventory (similar to `PowerUpItem.activateSkill`)
- Weapon factory already exists: `gameplayer.NewHeavyCannonWeapon(manager)`
- Inventory interface supports: `AddWeapon()`, `SetAmmo()`, `HasAmmo()`
