package scene

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

type SceneFactory interface {
	Create(sceneType navigation.SceneType, freshInstance bool) (navigation.Scene, error)
}

type DefaultSceneFactory struct {
	core.AppContextHolder

	sceneMap     navigation.SceneMap
	cachedScenes map[navigation.SceneType]navigation.Scene
}

func NewDefaultSceneFactory(sceneMap navigation.SceneMap) *DefaultSceneFactory {
	return &DefaultSceneFactory{
		sceneMap:     sceneMap,
		cachedScenes: make(map[navigation.SceneType]navigation.Scene),
	}
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
	scene.SetAppContext(f.AppContext())

	if !freshInstance {
		f.cachedScenes[sceneType] = scene
	}

	return scene, nil
}
