// Package actors implements the registry-based state machine for game entities.
// Each actor state is registered at init time via RegisterState, which assigns a
// unique ActorStateEnum value and a constructor. The Character drives transitions
// through handleState, delegating per-state behaviour to ActorState implementations.
package actors

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity"
)

// ActorState is the interface every concrete state must satisfy.
type ActorState interface {
	State() ActorStateEnum
	OnStart(currentCount int)
	OnFinish()
	GetAnimationCount(currentCount int) int
	IsAnimationFinished() bool
}

// ActorStateEnum is an integer token that uniquely identifies a registered state.
type ActorStateEnum int

// State enum: part of engine public API
//
//nolint:gochecknoglobals
var (
	Idle    ActorStateEnum
	Walking ActorStateEnum
	Jumping ActorStateEnum
	Falling ActorStateEnum
	Landing ActorStateEnum
	Hurted  ActorStateEnum
	Dying   ActorStateEnum
	Dead    ActorStateEnum
	Exiting ActorStateEnum

	IdleShooting    ActorStateEnum
	WalkingShooting ActorStateEnum
	JumpingShooting ActorStateEnum
	FallingShooting ActorStateEnum
)

func init() {
	Idle = RegisterState("idle", func(b BaseState) ActorState { return &IdleState{BaseState: b} })
	Walking = RegisterState("walk", func(b BaseState) ActorState { return &WalkState{BaseState: b} })
	Jumping = RegisterState("jump", func(b BaseState) ActorState { return &JumpState{BaseState: b} })
	Falling = RegisterState("fall", func(b BaseState) ActorState { return &FallState{BaseState: b} })
	Landing = RegisterState("land", func(b BaseState) ActorState { return &LandingState{BaseState: b} })
	Hurted = RegisterState("hurt", func(b BaseState) ActorState { return &HurtState{BaseState: b} })
	Dying = RegisterState("die", func(b BaseState) ActorState { return &DyingState{BaseState: b} })
	Dead = RegisterState("dead", func(b BaseState) ActorState { return &DeadState{BaseState: b} })
	Exiting = RegisterState("exit", func(b BaseState) ActorState { return &ExitingState{BaseState: b} })

	IdleShooting = RegisterState("idle_shoot", func(b BaseState) ActorState { return &IdleShootingState{BaseState: b} })
	WalkingShooting = RegisterState("walk_shoot", func(b BaseState) ActorState { return &WalkingShootingState{BaseState: b} })
	JumpingShooting = RegisterState("jump_shoot", func(b BaseState) ActorState { return &JumpingShootingState{BaseState: b} })
	FallingShooting = RegisterState("fall_shoot", func(b BaseState) ActorState { return &FallingShootingState{BaseState: b} })
}

// BaseState provides shared bookkeeping (entry count, tick) for all concrete states.
type BaseState struct {
	actor      ActorEntity
	state      ActorStateEnum
	entryCount int
	tick       int
}

func NewBaseState(actor ActorEntity, state ActorStateEnum) BaseState {
	return BaseState{actor: actor, state: state}
}

func (s *BaseState) State() ActorStateEnum {
	return s.state
}
func (s *BaseState) GetActor() ActorEntity {
	return s.actor
}

func (s *BaseState) GetRootOwner() interface{} {
	actor := s.GetActor()
	var root interface{} = actor
	if lastOwner := actor.LastOwner(); lastOwner != nil {
		root = lastOwner
	}
	return root
}

func (s *BaseState) OnStart(currentCount int) {
	s.entryCount = currentCount
	s.tick = 0
}

func (s *BaseState) GetAnimationCount(currentCount int) int {
	return currentCount - s.entryCount
}

func (s *BaseState) OnFinish() {}

func (s *BaseState) IsAnimationFinished() bool {
	s.tick++

	character := s.GetActor().GetCharacter()
	return entity.IsAnimationFinished(s.tick, character, s.State())
}
