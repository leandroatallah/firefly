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
			Name:         "Area 1 - Story",
			NextPhaseID:  2,
			SequencePath: "assets/sequences/area-1-story.json",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.SceneStory,
		},
		{
			ID:          2,
			Name:        "Area 1 - Title",
			Title:       "Area 1",
			NextPhaseID: 3,
			GoalType:    gamescenephases.NoGoalType,
			SceneType:   scenestypes.ScenePhaseTitle,
		},
		{
			ID:           3,
			Name:         "Area 1 - Phase 1",
			TilemapPath:  "assets/tilemap/phase-000.tmj",
			NextPhaseID:  4,
			SequencePath: "assets/sequences/phase-start.json",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.ScenePhases,
		},
		{
			ID:           4,
			Name:         "Area 1 - Phase 2",
			TilemapPath:  "assets/tilemap/phase-001.tmj",
			NextPhaseID:  5,
			SequencePath: "assets/sequences/phase-start.json",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.ScenePhases,
		},
		{
			ID:           5,
			Name:         "Area 1 - Phase 2",
			TilemapPath:  "assets/tilemap/phase-002.tmj",
			NextPhaseID:  6,
			SequencePath: "assets/sequences/phase-start.json",
			GoalType:     gamescenephases.ReactEndpointType,
			SceneType:    scenestypes.ScenePhases,
		},
		{
			ID:          6,
			Name:        "Credits",
			NextPhaseID: 1,
			GoalType:    gamescenephases.NoGoalType,
			SceneType:   scenestypes.SceneCredits,
		},
	}
}
