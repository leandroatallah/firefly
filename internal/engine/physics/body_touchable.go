package physics

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// TouchTrigger implements the physics.Touchable interface to handle body contact.
type TouchTrigger struct {
	execute func()
	subject body.Body
}

func NewTouchTrigger(execute func(), subject body.Body) *TouchTrigger {
	return &TouchTrigger{execute: execute, subject: subject}
}

// OnTouch is called by the physics engine when a non-obstructive collision occurs.
func (e *TouchTrigger) OnTouch(other body.Collidable) {
	if e.subject == nil {
		return
	}
	// Check if the body that touched the endpoint is the subject.
	if other.ID() == e.subject.ID() {
		if e.execute != nil {
			e.execute()
		}
	}
}

// OnBlock is called for obstructive collisions, which won't happen for a sensor.
func (e *TouchTrigger) OnBlock(other body.Collidable) {}
