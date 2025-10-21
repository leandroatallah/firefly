package actors

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

type ShapeRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (s *ShapeRect) Rect() (x, y, width, height int) {
	return s.X, s.Y, s.Width, s.Height
}

type SpriteData struct {
	BodyRect      ShapeRect         `json:"body_rect"`
	CollisionRect ShapeRect         `json:"collision_rect"`
	Assets        map[string]string `json:"assets"`
	FrameRate     int               `json:"frame_rate"`
	// TODO: Make facing direction work
	FacingDirection physics.FacingDirectionEnum `json:"facing_direction"` // 0 - right, 1 - left
}

type StatData struct {
	Health   int `json:"health"`
	Speed    int `json:"speed"`
	MaxSpeed int `json:"max_speed"`
}

type PlayerData struct {
	SpriteData SpriteData `json:"sprites"`
	StatData   StatData   `json:"stats"`
}

func ParseJsonPlayer(path string) (SpriteData, StatData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SpriteData{}, StatData{}, err
	}

	var playerData PlayerData
	if err := json.Unmarshal(data, &playerData); err != nil {
		return SpriteData{}, StatData{}, err
	}

	return playerData.SpriteData, playerData.StatData, nil
}
