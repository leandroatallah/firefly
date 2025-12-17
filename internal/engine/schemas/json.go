package schemas

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

// ShapeRect defines a rectangular shape with position and dimensions.
type ShapeRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Rect returns the coordinates and dimensions of the rectangle.
func (s ShapeRect) Rect() (x, y, width, height int) {
	return s.X, s.Y, s.Width, s.Height
}

// AssetData holds information about a single asset, including its path and collision areas.
type AssetData struct {
	Path           string      `json:"path"`
	CollisionRects []ShapeRect `json:"collision_rect"`
}

// SpriteData contains all data related to a sprite's appearance and behavior,
// including its body rectangle, assets for different states, animation frame rate, and initial facing direction.
type SpriteData struct {
	BodyRect        ShapeRect                `json:"body_rect"`
	Assets          map[string]AssetData     `json:"assets"`
	FrameRate       int                      `json:"frame_rate"`
	FacingDirection body.FacingDirectionEnum `json:"facing_direction"` // 0 - right, 1 - left
}
