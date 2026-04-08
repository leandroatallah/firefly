# Combat Inventory

The `inventory` package manages a collection of weapons and tracks their ammunition counts.

## Key Features

- **Weapon Collection**: Store multiple `combat.Weapon` implementations.
- **Active Weapon Switching**: Supports cycling through weapons (`SwitchNext`, `SwitchPrev`) or jumping to a specific index (`SwitchTo`).
- **Ammo Tracking**: Maps weapon IDs to ammo counts.
- **Unlimited Ammo**: Uses a count of `-1` to represent infinite ammunition.

## Implementation Details

The `Inventory` struct implements the `combat.Inventory` interface. It keeps track of an `activeIndex` and a map of ammo counts indexed by weapon ID.

### Usage Example

```go
import (
    "github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
    "github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
)

// Create a new inventory
inv := inventory.New()

// Add a weapon (defaults to unlimited ammo)
inv.AddWeapon(myWeapon)

// Set limited ammo
inv.SetAmmo(myWeapon.ID(), 100)

// Check if we can fire
if inv.HasAmmo(inv.ActiveWeapon().ID()) {
    // Fire weapon
    inv.ActiveWeapon().Fire(...)
    inv.ConsumeAmmo(inv.ActiveWeapon().ID(), 1)
}
```

## Methods

- `ActiveWeapon()`: Returns the currently selected weapon.
- `SwitchNext() / SwitchPrev()`: Cycles through the available weapons.
- `Update()`: Calls `Update()` on all weapons in the inventory (useful for cooldowns).
- `HasAmmo(id)`: Checks if the weapon has ammo (>= 1 or -1).
- `ConsumeAmmo(id, amount)`: Decrements ammo unless it is unlimited.
