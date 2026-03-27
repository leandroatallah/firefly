package vfx

import (
	"image/color"
	"math/rand"

	contractvfx "github.com/leandroatallah/firefly/internal/engine/contracts/vfx"
	"github.com/leandroatallah/firefly/internal/engine/render/particles"
)

// SpawnAuraParticles spawns flame-shaped aura particles with a wave-like motion.
// Particles initially move outward from center, then flow inward as they rise.
// colors: palette to choose from for each particle
// count: number of particles to spawn
// width: horizontal spread width from center
// speed: upward velocity magnitude
func SpawnAuraParticles(m contractvfx.Manager, x, y float64, colors []color.RGBA, count int, width float64, speed float64) {
	if m == nil {
		return
	}

	for i := 0; i < count; i++ {
		// Spread particles horizontally around the center (narrow spread)
		rx := x + (rand.Float64()-0.5)*width*0.6
		ry := y + (rand.Float64()-0.5)*width*0.3

		// Initial outward velocity from center
		direction := -1.0
		if rx > x {
			direction = 1.0
		}
		initialVelX := direction * (0.3 + rand.Float64()*0.3)

		// Upward velocity
		velY := -speed - rand.Float64()*speed*0.3

		// Choose color from palette
		c := colors[rand.Intn(len(colors))]

		p := &particles.Particle{
			X:           rx,
			Y:           ry,
			VelX:        initialVelX,
			VelY:        velY,
			AccX:        -direction * 0.04, // Acceleration inward (creates wave: out → slow → in)
			AccY:        -0.02,             // Slight upward acceleration for flame effect
			Duration:    12 + rand.Intn(8),
			MaxDuration: 22,
			Scale:       1.0 + rand.Float64()*0.5,
			ScaleSpeed:  -0.04, // Shrink faster over time
			Config:      m.PixelConfig(),
		}
		p.ColorScale.ScaleWithColor(c)
		m.AddParticle(p)
	}
}
