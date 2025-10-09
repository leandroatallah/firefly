package items

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type Item interface {
	// TODO: Move al these repeated Character method to a new struct/interface
	physics.Body
	Update(space *physics.Space) error
	SetBody(rect *physics.Rect) Item
	SetCollisionArea(rect *physics.Rect) Item
	SetTouchable(t physics.Touchable)
	OnTouch(other physics.Body)
	Image() *ebiten.Image
	ImageCollisionBox() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions

	IsRemoved() bool
	SetRemoved(value bool)
}

type BaseItem struct {
	// TODO: Item is not alive. Try to split PhysicsBody to not inherit AliveBody methods
	physics.PhysicsBody
	actors.SpriteEntity

	removed      bool
	imageOptions *ebiten.DrawImageOptions
}

func NewBaseItem() *BaseItem {
	return &BaseItem{}
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

func (b *BaseItem) SetTouchable(t physics.Touchable) {
	b.PhysicsBody.Touchable = t
}

func (b *BaseItem) Update(space *physics.Space) error {
	fmt.Println("UPDATE ITEM")
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

func (b *BaseItem) OnBlock(other physics.Body) {}

func (b *BaseItem) OnTouch(other physics.Body) {}

func (b *BaseItem) Image() *ebiten.Image {
	img := ebiten.NewImage(b.Position().Dx(), b.Position().Dy())
	img.Fill(color.RGBA{0xff, 0xff, 0, 0xff})
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
