package state

import (
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/navigation"
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
