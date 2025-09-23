package physics

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

type Obstacle interface {
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

// Body methods
func (o *ObstacleRect) Position() (minX, minY, maxX, maxY int) {
	return o.PhysicsBody.Position()
}

func (o *ObstacleRect) DrawCollisionBox(screen *ebiten.Image) {
	o.PhysicsBody.DrawCollisionBox(screen)
}

func (o *ObstacleRect) CollisionPosition() []image.Rectangle {
	return o.PhysicsBody.CollisionPosition()
}
func (o *ObstacleRect) IsColliding(boundaries []Body) bool {
	return o.PhysicsBody.IsColliding(boundaries)
}

func (o *ObstacleRect) Draw(screen *ebiten.Image) {
	rect := o.Shape.(*Rect)
	vector.DrawFilledRect(
		screen,
		float32(rect.x16)/config.Unit,
		float32(rect.y16)/config.Unit,
		float32(rect.width),
		float32(rect.height),
		color.RGBA{0xff, 0, 0, 0xff},
		false,
	)
}

type ObstacleCircle struct {
	PhysicsBody
}
