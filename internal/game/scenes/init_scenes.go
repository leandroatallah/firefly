package gamescene

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	gamescenephases "github.com/boilerplate/ebiten-template/internal/game/scenes/phases"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
)

func InitSceneMap(ctx *app.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		scenestypes.SceneIntro: func() navigation.Scene {
			return NewIntroScene(ctx)
		},
		scenestypes.ScenePhaseTitle: func() navigation.Scene {
			return NewPhaseTitleScene(ctx)
		},
		scenestypes.SceneMenu: func() navigation.Scene {
			return NewMenuScene(ctx)
		},
		scenestypes.SceneStory: func() navigation.Scene {
			return NewStoryScene(ctx)
		},
		scenestypes.ScenePhases: func() navigation.Scene {
			return gamescenephases.NewPhasesScene(ctx)
		},
		scenestypes.SceneSummary: func() navigation.Scene {
			return NewSummaryScene(ctx)
		},
		scenestypes.ScenePhaseReboot: func() navigation.Scene {
			return NewPhaseRebootScene(ctx)
		},
		scenestypes.SceneCredits: func() navigation.Scene {
			return NewCreditsScene(ctx)
		},
	}
	return sceneMap
}
