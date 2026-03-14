package gameitems

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/jsonutil"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// Concrete
type CollectiblePowerItem struct {
	items.BaseItem
}

func NewCollectiblePowerItem(ctx *app.AppContext, x, y int, id string) (*CollectiblePowerItem, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[items.StatData]("internal/game/entity/items/power.json")
	if err != nil {
		return nil, err
	}

	base, err := CreateAnimatedItem(id, spriteData, nil)
	if err != nil {
		return nil, err
	}

	powerItem := &CollectiblePowerItem{
		BaseItem: *base,
	}

	// SetPosition must be before SetItemBodies
	powerItem.SetPosition(x, y)
	powerItem.SetAppContext(ctx)
	powerItem.SetOwner(powerItem)

	if err = SetItemBodies(powerItem, spriteData, nil); err != nil {
		return nil, fmt.Errorf("SetItemBodies: %w", err)
	}
	if err = SetItemStats(powerItem, statData); err != nil {
		return nil, fmt.Errorf("SetItemStats: %w", err)
	}

	powerItem.StateCollisionManager.RefreshCollisions()

	return powerItem, nil
}

func (c *CollectiblePowerItem) OnTouch(other body.Collidable) {
	if c.IsRemoved() {
		return
	}

	player, found := c.AppContext().ActorManager.GetPlayer()
	if !found {
		return
	}

	if other.ID() != player.ID() {
		return
	}

	skillUser, ok := player.(gameentitytypes.FreezeSkillUser)
	if ok {
		c.SetRemoved(true)
		skillUser.ActivateFreezeSkill()
	}
}
