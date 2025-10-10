package state

import (
	"fmt"
)

type StateFactory interface {
	Create(state GameStateEnum) (GameState, error)
}

type DefaultStateFactory struct {
	stateMap StateMap
}

func NewDefaultSceneFactory(stateMap StateMap) *DefaultStateFactory {
	return &DefaultStateFactory{stateMap: stateMap}
}

func (f *DefaultStateFactory) Create(state GameStateEnum) (GameState, error) {
	s, ok := f.stateMap[state]
	if !ok {
		return nil, fmt.Errorf("unknown state type")
	}

	return s, nil
}
