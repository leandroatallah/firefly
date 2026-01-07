package gameitems

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// Concrete
type CollectibleCoinItem struct {
	items.BaseItem
}

func NewCollectibleCoinItem(ctx *app.AppContext, x, y int, id string) (*CollectibleCoinItem, error) {
	spriteData, statData, err := items.ParseJsonItem("internal/game/entity/items/coin.json")
	if err != nil {
		return nil, err
	}

	base, err := CreateAnimatedItem(id, spriteData)
	if err != nil {
		return nil, err
	}

	coinItem := &CollectibleCoinItem{
		BaseItem: *base,
	}

	// SetPosition must be before SetItemBodies
	coinItem.SetPosition(x, y)
	coinItem.SetAppContext(ctx)

	if err = SetItemBodies(coinItem, spriteData); err != nil {
		return nil, fmt.Errorf("SetItemBodies: %w", err)
	}
	if err = SetItemStats(coinItem, statData); err != nil {
		return nil, fmt.Errorf("SetItemStats: %w", err)
	}

	coinItem.StateCollisionManager.RefreshCollisions()

	return coinItem, nil
}

func (c *CollectibleCoinItem) OnTouch(other body.Collidable) {
	if c.IsRemoved() {
		return
	}

	player, found := c.AppContext().ActorManager.GetPlayer()
	if !found {
		return
	}
	coinCollector, ok := player.(*gameentitytypes.PlatformerCharacter)
	if ok {
		c.SetRemoved(true)
		coinCollector.AddCoinCount(1)
	}
}
