package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// NewClimberInventory creates a new inventory with two weapons for the climber player.
func NewClimberInventory(manager combat.ProjectileManager) *inventory.Inventory {
	inv := inventory.New()
	// Speeds in fp16 (scale factor 16): 6 pixels/frame = 96, 9 pixels/frame = 144
	inv.AddWeapon(weapon.NewProjectileWeapon("light_blaster", 8, "bullet_small", 96, manager))
	inv.AddWeapon(weapon.NewProjectileWeapon("heavy_cannon", 30, "bullet_large", 144, manager))
	return inv
}
