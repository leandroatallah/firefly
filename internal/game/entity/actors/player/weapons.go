package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
)

// NewClimberInventory creates a new inventory with light_blaster and heavy_cannon.
func NewClimberInventory(projectileManager combat.ProjectileManager, vfxManager vfx.Manager) *inventory.Inventory {
	inv := inventory.New()
	// Speeds in fp16 (scale factor 16): 6 pixels/frame = 96
	light := weapon.NewProjectileWeapon("light_blaster", 8, "bullet_small", 96, projectileManager, "muzzle_flash")
	light.SetVFXManager(vfxManager)
	inv.AddWeapon(light)
	inv.AddWeapon(NewHeavyCannonWeapon(projectileManager, vfxManager))
	return inv
}

// NewHeavyCannonWeapon creates the heavy_cannon weapon.
func NewHeavyCannonWeapon(projectileManager combat.ProjectileManager, vfxManager vfx.Manager) combat.Weapon {
	// Speed in fp16: 9 pixels/frame = 144
	w := weapon.NewProjectileWeapon("heavy_cannon", 30, "bullet_large", 144, projectileManager, "muzzle_flash")
	w.SetVFXManager(vfxManager)
	return w
}
