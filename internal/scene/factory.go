package scene

import "fmt"

type SceneType int

type SceneFactory interface {
	Create(sceneType SceneType) (Scene, error)
}

type DefaultSceneFactory struct{}

func NewDefaultSceneFactory() *DefaultSceneFactory {
	return &DefaultSceneFactory{}
}

const (
	SceneMenu SceneType = iota
	SceneSandbox
)

func (f *DefaultSceneFactory) Create(sceneType SceneType) (Scene, error) {
	switch sceneType {
	case SceneMenu:
		return &MenuScene{}, nil
	case SceneSandbox:
		return &SandboxScene{}, nil
	default:
		return nil, fmt.Errorf("unknown scene type")
	}
}
