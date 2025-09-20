package object

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

type Obstacle interface {
	Draw(screen *ebiten.Image)
}

type ObstacleRect struct {
	BaseObject
}

func NewObstacleRect(element Element, collisionList []*CollisionArea) *ObstacleRect {
	if len(collisionList) == 0 {
		collisionList = []*CollisionArea{elementToCollisionArea(element)}
	}
	return &ObstacleRect{
		BaseObject: NewBaseObject(
			element,
			collisionList,
		),
	}
}

// Object methods
func (o *ObstacleRect) Position() (minX, minY, maxX, maxY int) {
	return o.BaseObject.Position()
}

func (o *ObstacleRect) DrawCollisionBox(screen *ebiten.Image) {
	o.BaseObject.DrawCollisionBox(screen)
}

func (o *ObstacleRect) CollisionPosition() []image.Rectangle {
	return o.BaseObject.CollisionPosition()
}

func (c *ObstacleRect) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(c.x16)/config.Unit,
		float32(c.y16)/config.Unit,
		float32(c.width),
		float32(c.height),
		color.RGBA{0xff, 0, 0, 0xff},
		false,
	)
}

type ObstacleCircle struct {
	BaseObject
}
