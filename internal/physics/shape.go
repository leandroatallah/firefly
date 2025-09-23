package physics

import (
	"image"
)

type Shape interface {
	// TODO: Make all Position return the same data type
	Position() image.Rectangle
}
