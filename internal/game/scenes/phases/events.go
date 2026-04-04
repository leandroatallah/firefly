package gamescenephases

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	actorevents "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
)

func subscribeEvents(ctx *app.AppContext, _ *PhasesScene) {

	em := ctx.EventManager

	// Common events
	em.Subscribe(actorevents.ActorJumpedType, func(e event.Event) {
		if ctx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorJumpedEvent); ok {
			yOffset := 1.0
			ctx.VFX.SpawnJumpPuff(evt.X, evt.Y+yOffset, 1)
			ctx.AudioManager.PlaySoundAtVolume("assets/audio/Menu_Select.ogg", 0.3)
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
