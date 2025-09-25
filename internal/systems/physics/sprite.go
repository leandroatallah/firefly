package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO: Should it be a package?
type SpriteEntity struct {
	sprites    SpriteMap
	isMirrored bool
}

type SpriteMap map[CharacterState]*ebiten.Image

func NewSpriteEntity(sprites SpriteMap) SpriteEntity {
	return SpriteEntity{sprites: sprites}
}

type spriteAssets map[CharacterState]string

func (s spriteAssets) addSprite(state CharacterState, path string) spriteAssets {
	if len(s) == 0 {
		s = make(spriteAssets)
	}
	s[state] = path
	return s
}

func loadSprites(list spriteAssets) (SpriteMap, error) {
	res := make(map[CharacterState]*ebiten.Image)
	var err error

	for state, path := range list {
		res[state], _, err = ebitenutil.NewImageFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
