package movement

type IdleMovementState struct {
	BaseMovementState
}

func NewIdleMovementState(base BaseMovementState) *IdleMovementState {
	return &IdleMovementState{BaseMovementState: base}
}

func (s *IdleMovementState) Move() {}
