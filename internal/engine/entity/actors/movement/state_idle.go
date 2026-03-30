package movement

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

type IdleMovementState struct {
	BaseMovementState
}

func NewIdleMovementState(base BaseMovementState) *IdleMovementState {
	return &IdleMovementState{BaseMovementState: base}
}

func (s *IdleMovementState) Move(space body.BodiesSpace) {}
