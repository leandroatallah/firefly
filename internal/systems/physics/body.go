package physics

import (
	"github.com/google/uuid"
)

// Body is a Shape with collision, movable and alive
type Body interface {
	Shape
	Movable
	Collidable
	Alive

	ID() string
	SetMovementModel(model MovementModel)
	MovementModel() MovementModel
}

type FacingDirectionEnum int

const (
	FaceDirectionRight FacingDirectionEnum = iota
	FaceDirectionLeft
)

type PhysicsBody struct {
	Shape

	MovableBody
	CollidableBody
	AliveBody

	id            string
	invulnerable  bool
	movementModel MovementModel
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
	if b.movementModel != nil {
		b.movementModel.Update(b, space)
	}
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

// Movement Model methods
func (b *PhysicsBody) SetMovementModel(model MovementModel) {
	b.movementModel = model
}

func (b *PhysicsBody) MovementModel() MovementModel {
	return b.movementModel
}

// Platform methods
func (b *PhysicsBody) TryJump(force int) {
	b.MovableBody.TryJump(force)
}
