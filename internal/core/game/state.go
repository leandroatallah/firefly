package game

import (
	"fmt"
)

type GameState interface {
	SetContext(game *Game)
	SetState(state GameState)
	OnStart()
}

type GameStateEnum int

const (
	MainMenu GameStateEnum = iota
	Playing
	Paused
	GameOver
)

type BaseState struct {
	game *Game
}

func (s *BaseState) SetState(state GameState) {
	s.game.state = state
}

func (s *BaseState) OnStart() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseState) SetContext(game *Game) {
	s.game = game
}

// Concrete States
// TODO: Move to different files as MainMenuState
type PlayingState struct {
	BaseState
}

func (s *PlayingState) OnStart() {}

type PausedState struct {
	BaseState
}

func (s *PausedState) OnStart() {}

type GameOverState struct {
	BaseState
}

func (s *GameOverState) OnStart() {}

// State factory method
func NewGameState(state GameStateEnum) (GameState, error) {
	switch state {
	case MainMenu:
		return NewMainMenuState(), nil
	case Playing:
		return &PlayingState{}, nil
	case Paused:
		return &PausedState{}, nil
	case GameOver:
		return &GameOverState{}, nil
	default:
		return nil, fmt.Errorf("unknown scene type")
	}
}
