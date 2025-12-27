package items

import (
	"fmt"
)

type ItemState interface {
	State() ItemStateEnum
	OnStart()
}

type ItemStateEnum int

const (
	Idle ItemStateEnum = iota
	Walking
	Falling
	Hurted
)

type BaseState struct {
	item  Item
	state ItemStateEnum
}

func (s *BaseState) State() ItemStateEnum {
	return s.state
}

func (s *BaseState) OnStart() {}

// State factory method
func NewItemState(item Item, state ItemStateEnum) (ItemState, error) {
	b := BaseState{item: item, state: state}
	switch state {
	case Idle:
		return &IdleState{BaseState: b}, nil
	default:
		return nil, fmt.Errorf("unknown item state")
	}
}
