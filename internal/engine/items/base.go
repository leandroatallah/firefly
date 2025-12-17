package items

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type BaseItem struct {
	core.AppContextHolder
	sprites.SpriteEntity
	*physics.CollidableBody
	*physics.MovableBody

	count           int
	frameRate       int
	removed         bool
	imageOptions    *ebiten.DrawImageOptions
	state           ItemState
	collisionBodies map[ItemStateEnum][]body.Collidable
}

func NewBaseItem(id string, s sprites.SpriteMap, bodyRect *physics.Rect) *BaseItem {
	spriteEntity := sprites.NewSpriteEntity(s)
	b := physics.NewBody(bodyRect)
	movable := physics.NewMovableBody(b)
	collidable := physics.NewCollidableBody(b)

	base := &BaseItem{
		MovableBody:     movable,
		CollidableBody:  collidable,
		imageOptions:    &ebiten.DrawImageOptions{},
		SpriteEntity:    spriteEntity,
		collisionBodies: make(map[ItemStateEnum][]body.Collidable), // Character collisions based on state
	}
	base.SetID(id)

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

func (b *BaseItem) State() ItemStateEnum {
	return b.state.State()
}

// SetState set a new Character state and update current collision shapes.
func (b *BaseItem) SetState(state ItemState) {
	b.state = state
	b.RefreshCollisionBasedOnState()
	b.state.OnStart()
}

func (b *BaseItem) RefreshCollisionBasedOnState() {
	// TODO: Duplicated
	if rects, ok := b.collisionBodies[b.state.State()]; ok {
		b.ClearCollisions()
		x, y := b.GetPositionMin()
		for _, r := range rects {
			// Create a deep copy of the collision body to avoid mutating the template
			template := r.(*physics.CollidableBody)
			newCollisionBody := physics.NewCollidableBody(
				physics.NewBody(template.GetShape()),
			)
			relativePos := template.Position()
			newPos := image.Rect(
				x+relativePos.Min.X,
				y+relativePos.Min.Y,
				x+relativePos.Max.X,
				y+relativePos.Max.Y,
			)
			newCollisionBody.SetPosition(newPos.Min.X, newPos.Min.Y)
			// FIX: It should not use Nanosecond
			newCollisionBody.SetID(fmt.Sprintf("%v_COLLISION_%v", b.ID(), time.Now().Nanosecond()))
			b.AddCollision(newCollisionBody)
		}
	}
}

// TODO: Duplicated
func (b *BaseItem) AddCollisionRect(state ItemStateEnum, rect body.Collidable) {
	b.collisionBodies[state] = append(b.collisionBodies[state], rect)
}

// TODO: Duplicated
func (b *BaseItem) SetFrameRate(value int) {
	b.frameRate = value
}
