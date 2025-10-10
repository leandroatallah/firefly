package screenutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/config"
)

func GetCenterOfScreenPosition(width, height int) (int, int) {
	x := config.ScreenWidth/2 - width/2
	y := config.ScreenHeight/2 - height/2
	return x, y
}

func DrawCenteredText(screen *ebiten.Image, fontText *font.FontText, str string, size float64, c color.Color) {
	textOp := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}

	textOp.GeoM.Translate(config.ScreenWidth/2, config.ScreenHeight/2)
	textOp.ColorScale.ScaleWithColor(c)

	fontText.Draw(screen, str, size, textOp)
}

