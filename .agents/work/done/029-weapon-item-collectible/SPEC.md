# Technical Specification — 029 Weapon Item Collectible

**Branch**: `029-weapon-item-collectible`

## Overview

Implement a collectible weapon item (`ITEM_WEAPON_CANNON`) that grants the player the `heavy_cannon` weapon with 10 ammo. The item reuses the `PowerUpItem` pattern with a callback to access player inventory.

## Bounded Context

- **Entity** (`internal/game/entity/items/`)
- **Combat** (`internal/engine/combat/inventory/`)

## Technical Requirements

### 1. Item Configuration

**File**: `internal/game/entity/items/item-weapon-cannon.json`

```json
{
  "sprites": {
    "body_rect": {
      "x": 0,
      "y": 0,
      "width": 16,
      "height": 16
    },
    "assets": {
      "idle": {
        "path": "assets/images/item-power-grow.png",
        "collision_rect": [
          {
            "x": 0,
            "y": 0,
            "width": 16,
            "height": 16
          }
        ]
      }
    },
    "facing_direction": 0,
    "frame_rate": 10
  },
  "stats": {}
}
```

### 2. Item Implementation

**File**: `internal/game/entity/items/item_weapon_cannon.go`

- Struct: `WeaponCannonItem` embedding `PowerUpItem`
- Constructor: `NewWeaponCannonItem(ctx *app.AppContext, x, y int, id string) (*WeaponCannonItem, error)`
- Logic:
  - Use `NewPowerUpItem` with `activateSkill` callback
  - Callback accesses player via `ctx.ActorManager.GetPlayer()`
  - Check if player has `heavy_cannon` via inventory lookup
  - If weapon exists: add 10 ammo via `SetAmmo(weaponID, currentAmmo + 10)`
  - If weapon missing: create via `NewHeavyCannonWeapon(ctx.ProjectileManager)`, call `AddWeapon()`, set ammo to 10

### 3. Inventory Extension

**File**: `internal/engine/combat/inventory/inventory.go`

Add methods:

```go
// GetAmmo returns the current ammo count for a weapon (-1 = unlimited, 0 = none).
func (i *Inventory) GetAmmo(weaponID string) int

// HasWeapon returns true if the weapon exists in the inventory.
func (i *Inventory) HasWeapon(weaponID string) bool

// RemoveWeapon removes a weapon from the inventory if ammo reaches 0.
func (i *Inventory) RemoveWeapon(weaponID string)
```

### 4. Item Registration

**File**: `internal/game/entity/items/init_items.go`

- Add constant: `WeaponCannonType items.ItemType = "ITEM_WEAPON_CANNON"`
- Register factory in `InitItemMap`:

```go
WeaponCannonType: func(x, y int, id string) items.Item {
    return itemFactoryOrFatal(NewWeaponCannonItem(ctx, x, y, id))
},
```

## Integration Points

### Existing Contracts

- `items.Item` (from `internal/engine/entity/items/`)
- `combat.Weapon` (from `internal/engine/contracts/combat/`)
- `combat.ProjectileManager` (from `internal/engine/contracts/combat/`)

### Dependencies

- `PowerUpItem.activateSkill` callback pattern
- `gameplayer.NewHeavyCannonWeapon(manager)` factory
- `inventory.AddWeapon()`, `inventory.SetAmmo()`, `inventory.HasAmmo()`

## Pre-Conditions

- Player starts with only `light_blaster` (infinite ammo)
- `heavy_cannon` weapon factory exists
- `PowerUpItem` base class handles collision and removal

## Post-Conditions

- Item with type `ITEM_WEAPON_CANNON` can be placed on tilemap
- Player collects item → gains `heavy_cannon` with 10 ammo
- Re-collecting item → adds 10 ammo (no duplicate weapon)
- When ammo reaches 0 → weapon removed from inventory
- Item removed from scene after collection

## Red Phase (Failing Test Scenario)

**Test File**: `internal/game/entity/items/item_weapon_cannon_test.go`

### Test 1: Collect weapon when not owned
```
Given: Player has only light_blaster
When: Player touches ITEM_WEAPON_CANNON
Then: 
  - Player inventory contains heavy_cannon
  - heavy_cannon ammo = 10
  - Item is removed
```

### Test 2: Collect weapon when already owned
```
Given: Player has heavy_cannon with 5 ammo
When: Player touches ITEM_WEAPON_CANNON
Then:
  - heavy_cannon ammo = 15
  - Inventory weapon count unchanged
  - Item is removed
```

### Test 3: Weapon removal at zero ammo
```
Given: Player has heavy_cannon with 1 ammo
When: Player fires once (ammo → 0)
Then: heavy_cannon removed from inventory
```

## Design Decisions

- **Reuse `PowerUpItem`**: Avoids duplicating collision/removal logic
- **Callback pattern**: Decouples item from player implementation
- **Ammo stacking**: Prevents inventory clutter from duplicate weapons
- **Auto-removal at 0 ammo**: Requires `ConsumeAmmo` to trigger cleanup (handled by weapon firing logic)
