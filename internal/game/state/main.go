package gamestate

import "github.com/leandroatallah/firefly/internal/engine/core/game/state"

const (
	Intro state.GameStateEnum = iota
	MainMenu
	Playing
	Paused
	GameOver
)
