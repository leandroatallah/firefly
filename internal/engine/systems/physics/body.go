package physics

import (
	"github.com/google/uuid"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type FacingDirectionEnum int

const (
	FaceDirectionRight FacingDirectionEnum = iota
	FaceDirectionLeft
)

// TODO: Split PhysicsBody in a complex MovableCollidableAlive and MovableCollidable for items and move the methods to the right one.
type PhysicsBody struct {
	body.Shape

	MovableBody
	CollidableBody
	AliveBody

	id            string
	invulnerable  bool
	movementModel MovementModel
}

func NewPhysicsBody(shape body.Shape) *PhysicsBody {
	return &PhysicsBody{
		Shape: shape,
		id:    uuid.New().String(),
	}
}

// Attribute methods
func (b *PhysicsBody) ID() string {
	return b.id
}
func (b *PhysicsBody) SetID(id string) {
	b.id = id
}
func (b *PhysicsBody) Invulnerable() bool {
	return b.invulnerable
}
func (b *PhysicsBody) SetInvulnerable(value bool) {
	b.invulnerable = value
}

// Collision methods
func (b *PhysicsBody) OnTouch(other body.Body) {
	if b.Touchable != nil {
		b.Touchable.OnTouch(other)
	}
}

func (b *PhysicsBody) OnBlock(other body.Body) {
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

func (b *PhysicsBody) CollisionShapes() []*CollisionArea {
	return b.collisionList
}

// Movement methods
func (b *PhysicsBody) ApplyValidMovement(distance int, isXAxis bool, space body.BodiesSpace) {
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
		b.SetFaceDirection(FaceDirectionRight)
	} else if b.accelerationX < 0 {
		b.SetFaceDirection(FaceDirectionLeft)
	}
}

func (b *PhysicsBody) UpdateMovement(space body.BodiesSpace) {
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

func (b *PhysicsBody) SetPosition(x, y int) {
	switch b.Shape.(type) {
	case *Rect:
		rect := b.Shape.(*Rect)
		// Calculate the difference to move the collision areas as well
		diffX := x - rect.x16
		diffY := y - rect.y16

		rect.x16 = x
		rect.y16 = y

		for _, c := range b.collisionList {
			c.Shape.(*Rect).x16 += diffX
			c.Shape.(*Rect).y16 += diffY
		}
	}
}
