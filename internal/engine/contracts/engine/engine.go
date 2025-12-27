package engine

import "github.com/hajimehoshi/ebiten/v2"

// TODO: Remove if unused
type Game interface {
	Draw(screen *ebiten.Image)
	Layout(outWidth int, outHeight int) (screenWidth int, screenHeight int)
	Update() error
}
