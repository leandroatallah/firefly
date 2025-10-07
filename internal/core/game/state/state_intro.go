package state

import (
	"github.com/leandroatallah/firefly/internal/core"
	"github.com/leandroatallah/firefly/internal/navigation"
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
