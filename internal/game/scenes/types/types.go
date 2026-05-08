package scenestypes

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"

const (
	SceneMenu navigation.SceneType = iota
	ScenePlatformerPhase
	SceneBeatemupPhase
	ScenePhaseReboot
)
