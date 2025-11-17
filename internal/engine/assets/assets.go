package assets

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

func LoadImageFromFs(ctx *core.AppContext, path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(ctx.Assets, path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
