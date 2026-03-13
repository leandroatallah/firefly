package gamescene

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	gamescenephases "github.com/leandroatallah/firefly/internal/game/scenes/phases"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

func InitSceneMap(context *app.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		scenestypes.SceneIntro: func() navigation.Scene {
			return NewIntroScene(context)
		},
		scenestypes.ScenePhaseTitle: func() navigation.Scene {
			return NewPhaseTitleScene(context)
		},
		scenestypes.SceneMenu: func() navigation.Scene {
			return NewMenuScene(context)
		},
		scenestypes.SceneStory: func() navigation.Scene {
			return NewStoryScene(context)
		},
		scenestypes.ScenePhases: func() navigation.Scene {
			return gamescenephases.NewPhasesScene(context)
		},
		scenestypes.SceneSummary: func() navigation.Scene {
			return NewSummaryScene(context)
		},
		scenestypes.ScenePhaseReboot: func() navigation.Scene {
			return NewPhaseRebootScene(context)
		},
	}
	return sceneMap
}
