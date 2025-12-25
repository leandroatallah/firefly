package movement

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type MovementState interface {
	State() MovementStateEnum
	OnStart()
	Move()
	Target() body.MovableCollidable
}

type MovementStateEnum int

const (
	Input MovementStateEnum = iota
	Idle
	Rand
	Chase
	DumbChase
	Patrol
	Avoid
)

type BaseMovementState struct {
	state  MovementStateEnum
	actor  body.MovableCollidable
	target body.MovableCollidable
}

func NewBaseMovementState(
	state MovementStateEnum,
	actor body.MovableCollidable,
	target body.MovableCollidable,
) BaseMovementState {
	return BaseMovementState{state: state, actor: actor, target: target}
}

func (s *BaseMovementState) State() MovementStateEnum {
	return s.state
}

func (s *BaseMovementState) OnStart() {}

func (s *BaseMovementState) Target() body.MovableCollidable {
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
func calculateMovementDirections(actorPos, targetPos body.Body, isAvoid bool) MovementDirections {
	actorRect := actorPos.Position()
	targetRect := targetPos.Position()
	p0x, p0y := actorRect.Min.X, actorRect.Min.Y
	p1x, p1y := actorRect.Max.X, actorRect.Max.Y
	e0x, e0y := targetRect.Min.X, targetRect.Min.Y
	e1x, e1y := targetRect.Max.X, targetRect.Max.Y
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

func executeMovement(actor body.MovableCollidable, directions MovementDirections) {
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
