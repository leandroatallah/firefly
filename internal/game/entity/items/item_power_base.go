package gameitems

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/jsonutil"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
)

// PowerUpItem is a base struct for collectible power-up items.
// It provides common functionality for all power-up items,
// with skill activation delegated to a callback.
type PowerUpItem struct {
	items.BaseItem
	activateSkill func()
}

// NewPowerUpItem creates a new power-up item with the given sprite config and skill activation callback.
func NewPowerUpItem(ctx *app.AppContext, x, y int, id string, spriteConfigPath string, activateSkill func()) (*PowerUpItem, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[items.StatData](spriteConfigPath)
	if err != nil {
		return nil, err
	}

	base, err := CreateAnimatedItem(id, spriteData, nil)
	if err != nil {
		return nil, err
	}

	powerItem := &PowerUpItem{
		BaseItem:      *base,
		activateSkill: activateSkill,
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

// OnTouch handles collision with the player and activates the skill.
func (p *PowerUpItem) OnTouch(other body.Collidable) {
	if p.IsRemoved() {
		return
	}

	player, found := p.AppContext().ActorManager.GetPlayer()
	if !found {
		return
	}

	if other.ID() != player.ID() {
		return
	}

	// Activate the skill via callback
	p.SetRemoved(true)
	if p.activateSkill != nil {
		p.activateSkill()
	}
}
