package object

import "github.com/leandroatallah/firefly/internal/config"

type Element struct {
	x16, y16      int
	width, height int
}

func NewElement(x, y, width, height int) Element {
	return Element{
		x16:    x * config.Unit,
		y16:    y * config.Unit,
		width:  width,
		height: height,
	}
}

func (e *Element) Position() (minX, minY, maxX, maxY int) {
	minX = e.x16 / config.Unit
	minY = e.y16 / config.Unit
	maxX = minX + e.width
	maxY = minY + e.height
	return
}
