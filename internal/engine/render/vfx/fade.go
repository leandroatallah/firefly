package vfx

import (
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// FadeOverlay fades the screen to a color over a given number of frames.
type FadeOverlay struct {
	animating bool // true while animating
	alpha     float64
	speed     float64
	direction int // +1 = fade-out (0→255), -1 = fade-in (255→0)
	color     color.RGBA
}

// NewFadeOverlay creates a new fade overlay (inactive by default).
func NewFadeOverlay() *FadeOverlay {
	return &FadeOverlay{
		animating: false,
		alpha:     0,
		speed:     0,
		color:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

// FadeOut starts a fade-out (alpha 0→255) over the given number of frames.
func (f *FadeOverlay) FadeOut(frames int) {
	if frames <= 0 {
		frames = 17
	}
	f.speed = 255.0 / float64(frames)
	f.alpha = 0
	f.direction = 1
	f.animating = true
}

// FadeIn starts a fade-in (alpha 255→0) over the given number of frames.
func (f *FadeOverlay) FadeIn(frames int) {
	if frames <= 0 {
		frames = 17
	}
	f.speed = 255.0 / float64(frames)
	f.alpha = 255
	f.direction = -1
	f.animating = true
}

// Update advances the alpha by one frame. Returns true when the animation completes.
// On fade-out, the overlay persists at alpha=255 until Reset() is called.
// On fade-in, the overlay clears (alpha=0) when done.
func (f *FadeOverlay) Update() bool {
	if !f.animating {
		return false
	}
	f.alpha += f.speed * float64(f.direction)
	if f.direction >= 0 && f.alpha >= 255 {
		f.alpha = 255
		f.animating = false
		return true
	}
	if f.direction < 0 && f.alpha <= 0 {
		f.alpha = 0
		f.animating = false
		return true
	}
	return false
}

// IsActive returns true if the fade is currently animating.
func (f *FadeOverlay) IsActive() bool {
	return f.animating
}

// IsPersisting returns true if fade is drawn (animating or at target).
func (f *FadeOverlay) IsPersisting() bool {
	return f.alpha > 0
}

// Reset clears the fade completely.
func (f *FadeOverlay) Reset() {
	f.animating = false
	f.alpha = 0
	f.speed = 0
}

// Draw renders the fade overlay to the screen.
func (f *FadeOverlay) Draw(screen *ebiten.Image) {
	if !f.IsPersisting() {
		return
	}
	cfg := config.Get()
	img := ebiten.NewImage(cfg.ScreenWidth, cfg.ScreenHeight)
	rgba := f.color
	rgba.A = uint8(f.alpha)
	img.Fill(rgba)
	screen.DrawImage(img, nil)
}

// Alpha returns the current alpha value (for testing).
func (f *FadeOverlay) Alpha() float64 {
	return f.alpha
}

// SetColor define the base color
func (f *FadeOverlay) SetColor(c color.Color) {
	if nc, ok := c.(color.RGBA); ok {
		f.color = nc
	}
}
