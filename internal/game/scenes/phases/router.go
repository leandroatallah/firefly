package gamescenephases

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
	phaseskit "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases"
)

// SceneTypeForGenre maps a phase Genre to the navigation SceneType used by the
// scene factory. Panics on unknown genres so a misconfigured phase is caught early.
func SceneTypeForGenre(g phases.Genre) navigation.SceneType {
	switch g {
	case phaseskit.GenrePlatformer:
		return scenestypes.ScenePlatformerPhase
	case phaseskit.GenreBeatemup:
		return scenestypes.SceneBeatemupPhase
	}
	panic(fmt.Sprintf("SceneTypeForGenre: unknown genre %d", g))
}
