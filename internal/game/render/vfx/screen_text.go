package vfx

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/vfx/text"
)

// ScreenText displays text at fixed screen coordinates.
type ScreenText struct {
	*text.SimpleFloatingText
}

// NewScreenText creates text at fixed screen position.
func NewScreenText(x, y float64, msg string, duration int) *ScreenText {
	return &ScreenText{
		SimpleFloatingText: text.NewFloatingText(msg, x, y, duration),
	}
}

func (st *ScreenText) Draw(screen *ebiten.Image, cam *camera.Controller) {
	// Screen coordinates - no camera transformation
	st.SimpleFloatingText.Draw(screen, nil)
}
