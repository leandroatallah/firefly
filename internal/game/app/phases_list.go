package gamesetup

import (
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	gamescenephases "github.com/leandroatallah/firefly/internal/game/scenes/phases"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

func GetPhases() []phases.Phase {
	return []phases.Phase{
		{
			ID:           1,
			Name:         "Area 1 - Phase 1",
			TilemapPath:  "assets/tilemap/phase-000.tmj",
			NextPhaseID:  2,
			SequencePath: "assets/sequences/phase-start.json",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.ScenePhases,
		},
	}
}
