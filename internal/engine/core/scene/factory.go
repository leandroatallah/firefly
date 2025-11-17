package scene

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

type SceneFactory interface {
	Create(sceneType navigation.SceneType, freshInstance bool) (navigation.Scene, error)
	SetAppContext(appContext any)
}

type DefaultSceneFactory struct {
	manager      navigation.SceneManager
	sceneMap     navigation.SceneMap
	cachedScenes map[navigation.SceneType]navigation.Scene
	appContext   *core.AppContext
}

func NewDefaultSceneFactory(sceneMap navigation.SceneMap) *DefaultSceneFactory {
	return &DefaultSceneFactory{
		sceneMap:     sceneMap,
		cachedScenes: make(map[navigation.SceneType]navigation.Scene),
	}
}

func (f *DefaultSceneFactory) SetAppContext(appContext any) {
	f.appContext = appContext.(*core.AppContext)
	f.manager = f.appContext.SceneManager
}

func (f *DefaultSceneFactory) Create(sceneType navigation.SceneType, freshInstance bool) (navigation.Scene, error) {
	if !freshInstance {
		if scene, ok := f.cachedScenes[sceneType]; ok {
			return scene, nil
		}
	}

	sceneFunc, ok := f.sceneMap[sceneType]
	if !ok {
		return nil, fmt.Errorf("unknown scene type")
	}

	scene := sceneFunc()
	scene.SetAppContext(f.appContext)

	if !freshInstance {
		f.cachedScenes[sceneType] = scene
	}

	return scene, nil
}
