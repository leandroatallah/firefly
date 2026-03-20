package gamestates

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
)

// Dying
type DyingState struct {
	actors.BaseState
}

func (s *DyingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if p, ok := s.GetRootOwner().(platformer.PlatformerActorEntity); ok {
		p.OnDie()
	}
}

// Dead
type DeadState struct {
	actors.BaseState
}

func (s *DeadState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	actor := s.GetActor()
	actor.SetFreeze(true)
	actor.SetImmobile(true)
}

// Exiting
type ExitingState struct {
	actors.BaseState
}

func (s *ExitingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

}

// Lying
type LyingState struct {
	actors.BaseState
}

func (s *LyingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Rising
type RisingState struct {
	actors.BaseState
}

func (s *RisingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	actor := s.GetActor()
	actor.SetHealth(actor.MaxHealth())
	actor.SetFreeze(false)
	actor.SetVelocity(0, 0)
	actor.SetAcceleration(0, 0)
}

func (s *RisingState) OnFinish() {
	s.GetActor().SetImmobile(false)
}

var (
	Dying   actors.ActorStateEnum
	Dead    actors.ActorStateEnum
	Exiting actors.ActorStateEnum
	Lying   actors.ActorStateEnum
	Rising  actors.ActorStateEnum
)

func init() {
	Dying = actors.RegisterState("die", func(b actors.BaseState) actors.ActorState { return &DyingState{BaseState: b} })
	Dead = actors.RegisterState("dead", func(b actors.BaseState) actors.ActorState { return &DeadState{BaseState: b} })
	Exiting = actors.RegisterState("exit", func(b actors.BaseState) actors.ActorState { return &ExitingState{BaseState: b} })
	Lying = actors.RegisterState("lie", func(b actors.BaseState) actors.ActorState { return &LyingState{BaseState: b} })
	Rising = actors.RegisterState("rise", func(b actors.BaseState) actors.ActorState { return &RisingState{BaseState: b} })
}
