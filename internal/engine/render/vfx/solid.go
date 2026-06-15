package vfx

import (
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// SolidColor shows a solid color over a given number of frames.
type SolidColor struct {
	animating bool
	duration  int
	color     color.RGBA
}

// NewSolidColor creates a new solid color overlay (inactive by default).
func NewSolidColor() *SolidColor {
	return &SolidColor{
		color: color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

// FadeOut starts the overlay.
func (f *SolidColor) FadeOut(frames int) {
	f.duration = frames
	f.animating = true
}

// FadeIn is a no-op for SolidColor (no fade-in concept).
func (f *SolidColor) FadeIn(_ int) {}

// Update returns true when the duration ends.
func (f *SolidColor) Update() bool {
	if !f.IsActive() {
		return true
	}

	f.duration--
	if f.duration <= 0 {
		f.animating = false
		return true
	}
	return false
}

// IsActive returns true if the fade is currently animating.
func (f *SolidColor) IsActive() bool {
	return f.animating
}

// Reset clears the overlay completely.
func (f *SolidColor) Reset() {
	f.animating = false
}

// Draw renders the color overlay to the screen.
func (f *SolidColor) Draw(screen *ebiten.Image) {
	if !f.IsActive() {
		return
	}

	cfg := config.Get()
	img := ebiten.NewImage(cfg.ScreenWidth, cfg.ScreenHeight)
	rgba := f.color
	img.Fill(rgba)
	screen.DrawImage(img, nil)
}

// SetColor define the base color
func (f *SolidColor) SetColor(c color.Color) {
	if nc, ok := c.(color.RGBA); ok {
		f.color = nc
	}
}
