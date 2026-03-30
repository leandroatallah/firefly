package gamesetup

import (
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	gamescenephases "github.com/boilerplate/ebiten-template/internal/game/scenes/phases"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
)

func GetPhases() []phases.Phase {
	return []phases.Phase{
		{
			ID:           1,
			Name:         "Sample Phase",
			TilemapPath:  "assets/tilemap/phase-000.tmj",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.ScenePhases,
			SequencePath: "assets/sequences/sample_phase.json",
		},
	}
}
