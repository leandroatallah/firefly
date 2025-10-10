package actors

import (
	"fmt"
)

type ActorState interface {
	State() ActorStateEnum
	OnStart()
}

type ActorStateEnum int

// TODO: Rename to IdleState, etc.
const (
	Idle ActorStateEnum = iota
	Walk
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

type IdleState struct {
	BaseState
}

func (s *IdleState) OnStart() {}

type WalkState struct {
	BaseState
}

func (s *WalkState) OnStart() {}

type HurtState struct {
	BaseState
	count         int
	durationLimit int
}

func (s *HurtState) OnStart() {
	s.durationLimit = 30 // 0.5 sec
}

func (s *HurtState) CheckRecovery() bool {
	s.count++

	if s.count > s.durationLimit {
		return true
	}

	return false
}

// State factory method
func NewActorState(actor ActorEntity, state ActorStateEnum) (ActorState, error) {
	b := BaseState{actor: actor, state: state}
	switch state {
	case Idle:
		return &IdleState{BaseState: b}, nil
	case Walk:
		return &WalkState{BaseState: b}, nil
	case Hurted:
		return &HurtState{BaseState: b}, nil
	default:
		return nil, fmt.Errorf("unknown scene type")
	}
}
