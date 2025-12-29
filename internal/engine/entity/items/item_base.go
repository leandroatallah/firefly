package items

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

type BaseItem struct {
	app.AppContextHolder
	sprites.SpriteEntity
	*bodyphysics.CollidableBody
	*bodyphysics.MovableBody
	*space.StateCollisionManager[ItemStateEnum]

	count        int
	removed      bool
	imageOptions *ebiten.DrawImageOptions
	state        ItemState
}

func NewBaseItem(id string, s sprites.SpriteMap, bodyRect *bodyphysics.Rect) *BaseItem {
	spriteEntity := sprites.NewSpriteEntity(s)
	b := bodyphysics.NewBody(bodyRect)
	movable := bodyphysics.NewMovableBody(b)
	collidable := bodyphysics.NewCollidableBody(b)

	base := &BaseItem{
		MovableBody:    movable,
		CollidableBody: collidable,
		imageOptions:   &ebiten.DrawImageOptions{},
		SpriteEntity:   spriteEntity,
	}
	base.SetID(id)
	base.StateCollisionManager = space.NewStateCollisionManager[ItemStateEnum](base)

	state, err := NewItemState(base, Idle)
	if err != nil {
		log.Fatal(err)
	}
	base.SetState(state)

	return base
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
	img := b.GetSpriteByState(b.state.State())
	if img == nil {
		// Try to fallback to idle sprite
		img = b.GetSpriteByState(Idle)
	}
	if img == nil {
		img = b.GetFirstSprite()
	}

	pos := b.Position()
	return b.AnimatedSpriteImage(img, pos, b.count, b.SpriteEntity.FrameRate())
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

func (b *BaseItem) State() ItemStateEnum {
	return b.state.State()
}

// SetState set a new Character state and update current collision shapes.
func (b *BaseItem) SetState(state ItemState) {
	b.state = state
	b.StateCollisionManager.RefreshCollisions()
	b.state.OnStart()
}
