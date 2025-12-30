package body

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
)

type CollidableBody struct {
	*Body

	Touchable     body.Touchable
	collisionList []body.Collidable
	isObstructive bool
}

func NewCollidableBody(body *Body) *CollidableBody {
	if body == nil {
		panic("NewCollidableBody: body must not be nil")
	}
	return &CollidableBody{Body: body}
}

func NewCollidableBodyFromRect(rect body.Shape) *CollidableBody {
	b := NewBody(rect)
	return NewCollidableBody(b)
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

func (b *CollidableBody) OnTouch(other body.Collidable) {
	if b.Touchable != nil {
		b.Touchable.OnTouch(other)
	}
}
func (b *CollidableBody) OnBlock(other body.Collidable) {
	if b.Touchable != nil {
		b.Touchable.OnBlock(other)
	}
}
func (b *CollidableBody) AddCollision(list ...body.Collidable) {
	if b.ID() == "" {
		debug.PrintStack()
		log.Fatal("(*CollidableBody).AddCollision: A body must have an ID")
	}

	if b.GetShape() == nil {
		log.Fatal("(*CollidableBody).AddCollision: A body must have a shape")
	}

	for d, i := range list {
		i.SetID(fmt.Sprintf("%v_COLLISION_%d", b.ID(), d))
		b.collisionList = append(b.collisionList, i)
	}
}
func (b *CollidableBody) CollisionShapes() []body.Collidable {
	return b.collisionList
}

func (b *CollidableBody) ClearCollisions() {
	b.collisionList = nil
}

// SetPosition overrides Body.SetPosition method to updates the body position and its collisions
func (b *CollidableBody) SetPosition(x, y int) {
	cfg := config.Get()

	// Calculate the difference to move the collision areas as well
	diffX16 := cfg.To16(x) - b.Body.x16
	diffY16 := cfg.To16(y) - b.Body.y16

	b.Body.SetPosition(x, y)

	for _, c := range b.collisionList {
		x, y := c.GetPositionMin()
		x16, y16 := cfg.To16(x), cfg.To16(y)
		c.SetPosition(cfg.From16(x16+diffX16), cfg.From16(y16+diffY16))
	}
}

// ApplyValidPosition moves the body by a given distance, ensuring it stops at the first collision.
// It works by moving one pixel at a time to avoid truncation issues with sub-pixel velocities
// and guarantees the body will be perfectly flush with any obstacle it collides with.
// Might update: Body x16 and y16
func (b *CollidableBody) ApplyValidPosition(distance16 int, isXAxis bool, space body.BodiesSpace) (int, int, bool) {
	if distance16 == 0 || space == nil {
		x, y := b.GetPositionMin()
		return x, y, false
	}

	cfg := config.Get()
	var isBlocking bool

	// Determine the direction of movement (step is one pixel).
	step := 1
	if distance16 < 0 {
		step = -1
	}

	// Calculate how many pixels we need to move.
	// We use ceiling division to ensure that any velocity, no matter how small,
	// results in at least a 1-pixel check.
	pixelDistance := (abs(distance16) + cfg.Unit - 1) / cfg.Unit

	// Move pixel by pixel and check for collisions.
	for i := 0; i < pixelDistance; i++ {
		lastX, lastY := b.GetPositionMin()

		// Move one pixel.
		if isXAxis {
			b.SetPosition(lastX+step, lastY)
		} else {
			b.SetPosition(lastX, lastY+step)
		}

		// Check if this new position causes a collision.
		_, blocking := space.ResolveCollisions(b)
		if blocking {
			// A collision occurred. We should only undo the movement if the collision
			// is on the same axis we are currently moving on.
			// It does. The last position was the correct one.
			b.SetPosition(lastX, lastY)
			isBlocking = true
			break
		}
	}

	x, y := b.GetPositionMin()
	return x, y, isBlocking
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
