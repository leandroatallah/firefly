package items

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// TODO: Duplicated
type ShapeRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// TODO: Duplicated
func (s ShapeRect) Rect() (x, y, width, height int) {
	return s.X, s.Y, s.Width, s.Height
}

// TODO: Duplicated
type AssetData struct {
	Path           string      `json:"path"`
	CollisionRects []ShapeRect `json:"collision_rect"`
}

// TODO: Duplicated
type SpriteData struct {
	BodyRect        ShapeRect                `json:"body_rect"`
	Assets          map[string]AssetData     `json:"assets"`
	FrameRate       int                      `json:"frame_rate"`
	FacingDirection body.FacingDirectionEnum `json:"facing_direction"` // 0 - right, 1 - left
}

type StatData struct {
	Id string `json:"id"`
}

type ItemData struct {
	SpriteData SpriteData `json:"sprites"`
	StatData   StatData   `json:"stats"`
}

func ParseJsonItem(path string) (SpriteData, StatData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SpriteData{}, StatData{}, err
	}

	var itemData ItemData
	if err := json.Unmarshal(data, &itemData); err != nil {
		return SpriteData{}, StatData{}, err
	}

	return itemData.SpriteData, itemData.StatData, nil
}
