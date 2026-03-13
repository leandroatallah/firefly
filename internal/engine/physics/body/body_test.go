package body

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

func TestNewBody(t *testing.T) {
	shape := NewRect(0, 0, 10, 10)
	b := NewBody(shape)

	if b == nil {
		t.Fatal("NewBody returned nil")
	}
	if b.shape != shape {
		t.Errorf("expected shape %v; got %v", shape, b.shape)
	}
}

func TestBody_GetShape(t *testing.T) {
	shape := NewRect(0, 0, 20, 30)
	b := NewBody(shape)

	got := b.GetShape()
	if got != shape {
		t.Errorf("expected shape %v; got %v", shape, got)
	}
}

func TestBody_SetID(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))

	b.SetID("test-id")
	if b.id != "test-id" {
		t.Errorf("expected id 'test-id'; got '%s'", b.id)
	}
}

func TestBody_ID(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))
	b.SetID("my-body")

	got := b.ID()
	if got != "my-body" {
		t.Errorf("expected ID 'my-body'; got '%s'", got)
	}
}

func TestBody_Position(t *testing.T) {
	tests := []struct {
		name          string
		x, y          int
		width, height int
		wantMinX      int
		wantMinY      int
		wantMaxX      int
		wantMaxY      int
	}{
		{"origin", 0, 0, 10, 10, 0, 0, 10, 10},
		{"offset", 5, 7, 20, 15, 5, 7, 25, 22},
		{"negative", -10, -5, 10, 10, -10, -5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shape := NewRect(tt.x, tt.y, tt.width, tt.height)
			b := NewBody(shape)
			b.SetPosition(tt.x, tt.y)

			pos := b.Position()
			if pos.Min.X != tt.wantMinX || pos.Min.Y != tt.wantMinY ||
				pos.Max.X != tt.wantMaxX || pos.Max.Y != tt.wantMaxY {
				t.Errorf("expected rect (%d,%d)-(%d,%d); got %v",
					tt.wantMinX, tt.wantMinY, tt.wantMaxX, tt.wantMaxY, pos)
			}
		})
	}
}

func TestBody_GetPositionMin(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))
	b.SetPosition(15, 25)

	x, y := b.GetPositionMin()
	if x != 15 || y != 25 {
		t.Errorf("expected (15, 25); got (%d, %d)", x, y)
	}
}

func TestBody_SetPosition(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))

	b.SetPosition(100, 200)
	if b.x16 != fp16.To16(100) || b.y16 != fp16.To16(200) {
		t.Errorf("expected position (100, 200) in fp16; got (%d, %d)", b.x16, b.y16)
	}

	// Verify actual position
	x, y := b.GetPositionMin()
	if x != 100 || y != 200 {
		t.Errorf("expected GetPositionMin to return (100, 200); got (%d, %d)", x, y)
	}
}

func TestBody_SetPosition16(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))

	x16 := fp16.To16(50)
	y16 := fp16.To16(75)
	b.SetPosition16(x16, y16)

	gotX16, gotY16 := b.GetPosition16()
	if gotX16 != x16 || gotY16 != y16 {
		t.Errorf("expected (%d, %d); got (%d, %d)", x16, y16, gotX16, gotY16)
	}
}

func TestBody_GetPosition16(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))
	b.SetPosition16(100, 200)

	x16, y16 := b.GetPosition16()
	if x16 != 100 || y16 != 200 {
		t.Errorf("expected (100, 200); got (%d, %d)", x16, y16)
	}
}

func TestBody_Position_WithFixedPoint(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))

	// Set position with fractional part in fixed-point
	// fp16 uses scale of 16, so 0.5 = 8, 0.25 = 4
	x16 := fp16.To16(10) + 8  // 10.5 in fp16
	y16 := fp16.To16(20) + 4  // 20.25 in fp16
	b.SetPosition16(x16, y16)

	pos := b.Position()
	// From16 should truncate to integer
	if pos.Min.X != 10 || pos.Min.Y != 20 {
		t.Errorf("expected (10, 20) after truncation; got (%d, %d)", pos.Min.X, pos.Min.Y)
	}
}

func TestBody_RectangleDimensions(t *testing.T) {
	tests := []struct {
		name          string
		x, y          int
		width, height int
		wantWidth     int
		wantHeight    int
	}{
		{"square", 0, 0, 10, 10, 10, 10},
		{"wide", 0, 0, 50, 10, 50, 10},
		{"tall", 0, 0, 10, 50, 10, 50},
		{"large", 100, 100, 200, 150, 200, 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shape := NewRect(tt.x, tt.y, tt.width, tt.height)
			b := NewBody(shape)
			b.SetPosition(tt.x, tt.y)

			pos := b.Position()
			gotWidth := pos.Dx()
			gotHeight := pos.Dy()

			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Errorf("expected %dx%d; got %dx%d",
					tt.wantWidth, tt.wantHeight, gotWidth, gotHeight)
			}
		})
	}
}

func TestBody_SetPositionOverwrites(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))

	b.SetPosition(10, 20)
	b.SetPosition(30, 40)

	x, y := b.GetPositionMin()
	if x != 30 || y != 40 {
		t.Errorf("expected (30, 40); got (%d, %d)", x, y)
	}
}

func TestBody_EmptyID(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))

	if b.ID() != "" {
		t.Errorf("expected empty ID for new body; got '%s'", b.ID())
	}
}

type mockShape struct{}

func (s *mockShape) Width() int  { return 10 }
func (s *mockShape) Height() int { return 10 }

func TestBody_SetPosition_InvalidShape(t *testing.T) {
	// log.Fatal path
	t.Skip("SetPosition calls log.Fatal for non-Rect shape")
}

func TestBody_SetPosition16_InvalidShape(t *testing.T) {
	// log.Fatal path
	t.Skip("SetPosition16 calls log.Fatal for non-Rect shape")
}
