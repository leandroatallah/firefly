package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	actorevents "github.com/leandroatallah/firefly/internal/engine/entity/actors/events"
	"github.com/leandroatallah/firefly/internal/engine/event"
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
		}
	})
	em.Subscribe(actorevents.ActorLandedType, func(e event.Event) {
		if ctx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorLandedEvent); ok {
			yOffset := 1.0
			ctx.VFX.SpawnLandingPuff(evt.X, evt.Y+yOffset, 1)
		}
	})
}
