package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// StarPowerItem is a collectible power-up item that activates the star skill.
type StarPowerItem struct {
	PowerUpItem
}

// NewStarPowerItem creates a new star power-up item.
func NewStarPowerItem(ctx *app.AppContext, x, y int, id string) (*StarPowerItem, error) {
	powerItem, err := NewPowerUpItem(
		ctx, x, y, id,
		"internal/game/entity/items/star.json",
		func() {
			// Activate star skill when collected
			player, found := ctx.ActorManager.GetPlayer()
			if !found {
				return
			}
			if skillUser, ok := player.(gameentitytypes.StarSkillUser); ok {
				skillUser.ActivateStarSkill()
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return &StarPowerItem{
		PowerUpItem: *powerItem,
	}, nil
}
