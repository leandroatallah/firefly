package transition

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type Fader struct {
	BaseTransition
	alpha           float64
	waitFrames      int
	holdDuration    time.Duration // wait with black screen BEFORE callback
	fadeSpeed       float64
	visibleDuration time.Duration // wait with black screen AFTER callback (before fade in)
}

func NewFader(holdDuration, visibleDuration time.Duration) *Fader {
	return &Fader{
		holdDuration:    holdDuration,
		visibleDuration: visibleDuration,
		fadeSpeed:       15, // default: ~17 frames for full fade (255/15)
	}
}

// Transition methods
func (f *Fader) Update() {
	if !f.active {
		return
	}

	if f.exiting {
		f.alpha += f.fadeSpeed
		if f.alpha >= 255 {
			f.alpha = 255
			f.exiting = false
			// Call callback immediately when fade-out completes
			// This ensures scene.OnStart() runs and sequences start instantly
			if f.onExitCb != nil {
				f.onExitCb()
			}
			f.waitFrames = 0
			// State transition handled below
		}
		return
	}

	// State transition from black screen to hold/visible/starting
	if !f.exiting && !f.starting && f.alpha == 255 {
		if f.holdDuration > 0 {
			f.waitFrames++
			if timing.ToDuration(f.waitFrames) >= f.holdDuration {
				// Hold complete, move to next stage
				if f.visibleDuration > 0 {
					f.waitFrames = 0 // Wait for visible duration
					return
				}
				f.starting = true
				return
			}
			return
		}

		if f.visibleDuration > 0 {
			f.waitFrames++
			if timing.ToDuration(f.waitFrames) >= f.visibleDuration {
				f.starting = true
			}
			return
		}

		// No hold or visible duration, start fade in immediately
		f.starting = true
		return
	}

	if f.starting {
		f.alpha -= f.fadeSpeed
		if f.alpha <= 0 {
			f.alpha = 0
			f.starting = false
			f.active = false
		}
		return
	}
}

func (f *Fader) Draw(screen *ebiten.Image) {
	if !f.active {
		return
	}
	c := color.RGBA{A: uint8(f.alpha)}
	img := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)
	img.Fill(c)
	screen.DrawImage(img, nil)
}

func (f *Fader) StartTransition(cb func()) {
	f.fadeOut(func() {
		f.fadeIn(cb)
	})
}

func (f *Fader) EndTransition(cb func()) {}

// Custom methods
func (f *Fader) fadeOut(cb func()) {
	if f.active {
		return
	}
	f.active = true
	f.exiting = true
	f.alpha = 0
	f.onExitCb = cb
}

func (f *Fader) fadeIn(cb func()) {
	// fadeIn now only calls the callback. Update handles setting f.starting = true
	// so that holdDuration and visibleDuration are respected.
	cb()
}
