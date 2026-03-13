package movement

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

type IdleMovementState struct {
	BaseMovementState
}

func NewIdleMovementState(base BaseMovementState) *IdleMovementState {
	return &IdleMovementState{BaseMovementState: base}
}

func (s *IdleMovementState) Move(space body.BodiesSpace) {}
