package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// NewClimberInventory creates a new inventory with light_blaster and heavy_cannon.
func NewClimberInventory(manager combat.ProjectileManager) *inventory.Inventory {
	inv := inventory.New()
	// Speeds in fp16 (scale factor 16): 6 pixels/frame = 96
	inv.AddWeapon(weapon.NewProjectileWeapon("light_blaster", 8, "bullet_small", 96, manager, "muzzle_flash"))
	inv.AddWeapon(NewHeavyCannonWeapon(manager))
	return inv
}

// NewHeavyCannonWeapon creates the heavy_cannon weapon.
func NewHeavyCannonWeapon(manager combat.ProjectileManager) combat.Weapon {
	// Speed in fp16: 9 pixels/frame = 144
	return weapon.NewProjectileWeapon("heavy_cannon", 30, "bullet_large", 144, manager, "muzzle_flash")
}
