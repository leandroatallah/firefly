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
}

var (
	Dying   actors.ActorStateEnum
	Exiting actors.ActorStateEnum
	Lying   actors.ActorStateEnum
	Rising  actors.ActorStateEnum
)

func init() {
	Dying = actors.RegisterState("die", func(b actors.BaseState) actors.ActorState { return &DyingState{BaseState: b} })
	Exiting = actors.RegisterState("exit", func(b actors.BaseState) actors.ActorState { return &ExitingState{BaseState: b} })
	Lying = actors.RegisterState("lie", func(b actors.BaseState) actors.ActorState { return &LyingState{BaseState: b} })
	Rising = actors.RegisterState("rise", func(b actors.BaseState) actors.ActorState { return &RisingState{BaseState: b} })
}
