package gamestates

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
)

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

// Growing
type GrowingState struct {
	actors.BaseState
}

func (s *GrowingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
	if scaler, ok := s.GetRootOwner().(interface{ SetScale(float64) }); ok {
		scaler.SetScale(1.0)
	}
}

// Shrinking
type ShrinkingState struct {
	actors.BaseState
}

func (s *ShrinkingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
	if scaler, ok := s.GetRootOwner().(interface{ SetScale(float64) }); ok {
		scaler.SetScale(1.0)
	}
}

var (
	Dying     actors.ActorStateEnum
	Dead      actors.ActorStateEnum
	Exiting   actors.ActorStateEnum
	Lying     actors.ActorStateEnum
	Rising    actors.ActorStateEnum
	Growing   actors.ActorStateEnum
	Shrinking actors.ActorStateEnum
)

func init() {
	Dying = actors.Dying
	Dead = actors.Dead
	Exiting = actors.Exiting
	Lying = actors.RegisterState("lie", func(b actors.BaseState) actors.ActorState { return &LyingState{BaseState: b} })
	Rising = actors.RegisterState("rise", func(b actors.BaseState) actors.ActorState { return &RisingState{BaseState: b} })
	Growing = actors.RegisterState("grow", func(b actors.BaseState) actors.ActorState { return &GrowingState{BaseState: b} })
	Shrinking = actors.RegisterState("shrink", func(b actors.BaseState) actors.ActorState { return &ShrinkingState{BaseState: b} })
}
