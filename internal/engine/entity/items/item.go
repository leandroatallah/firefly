package items

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type Item interface {
	body.MovableCollidable
	body.Drawable

	Update(space body.BodiesSpace) error
	IsRemoved() bool
	SetRemoved(value bool)
}
