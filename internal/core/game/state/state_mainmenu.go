package state

import (
	"github.com/leandroatallah/firefly/internal/core"
	"github.com/leandroatallah/firefly/internal/navigation"
)

type MainMenuState struct {
	BaseState
}

func NewMainMenuState(ctx *core.AppContext) *MainMenuState {
	return &MainMenuState{BaseState: BaseState{ctx: ctx}}
}

func (s *MainMenuState) OnStart() {
	s.ctx.SceneManager.NavigateTo(navigation.SceneMenu, nil)
}
