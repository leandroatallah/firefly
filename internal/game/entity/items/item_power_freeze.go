package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// FreezePowerItem is a collectible power-up item that activates the freeze skill.
type FreezePowerItem struct {
	PowerUpItem
}

// NewFreezePowerItem creates a new freeze power-up item.
func NewFreezePowerItem(ctx *app.AppContext, x, y int, id string) (*FreezePowerItem, error) {
	powerItem, err := NewPowerUpItem(
		ctx, x, y, id,
		"internal/game/entity/items/power.json",
		func() {
			// Activate freeze skill when collected
			player, found := ctx.ActorManager.GetPlayer()
			if !found {
				return
			}
			if skillUser, ok := player.(gameentitytypes.FreezeSkillUser); ok {
				skillUser.ActivateFreezeSkill()
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return &FreezePowerItem{
		PowerUpItem: *powerItem,
	}, nil
}
