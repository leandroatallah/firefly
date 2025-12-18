package sprites

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/schemas"
)

// GetSpritesFromAssets converts asset data from a JSON schema into a SpriteMap,
// using a provided mapping from string keys to sprite states.
func GetSpritesFromAssets(assets map[string]schemas.AssetData, stateMap map[string]animation.SpriteState) (SpriteMap, error) {
	s := make(SpriteAssets)
	for key, value := range assets {
		if state, ok := stateMap[key]; ok {
			s = s.AddSprite(state, value.Path)
		}
	}
	return LoadSprites(s)
}
