package physics

import (
	"math"

	"github.com/google/uuid"
	"github.com/leandroatallah/firefly/internal/config"
)

// Body is a Shape with collision, movable and alive
type Body interface {
	// TODO: Check this Body, Movable and Collidable usage
	Shape
	Movable
	Collidable
	Alive

	ID() string
}

type FacingDirectionEnum int

const (
	FaceDirectionLeft FacingDirectionEnum = iota
	FaceDirectionRight
)

type PhysicsBody struct {
	Shape

	MovableBody
	CollidableBody
	AliveBody

	id           string
	invulnerable bool
}

func NewPhysicsBody(shape Shape) *PhysicsBody {
	return &PhysicsBody{
		Shape: shape,
		id:    uuid.New().String(),
	}
}

// Attribute methods
func (b *PhysicsBody) ID() string {
	return b.id
}
func (b *PhysicsBody) Invulnerable() bool {
	return b.invulnerable
}
func (b *PhysicsBody) SetInvulnerable(value bool) {
	b.invulnerable = value
}

// Collision methods
func (b *PhysicsBody) OnTouch(other Body) {
	if b.Touchable != nil {
		b.Touchable.OnTouch(other)
	}
}

func (b *PhysicsBody) OnBlock(other Body) {
	if b.Touchable != nil {
		b.Touchable.OnBlock(other)
	}
}

func (b *PhysicsBody) AddCollision(list ...*CollisionArea) *PhysicsBody {
	for _, i := range list {
		b.collisionList = append(b.collisionList, i)
	}
	return b
}

// Movement methods
func (b *PhysicsBody) ApplyValidMovement(distance int, isXAxis bool, space *Space) {
	if distance == 0 {
		return
	}

	b.updatePosition(distance, isXAxis)

	if space == nil {
		return
	}

	_, isBlocking := space.ResolveCollisions(b)
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

func (b *PhysicsBody) UpdateMovement(space *Space) {
	// Apply physics to player's position based on the velocity from previous frame.
	// This is a simple Euler integration step: position += velocity * deltaTime (where deltaTime=1 frame).
	b.ApplyValidMovement(b.vx16, true, space)
	b.ApplyValidMovement(b.vy16, false, space)

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
