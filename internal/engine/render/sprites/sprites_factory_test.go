package sprites

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
)

func TestGetSpritesFromAssetsSmoke(t *testing.T) {
	// Since GetSpritesFromAssets calls LoadSprites which hits the disk,
	// we test the mapping logic but expect an error unless we have a real image.
	
	assets := map[string]schemas.AssetData{
		"idle": {Path: "non_existent.png"},
	}
	stateMap := map[string]animation.SpriteState{
		"idle": "idle_state",
	}

	_, err := GetSpritesFromAssetsOS(assets, stateMap)
	if err == nil {
		t.Error("expected error for non-existent image path")
	}
	
	// Test with no matching states (should be empty map but no error)
	emptyStateMap := map[string]animation.SpriteState{}
	res, err := GetSpritesFromAssetsOS(assets, emptyStateMap)
	if err != nil {
		t.Errorf("unexpected error for empty state map: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected empty result, got %v", res)
	}
}
