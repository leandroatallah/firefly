package object

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

type Object interface {
	Position() (minX, minY, maxX, maxY int)
	DrawCollisionBox(screen *ebiten.Image)
	CollisionPosition() []image.Rectangle
}

// TODO: Rename
type BaseObject struct {
	Element
	vx16 int
	vy16 int
	// TODO: Convert to a list
	collisionList []*CollisionArea
}

func NewBaseObject(element Element, collisionList []*CollisionArea) BaseObject {
	return BaseObject{Element: element, collisionList: collisionList}
}

func (b *BaseObject) Move() {
	panic("You should implement this method in derivated structs")
}

// TODO: Implement ease in movement
func (b *BaseObject) MoveY(distance int) {
	b.vy16 += distance * config.Unit
}

// TODO: Implement ease in movement
func (b *BaseObject) MoveX(distance int) {
	b.vx16 += distance * config.Unit
}

func (b *BaseObject) Position() (minX, minY, maxX, maxY int) {
	minX = b.x16 / config.Unit
	minY = b.y16 / config.Unit
	maxX = minX + b.width
	maxY = minY + b.height
	return
}

func (b *BaseObject) DrawCollisionBox(screen *ebiten.Image) {
	for _, c := range b.CollisionPosition() {
		minX := c.Min.X
		minY := c.Min.Y
		maxX := c.Max.X
		maxY := c.Max.Y

		width := float32(maxX - minX)
		height := float32(maxY - minY)
		vector.DrawFilledRect(
			screen,
			float32(minX), float32(minY), width, height,
			color.RGBA{0, 0xaa, 0, 0xff}, false)
		vector.DrawFilledRect(
			screen,
			float32(minX)+1, float32(minY)+1, width-2, height-2,
			color.RGBA{0, 0xff, 0, 0xff}, false)
	}
}

func (b *BaseObject) CollisionPosition() []image.Rectangle {
	res := []image.Rectangle{}
	for _, c := range b.collisionList {
		minX, minY, maxX, maxY := c.Position()
		rect := image.Rect(minX, minY, maxX, maxY)
		res = append(res, rect)
	}
	return res
}
