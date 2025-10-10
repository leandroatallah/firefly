package state

import (
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/navigation"
)

type IntroState struct {
	BaseState
}

func NewIntroState(ctx *core.AppContext) *IntroState {
	return &IntroState{BaseState: BaseState{ctx: ctx}}
}

func (s *IntroState) OnStart() {
	s.ctx.SceneManager.NavigateTo(navigation.SceneIntro, nil)
}
