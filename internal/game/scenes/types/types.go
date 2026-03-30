package scenestypes

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"

const (
	SceneIntro navigation.SceneType = iota
	SceneMenu
	ScenePhases
	SceneSummary
	ScenePhaseReboot
	SceneStory
	ScenePhaseTitle
	SceneCredits
)
