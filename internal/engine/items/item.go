package items

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

type Item interface {
	// TODO: Move al these repeated Character method to a new struct/interface
	body.MovableCollidable
	body.Drawable

	Update(space body.BodiesSpace) error
	SetCollisionArea(rect *physics.Rect)
	Image() *ebiten.Image
	ImageCollisionBox() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions

	IsRemoved() bool
	SetRemoved(value bool)
}
