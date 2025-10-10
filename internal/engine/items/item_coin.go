package items

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

// Concrete
type CollectibleCoinItem struct {
	BaseItem
}

func NewCollectibleCoinItem(x, y int) (*CollectibleCoinItem, error) {
	frameWidth, frameHeight := 16, 16

	var assets sprites.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/collectible-coin.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		return nil, err
	}

	base := NewBaseItem(sprites)
	// TODO: It should be set in a better place
	base.frameRate = 10
	rect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := rect
	base.SetBody(rect)
	base.SetCollisionArea(collisionRect)
	base.SetTouchable(base)

	return &CollectibleCoinItem{BaseItem: *base}, nil
}

func (c *CollectibleCoinItem) OnTouch(other physics.Body) {
	if c.IsRemoved() {
		return
	}

	if p, ok := other.GetTouchable().(*actors.PlayerPlatform); ok {
		c.SetRemoved(true)
		p.AddCoinCount(1)
	}
}
