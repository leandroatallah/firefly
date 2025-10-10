package hud

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type HUD interface {
	Draw(screen *ebiten.Image)
	Update() error
}

// Create base HUD
