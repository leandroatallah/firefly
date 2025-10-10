package gamestate

import "github.com/leandroatallah/firefly/internal/engine/core/game/state"

type PausedState struct {
	state.BaseState
}

func (s *PausedState) OnStart() {}
