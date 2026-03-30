package screenutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
)

func DrawScreenFlash(screen *ebiten.Image) {
	cfg := config.Get()
	bg := ebiten.NewImage(cfg.ScreenWidth, cfg.ScreenHeight)
	bg.Fill(color.RGBA{255, 255, 255, 255})
	screen.DrawImage(bg, nil)
}
