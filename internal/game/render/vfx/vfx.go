package vfx

import (
	"image/color"

	contractvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	enginevfx "github.com/boilerplate/ebiten-template/internal/engine/render/vfx"
)

// SpawnAuraParticles spawns aura particles with the specified colors.
func SpawnAuraParticles(m contractvfx.Manager, x, y float64, colors []color.RGBA, count int) {
	if m == nil {
		return
	}

	enginevfx.SpawnAuraParticles(m, x, y, colors, count, 16.0, 1.5)
}
