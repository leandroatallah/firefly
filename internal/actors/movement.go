package actors

import (
	"fmt"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/systems/input"
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
	actor  ActorEntity
	target physics.Body
}

func NewBaseMovementState(
	state MovementStateEnum,
	actor ActorEntity,
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

// Movement States
type InputMovementState struct {
	BaseMovementState
}

func (s *InputMovementState) Move() {
	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		s.actor.OnMoveLeft(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		s.actor.OnMoveRight(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		s.actor.OnMoveUp(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		s.actor.OnMoveDown(s.actor.Speed())
	}
}

type RandMovementState struct {
	BaseMovementState
}

func (s *RandMovementState) Move() {
	// if s..count%60 != 0 {
	// 	return
	// }
	a := []func(){
		func() { s.actor.OnMoveLeft(s.actor.Speed()) },
		func() { s.actor.OnMoveRight(s.actor.Speed()) },
		func() { s.actor.OnMoveUp(s.actor.Speed()) },
		func() { s.actor.OnMoveDown(s.actor.Speed()) },
	}

	i := rand.IntN(4)
	a[i]()
}

type ChaseMovementState struct {
	BaseMovementState
}

func (s *ChaseMovementState) Move() {}

type DumbChaseMovementState struct {
	BaseMovementState
}

func (s *DumbChaseMovementState) Move() {
	p0x, p0y, p1x, p1y := s.actor.Position()
	e0x, e0y, e1x, e1y := s.target.Position()
	var up, down, left, right bool

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

	if !up && !down && !left && !right {
		return
	}

	speed := s.actor.Speed()

	if up {
		if left {
			s.actor.OnMoveUpLeft(speed)
		} else if right {
			s.actor.OnMoveUpRight(speed)
		} else {
			s.actor.OnMoveUp(speed)
		}
	} else if down {
		if left {
			s.actor.OnMoveDownLeft(speed)
		} else if right {
			s.actor.OnMoveDownRight(speed)
		} else {
			s.actor.OnMoveDown(speed)
		}
	} else if left {
		s.actor.OnMoveLeft(speed)
	} else if right {
		s.actor.OnMoveRight(speed)
	}
}

type PatrolMovementState struct {
	BaseMovementState
}

type AvoidMovementState struct {
	BaseMovementState
}

// State factory method
// TODO: It should be a method
func NewMovementState(actor ActorEntity, state MovementStateEnum, target physics.Body) (MovementState, error) {
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
	default:
		return nil, fmt.Errorf("unknown movement state type")
	}
}
