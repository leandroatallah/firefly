package body

import (
	"image"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type Body struct {
	body.Body

	shape body.Shape

	id       string
	x16, y16 int
}

func NewBody(shape body.Shape) *Body {
	return &Body{
		shape: shape,
	}
}

func (b *Body) ID() string {
	return b.id
}

func (b *Body) SetID(id string) {
	b.id = id
}

// Position() returns the body coordinates as a image.Rectangle.
func (b *Body) Position() image.Rectangle {
	minX := fp16.From16(b.x16)
	minY := fp16.From16(b.y16)
	maxX := minX + b.shape.Width()
	maxY := minY + b.shape.Height()
	return image.Rect(minX, minY, maxX, maxY)
}

func (b *Body) GetPositionMin() (int, int) {
	pos := b.Position()
	return pos.Min.X, pos.Min.Y
}

// SetPosition updates the body position.
func (b *Body) SetPosition(x, y int) {
	// NOTE: For now, it only accepts rect shape.
	_, ok := b.GetShape().(*Rect)
	if !ok {
		log.Fatal("SetPosition expects a *Rect instance")
	}
	b.x16 = fp16.To16(x)
	b.y16 = fp16.To16(y)
}

func (b *Body) GetShape() body.Shape {
	return b.shape
}
