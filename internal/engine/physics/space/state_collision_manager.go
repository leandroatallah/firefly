package space

import (
	"fmt"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
)

// StateEnum is a constraint for generic types that represent states, typically as integers.
type StateEnum interface {
	~int
}

// StateBasedCollisioner defines the interface for an entity that manages collisions based on its state.
type StateBasedCollisioner[T StateEnum] interface {
	State() T
	GetPositionMin() (int, int)
	ClearCollisions()
	AddCollision(...body.Collidable)
	ID() string
	Scale() float64
}

// StateCollisionManager manages state-based collision bodies for an entity.
type StateCollisionManager[T StateEnum] struct {
	owner           StateBasedCollisioner[T]
	collisionBodies map[T][]body.Collidable
}

// NewStateCollisionManager creates a new manager for state-based collisions.
func NewStateCollisionManager[T StateEnum](owner StateBasedCollisioner[T]) *StateCollisionManager[T] {
	return &StateCollisionManager[T]{
		owner:           owner,
		collisionBodies: make(map[T][]body.Collidable),
	}
}

// AddCollisionRect associates a collision rectangle with a specific state.
func (m *StateCollisionManager[T]) AddCollisionRect(state T, rect body.Collidable) {
	m.collisionBodies[state] = append(m.collisionBodies[state], rect)
}

// RefreshCollisions updates the entity's collision bodies based on its current state.
func (m *StateCollisionManager[T]) RefreshCollisions() {
	currentState := m.owner.State()
	if rects, ok := m.collisionBodies[currentState]; ok {
		m.owner.ClearCollisions()
		x, y := m.owner.GetPositionMin()
		scale := m.owner.Scale()
		for _, r := range rects {
			template, ok := r.(*bodyphysics.CollidableBody)
			if !ok {
				continue
			}

			shape := template.GetShape()
			if scale != 0 && scale != 1.0 {
				if rect, ok := shape.(*bodyphysics.Rect); ok {
					shape = bodyphysics.NewRect(
						0, 0,
						int(float64(rect.Width())*scale),
						int(float64(rect.Height())*scale),
					)
				}
			}

			newBody := bodyphysics.NewBody(shape)
			// Set owner to movable? We don't have movable here.
			// But m.owner is typically a Character.
			// If we set newBody.Owner = m.owner directly, newBody.TopOwner() -> m.owner.TopOwner().
			// This shortcuts MovableBody but it's fine for collision bodies generated for states.
			newBody.SetOwner(m.owner)
			newCollisionBody := bodyphysics.NewCollidableBody(newBody)
			newCollisionBody.SetOwner(m.owner)
			relativePos := template.Position()
			relMinX := relativePos.Min.X
			relMinY := relativePos.Min.Y
			if scale != 0 && scale != 1.0 {
				relMinX = int(float64(relMinX) * scale)
				relMinY = int(float64(relMinY) * scale)
			}
			newCollisionBody.SetPosition(x+relMinX, y+relMinY)
			newCollisionBody.SetID(fmt.Sprintf("%s_collision_%d", m.owner.ID(), time.Now().UnixNano()))
			m.owner.AddCollision(newCollisionBody)
		}
	}
}
