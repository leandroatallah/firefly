package phases

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/sequences"

// Goal defines the interface for phase completion criteria
type Goal interface {
	IsCompleted() bool
	OnCompletion()
}

// GoalType constants for identifying the completion criteria of a phase.
// These live in the engine package so both kit and game layers can reference them.
//
//nolint:gochecknoglobals
var (
	ReactEndpointType GoalType = "reach_endpoint"
	SequenceGoalType  GoalType = "sequence"
	NoGoalType        GoalType = "no_goal"
)

// SequenceGoal: Complete when sequence finishes
type SequenceGoal struct {
	Player         sequences.Player
	OnCompleteFunc func()
}

func (g *SequenceGoal) IsCompleted() bool {
	return g.Player != nil && !g.Player.IsPlaying()
}

func (g *SequenceGoal) OnCompletion() {
	if g.OnCompleteFunc != nil {
		g.OnCompleteFunc()
	}
}

// NoGoal: Never completes
type NoGoal struct{}

func (g *NoGoal) IsCompleted() bool {
	return false
}

func (g *NoGoal) OnCompletion() {}

// ReachEndpointGoal completes when a flag is flipped via Reach().
// The optional OnCompletion_ callback is invoked by OnCompletion when set.
type ReachEndpointGoal struct {
	reached       bool
	OnCompletion_ func() // optional callback (e.g., game-layer freeze/audio fade)
}

func (g *ReachEndpointGoal) IsCompleted() bool { return g.reached }
func (g *ReachEndpointGoal) OnCompletion() {
	if g.OnCompletion_ != nil {
		g.OnCompletion_()
	}
}
func (g *ReachEndpointGoal) Reach() { g.reached = true }
