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

func (b *CollidableBody) DrawCollisionBox(screen *ebiten.Image, position image.Rectangle) {
	for _, collisionRect := range b.CollisionPosition() {
		// Calculate top-left corner of collision box relative to the character's body origin.
		offsetX := float32(collisionRect.Min.X - position.Min.X)
		offsetY := float32(collisionRect.Min.Y - position.Min.Y)

		width := float32(collisionRect.Dx())
		height := float32(collisionRect.Dy())

		// Draw on the 'screen' (which is the sprite) at the relative offset.
		vector.DrawFilledRect(
			screen,
			offsetX, offsetY, width, height,
			color.RGBA{0, 0xaa, 0, 0xff}, false)
		vector.DrawFilledRect(
			screen,
			offsetX+1, offsetY+1, width-2, height-2,
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
