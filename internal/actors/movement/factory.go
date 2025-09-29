package movement

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/systems/physics"
)

func NewMovementState(actor physics.Body, state MovementStateEnum, target physics.Body) (MovementState, error) {
	b := NewBaseMovementState(state, actor, target)

	switch state {
	case Input:
		return &InputMovementState{BaseMovementState: *b}, nil
	case Rand:
		return &RandMovementState{BaseMovementState: *b}, nil
	case Chase:
		return &ChaseMovementState{BaseMovementState: *b}, nil
	case DumbChase:
		return &DumbChaseMovementState{BaseMovementState: *b}, nil
	case Avoid:
		return &AvoidMovementState{BaseMovementState: *b}, nil
	case Patrol:
		return &PatrolMovementState{BaseMovementState: *b}, nil
	default:
		return nil, fmt.Errorf("unknown movement state type")
	}
}
