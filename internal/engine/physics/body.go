package physics

import (
	"image"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type Body struct {
	body.Body

	shape body.Shape

	id       string
	x16, y16 int
}

func NewBody(shape body.Shape) *Body {
	// cfg := config.Get()
	return &Body{
		shape: shape,
		// x16:    x * cfg.Unit,
		// y16:    y * cfg.Unit,
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
	minX := b.x16 / config.Get().Unit
	minY := b.y16 / config.Get().Unit
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
	cfg := config.Get()
	b.x16 = cfg.To16(x)
	b.y16 = cfg.To16(y)
}

func (b *Body) GetShape() body.Shape {
	return b.shape
}
