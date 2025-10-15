package items

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type BaseItem struct {
	// TODO: Item is not alive. Try to split PhysicsBody to not inherit AliveBody methods
	physics.PhysicsBody
	sprites.SpriteEntity

	count        int
	frameRate    int
	removed      bool
	imageOptions *ebiten.DrawImageOptions
}

func NewBaseItem(s sprites.SpriteMap, frameRate int) *BaseItem {
	spriteEntity := sprites.NewSpriteEntity(s)
	return &BaseItem{
		imageOptions: &ebiten.DrawImageOptions{},
		SpriteEntity: spriteEntity,
		frameRate:    frameRate,
	}
}

func (b *BaseItem) SetBody(rect *physics.Rect) Item {
	b.PhysicsBody = *physics.NewPhysicsBody(rect)
	b.PhysicsBody.SetTouchable(b)
	return b
}

func (b *BaseItem) SetCollisionArea(rect *physics.Rect) Item {
	collisionArea := &physics.CollisionArea{Shape: rect}
	b.PhysicsBody.AddCollision(collisionArea)
	return b
}

func (b *BaseItem) SetTouchable(t body.Touchable) {
	b.PhysicsBody.Touchable = t
}

func (b *BaseItem) Update(space body.BodiesSpace) error {
	b.count++

	return nil
}

func (b *BaseItem) UpdateImageOptions() {
	b.imageOptions.GeoM.Reset()

	pos := b.Position()
	minX, minY := pos.Min.X, pos.Min.Y

	// Apply character position
	b.imageOptions.GeoM.Translate(
		float64(minX),
		float64(minY),
	)
}

func (b *BaseItem) OnBlock(other body.Body) {}

func (b *BaseItem) OnTouch(other body.Body) {}

func (b *BaseItem) Image() *ebiten.Image {
	pos := b.Position()
	img := b.GetFirstSprite()
	img = b.AnimatedSpriteImage(img, pos, b.count, b.frameRate)
	return img
}

func (b *BaseItem) ImageCollisionBox() *ebiten.Image {
	img := b.Image()
	if b.IsObstructive() {
		img.Fill(color.RGBA{G: 255, A: 255})
	} else {
		img.Fill(color.RGBA{R: 255, A: 255})
	}
	return img
}

func (b *BaseItem) ImageOptions() *ebiten.DrawImageOptions {
	return b.imageOptions
}

func (b *BaseItem) IsRemoved() bool {
	return b.removed
}

func (b *BaseItem) SetRemoved(value bool) {
	b.removed = value
}
