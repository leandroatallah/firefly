package movement

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

// DumbChaseMovementState provides a simple, direct chase behavior.
type DumbChaseMovementState struct {
	BaseMovementState
}

func NewDumbChaseMovementState(base BaseMovementState) *DumbChaseMovementState {
	return &DumbChaseMovementState{BaseMovementState: base}
}

func (s *DumbChaseMovementState) Move(space body.BodiesSpace) {
	if s.actor.Immobile() {
		return
	}

	directions := calculateMovementDirections(s.actor, s.target, false)
	executeMovement(s.actor, directions)
}
