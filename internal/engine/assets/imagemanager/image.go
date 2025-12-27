package imagemanager

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageItem struct {
	name  string
	image *ebiten.Image
}

func (a *ImageItem) Name() string {
	return a.name
}
func (a *ImageItem) Image() *ebiten.Image {
	return a.image
}

type ImageManager struct {
	images map[string]*ebiten.Image
}

func NewImageManager() *ImageManager {
	return &ImageManager{
		images: make(map[string]*ebiten.Image),
	}
}

func (am *ImageManager) Add(name string, source *ebiten.Image) {
	am.images[name] = source
}

func (am *ImageManager) GetImage(name string) *ebiten.Image {
	img, ok := am.images[name]
	if !ok {
		log.Printf("image not found: %s", name)
		return nil
	}
	return img
}
