package mocks

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// MockOneWayPlatform implements body.OneWayPlatform for testing.
type MockOneWayPlatform struct {
	IsOneWayFunc       func() bool
	SetPassThroughFunc func(actor body.Collidable, frames int)
	IsPassThroughFunc  func(actor body.Collidable) bool
	UpdateFunc         func()

	// Embedded minimal Body/Collidable state
	id  string
	pos image.Rectangle
}

func NewMockOneWayPlatform(id string, pos image.Rectangle) *MockOneWayPlatform {
	return &MockOneWayPlatform{id: id, pos: pos}
}

func (m *MockOneWayPlatform) IsOneWay() bool {
	if m.IsOneWayFunc != nil {
		return m.IsOneWayFunc()
	}
	return true
}

func (m *MockOneWayPlatform) SetPassThrough(actor body.Collidable, frames int) {
	if m.SetPassThroughFunc != nil {
		m.SetPassThroughFunc(actor, frames)
	}
}

func (m *MockOneWayPlatform) IsPassThrough(actor body.Collidable) bool {
	if m.IsPassThroughFunc != nil {
		return m.IsPassThroughFunc(actor)
	}
	return false
}

func (m *MockOneWayPlatform) Update() {
	if m.UpdateFunc != nil {
		m.UpdateFunc()
	}
}

// body.Body stubs
func (m *MockOneWayPlatform) ID() string                          { return m.id }
func (m *MockOneWayPlatform) SetID(id string)                     { m.id = id }
func (m *MockOneWayPlatform) Position() image.Rectangle           { return m.pos }
func (m *MockOneWayPlatform) SetPosition(x, y int)                { m.pos = image.Rect(x, y, x+m.pos.Dx(), y+m.pos.Dy()) }
func (m *MockOneWayPlatform) SetPosition16(x16, y16 int)          { m.SetPosition(x16/16, y16/16) }
func (m *MockOneWayPlatform) SetSize(w, h int)                    { m.pos.Max = m.pos.Min.Add(image.Pt(w, h)) }
func (m *MockOneWayPlatform) Scale() float64                      { return 1 }
func (m *MockOneWayPlatform) SetScale(float64)                    {}
func (m *MockOneWayPlatform) GetPosition16() (int, int)           { return m.pos.Min.X * 16, m.pos.Min.Y * 16 }
func (m *MockOneWayPlatform) GetPositionMin() (int, int)          { return m.pos.Min.X, m.pos.Min.Y }
func (m *MockOneWayPlatform) GetShape() body.Shape                { return m }
func (m *MockOneWayPlatform) Width() int                          { return m.pos.Dx() }
func (m *MockOneWayPlatform) Height() int                         { return m.pos.Dy() }
func (m *MockOneWayPlatform) Owner() interface{}                  { return nil }
func (m *MockOneWayPlatform) SetOwner(interface{})                {}
func (m *MockOneWayPlatform) LastOwner() interface{}              { return nil }

// body.Collidable stubs
func (m *MockOneWayPlatform) GetTouchable() body.Touchable                        { return m }
func (m *MockOneWayPlatform) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (m *MockOneWayPlatform) CollisionPosition() []image.Rectangle                { return []image.Rectangle{m.pos} }
func (m *MockOneWayPlatform) CollisionShapes() []body.Collidable                  { return nil }
func (m *MockOneWayPlatform) IsObstructive() bool                                 { return true }
func (m *MockOneWayPlatform) SetIsObstructive(bool)                               {}
func (m *MockOneWayPlatform) AddCollision(...body.Collidable)                     {}
func (m *MockOneWayPlatform) ClearCollisions()                                    {}
func (m *MockOneWayPlatform) SetTouchable(body.Touchable)                         {}
func (m *MockOneWayPlatform) OnTouch(body.Collidable)                             {}
func (m *MockOneWayPlatform) OnBlock(body.Collidable)                             {}
func (m *MockOneWayPlatform) ApplyValidPosition(d int, ax bool, sp body.BodiesSpace) (int, int, bool) {
	return m.pos.Min.X, m.pos.Min.Y, false
}
