package navigation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/systems/audiomanager"
)

type SceneType int

// TODO: should it be initialized in Game?
const (
	SceneMenu SceneType = iota
	SceneSandbox
	ScenePlatform
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update() error
	OnStart()
	OnFinish()
	// Use any to prevent import cycle error
	SetAppContext(appContext any)
}

type SceneManager interface {
	AudioManager() *audiomanager.AudioManager
	Draw(screen *ebiten.Image)
	NavigateTo(sceneType SceneType, sceneTransition Transition)
	// SetFactory(factory SceneFactory)
	SwitchTo(scene Scene)
	Update() error
}

type Transition interface {
	Update()
	Draw(screen *ebiten.Image)
	StartTransition(func())
	EndTransition(func())
}
