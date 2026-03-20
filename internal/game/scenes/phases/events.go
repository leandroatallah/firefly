package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	actorevents "github.com/leandroatallah/firefly/internal/engine/entity/actors/events"
	"github.com/leandroatallah/firefly/internal/engine/event"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
)

func subscribeEvents(ctx *app.AppContext, scene *PhasesScene) {

	em := ctx.EventManager

	// Common events
	em.Subscribe(actorevents.ActorJumpedType, func(e event.Event) {
		if ctx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorJumpedEvent); ok {
			yOffset := 1.0
			ctx.VFX.SpawnJumpPuff(evt.X, evt.Y+yOffset, 1)
			ctx.AudioManager.PlaySoundAtVolume("assets/audio/Jump.ogg", 0.2)
		}
	})
	em.Subscribe(actorevents.ActorLandedType, func(e event.Event) {
		if ctx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorLandedEvent); ok {
			yOffset := 1.0
			ctx.VFX.SpawnLandingPuff(evt.X, evt.Y+yOffset, 1)

			// Trigger screen shake if player is grown
			if player, found := ctx.ActorManager.GetPlayer(); found {
				if climber, ok := player.(*gameplayer.ClimberPlayer); ok && climber.IsGrowActive() {
					scene.Camera().Base().AddTrauma(0.3)
				}
			}
		}
	})
}
