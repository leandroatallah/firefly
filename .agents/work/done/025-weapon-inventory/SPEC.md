# SPEC-025 — Weapon Inventory System

**Branch:** `025-weapon-inventory`  
**Bounded Context:** Engine (Combat)  
**Package:** `internal/engine/combat/inventory/`

## Overview

Implements a weapon inventory system that holds multiple weapons, tracks per-weapon ammo, and supports switching between weapons. Enables Megaman-style weapon selection and Metal Slug weapon pickups.

## Technical Requirements

### New Contract

**File:** `internal/engine/contracts/combat/inventory.go`

```go
package combat

// Inventory manages a collection of weapons with ammo tracking.
type Inventory interface {
	AddWeapon(weapon Weapon)
	ActiveWeapon() Weapon
	SwitchNext()
	SwitchPrev()
	SwitchTo(index int)
	HasAmmo(weaponID string) bool
	ConsumeAmmo(weaponID string, amount int)
	SetAmmo(weaponID string, amount int)
}
```

### Implementation

**File:** `internal/engine/combat/inventory/inventory.go`

**Struct:**
```go
type Inventory struct {
	weapons      []combat.Weapon
	activeIndex  int
	ammo         map[string]int // key: weapon.ID(), value: ammo count (-1 = unlimited)
}
```

**Methods:**
- `New() *Inventory` — constructor, initializes empty weapons slice, ammo map, activeIndex = 0
- `AddWeapon(weapon combat.Weapon)` — appends weapon, initializes ammo to -1 (unlimited) if not set
- `ActiveWeapon() combat.Weapon` — returns `weapons[activeIndex]` or `nil` if empty
- `SwitchNext()` — increments activeIndex with wrap-around: `(activeIndex + 1) % len(weapons)`
- `SwitchPrev()` — decrements activeIndex with wrap-around: `(activeIndex - 1 + len(weapons)) % len(weapons)`
- `SwitchTo(index int)` — sets activeIndex if `0 <= index < len(weapons)`, else no-op
- `HasAmmo(weaponID string) bool` — returns `true` if ammo[weaponID] == -1 or > 0
- `ConsumeAmmo(weaponID string, amount int)` — decrements ammo[weaponID] by amount if not -1
- `SetAmmo(weaponID string, amount int)` — sets ammo[weaponID] = amount

## Pre-conditions

- `combat.Weapon` interface exists (US-022)
- Weapons have unique `ID()` strings

## Post-conditions

- Inventory can hold multiple weapons
- Active weapon can be retrieved and switched
- Ammo is tracked per weapon ID
- Empty inventory returns `nil` without panic
- Switching wraps around at boundaries

## Integration Points

- **Weapon Contract:** `internal/engine/contracts/combat/weapon.go` — consumed by inventory
- **Future:** US-026 will wire input handling for weapon switching in shooting skill

## Red Phase Scenario

**Test File:** `internal/engine/combat/inventory/inventory_test.go`

**Failing Test Cases:**

1. **Empty inventory guard:**
   - Given: new empty inventory
   - When: `ActiveWeapon()` called
   - Then: returns `nil` (no panic)

2. **Add and retrieve weapon:**
   - Given: inventory with one weapon added
   - When: `ActiveWeapon()` called
   - Then: returns the added weapon

3. **Switch next with wrap-around:**
   - Given: inventory with 3 weapons, activeIndex = 2
   - When: `SwitchNext()` called
   - Then: activeIndex = 0 (wraps to first)

4. **Switch prev with wrap-around:**
   - Given: inventory with 3 weapons, activeIndex = 0
   - When: `SwitchPrev()` called
   - Then: activeIndex = 2 (wraps to last)

5. **Unlimited ammo:**
   - Given: weapon with ammo = -1
   - When: `HasAmmo()` called
   - Then: returns `true`
   - When: `ConsumeAmmo(10)` called
   - Then: ammo remains -1

6. **Limited ammo consumption:**
   - Given: weapon with ammo = 5
   - When: `ConsumeAmmo(2)` called
   - Then: ammo = 3
   - When: `HasAmmo()` called
   - Then: returns `true`
   - When: `ConsumeAmmo(3)` called
   - Then: ammo = 0
   - When: `HasAmmo()` called
   - Then: returns `false`

7. **SwitchTo bounds check:**
   - Given: inventory with 3 weapons
   - When: `SwitchTo(1)` called
   - Then: activeIndex = 1
   - When: `SwitchTo(5)` called (out of bounds)
   - Then: activeIndex remains 1 (no-op)

## Design Decisions

- **Ammo map keyed by weapon ID:** Allows flexible ammo tracking independent of weapon order in slice
- **-1 for unlimited ammo:** Simple sentinel value, avoids separate boolean flag
- **Wrap-around switching:** Improves UX for cycling through weapons
- **No-op for invalid SwitchTo:** Prevents index out of bounds, fails silently (caller should validate)
- **Default unlimited ammo on AddWeapon:** Simplifies common case where weapons don't need ammo limits
