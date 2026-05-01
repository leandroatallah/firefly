package inventory

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// Inventory manages a collection of weapons with ammo tracking.
type Inventory struct {
	weapons     []combat.Weapon
	activeIndex int
	ammo        map[string]int // key: weapon.ID(), value: ammo count (-1 = unlimited)
}

// New creates a new empty inventory.
func New() *Inventory {
	return &Inventory{
		weapons: []combat.Weapon{},
		ammo:    make(map[string]int),
	}
}

// AddWeapon appends a weapon to the inventory and initializes ammo to -1 (unlimited).
func (i *Inventory) AddWeapon(weapon combat.Weapon) {
	i.weapons = append(i.weapons, weapon)
	i.ammo[weapon.ID()] = -1
}

// ActiveWeapon returns the currently active weapon, or nil if inventory is empty.
func (i *Inventory) ActiveWeapon() combat.Weapon {
	if len(i.weapons) == 0 {
		return nil
	}
	return i.weapons[i.activeIndex]
}

// SwitchNext increments the active weapon index with wrap-around.
func (i *Inventory) SwitchNext() {
	if len(i.weapons) > 0 {
		i.activeIndex = (i.activeIndex + 1) % len(i.weapons)
	}
}

// SwitchPrev decrements the active weapon index with wrap-around.
func (i *Inventory) SwitchPrev() {
	if len(i.weapons) > 0 {
		i.activeIndex = (i.activeIndex - 1 + len(i.weapons)) % len(i.weapons)
	}
}

// SwitchTo sets the active weapon index if valid, otherwise no-op.
func (i *Inventory) SwitchTo(index int) {
	if index >= 0 && index < len(i.weapons) {
		i.activeIndex = index
	}
}

// HasAmmo returns true if the weapon has ammo (-1 = unlimited, or > 0).
func (i *Inventory) HasAmmo(weaponID string) bool {
	ammo, exists := i.ammo[weaponID]
	return exists && (ammo == -1 || ammo > 0)
}

// ConsumeAmmo decrements ammo for a weapon if not unlimited (-1).
func (i *Inventory) ConsumeAmmo(weaponID string, amount int) {
	if ammo, exists := i.ammo[weaponID]; exists && ammo != -1 {
		i.ammo[weaponID] = ammo - amount
	}
}

// SetAmmo sets the ammo count for a weapon.
func (i *Inventory) SetAmmo(weaponID string, amount int) {
	i.ammo[weaponID] = amount
}

// GetAmmo returns the current ammo count for a weapon (-1 = unlimited, 0 = none).
func (i *Inventory) GetAmmo(weaponID string) int {
	if ammo, exists := i.ammo[weaponID]; exists {
		return ammo
	}
	return 0
}

// HasWeapon returns true if the weapon exists in the inventory.
func (i *Inventory) HasWeapon(weaponID string) bool {
	for _, w := range i.weapons {
		if w.ID() == weaponID {
			return true
		}
	}
	return false
}

// RemoveWeapon removes a weapon from the inventory if ammo reaches 0.
func (i *Inventory) RemoveWeapon(weaponID string) {
	for idx, w := range i.weapons {
		if w.ID() == weaponID {
			i.weapons = append(i.weapons[:idx], i.weapons[idx+1:]...)
			delete(i.ammo, weaponID)
			if i.activeIndex >= len(i.weapons) && len(i.weapons) > 0 {
				i.activeIndex = len(i.weapons) - 1
			}
			break
		}
	}
}

// Update updates all weapons in the inventory.
func (i *Inventory) Update() {
	for _, weapon := range i.weapons {
		weapon.Update()
	}
}
