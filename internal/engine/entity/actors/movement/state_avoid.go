package movement

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

type AvoidMovementState struct {
	BaseMovementState
}

func NewAvoidMovementState(base BaseMovementState) *AvoidMovementState {
	return &AvoidMovementState{BaseMovementState: base}
}

func (s *AvoidMovementState) Move(space body.BodiesSpace) {
	if s.actor.Immobile() {
		return
	}

	directions := calculateMovementDirections(s.actor, s.target, true)
	executeMovement(s.actor, directions)
}
