package assets

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/app"
)

func LoadImageFromFs(ctx *app.AppContext, path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(ctx.Assets, path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
