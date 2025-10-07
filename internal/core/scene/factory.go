package scene

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/core"
	"github.com/leandroatallah/firefly/internal/navigation"
)

type SceneFactory interface {
	Create(sceneType navigation.SceneType) (navigation.Scene, error)
	SetAppContext(appContext any)
}

type DefaultSceneFactory struct {
	manager    navigation.SceneManager
	appContext *core.AppContext
}

func NewDefaultSceneFactory() *DefaultSceneFactory {
	return &DefaultSceneFactory{}
}

func (f *DefaultSceneFactory) SetAppContext(appContext any) {
	f.appContext = appContext.(*core.AppContext)
	f.manager = f.appContext.SceneManager
}

func (f *DefaultSceneFactory) Create(sceneType navigation.SceneType) (navigation.Scene, error) {
	var scene navigation.Scene
	var err error

	switch sceneType {
	case navigation.SceneIntro:
		scene = NewIntroScene()
	case navigation.SceneMenu:
		scene = &MenuScene{}
	case navigation.SceneSandbox:
		scene = &SandboxScene{}
	case navigation.ScenePlatform:
		scene = &PlatformScene{}
	default:
		err = fmt.Errorf("unknown scene type")
	}

	if err == nil {
		scene.SetAppContext(f.appContext)
	}

	return scene, err
}
