package movement

import (
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type MovementState interface {
	State() MovementStateEnum
	OnStart()
	Move()
	Target() physics.Body
}

type MovementStateEnum int

const (
	Input MovementStateEnum = iota
	Rand
	Chase
	DumbChase
	Patrol
	Avoid
)

type BaseMovementState struct {
	state  MovementStateEnum
	actor  physics.Body
	target physics.Body
}

func NewBaseMovementState(
	state MovementStateEnum,
	actor physics.Body,
	target physics.Body,
) *BaseMovementState {
	return &BaseMovementState{state: state, actor: actor, target: target}
}

func (s *BaseMovementState) State() MovementStateEnum {
	return s.state
}

func (s *BaseMovementState) OnStart() {}

func (s *BaseMovementState) Target() physics.Body {
	return s.target
}

// Movement utility functions
type MovementDirections struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

// calculateMovementDirections determines which directions to move based on actor and target positions
func calculateMovementDirections(actorPos, targetPos physics.Body, isAvoid bool) MovementDirections {
	p0x, p0y, p1x, p1y := actorPos.Position()
	e0x, e0y, e1x, e1y := targetPos.Position()
	var up, down, left, right bool

	// Check direction to chase destination
	if p1x < e0x {
		right = true
	} else if p0x > e1x {
		left = true
	}

	if p1y < e0y {
		down = true
	} else if p0y > e1y {
		up = true
	}

	if isAvoid {
		// Invert to  move away from target
		up, down, left, right = !up, !down, !left, !right
	}

	return MovementDirections{Up: up, Down: down, Left: left, Right: right}
}

func executeMovement(actor physics.Body, directions MovementDirections) {
	if !directions.Up && !directions.Down && !directions.Left && !directions.Right {
		return
	}

	speed := actor.Speed()

	if directions.Up {
		if directions.Left {
			actor.OnMoveUpLeft(speed)
		} else if directions.Right {
			actor.OnMoveUpRight(speed)
		} else {
			actor.OnMoveUp(speed)
		}
	} else if directions.Down {
		if directions.Left {
			actor.OnMoveDownLeft(speed)
		} else if directions.Right {
			actor.OnMoveDownRight(speed)
		} else {
			actor.OnMoveDown(speed)
		}
	} else if directions.Left {
		actor.OnMoveLeft(speed)
	} else if directions.Right {
		actor.OnMoveRight(speed)
	}
}
