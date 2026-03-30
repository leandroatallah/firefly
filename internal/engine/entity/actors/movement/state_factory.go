package movement

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// MovementStateOption defines a function that configures a movement state
type MovementStateOption func(MovementState)

func NewMovementState(
	actor body.MovableCollidable,
	state MovementStateEnum,
	target body.MovableCollidable,
	options ...MovementStateOption,
) (MovementState, error) {
	b := NewBaseMovementState(state, actor, target)

	var movementState MovementState

	switch state {
	case Idle:
		movementState = NewIdleMovementState(b)
	case Chase:
		movementState = NewChaseMovementState(b)
	case DumbChase:
		movementState = NewDumbChaseMovementState(b)
	case Avoid:
		movementState = NewAvoidMovementState(b)
	case Patrol:
		movementState = NewPatrolMovementState(b)
	case SideToSide:
		movementState = NewSideToSideMovementState(b)
	case Follow:
		movementState = NewFollowMovementState(b)
	default:
		// Check registry
		constructor, err := GetMovementStateConstructor(state)
		if err != nil {
			return nil, fmt.Errorf("unknown movement state type: %d", state)
		}
		movementState = constructor(b)
	}

	// Apply options
	for _, option := range options {
		if option != nil {
			option(movementState)
		}
	}

	return movementState, nil
}
