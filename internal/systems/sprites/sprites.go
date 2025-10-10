package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/contracts/animation"
)

// ?
type SpriteEntity struct {
	sprites SpriteMap
}

func NewSpriteEntity(sprites SpriteMap) SpriteEntity {
	return SpriteEntity{sprites: sprites}
}

// GetFirstSprite returns the first sprite. Useful when have only one sprite.
func (s *SpriteEntity) GetFirstSprite() *ebiten.Image {
	for _, img := range s.sprites {
		return img
	}

	return nil
}

func (s *SpriteEntity) GetSpriteByState(state animation.SpriteState) *ebiten.Image {
	return s.sprites[state]
}

func (s *SpriteEntity) Sprites() SpriteMap {
	return s.sprites
}

func (s *SpriteEntity) AnimatedSpriteImage(sprite *ebiten.Image, rect image.Rectangle, count int, frameRate int) *ebiten.Image {
	frameOX, frameOY := 0, 0
	width := rect.Dx()
	height := rect.Dy()

	elementWidth := sprite.Bounds().Dx()
	frameCount := elementWidth / width
	i := (count / frameRate) % frameCount
	sx, sy := frameOX+i*width, frameOY

	return sprite.SubImage(
		image.Rect(sx, sy, sx+width, sy+height),
	).(*ebiten.Image)

}

// ?
type SpriteMap map[animation.SpriteState]*ebiten.Image

// ?
type SpriteAssets map[animation.SpriteState]string

func (s SpriteAssets) AddSprite(state animation.SpriteState, path string) SpriteAssets {
	if len(s) == 0 {
		s = make(SpriteAssets)
	}
	s[state] = path
	return s
}

func LoadSprites(list SpriteAssets) (SpriteMap, error) {
	res := make(map[animation.SpriteState]*ebiten.Image)
	var err error

	for state, path := range list {
		res[state], _, err = ebitenutil.NewImageFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
