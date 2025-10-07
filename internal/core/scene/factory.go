package scene

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/levels"
	"github.com/leandroatallah/firefly/internal/navigation"
)

type SceneFactory interface {
	Create(sceneType navigation.SceneType) (navigation.Scene, error)
	SetManager(manager navigation.SceneManager)
	SetLevelManager(manager *levels.Manager)
}

type DefaultSceneFactory struct {
	manager      navigation.SceneManager
	levelManager *levels.Manager
}

func NewDefaultSceneFactory() *DefaultSceneFactory {
	return &DefaultSceneFactory{}
}

func (f *DefaultSceneFactory) SetManager(manager navigation.SceneManager) {
	f.manager = manager
}

func (f *DefaultSceneFactory) SetLevelManager(manager *levels.Manager) {
	f.levelManager = manager
}

func (f *DefaultSceneFactory) Create(sceneType navigation.SceneType) (navigation.Scene, error) {
	var scene navigation.Scene
	var err error

	switch sceneType {
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
		scene.SetManager(f.manager)
		scene.SetLevelManager(f.levelManager)
	}

	return scene, err
}
