package body

import (
	"testing"
)

func TestNewRect(t *testing.T) {
	r := NewRect(0, 0, 10, 20)

	if r == nil {
		t.Fatal("NewRect returned nil")
	}
	if r.width != 10 {
		t.Errorf("expected width 10; got %d", r.width)
	}
	if r.height != 20 {
		t.Errorf("expected height 20; got %d", r.height)
	}
}

func TestRect_Width(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   int
	}{
		{"zero", 0, 10, 0},
		{"positive", 10, 20, 10},
		{"large", 1000, 500, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRect(0, 0, tt.width, tt.height)
			if r.Width() != tt.want {
				t.Errorf("Width() = %d; want %d", r.Width(), tt.want)
			}
		})
	}
}

func TestRect_Height(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   int
	}{
		{"zero", 10, 0, 0},
		{"positive", 10, 20, 20},
		{"large", 500, 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRect(0, 0, tt.width, tt.height)
			if r.Height() != tt.want {
				t.Errorf("Height() = %d; want %d", r.Height(), tt.want)
			}
		})
	}
}

func TestRect_Dimensions(t *testing.T) {
	tests := []struct {
		name          string
		x, y          int
		width, height int
		wantWidth     int
		wantHeight    int
	}{
		{"square", 0, 0, 10, 10, 10, 10},
		{"wide", 0, 0, 100, 10, 100, 10},
		{"tall", 0, 0, 10, 100, 10, 100},
		{"small", 0, 0, 1, 1, 1, 1},
		{"origin ignored", 50, 50, 25, 30, 25, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRect(tt.x, tt.y, tt.width, tt.height)
			if r.Width() != tt.wantWidth || r.Height() != tt.wantHeight {
				t.Errorf("expected %dx%d; got %dx%d",
					tt.wantWidth, tt.wantHeight, r.Width(), r.Height())
			}
		})
	}
}

func TestRect_ImplementsShape(t *testing.T) {
	r := NewRect(0, 0, 10, 10)

	// Verify it implements body.Shape interface
	var _ interface {
		Width() int
		Height() int
	} = r
}

func TestRect_ZeroDimensions(t *testing.T) {
	r := NewRect(0, 0, 0, 0)

	if r.Width() != 0 {
		t.Errorf("expected width 0; got %d", r.Width())
	}
	if r.Height() != 0 {
		t.Errorf("expected height 0; got %d", r.Height())
	}
}

func TestRect_LargeDimensions(t *testing.T) {
	r := NewRect(0, 0, 10000, 20000)

	if r.Width() != 10000 {
		t.Errorf("expected width 10000; got %d", r.Width())
	}
	if r.Height() != 20000 {
		t.Errorf("expected height 20000; got %d", r.Height())
	}
}

func TestRect_MultipleInstances(t *testing.T) {
	r1 := NewRect(0, 0, 10, 20)
	r2 := NewRect(0, 0, 30, 40)

	if r1.Width() == r2.Width() {
		t.Error("expected different widths")
	}
	if r1.Height() == r2.Height() {
		t.Error("expected different heights")
	}
}

func TestRect_PositionIgnored(t *testing.T) {
	// Rect only stores width and height, x and y are ignored
	r1 := NewRect(0, 0, 10, 10)
	r2 := NewRect(100, 200, 10, 10)

	if r1.Width() != r2.Width() {
		t.Errorf("expected same width; got %d vs %d", r1.Width(), r2.Width())
	}
	if r1.Height() != r2.Height() {
		t.Errorf("expected same height; got %d vs %d", r1.Height(), r2.Height())
	}
}
