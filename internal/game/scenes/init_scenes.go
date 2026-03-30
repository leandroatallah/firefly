package gamescene

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	gamescenephases "github.com/boilerplate/ebiten-template/internal/game/scenes/phases"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
)

func InitSceneMap(ctx *app.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		scenestypes.SceneMenu: func() navigation.Scene {
			return NewMenuScene(ctx)
		},
		scenestypes.ScenePhases: func() navigation.Scene {
			return gamescenephases.NewPhasesScene(ctx)
		},
		scenestypes.ScenePhaseReboot: func() navigation.Scene {
			return NewPhaseRebootScene(ctx)
		},
	}
	return sceneMap
}
