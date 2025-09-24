package transition

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
)

const (
	transitionSpeed = 15
)

type Fader struct {
	BaseTransition
	alpha float64
}

func NewFader() *Fader {
	return &Fader{}
}

// Transition methods
func (f *Fader) Update() {
	if !f.active {
		return
	}

	if f.exiting {
		f.alpha += transitionSpeed
		if f.alpha >= 255 {
			f.alpha = 255
			f.exiting = false
			if f.onExitCb != nil {
				f.onExitCb()
			}
		}
	} else if f.starting {
		f.alpha -= transitionSpeed
		if f.alpha <= 0 {
			f.alpha = 0
			f.starting = false
			f.active = false
		}
	}
}

func (f *Fader) Draw(screen *ebiten.Image) {
	if !f.active {
		return
	}
	c := color.RGBA{A: uint8(f.alpha)}
	img := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
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
	f.starting = true
	cb()
}
