package items

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type BaseItem struct {
	sprites.SpriteEntity
	*physics.CollidableBody
	*physics.MovableBody

	count        int
	frameRate    int
	removed      bool
	imageOptions *ebiten.DrawImageOptions
	appContext   *core.AppContext
}

func NewBaseItem(s sprites.SpriteMap, frameRate int, bodyRect *physics.Rect) *BaseItem {
	spriteEntity := sprites.NewSpriteEntity(s)
	b := physics.NewBody(bodyRect)
	movable := physics.NewMovableBody(b)
	collidable := physics.NewCollidableBody(b)
	return &BaseItem{
		MovableBody:    movable,
		CollidableBody: collidable,

		imageOptions: &ebiten.DrawImageOptions{},
		SpriteEntity: spriteEntity,
		frameRate:    frameRate,
	}
}

// Forwarding methods for Body to avoid ambiguous selector
// Always route via the MovableBody component
func (b *BaseItem) ID() string {
	return b.MovableBody.ID()
}
func (b *BaseItem) SetID(id string) {
	b.MovableBody.SetID(id)
}
func (b *BaseItem) Position() image.Rectangle {
	return b.MovableBody.Position()
}
func (b *BaseItem) SetPosition(x, y int) {
	b.MovableBody.SetPosition(x, y)
}
func (b *BaseItem) GetPositionMin() (int, int) {
	return b.MovableBody.GetPositionMin()
}
func (b *BaseItem) GetShape() body.Shape {
	return b.MovableBody.GetShape()
}

func (b *BaseItem) SetCollisionArea(rect *physics.Rect) {
	collision := physics.NewCollidableBodyFromRect(rect)
	x, y := b.GetPositionMin()
	collision.SetID(fmt.Sprintf("%v_COLLISION_0", b.ID()))
	collision.SetPosition(x, y)
	b.AddCollision(collision)
}

func (b *BaseItem) SetTouchable(t body.Touchable) {
	b.Touchable = t
}

func (b *BaseItem) Update(space body.BodiesSpace) error {
	b.count++

	return nil
}

func (b *BaseItem) UpdateImageOptions() {
	b.imageOptions.GeoM.Reset()

	x, y := b.GetPositionMin()
	b.imageOptions.GeoM.Translate(float64(x), float64(y))
}

func (b *BaseItem) OnBlock(other body.Collidable) {}

func (b *BaseItem) OnTouch(other body.Collidable) {}

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

func (b *BaseItem) SetAppContext(appContext *core.AppContext) {
	b.appContext = appContext
}

func (b *BaseItem) AppContext() *core.AppContext {
	return b.appContext
}
