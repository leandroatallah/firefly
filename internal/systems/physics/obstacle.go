package physics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

// TODO: Should it be merge with Collidable?
type Obstacle interface {
	Body
	Draw(screen *ebiten.Image)
	DrawCollisionBox(screen *ebiten.Image)
}

type ObstacleRect struct {
	PhysicsBody
}

func NewObstacleRect(rect *Rect) *ObstacleRect {
	return &ObstacleRect{PhysicsBody: *NewPhysicsBody(rect)}
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
		float32(rect.x16)/config.Unit,
		float32(rect.y16)/config.Unit,
		float32(rect.width),
		float32(rect.height),
		color.RGBA{0x77, 0x55, 0x55, 0xff},
		false,
	)
}
