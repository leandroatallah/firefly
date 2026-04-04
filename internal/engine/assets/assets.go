package assets

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadImageFromFs(ctx *app.AppContext, path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(ctx.Assets, path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
