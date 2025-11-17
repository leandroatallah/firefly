package navigation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
)

type SceneType int

type Scene interface {
	Draw(screen *ebiten.Image)
	Update() error
	OnStart()
	OnFinish()
	// Use any to prevent import cycle error
	SetAppContext(appContext any)
}

type SceneFactory interface {
	Create(sceneType SceneType, freshInstance bool) (Scene, error)
	SetAppContext(appContext any)
}

type SceneMap map[SceneType]func() Scene

type SceneManager interface {
	AudioManager() *audiomanager.AudioManager
	Draw(screen *ebiten.Image)
	NavigateTo(sceneType SceneType, sceneTransition Transition, freshInstance bool)
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
