package gamescene

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

const (
	SceneIntro navigation.SceneType = iota
	SceneMenu
	SceneSandbox
	SceneLevels
	SceneSummary
)

func InitSceneMap(context *core.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		SceneIntro: func() navigation.Scene {
			return NewIntroScene(context)
		},
		SceneMenu: func() navigation.Scene {
			return NewMenuScene(context)
		},
		SceneLevels: func() navigation.Scene {
			return NewLevelsScene(context)
		},
		SceneSummary: func() navigation.Scene {
			return NewSummaryScene(context)
		},
	}
	return sceneMap
}
