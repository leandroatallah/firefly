package gamecamera

import (
	enginecamera "github.com/leandroatallah/firefly/internal/engine/camera"
)

// New creates a new camera controller with game-specific settings.
func New(x, y int) *enginecamera.Controller {
	cam := enginecamera.NewController(float64(x), float64(y))
	// Should it be removed?
	cam.DeadZoneRadius = 10.0
	cam.SmoothingFactor = 0.08
	return cam
}
