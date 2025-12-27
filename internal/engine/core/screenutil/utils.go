package screenutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
)

// TODO: Move
func GetCenterOfScreenPosition(width, height int) (int, int) {
	x := config.Get().ScreenWidth/2 - width/2
	y := config.Get().ScreenHeight/2 - height/2
	return x, y
}

// TODO: Move
func DrawCenteredText(screen *ebiten.Image, fontText *font.FontText, str string, size float64, c color.Color) {
	textOp := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}

	textOp.GeoM.Translate(
		float64(config.Get().ScreenWidth/2),
		float64(config.Get().ScreenHeight/2),
	)
	textOp.ColorScale.ScaleWithColor(c)

	fontText.Draw(screen, str, size, textOp)
}
