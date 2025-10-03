package movement

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/systems/physics"
)

// MovementStateOption defines a function that configures a movement state
type MovementStateOption func(MovementState)

func NewMovementState(
	actor physics.Body,
	state MovementStateEnum,
	target physics.Body,
	options ...MovementStateOption,
) (MovementState, error) {
	b := NewBaseMovementState(state, actor, target)

	var movementState MovementState

	switch state {
	case Chase:
		movementState = NewChaseMovementState(b)
	case DumbChase:
		movementState = NewDumbChaseMovementState(b)
	case Avoid:
		movementState = NewAvoidMovementState(b)
	case Patrol:
		movementState = NewPatrolMovementState(b)
	default:
		return nil, fmt.Errorf("unknown movement state type")
	}

	// Apply options
	// TODO: Improve this to handle variadict/optitonbal parameter
	for _, option := range options {
		if option != nil {
			option(movementState)
		}
	}

	return movementState, nil
}
