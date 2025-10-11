package physics

import (
	"image"

	"github.com/leandroatallah/firefly/internal/game/constants"
)

type Rect struct {
	x16, y16      int
	width, height int
}

func NewRect(x, y, width, height int) *Rect {
	return &Rect{
		x16:    x * constants.Unit,
		y16:    y * constants.Unit,
		width:  width,
		height: height,
	}
}

func (e *Rect) Position() image.Rectangle {
	minX := e.x16 / constants.Unit
	minY := e.y16 / constants.Unit
	maxX := minX + e.width
	maxY := minY + e.height
	return image.Rect(minX, minY, maxX, maxY)
}
