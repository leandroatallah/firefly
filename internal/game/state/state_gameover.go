package gamestate

import "github.com/leandroatallah/firefly/internal/engine/core/game/state"

type GameOverState struct {
	state.BaseState
}

func (s *GameOverState) OnStart() {}
