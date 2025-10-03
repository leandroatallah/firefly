package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core"
	"github.com/leandroatallah/firefly/internal/core/game/state"
)

type Game struct {
	AppContext *core.AppContext
	state      state.GameState
}

func NewGame(ctx *core.AppContext) *Game {
	return &Game{AppContext: ctx}
}

func (g *Game) Update() error {
	// First, update the input manager
	g.AppContext.InputManager.Update()

	// Then, update the current scene
	g.AppContext.SceneManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.AppContext.SceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) SetState(stateID state.GameStateEnum) error {
	state, err := state.NewGameState(stateID, g.AppContext)
	if err != nil {
		return err
	}

	g.state = state
	g.state.OnStart()

	return nil
}
