package gameitems

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/jsonutil"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
)

// PowerUpItem is a base struct for collectible power-up items.
// It provides common functionality for all power-up items,
// with skill activation delegated to a callback.
type PowerUpItem struct {
	items.BaseItem
	activateSkill func()
	onCollect     func()
}

// NewPowerUpItem creates a new power-up item with the given sprite config and skill activation callback.
func NewPowerUpItem(ctx *app.AppContext, x, y int, id string, spriteConfigPath string, activateSkill func()) (*PowerUpItem, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[items.StatData](ctx.Assets, spriteConfigPath)
	if err != nil {
		return nil, err
	}

	base, err := CreateAnimatedItem(ctx.Assets, id, spriteData, nil)
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

// SetOnCollect sets the callback to be called when the item is collected.
func (p *PowerUpItem) SetOnCollect(fn func()) {
	p.onCollect = fn
}

// createPowerUpBase creates a base item for power-ups with the given sprite config and callback.
// This is a helper function for power-up items that need to reference themselves in the callback.
func createPowerUpBase(ctx *app.AppContext, x, y int, id string, spriteConfigPath string, activateSkill func()) (*items.BaseItem, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[items.StatData](ctx.Assets, spriteConfigPath)
	if err != nil {
		return nil, err
	}

	base, err := CreateAnimatedItem(ctx.Assets, id, spriteData, nil)
	if err != nil {
		return nil, err
	}

	// SetPosition must be before SetItemBodies
	base.SetPosition(x, y)
	base.SetAppContext(ctx)
	base.SetOwner(base)

	if err = SetItemBodies(base, spriteData, nil); err != nil {
		return nil, fmt.Errorf("SetItemBodies: %w", err)
	}
	if err = SetItemStats(base, statData); err != nil {
		return nil, fmt.Errorf("SetItemStats: %w", err)
	}

	base.StateCollisionManager.RefreshCollisions()

	return base, nil
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
	
	// Call collection callback for feedback (sound, vfx)
	if p.onCollect != nil {
		p.onCollect()
	}
}
