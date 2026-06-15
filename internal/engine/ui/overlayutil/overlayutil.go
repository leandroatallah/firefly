// Package overlayutil provides shared rendering helpers for in-game overlays.
package overlayutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Base holds the open/close state shared by all overlays.
type Base struct{ open bool }

func (b *Base) Open()        { b.open = true }
func (b *Base) Close()       { b.open = false }
func (b *Base) IsOpen() bool { return b.open }

// DrawDimPanel fills screen with a semi-transparent black panel.
func DrawDimPanel(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	panel := ebiten.NewImage(w, h)
	panel.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(panel, nil)
}

// DrawText draws s at (x, y) on screen using face and c.
func DrawText(screen *ebiten.Image, face *text.GoTextFace, s string, x, y float64, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(screen, s, face, op)
}
