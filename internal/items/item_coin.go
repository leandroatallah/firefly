package items

import (
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

// Concrete
type CollectibleCoinItem struct {
	BaseItem
}

func NewCollectibleCoinItem() *CollectibleCoinItem {
	x, y := 220, -140
	frameWidth, frameHeight := 16, 16

	base := NewBaseItem()
	rect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := rect
	base.SetBody(rect)
	base.SetCollisionArea(collisionRect)
	base.SetTouchable(base)

	return &CollectibleCoinItem{BaseItem: *base}
}

func (c *CollectibleCoinItem) OnTouch(other physics.Body) {
	if c.IsRemoved() {
		return
	}

	c.SetRemoved(true)

	// TODO: Handle player reward
}
