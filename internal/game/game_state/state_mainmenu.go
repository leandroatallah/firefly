package gamestate

import (
	"github.com/leandroatallah/firefly/internal/engine/core/game/state"
)

type MainMenuState struct {
	state.BaseState
}

// func NewMainMenuState(ctx *core.AppContext) *MainMenuState {
// 	return &MainMenuState{BaseState: BaseState{ctx: ctx}}
// }

func (s *MainMenuState) OnStart() {
	// s.ctx.SceneManager.NavigateTo(navigation.SceneMenu, nil)
}
