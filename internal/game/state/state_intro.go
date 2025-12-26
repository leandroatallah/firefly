package gamestate

import (
	"github.com/leandroatallah/firefly/internal/engine/core/game/state"
)

type IntroState struct {
	state.BaseState
}

// func NewIntroState(ctx *app.AppContext) *IntroState {
// 	return &IntroState{BaseState: state.BaseState{ctx: ctx}}
// }

func (s *IntroState) OnStart() {
	// s.ctx.SceneManager.NavigateTo(navigation.SceneIntro, nil)
}
