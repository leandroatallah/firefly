package physics

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

type Body interface {
	Position() (minX, minY, maxX, maxY int)
	DrawCollisionBox(screen *ebiten.Image)
	CollisionPosition() []image.Rectangle
}

type PhysicsBody struct {
	Rect
	vx16 int
	vy16 int
	// TODO: Convert to a list
	collisionList []*CollisionArea
}

func NewPhysicsBody(element Rect, collisionList []*CollisionArea) PhysicsBody {
	return PhysicsBody{Rect: element, collisionList: collisionList}
}

func (b *PhysicsBody) Move() {
	panic("You should implement this method in derivated structs")
}

// TODO: Implement ease in movement
func (b *PhysicsBody) MoveY(distance int) {
	b.vy16 += distance * config.Unit
}

// TODO: Implement ease in movement
func (b *PhysicsBody) MoveX(distance int) {
	b.vx16 += distance * config.Unit
}

func (b *PhysicsBody) Position() (minX, minY, maxX, maxY int) {
	minX = b.x16 / config.Unit
	minY = b.y16 / config.Unit
	maxX = minX + b.width
	maxY = minY + b.height
	return
}

func (b *PhysicsBody) DrawCollisionBox(screen *ebiten.Image) {
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

func (b *PhysicsBody) CollisionPosition() []image.Rectangle {
	res := []image.Rectangle{}
	for _, c := range b.collisionList {
		res = append(res, c.Position())
	}
	return res
}
