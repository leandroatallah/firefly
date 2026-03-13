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
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type CollidableBody struct {
	*Body

	Ownership
	Touchable     body.Touchable
	collisionList []body.Collidable
	isObstructive bool

	// Accumulators for sub-pixel movement
	accumulatorX16 int
	accumulatorY16 int
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
	b.SetPosition16(fp16.To16(x), fp16.To16(y))
}

func (b *CollidableBody) SetPosition16(x16, y16 int) {
	// Calculate the difference to move the collision areas as well
	diffX16 := x16 - b.Body.x16
	diffY16 := y16 - b.Body.y16

	b.Body.SetPosition16(x16, y16)

	for _, c := range b.collisionList {
		cx16, cy16 := c.GetPosition16()
		c.SetPosition16(cx16+diffX16, cy16+diffY16)
	}
}

// ApplyValidPosition moves the body by a given distance, ensuring it stops at the first collision.
// It accumulates sub-pixel movements and only moves when a full pixel step is reached.
// Might update: Body x16 and y16
func (b *CollidableBody) ApplyValidPosition(distance16 int, isXAxis bool, space body.BodiesSpace) (int, int, bool) {
	if distance16 == 0 || space == nil {
		x, y := b.GetPositionMin()
		return x, y, false
	}

	var totalDistance16 int
	if isXAxis {
		b.accumulatorX16 += distance16
		totalDistance16 = b.accumulatorX16
	} else {
		b.accumulatorY16 += distance16
		totalDistance16 = b.accumulatorY16
	}

	// Calculate whole pixels to move
	pixelSteps := fp16.From16(totalDistance16)

	if pixelSteps == 0 {
		// Not enough accumulated distance to move a full pixel yet
		x, y := b.GetPositionMin()
		return x, y, false
	}

	var isBlocking bool

	// Determine the direction of movement (step is one pixel).
	step := 1
	if pixelSteps < 0 {
		step = -1
		pixelSteps = -pixelSteps
	}

	// Move pixel by pixel and check for collisions.
	pixelsMoved := 0
	for i := 0; i < pixelSteps; i++ {
		lastX16, lastY16 := b.GetPosition16()

		// Move one pixel (in 16.16 fixed point)
		step16 := fp16.To16(step)
		if isXAxis {
			b.SetPosition16(lastX16+step16, lastY16)
		} else {
			b.SetPosition16(lastX16, lastY16+step16)
		}

		// Check if this new position causes a collision.
		_, blocking := space.ResolveCollisions(b)
		if blocking {
			// A collision occurred. Revert to the last valid position.
			b.SetPosition16(lastX16, lastY16)
			isBlocking = true
			break
		}
		pixelsMoved += step
	}

	// Update the accumulator by subtracting the distance actually moved
	movedDistance16 := fp16.To16(pixelsMoved)
	if isXAxis {
		b.accumulatorX16 -= movedDistance16
		if isBlocking {
			b.accumulatorX16 = 0 // Clear momentum on collision
		}
	} else {
		b.accumulatorY16 -= movedDistance16
		if isBlocking {
			b.accumulatorY16 = 0 // Clear momentum on collision
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
