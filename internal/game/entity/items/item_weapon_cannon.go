package gameitems

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
)

// WeaponCannonItem is a collectible item that grants the heavy_cannon weapon.
type WeaponCannonItem struct {
	*PowerUpItem
}

// NewWeaponCannonItem creates a new weapon cannon item.
func NewWeaponCannonItem(ctx *app.AppContext, x, y int, id string) (*WeaponCannonItem, error) {
	activateSkill := func() {
		playerActor, found := ctx.ActorManager.GetPlayer()
		if !found {
			return
		}

		gamePlayer, ok := playerActor.(interface{ Inventory() interface{} })
		if !ok {
			return
		}

		inv, ok := gamePlayer.Inventory().(*inventory.Inventory)
		if !ok {
			return
		}

		if inv.HasWeapon("heavy_cannon") {
			currentAmmo := inv.GetAmmo("heavy_cannon")
			inv.SetAmmo("heavy_cannon", currentAmmo+10)
		} else {
			weapon := gameplayer.NewHeavyCannonWeapon(ctx.ProjectileManager, ctx.VFX)
			inv.AddWeapon(weapon)
			inv.SetAmmo("heavy_cannon", 10)
		}
	}

	powerItem, err := NewPowerUpItem(ctx, x, y, id, "assets/entities/items/item-weapon-cannon.json", activateSkill)
	if err != nil {
		return nil, fmt.Errorf("NewPowerUpItem: %w", err)
	}

	return &WeaponCannonItem{
		PowerUpItem: powerItem,
	}, nil
}
