package items

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type Item interface {
	// TODO: Move al these repeated Character method to a new struct/interface
	body.MovableCollidable
	body.Drawable

	Update(space body.BodiesSpace) error
	Image() *ebiten.Image
	ImageCollisionBox() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions

	IsRemoved() bool
	SetRemoved(value bool)
}
