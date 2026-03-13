package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Controller) CamDebug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		c.Kamera().Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		c.Kamera().Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		c.Kamera().Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		c.Kamera().ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		c.Kamera().ZoomFactor *= 1.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		c.Kamera().CenterOffsetX *= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		c.Kamera().CenterOffsetX /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		c.Kamera().CenterOffsetY /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		c.Kamera().CenterOffsetY *= 1.02
	}
}
