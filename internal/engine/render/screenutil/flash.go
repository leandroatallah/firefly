package screenutil

import (
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawScreenFlash(screen *ebiten.Image) {
	cfg := config.Get()
	bg := ebiten.NewImage(cfg.ScreenWidth, cfg.ScreenHeight)
	bg.Fill(color.RGBA{255, 255, 255, 255})
	screen.DrawImage(bg, nil)
}
