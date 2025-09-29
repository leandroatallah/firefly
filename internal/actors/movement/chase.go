package movement

type ChaseMovementState struct {
	BaseMovementState
}

func (s *ChaseMovementState) Move() {}

type DumbChaseMovementState struct {
	BaseMovementState
}

func (s *DumbChaseMovementState) Move() {
	directions := calculateMovementDirections(s.actor, s.target, false)
	executeMovement(s.actor, directions)
}
