package gamestate

import "github.com/leandroatallah/firefly/internal/engine/core/game/state"

type PlayingState struct {
	state.BaseState
}

func (s *PlayingState) OnStart() {}
