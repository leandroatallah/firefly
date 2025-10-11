package physics

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type CollidableBody struct {
	body.Shape
	body.Touchable

	isObstructive bool
	collisionList []*CollisionArea
	invulnerable  bool
}

func (b *CollidableBody) SetTouchable(t body.Touchable) {
	b.Touchable = t
}

func (b *CollidableBody) GetTouchable() body.Touchable {
	return b.Touchable
}

func (b *CollidableBody) DrawCollisionBox(screen *ebiten.Image) {
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

func (b *CollidableBody) CollisionPosition() []image.Rectangle {
	res := []image.Rectangle{}
	for _, c := range b.collisionList {
		res = append(res, c.Position())
	}
	return res
}

func (b *CollidableBody) SetIsObstructive(value bool) {
	b.isObstructive = value
}

func (b *CollidableBody) IsObstructive() bool {
	return b.isObstructive
}
