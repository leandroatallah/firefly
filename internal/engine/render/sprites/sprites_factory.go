package sprites

import (
	"io/fs"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
)

// GetSpritesFromAssets converts asset data from a JSON schema into a SpriteMap,
// using a provided mapping from string keys to sprite states.
func GetSpritesFromAssets(fsys fs.FS, assets map[string]schemas.AssetData, stateMap map[string]animation.SpriteState) (SpriteMap, error) {
	s := make(SpriteAssets)
	for key, value := range assets {
		if state, ok := stateMap[key]; ok {
			loop := true // Default to true
			if value.Loop != nil {
				loop = *value.Loop
			}
			s = s.AddSprite(state, value.Path, loop)
		}
	}
	return LoadSprites(fsys, s)
}

// GetSpritesFromAssetsOS is a convenience wrapper using the OS filesystem.
func GetSpritesFromAssetsOS(assets map[string]schemas.AssetData, stateMap map[string]animation.SpriteState) (SpriteMap, error) {
	return GetSpritesFromAssets(os.DirFS("."), assets, stateMap)
}
