package scene

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/navigation"
)

type SceneFactory interface {
	Create(sceneType navigation.SceneType) (navigation.Scene, error)
	SetAppContext(appContext any)
}

type DefaultSceneFactory struct {
	manager    navigation.SceneManager
	sceneMap   navigation.SceneMap
	appContext *core.AppContext
}

func NewDefaultSceneFactory(sceneMap navigation.SceneMap) *DefaultSceneFactory {
	return &DefaultSceneFactory{sceneMap: sceneMap}
}

func (f *DefaultSceneFactory) SetAppContext(appContext any) {
	f.appContext = appContext.(*core.AppContext)
	f.manager = f.appContext.SceneManager
}

func (f *DefaultSceneFactory) Create(sceneType navigation.SceneType) (navigation.Scene, error) {
	scene, ok := f.sceneMap[sceneType]
	if !ok {
		return nil, fmt.Errorf("unknown scene type")
	}

	scene.SetAppContext(f.appContext)

	return scene, nil
}
