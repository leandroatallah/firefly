package gamescenephases

import (
	"time"
)

// ReachEndpointGoal completes the phase when the player reaches the level endpoint.
type ReachEndpointGoal struct {
	scene *PhasesScene
}

func (g *ReachEndpointGoal) IsCompleted() bool {
	return g.scene.reachedEndpoint
}

func (g *ReachEndpointGoal) OnCompletion() {
	g.scene.freezeAllActors()
	if g.scene.TilemapScene != nil {
		if ctx := g.scene.AppContext(); ctx != nil && ctx.AudioManager != nil {
			ctx.AudioManager.FadeOutCurrentTrack(time.Second)
		}
	}
	g.scene.defaultCompletion()
}
