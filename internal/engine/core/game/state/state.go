package state

import (
	"github.com/leandroatallah/firefly/internal/engine/core"
)

type GameState interface {
	OnStart()
}

type GameStateEnum int

type StateMap map[GameStateEnum]GameState

type BaseState struct {
	ctx *core.AppContext

	stateMap map[GameStateEnum]GameState
}

func NewBaseState(stateMap map[GameStateEnum]GameState) *BaseState {
	return &BaseState{stateMap: stateMap}
}

func (s *BaseState) OnStart() {}
