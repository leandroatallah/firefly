package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO: Should it be a package?
type SpriteEntity struct {
	sprites    SpriteMap
	isMirrored bool
}

type SpriteMap map[ActorStateEnum]*ebiten.Image

func NewSpriteEntity(sprites SpriteMap) SpriteEntity {
	return SpriteEntity{sprites: sprites}
}

type SpriteAssets map[ActorStateEnum]string

func (s SpriteAssets) AddSprite(state ActorStateEnum, path string) SpriteAssets {
	if len(s) == 0 {
		s = make(SpriteAssets)
	}
	s[state] = path
	return s
}

func LoadSprites(list SpriteAssets) (SpriteMap, error) {
	res := make(map[ActorStateEnum]*ebiten.Image)
	var err error

	for state, path := range list {
		res[state], _, err = ebitenutil.NewImageFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *SpriteEntity) SetIsMirrored(value bool) {
	s.isMirrored = value
}
