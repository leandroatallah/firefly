package body

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type ObstacleRect struct {
	*MovableBody
	*CollidableBody

	imageOptions *ebiten.DrawImageOptions
}

func NewObstacleRect(bodyRect *Rect) *ObstacleRect {
	b := NewBody(bodyRect)
	movable := NewMovableBody(b)
	collidable := NewCollidableBody(b)
	return &ObstacleRect{
		MovableBody:    movable,
		CollidableBody: collidable,

		imageOptions: &ebiten.DrawImageOptions{},
	}
}

// Forwarding methods for Body to avoid ambiguous selector
// Always route via the MovableBody component
func (o *ObstacleRect) ID() string {
	return o.MovableBody.ID()
}
func (o *ObstacleRect) SetID(id string) {
	o.MovableBody.SetID(id)
}
func (o *ObstacleRect) Position() image.Rectangle {
	return o.MovableBody.Position()
}
func (o *ObstacleRect) SetPosition(x, y int) {
	o.MovableBody.SetPosition(x, y)
}
func (o *ObstacleRect) GetPositionMin() (int, int) {
	return o.MovableBody.GetPositionMin()
}
func (o *ObstacleRect) GetShape() body.Shape {
	return o.MovableBody.GetShape()
}

func (o *ObstacleRect) AddCollisionBodies(list ...body.Collidable) {
	if len(list) == 0 {
		b := NewCollidableBodyFromRect(o.GetShape())
		x, y := o.GetPositionMin()
		b.SetPosition(x, y)
		b.SetID(fmt.Sprintf("%v_COLLISION_0", o.ID()))
		list = []body.Collidable{b}
	}
	o.AddCollision(list...)
}

func (o *ObstacleRect) Draw(screen *ebiten.Image) {
	rect := o.GetShape().(*Rect)
	x, y := o.GetPositionMin()
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(rect.width),
		float32(rect.height),
		color.Transparent,
		false,
	)
}

func (o *ObstacleRect) Image() *ebiten.Image {
	w := o.Position().Dx()
	h := o.Position().Dy()
	i := ebiten.NewImage(w, h)
	return i
}

func (o *ObstacleRect) ImageCollisionBox() *ebiten.Image {
	img := o.Image()
	if o.IsObstructive() {
		img.Fill(color.RGBA{G: 255})
	} else {
		img.Fill(color.RGBA{R: 255})
	}
	return img
}

func (o *ObstacleRect) ImageOptions() *ebiten.DrawImageOptions {
	return o.imageOptions
}

func (o *ObstacleRect) UpdateImageOptions() {
	o.imageOptions.GeoM.Reset()
	x, y := o.GetPositionMin()
	o.imageOptions.GeoM.Translate(float64(x), float64(y))
}
