package physics

import (
	"image"
)

type Shape interface {
	Position() image.Rectangle
}
