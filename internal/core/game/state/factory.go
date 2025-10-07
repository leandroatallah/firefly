package state

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/core"
)

// State factory method
func NewGameState(state GameStateEnum, ctx *core.AppContext) (GameState, error) {
	base := BaseState{ctx: ctx}
	switch state {
	case Intro:
		return NewIntroState(ctx), nil
	case MainMenu:
		return NewMainMenuState(ctx), nil
	case Playing:
		return &PlayingState{BaseState: base}, nil
	case Paused:
		return &PausedState{BaseState: base}, nil
	case GameOver:
		return &GameOverState{BaseState: base}, nil
	default:
		return nil, fmt.Errorf("unknown scene type")
	}
}
