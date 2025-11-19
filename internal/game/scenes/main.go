package gamescene

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
	gamescenelevels "github.com/leandroatallah/firefly/internal/game/scenes/levels"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

func InitSceneMap(context *core.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		scenestypes.SceneIntro: func() navigation.Scene {
			return NewIntroScene(context)
		},
		scenestypes.SceneMenu: func() navigation.Scene {
			return NewMenuScene(context)
		},
		scenestypes.SceneLevels: func() navigation.Scene {
			return gamescenelevels.NewLevelsScene(context)
		},
		scenestypes.SceneSummary: func() navigation.Scene {
			return NewSummaryScene(context)
		},
	}
	return sceneMap
}
