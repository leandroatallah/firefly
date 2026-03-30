package items

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity"
)

type ItemState interface {
	State() ItemStateEnum
	OnStart()
	IsAnimationFinished() bool
}

type ItemStateEnum int

var (
	Idle    ItemStateEnum
	Walking ItemStateEnum
	Falling ItemStateEnum
	Hurted  ItemStateEnum
)

func init() {
	Idle = RegisterState("idle", func(b BaseState) ItemState { return &IdleState{BaseState: b} })
	Walking = RegisterState("walking", func(b BaseState) ItemState { return &IdleState{BaseState: b} }) // Placeholder
	Falling = RegisterState("falling", func(b BaseState) ItemState { return &IdleState{BaseState: b} }) // Placeholder
	Hurted = RegisterState("hurted", func(b BaseState) ItemState { return &IdleState{BaseState: b} })   // Placeholder
}

type BaseState struct {
	item  Item
	state ItemStateEnum
	tick  int
}

func (s *BaseState) Item() Item {
	return s.item
}

func (s *BaseState) State() ItemStateEnum {
	return s.state
}

func (s *BaseState) OnStart() {
	s.tick = 0
}

func (s *BaseState) IsAnimationFinished() bool {
	s.tick++

	animatable, ok := s.item.(entity.Animatable)
	if !ok {
		return true
	}

	return entity.IsAnimationFinished(s.tick, animatable, s.State())
}

// State factory method
func NewItemState(item Item, state ItemStateEnum) (ItemState, error) {
	return NewState(item, state)
}
