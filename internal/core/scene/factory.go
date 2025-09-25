package scene

import "fmt"

type SceneType int

type SceneFactory interface {
	Create(sceneType SceneType) (Scene, error)
	SetManager(manager *SceneManager)
}

type DefaultSceneFactory struct {
	manager *SceneManager
}

func NewDefaultSceneFactory() *DefaultSceneFactory {
	return &DefaultSceneFactory{}
}

func (f *DefaultSceneFactory) SetManager(manager *SceneManager) {
	f.manager = manager
}

const (
	SceneMenu SceneType = iota
	SceneSandbox
)

func (f *DefaultSceneFactory) Create(sceneType SceneType) (Scene, error) {
	var scene Scene
	var err error

	switch sceneType {
	case SceneMenu:
		scene = &MenuScene{}
	case SceneSandbox:
		scene = &SandboxScene{}
	default:
		err = fmt.Errorf("unknown scene type")
	}

	if err == nil {
		scene.SetManager(f.manager)
	}

	return scene, err
}
