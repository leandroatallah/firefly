package movement

type AvoidMovementState struct {
	BaseMovementState
}

func NewAvoidMovementState(base BaseMovementState) *AvoidMovementState {
	return &AvoidMovementState{BaseMovementState: base}
}

func (s *AvoidMovementState) Move() {
	directions := calculateMovementDirections(s.actor, s.target, true)
	executeMovement(s.actor, directions)
}
