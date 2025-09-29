package movement

type AvoidMovementState struct {
	BaseMovementState
}

func (s *AvoidMovementState) Move() {
	directions := calculateMovementDirections(s.actor, s.target, true)
	executeMovement(s.actor, directions)
}
