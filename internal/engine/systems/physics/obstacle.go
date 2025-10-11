package physics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/game/constants"
)

type ObstacleRect struct {
	PhysicsBody

	// TODO: Rename this
	op *ebiten.DrawImageOptions
}

func NewObstacleRect(rect *Rect) *ObstacleRect {
	return &ObstacleRect{
		PhysicsBody: *NewPhysicsBody(rect),
		op:          &ebiten.DrawImageOptions{},
	}
}

func (o *ObstacleRect) AddCollision(list ...*CollisionArea) *ObstacleRect {
	if len(list) == 0 {
		list = []*CollisionArea{rectToCollisionArea(o.Shape)}
	}
	o.PhysicsBody.AddCollision(list...)
	return o
}

func (o *ObstacleRect) Draw(screen *ebiten.Image) {
	rect := o.Shape.(*Rect)
	vector.DrawFilledRect(
		screen,
		float32(rect.x16)/constants.Unit,
		float32(rect.y16)/constants.Unit,
		float32(rect.width),
		float32(rect.height),
		color.Transparent,
		false,
	)
}

func (o *ObstacleRect) Image() *ebiten.Image {
	w := o.Position().Dx()
	h := o.Position().Dy()
	return ebiten.NewImage(w, h)
}

func (o *ObstacleRect) ImageCollisionBox() *ebiten.Image {
	img := o.Image()
	if o.IsObstructive() {
		img.Fill(color.RGBA{G: 255, A: 255})
	} else {
		img.Fill(color.RGBA{R: 255, A: 255})
	}
	return img
}

func (o *ObstacleRect) ImageOptions() *ebiten.DrawImageOptions {
	return o.op
}
