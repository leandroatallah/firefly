package engine

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Draw(screen *ebiten.Image)
	Layout(outWidth int, outHeight int) (screenWidth int, screenHeight int)
	Update() error
}
