package physics

import (
	"image"

	"github.com/leandroatallah/firefly/internal/config"
)

type Shape interface {
	Position() image.Rectangle
}

type Rect struct {
	x16, y16      int
	width, height int
}

func NewRect(x, y, width, height int) Rect {
	return Rect{
		x16:    x * config.Unit,
		y16:    y * config.Unit,
		width:  width,
		height: height,
	}
}

func (e *Rect) Position() image.Rectangle {
	minX := e.x16 / config.Unit
	minY := e.y16 / config.Unit
	maxX := minX + e.width
	maxY := minY + e.height
	return image.Rect(minX, minY, maxX, maxY)
}
