package pause

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type PauseScreen struct {
	isPaused   bool
	disable    bool
	count      int
	key        ebiten.Key
	disableFor time.Duration

	onStart  func(p *PauseScreen)
	onFinish func(p *PauseScreen)
}

func NewPauseScreen(key ebiten.Key, disableFor time.Duration) *PauseScreen {
	return &PauseScreen{
		key:        key,
		disableFor: disableFor,
	}
}

func (p *PauseScreen) Update() {
	p.handlePause()

	if p.isPaused {
		p.count++
	}
}

func (p *PauseScreen) handlePause() {
	if p.disable && p.disableFor > 0 && p.count > timing.FromDuration(p.disableFor) {
		p.disable = false
	}

	if inpututil.IsKeyJustPressed(p.key) {
		p.Toggle()
	}
}

func (p *PauseScreen) Toggle() {
	if p.disable {
		return
	}

	p.isPaused = !p.isPaused
	p.count = 0

	if p.isPaused {
		if p.disableFor > 0 {
			p.disable = true
		}
		if p.onStart != nil {
			p.onStart(p)
		}
	} else {
		if p.onFinish != nil {
			p.onFinish(p)
		}
	}
}

func (p *PauseScreen) IsPaused() bool {
	return p.isPaused
}

func (p *PauseScreen) Count() int {
	return p.count
}

func (p *PauseScreen) SetOnStart(fn func(p *PauseScreen)) {
	p.onStart = fn
}

func (p *PauseScreen) SetOnFinish(fn func(p *PauseScreen)) {
	p.onFinish = fn
}

