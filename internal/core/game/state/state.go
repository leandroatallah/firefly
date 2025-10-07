package state

import (
	"github.com/leandroatallah/firefly/internal/core"
)

type GameState interface {
	OnStart()
}

type GameStateEnum int

const (
	Intro GameStateEnum = iota
	MainMenu
	Playing
	Paused
	GameOver
)

type BaseState struct {
	ctx *core.AppContext
}

func (s *BaseState) OnStart() {}
