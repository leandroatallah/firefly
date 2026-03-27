package vfx

import (
	"image/color"

	contractvfx "github.com/leandroatallah/firefly/internal/engine/contracts/vfx"
	enginevfx "github.com/leandroatallah/firefly/internal/engine/render/vfx"
)

// SpawnFreezeAuraParticles spawns blue aura particles for the freeze power effect.
func SpawnFreezeAuraParticles(m contractvfx.Manager, x, y float64, count int) {
	if m == nil {
		return
	}

	// Blue color palette
	colors := []color.RGBA{
		{50, 100, 255, 255},  // Bright blue
		{100, 150, 255, 255}, // Light blue
		{50, 150, 255, 255},  // Sky blue
	}

	enginevfx.SpawnAuraParticles(m, x, y, colors, count, 16.0, 1.5)
}

// SpawnGrowAuraParticles spawns orange aura particles for the grow power effect.
func SpawnGrowAuraParticles(m contractvfx.Manager, x, y float64, count int) {
	if m == nil {
		return
	}

	// Orange color palette
	colors := []color.RGBA{
		{255, 140, 0, 255}, // Orange
		{255, 165, 0, 255}, // Standard orange
		{255, 100, 0, 255}, // Deep orange
	}

	// NOTE: For now, use only white
	colors = []color.RGBA{
		{255, 255, 255, 255},
	}

	enginevfx.SpawnAuraParticles(m, x, y, colors, count, 16.0, 1.5)
}

// SpawnStarParticles spawns rainbow aura particles for the star power effect.
func SpawnStarParticles(m contractvfx.Manager, x, y float64, count int) {
	if m == nil {
		return
	}

	// Vibrant rainbow colors
	colors := []color.RGBA{
		{255, 50, 50, 255},  // Red
		{50, 255, 50, 255},  // Green
		{50, 50, 255, 255},  // Blue
		{255, 255, 50, 255}, // Yellow
		{255, 50, 255, 255}, // Magenta
		{50, 255, 255, 255}, // Cyan
	}

	enginevfx.SpawnAuraParticles(m, x, y, colors, count, 18.0, 1.8)
}
