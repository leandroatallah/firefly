package phases

import "github.com/leandroatallah/firefly/internal/engine/contracts/sequences"

// Goal defines the interface for phase completion criteria
type Goal interface {
	IsCompleted() bool
	OnCompletion()
}

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
