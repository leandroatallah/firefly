package gamesetup

import (
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	gamescenephases "github.com/boilerplate/ebiten-template/internal/game/scenes/phases"
	phaseskit "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases"
)

func GetPhases() []phases.Phase {
	return []phases.Phase{
		{
			Genre:        phaseskit.GenrePlatformer,
			SceneType:    gamescenephases.SceneTypeForGenre(phaseskit.GenrePlatformer),
			ID:           1,
			Name:         "Sample Phase",
			TilemapPath:  "assets/tilemap/phase-000.tmj",
			GoalType:     gamescenephases.ReactEndpointType,
			SequencePath: "assets/sequences/sample_phase.json",
		},
		{
			Genre:        phaseskit.GenreBeatemup,
			SceneType:    gamescenephases.SceneTypeForGenre(phaseskit.GenreBeatemup),
			ID:           2,
			Name:         "Sample Phase",
			TilemapPath:  "assets/tilemap/beatemup-000.tmj",
			GoalType:     gamescenephases.ReactEndpointType,
			SequencePath: "assets/sequences/sample_phase.json",
		},
	}
}
