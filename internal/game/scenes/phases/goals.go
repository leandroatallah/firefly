package gamescenephases

import (
	"time"
)

// ReachEndpointGoal: Complete when player reaches endpoint
type ReachEndpointGoal struct {
	scene *PhasesScene
}

func (g *ReachEndpointGoal) IsCompleted() bool {
	return g.scene.reachedEndpoint
}

func (g *ReachEndpointGoal) OnCompletion() {
	g.scene.freezeAllActors()
	g.scene.Audiomanager().FadeOutCurrentTrack(time.Second)
	g.scene.defaultCompletion()
}
