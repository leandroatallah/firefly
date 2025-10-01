package movement

// DumbChaseMovementState provides a simple, direct chase behavior.
type DumbChaseMovementState struct {
	BaseMovementState
}

func NewDumbChaseMovementState(base BaseMovementState) *DumbChaseMovementState {
	return &DumbChaseMovementState{BaseMovementState: base}
}

func (s *DumbChaseMovementState) Move() {
	directions := calculateMovementDirections(s.actor, s.target, false)
	executeMovement(s.actor, directions)
}
