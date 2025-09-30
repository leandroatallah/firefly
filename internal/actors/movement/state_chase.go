package movement

type ChaseMovementState struct {
	BaseMovementState
}

func NewChaseMovementState(base BaseMovementState) *ChaseMovementState {
	return &ChaseMovementState{BaseMovementState: base}
}

func (s *ChaseMovementState) Move() {}

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
