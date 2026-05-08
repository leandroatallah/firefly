package gamescene

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	gamebeatemupphase "github.com/boilerplate/ebiten-template/internal/game/scenes/phases/beatemup"
	gameplatformerphase "github.com/boilerplate/ebiten-template/internal/game/scenes/phases/platformer"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
)

func InitSceneMap(ctx *app.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		scenestypes.SceneMenu: func() navigation.Scene {
			return NewMenuScene(ctx)
		},
		scenestypes.ScenePlatformerPhase: func() navigation.Scene {
			return gameplatformerphase.NewPlatformerPhaseScene(ctx)
		},
		scenestypes.SceneBeatemupPhase: func() navigation.Scene {
			return gamebeatemupphase.NewBeatemupPhaseScene(ctx)
		},
		scenestypes.ScenePhaseReboot: func() navigation.Scene {
			return NewPhaseRebootScene(ctx)
		},
	}
	return sceneMap
}
