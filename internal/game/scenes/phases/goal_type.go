package gamescenephases

import "github.com/boilerplate/ebiten-template/internal/engine/scene/phases"

// State enum: part of engine public API
//
//nolint:gochecknoglobals
var (
	ReactEndpointType phases.GoalType = "reach_endpoint"
	SequenceGoalType  phases.GoalType = "sequence"
	NoGoalType        phases.GoalType = "no_goal"
)
