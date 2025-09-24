package game

import (
	"github.com/leandroatallah/firefly/internal/scene"
)

type MainMenuState struct {
	BaseState
}

func NewMainMenuState() *MainMenuState {
	return &MainMenuState{}
}

func (s *MainMenuState) OnStart() {
	s.game.sceneManager.GoToScene(scene.SceneMenu, nil)

}
