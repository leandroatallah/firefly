package vfx

import (
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles"
	"github.com/hajimehoshi/ebiten/v2"
)

type Manager interface {
	SetAppContext(appContext any)

	Update()
	Draw(screen *ebiten.Image, cam *camera.Controller)

	AddParticle(p *particles.Particle)
	AddTrauma(cam *camera.Controller, amount float64)
	PixelConfig() *particles.Config
	SetDefaultFont(f *font.FontText)

	// Effects
	SpawnDeathExplosion(x float64, y float64, count int)
	SpawnFallingRocks(x float64, y float64, width float64, count int)
	SpawnFloatingText(msg string, x float64, y float64, duration int)
	SpawnFloatingTextAbove(actor actors.ActorEntity, msg string, duration int)
	SpawnJumpPuff(x float64, y float64, count int)
	SpawnLandingPuff(x float64, y float64, count int)
	SpawnPuff(typeKey string, x float64, y float64, count int, randRange float64)
	TriggerScreenFlash()
}
