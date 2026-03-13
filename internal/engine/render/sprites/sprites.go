package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
)

type Sprite struct {
	Image *ebiten.Image
	Loop  bool
}

type SpriteEntity struct {
	sprites   SpriteMap
	frameRate int
}

func NewSpriteEntity(sprites SpriteMap) SpriteEntity {
	return SpriteEntity{sprites: sprites, frameRate: 1} // Default frame rate to 1
}

// GetFirstSprite returns the first sprite. Useful when have only one sprite.
func (s *SpriteEntity) GetFirstSprite() *Sprite {
	for _, sprite := range s.sprites {
		return sprite
	}

	return nil
}

func (s *SpriteEntity) GetSpriteByState(state animation.SpriteState) *Sprite {
	return s.sprites[state]
}

func (s *SpriteEntity) Sprites() SpriteMap {
	return s.sprites
}

func (s *SpriteEntity) AnimatedSpriteImage(sprite *Sprite, rect image.Rectangle, count int, frameRate int) *ebiten.Image {
	if sprite == nil || sprite.Image == nil {
		return nil
	}

	frameOX, frameOY := 0, 0
	width := rect.Dx()
	height := rect.Dy()

	elementWidth := sprite.Image.Bounds().Dx()

	if width <= 0 {
		return sprite.Image
	}

	frameCount := elementWidth / width
	if frameCount <= 1 {
		return sprite.Image
	}

	frameNum := count / frameRate
	if sprite.Loop {
		frameNum = frameNum % frameCount
	} else {
		if frameNum >= frameCount {
			frameNum = frameCount - 1
		}
	}

	sx, sy := frameOX+frameNum*width, frameOY

	return sprite.Image.SubImage(
		image.Rect(sx, sy, sx+width, sy+height),
	).(*ebiten.Image)
}

func (s *SpriteEntity) SetFrameRate(value int) {
	s.frameRate = value
}

func (s *SpriteEntity) FrameRate() int {
	return s.frameRate
}

// SpriteMap represents a collection of Ebiten images, keyed by their animation state.
type SpriteMap map[animation.SpriteState]*Sprite

type AssetInfo struct {
	Path string
	Loop bool
}

// SpriteAssets maps an animation state to the file path of the corresponding sprite image.
type SpriteAssets map[animation.SpriteState]AssetInfo

func (s SpriteAssets) AddSprite(state animation.SpriteState, path string, loop bool) SpriteAssets {
	if len(s) == 0 {
		s = make(SpriteAssets)
	}
	s[state] = AssetInfo{Path: path, Loop: loop}
	return s
}

func LoadSprites(list SpriteAssets) (SpriteMap, error) {
	res := make(map[animation.SpriteState]*Sprite)

	for state, assetInfo := range list {
		img, _, err := ebitenutil.NewImageFromFile(assetInfo.Path)
		if err != nil {
			return nil, err
		}
		res[state] = &Sprite{Image: img, Loop: assetInfo.Loop}
	}

	return res, nil
}
