package items

import (
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

// Concrete
type CollectibleCoinItem struct {
	BaseItem
}

// TODO: Replace param with a rect
func NewCollectibleCoinItem(x, y int) *CollectibleCoinItem {
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
