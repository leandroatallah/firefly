package physics

import (
	"image"
	"image/color"
	"math"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

type Collidable interface {
	Shape
	DrawCollisionBox(screen *ebiten.Image)
	CollisionPosition() []image.Rectangle
	IsColliding(boundaries []Body) (isTouching, isBlocking bool)
	IsObstructive() bool
	SetIsObstructive(value bool)
}

// Movable is a Shape but with movement
type Movable interface {
	Shape
	Position() image.Rectangle
	ApplyValidMovement(velocity int, isXAxis bool, boundaries []Body)

	SetSpeedAndMaxSpeed(speed, maxSpeed int)
	Speed() int

	OnMoveUp(distance int)
	OnMoveDown(distance int)
	OnMoveLeft(distance int)
	OnMoveRight(distance int)
	OnMoveUpLeft(distance int)
	OnMoveUpRight(distance int)
	OnMoveDownLeft(distance int)
	OnMoveDownRight(distance int)
}

type Touchable interface {
	OnTouch(other Body)
	OnBlock(other Body)
}

// Body is a Shape but with collision
type Body interface {
	// TODO: Check this Body, Movable and Collidable usage
	Shape
	Movable
	Collidable
	// Touchable

	ID() string
}

type FacingDirectionEnum int

const (
	FaceDirectionLeft FacingDirectionEnum = iota
	FaceDirectionRight
)

// TODO: Split into Movable struct
type PhysicsBody struct {
	Shape
	Touchable     Touchable
	id            string
	vx16          int
	vy16          int
	accelerationX int
	accelerationY int
	speed         int
	maxSpeed      int
	faceDirection FacingDirectionEnum

	isObstructive bool
	collisionList []*CollisionArea
}

func NewPhysicsBody(shape Shape) *PhysicsBody {
	return &PhysicsBody{
		Shape: shape,
		id:    uuid.New().String(),
	}
}

func (b *PhysicsBody) Move() {
	panic("You should implement this method in derivated structs")
}

func (b *PhysicsBody) MoveX(distance int) {
	b.accelerationX = distance * config.Unit
}

func (b *PhysicsBody) MoveY(distance int) {
	b.accelerationY = distance * config.Unit
}

// TODO: Should it be moved to Movable?
func (b *PhysicsBody) OnMoveLeft(distance int) {
	b.MoveX(-distance)
}
func (b *PhysicsBody) OnMoveUpLeft(distance int) {
	b.MoveX(-distance)
	b.MoveY(-distance)
}
func (b *PhysicsBody) OnMoveDownLeft(distance int) {
	b.MoveX(-distance)
	b.MoveY(distance)
}
func (b *PhysicsBody) OnMoveRight(distance int) {
	b.MoveX(distance)
}
func (b *PhysicsBody) OnMoveUpRight(distance int) {
	b.MoveX(distance)
	b.MoveY(-distance)
}
func (b *PhysicsBody) OnMoveDownRight(distance int) {
	b.MoveX(distance)
	b.MoveY(distance)
}
func (b *PhysicsBody) OnMoveUp(distance int) {
	b.MoveY(-distance)
}
func (b *PhysicsBody) OnMoveDown(distance int) {
	b.MoveY(distance)
}

// TODO: Improve this method (split of find out a better approach)
func (b *PhysicsBody) SetSpeedAndMaxSpeed(speed, maxSpeed int) {
	b.speed = speed
	b.maxSpeed = maxSpeed
}

func (b *PhysicsBody) Speed() int {
	return b.speed
}

func (b *PhysicsBody) FaceDirection() FacingDirectionEnum {
	return b.faceDirection
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

func (b *PhysicsBody) AddCollision(list ...*CollisionArea) *PhysicsBody {
	for _, i := range list {
		b.collisionList = append(b.collisionList, i)
	}
	return b
}

func (b *PhysicsBody) CollisionPosition() []image.Rectangle {
	res := []image.Rectangle{}
	for _, c := range b.collisionList {
		res = append(res, c.Position())
	}
	return res
}

func (b *PhysicsBody) SetIsObstructive(value bool) {
	b.isObstructive = value
}

func (b *PhysicsBody) IsObstructive() bool {
	return b.isObstructive
}

// TODO: Needs to be updated when dealing with different shapes (e.g. circle)
func (b *PhysicsBody) checkRectIntersect(obj1, obj2 Body) bool {
	rects1 := obj1.CollisionPosition()
	rects2 := obj2.CollisionPosition()

	for _, r1 := range rects1 {
		for _, r2 := range rects2 {
			if r1.Overlaps(r2) {
				return true
			}
		}
	}

	return false
}

func (b *PhysicsBody) ID() string {
	return b.id
}

// TODO: Should return the collisions? If yes, collision need to become a struct
// TODO: Handle multiple simultaneos collision
func (b *PhysicsBody) IsColliding(boundaries []Body) (isTouching, isBlocking bool) {
	for _, o := range boundaries {
		if b.ID() == o.ID() {
			continue
		}

		if b.checkRectIntersect(b, o) {
			// A collision happened. Notify the touch handler if it's set.
			b.OnTouch(o)

			// Check if it's a blocking collision.
			if o.IsObstructive() {
				b.OnBlock(o)
				return true, true
			}

			return true, false
		}
	}
	return false, false
}

func (b *PhysicsBody) OnTouch(other Body) {}

func (b *PhysicsBody) OnBlock(other Body) {}

func (b *PhysicsBody) updatePosition(distance int, isXAxis bool) {
	// TODO: Replace switch with "polymorphism"
	switch b.Shape.(type) {
	case *Rect:
		rect := b.Shape.(*Rect)
		if isXAxis {
			rect.x16 += distance
			for _, c := range b.collisionList {
				c.Shape.(*Rect).x16 += distance
			}
		} else {
			rect.y16 += distance
			for _, c := range b.collisionList {
				c.Shape.(*Rect).y16 += distance
			}
		}
	}
}

func (b *PhysicsBody) ApplyValidMovement(distance int, isXAxis bool, boundaries []Body) {
	if distance == 0 {
		return
	}

	b.updatePosition(distance, isXAxis)

	_, isBlocking := b.IsColliding(boundaries)
	if isBlocking {
		b.updatePosition(-distance, isXAxis)
	}
}

func (b *PhysicsBody) CheckMovementDirectionX() {
	if b.accelerationX > 0 {
		b.faceDirection = FaceDirectionRight
	}
	if b.accelerationX < 0 {
		b.faceDirection = FaceDirectionLeft
	}
}

func (b *PhysicsBody) UpdateMovement(boundaries []Body) {
	// Apply physics to player's position based on the velocity from previous frame.
	// This is a simple Euler integration step: position += velocity * deltaTime (where deltaTime=1 frame).
	b.ApplyValidMovement(b.vx16, true, boundaries)
	b.ApplyValidMovement(b.vy16, false, boundaries)

	// Convert the raw input acceleration into a scaled and normalized vector.
	scaledAccX, scaledAccY := smoothDiagonalMovement(b.accelerationX, b.accelerationY)

	b.vx16 = increaseVelocity(b.vx16, scaledAccX)
	b.vy16 = increaseVelocity(b.vy16, scaledAccY)

	// Cap the magnitude of the velocity vector to enforce a maximum speed.
	// This is crucial for preventing faster movement on diagonals.
	// We need to check if the velocity magnitude `sqrt(vx² + vy²)` exceeds `speedMax16²`.
	// To avoid a costly square root, we can compare the squared values:
	speedMax16 := b.maxSpeed * config.Unit
	// Use int64 for squared values to prevent potential overflow.
	velSq := int64(b.vx16)*int64(b.vx16) + int64(b.vy16)*int64(b.vy16)
	maxSq := int64(speedMax16) * int64(speedMax16)

	if velSq > maxSq {
		// If the speed is too high, we need to scale the velocity vector down.
		// The scaling factor is `scale = speedMax16 / current_speed`.
		// `current_speed` is `sqrt(velSq)`.
		// So, `scale = speedMax16 / sqrt(velSq)`.
		scale := float64(speedMax16) / math.Sqrt(float64(velSq))
		b.vx16 = int(float64(b.vx16) * scale)
		b.vy16 = int(float64(b.vy16) * scale)
	}

	b.CheckMovementDirectionX()

	// Reset frame-specific acceleration.
	// It will be recalculated on the next frame from input.
	b.accelerationX, b.accelerationY = 0, 0

	// Apply friction to slow the player down when there is no input.
	b.vx16 = reduceVelocity(b.vx16)
	b.vy16 = reduceVelocity(b.vy16)

}

type BodyState int

const (
	Idle BodyState = iota
	Walk
)

func (b *PhysicsBody) CurrentBodyState() BodyState {
	isWalking := b.vx16 != 0 || b.vy16 != 0
	if isWalking {
		return Walk
	}
	return Idle
}
