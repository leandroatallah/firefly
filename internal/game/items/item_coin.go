package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

// Concrete
type CollectibleCoinItem struct {
	items.BaseItem
}

func NewCollectibleCoinItem(x, y int) (*CollectibleCoinItem, error) {
	frameWidth, frameHeight := 16, 16

	var assets sprites.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/collectible-coin.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		return nil, err
	}

	// TODO: It should be set in a better place (frameRate)
	frameRate := 10
	base := items.NewBaseItem(sprites, frameRate)
	rect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := rect
	base.SetBody(rect)
	base.SetCollisionArea(collisionRect)
	base.SetTouchable(base)

	return &CollectibleCoinItem{BaseItem: *base}, nil
}

func (c *CollectibleCoinItem) OnTouch(other body.Body) {
	if c.IsRemoved() {
		return
	}

	if p, ok := other.GetTouchable().(*actors.PlayerPlatform); ok {
		c.SetRemoved(true)
		p.AddCoinCount(1)
	}
}
