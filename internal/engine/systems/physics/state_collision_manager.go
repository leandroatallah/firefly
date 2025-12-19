package physics

import (
	"fmt"
	"image"
	"time"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
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
		for _, r := range rects {
			template, ok := r.(*CollidableBody)
			if !ok {
				continue
			}

			newCollisionBody := NewCollidableBody(NewBody(template.GetShape()))
			relativePos := template.Position()
			newPos := image.Rect(
				x+relativePos.Min.X,
				y+relativePos.Min.Y,
				x+relativePos.Max.X,
				y+relativePos.Max.Y,
			)
			newCollisionBody.SetPosition(newPos.Min.X, newPos.Min.Y)
			newCollisionBody.SetID(fmt.Sprintf("%s_collision_%d", m.owner.ID(), time.Now().UnixNano()))
			m.owner.AddCollision(newCollisionBody)
		}
	}
}
