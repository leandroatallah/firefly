package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// GrowPowerItem is a collectible power-up item that activates the grow skill.
type GrowPowerItem struct {
	PowerUpItem
}

// NewGrowPowerItem creates a new grow power-up item.
func NewGrowPowerItem(ctx *app.AppContext, x, y int, id string) (*GrowPowerItem, error) {
	powerItem, err := NewPowerUpItem(
		ctx, x, y, id,
		"internal/game/entity/items/item-power-grow.json",
		func() {
			// Activate grow skill when collected
			player, found := ctx.ActorManager.GetPlayer()
			if !found {
				return
			}
			if skillUser, ok := player.(gameentitytypes.GrowSkillUser); ok {
				skillUser.ActivateGrowSkill()
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return &GrowPowerItem{
		PowerUpItem: *powerItem,
	}, nil
}
