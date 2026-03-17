package vfx

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/particles"
	"github.com/leandroatallah/firefly/internal/engine/render/particles/vfx"
)

// SpawnStarParticles spawns rainbow aura particles for the star power effect.
func SpawnStarParticles(m *vfx.Manager, x, y float64, count int) {
	if m == nil {
		return
	}

	for i := 0; i < count; i++ {
		// Spread particles around the location
		rx := x + (rand.Float64()-0.5)*12
		ry := y + (rand.Float64()-0.5)*20

		// Vibrant rainbow colors
		colors := []color.RGBA{
			{255, 50, 50, 255},   // Red
			{50, 255, 50, 255},   // Green
			{50, 50, 255, 255},   // Blue
			{255, 255, 50, 255},  // Yellow
			{255, 50, 255, 255},  // Magenta
			{50, 255, 255, 255},  // Cyan
		}
		c := colors[rand.Intn(len(colors))]

		// Aura-like velocity: moving upward initially
		velX := (rand.Float64() - 0.5) * 1.5
		velY := -(1.0 + rand.Float64()*1.5)

		p := &particles.Particle{
			X:           rx,
			Y:           ry,
			VelX:        velX,
			VelY:        velY,
			AccY:        0.08, // Gravity will make them fall after the initial burst
			Duration:    20 + rand.Intn(15),
			MaxDuration: 40,
			Scale:       1.5 + rand.Float64()*1.0,
			ScaleSpeed:  -0.05, // Shrink over time
			Config:      m.PixelConfig(),
		}
		p.ColorScale.ScaleWithColor(c)
		m.AddParticle(p)
	}
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
