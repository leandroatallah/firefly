package vfx

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	contractvfx "github.com/leandroatallah/firefly/internal/engine/contracts/vfx"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
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
		direction := 1.0
		if rx > x {
			direction = 1.0
		} else {
			direction = -1.0
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

	SpawnAuraParticles(m, x, y, colors, count, 16.0, 1.5)
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

	SpawnAuraParticles(m, x, y, colors, count, 16.0, 1.5)
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

	SpawnAuraParticles(m, x, y, colors, count, 18.0, 1.8)
}

// Vignette draws a retro darkness overlay with a jagged circular opening.
// The opening is computed per-screen pixel (1px blocks) so it naturally looks imperfect.
type Vignette struct {
	enabled  bool
	radiusPx float64

	overlay *ebiten.Image
	pixels  []byte
	w, h    int
}

func NewVignette() *Vignette {
	return &Vignette{
		enabled:  false,
		radiusPx: 48,
	}
}

func (v *Vignette) Enabled() bool {
	return v.enabled
}

// Enable turns the effect on and sets the opening radius (in screen pixels).
func (v *Vignette) Enable(radiusPx float64) {
	v.enabled = true
	if radiusPx > 0 {
		v.radiusPx = radiusPx
	}
}

func (v *Vignette) Disable() {
	v.enabled = false
}

func (v *Vignette) ensureOverlay(w, h int) {
	if v.overlay != nil && v.w == w && v.h == h {
		return
	}
	v.w, v.h = w, h
	v.overlay = ebiten.NewImage(w, h)
	v.pixels = make([]byte, 4*w*h)
}

func (v *Vignette) buildMaskPixels(cam *enginecamera.Controller, target body.Body, w, h int) []byte {
	v.ensureOverlay(w, h)

	if !v.enabled {
		for i := 0; i < len(v.pixels); i += 4 {
			v.pixels[i+0] = 0
			v.pixels[i+1] = 0
			v.pixels[i+2] = 0
			v.pixels[i+3] = 0
		}
		return v.pixels
	}

	// Convert target world center to screen coordinates based on camera.
	centerX, centerY := cam.GetActualCenter()
	topLeftX := centerX - float64(w)/2
	topLeftY := centerY - float64(h)/2

	tx, ty := target.GetPositionMin()
	tw, th := target.GetShape().Width(), target.GetShape().Height()
	targetCenterWorldX := float64(tx) + float64(tw)/2
	targetCenterWorldY := float64(ty) + float64(th)/2

	targetScreenX := targetCenterWorldX - topLeftX
	targetScreenY := targetCenterWorldY - topLeftY

	// Clamp radius to sane values.
	r := v.radiusPx
	if r < 1 {
		r = 1
	}
	maxR := math.Max(float64(w), float64(h)) * 2
	if r > maxR {
		r = maxR
	}
	rSq := r * r

	// Build an RGBA mask: black outside the circle, transparent inside.
	for y := 0; y < h; y++ {
		dy := float64(y) - targetScreenY
		row := y * w * 4
		for x := 0; x < w; x++ {
			dx := float64(x) - targetScreenX
			i := row + x*4
			if dx*dx+dy*dy > rSq {
				v.pixels[i+0] = 0
				v.pixels[i+1] = 0
				v.pixels[i+2] = 0
				v.pixels[i+3] = 0xff
			} else {
				v.pixels[i+0] = 0
				v.pixels[i+1] = 0
				v.pixels[i+2] = 0
				v.pixels[i+3] = 0x00
			}
		}
	}

	return v.pixels
}

// Draw renders the darkness overlay centered around the given body (typically the player).
// It draws in screen space, so it should be called after the world layer has been rendered.
func (v *Vignette) Draw(screen *ebiten.Image, cam *enginecamera.Controller, target body.Body) {
	if !v.enabled || screen == nil || cam == nil || target == nil {
		return
	}

	cfg := config.Get()
	if cfg == nil || cfg.ScreenWidth <= 0 || cfg.ScreenHeight <= 0 {
		return
	}

	w, h := cfg.ScreenWidth, cfg.ScreenHeight
	pixels := v.buildMaskPixels(cam, target, w, h)
	v.overlay.ReplacePixels(pixels)

	op := &ebiten.DrawImageOptions{}
	// Defensive: make sure the overlay is drawn as pure black regardless of blending artifacts.
	op.ColorScale.ScaleWithColor(color.Black)
	screen.DrawImage(v.overlay, op)
}
