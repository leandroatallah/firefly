package body

import (
	"image"
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

type Body struct {
	shape body.Shape

	Ownership

	id         string
	x16, y16   int
	altitude16 int
	scale      float64
}

func NewBody(shape body.Shape) *Body {
	return &Body{
		shape: shape,
		scale: 1.0,
	}
}

func (b *Body) GetShape() body.Shape {
	return b.shape
}

func (b *Body) SetID(id string) {
	b.id = id
}

func (b *Body) ID() string {
	return b.id
}

func (b *Body) Scale() float64 {
	return b.scale
}

func (b *Body) SetScale(scale float64) {
	b.scale = scale
}

// Position() returns the body coordinates as a image.Rectangle.
func (b *Body) Position() image.Rectangle {
	minX := fp16.From16(b.x16)
	groundY := fp16.From16(b.y16)
	alt := fp16.From16(b.altitude16)
	minY := groundY - alt
	maxX := minX + b.shape.Width()
	maxY := minY + b.shape.Height()
	return image.Rect(minX, minY, maxX, maxY)
}

func (b *Body) Altitude() int       { return fp16.From16(b.altitude16) }
func (b *Body) Altitude16() int     { return b.altitude16 }
func (b *Body) SetAltitude(alt int) { b.altitude16 = fp16.To16(alt) }
func (b *Body) SetAltitude16(a int) { b.altitude16 = a }

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

func (b *Body) SetPosition16(x16, y16 int) {
	// NOTE: For now, it only accepts rect shape.
	_, ok := b.GetShape().(*Rect)
	if !ok {
		log.Fatal("SetPosition expects a *Rect instance")
	}
	b.x16 = x16
	b.y16 = y16
}

func (b *Body) GetPosition16() (int, int) {
	return b.x16, b.y16
}

func (b *Body) SetSize(width, height int) {
	if r, ok := b.shape.(*Rect); ok {
		r.SetSize(width, height)
	} else {
		log.Printf("Warning: SetSize called on body with non-Rect shape: %T", b.shape)
	}
}
