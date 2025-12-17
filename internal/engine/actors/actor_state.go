package actors

import (
	"fmt"
)

type ActorState interface {
	State() ActorStateEnum
	OnStart()
}

type ActorStateEnum int

const (
	Idle ActorStateEnum = iota
	Walking
	Falling
	Hurted
)

type BaseState struct {
	actor ActorEntity
	state ActorStateEnum
}

func (s *BaseState) State() ActorStateEnum {
	return s.state
}

func (s *BaseState) OnStart() {}

// State factory method
func NewActorState(actor ActorEntity, state ActorStateEnum) (ActorState, error) {
	b := BaseState{actor: actor, state: state}
	switch state {
	case Idle:
		return &IdleState{BaseState: b}, nil
	case Walking:
		return &WalkState{BaseState: b}, nil
	case Falling:
		return &FallState{BaseState: b}, nil
	case Hurted:
		return &HurtState{BaseState: b}, nil
	default:
		return nil, fmt.Errorf("unknown actor state")
	}
}
