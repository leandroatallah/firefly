package vfx

import (
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles"
	"github.com/hajimehoshi/ebiten/v2"
)

// Manager coordinates all visual effects including particles, screen shake, and floating text.
type Manager interface {
	// SetAppContext injects the shared application context.
	SetAppContext(appContext any)

	// Update advances all active effects by one game tick.
	Update()
	// Draw renders all active effects to the screen using the given camera.
	Draw(screen *ebiten.Image, cam *camera.Controller)

	// AddParticle registers a new particle to be updated and drawn.
	AddParticle(p *particles.Particle)
	// AddTrauma adds screen-shake trauma to the given camera.
	AddTrauma(cam *camera.Controller, amount float64)
	// PixelConfig returns the default particle configuration.
	PixelConfig() *particles.Config
	// SetDefaultFont sets the font used for floating text effects.
	SetDefaultFont(f *font.FontText)

	// Effects

	// SpawnDeathExplosion spawns an explosion particle burst at the given position.
	SpawnDeathExplosion(x float64, y float64, count int)
	// SpawnFallingRocks spawns falling rock particles across the given width.
	SpawnFallingRocks(x float64, y float64, width float64, count int)
	// SpawnFloatingText spawns a floating text label at the given position for the given duration.
	SpawnFloatingText(msg string, x float64, y float64, duration int)
	// SpawnFloatingTextAbove spawns a floating text label above the given body.
	SpawnFloatingTextAbove(target body.Body, msg string, duration int)
	// SpawnJumpPuff spawns a small puff of particles for a jump action.
	SpawnJumpPuff(x float64, y float64, count int)
	// SpawnLandingPuff spawns a small puff of particles for a landing action.
	SpawnLandingPuff(x float64, y float64, count int)
	// SpawnPuff spawns a named puff effect at the given position with optional random spread.
	SpawnPuff(typeKey string, x float64, y float64, count int, randRange float64)
	// SpawnDirectionalPuff spawns a named puff effect anchored to extend outward from the spawn point.
	// faceRight=true anchors the sprite left edge at x; faceRight=false anchors the right edge.
	SpawnDirectionalPuff(typeKey string, x float64, y float64, faceRight bool, count int, randRange float64)
	// TriggerScreenFlash triggers a full-screen flash effect.
	TriggerScreenFlash()
}
