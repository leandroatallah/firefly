package body

import (
	"image"
	"testing"
)

func TestResizeFixedBottom(t *testing.T) {
	base := image.Rect(0, 0, 32, 64)

	tests := []struct {
		name      string
		newHeight int
		wantMinY  int
		wantMaxY  int
	}{
		{"normal shrink", 32, 32, 64},
		{"normal grow", 80, -16, 64},
		{"zero height", 0, 64, 64},
		{"negative height", -5, 64, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResizeFixedBottom(base, tt.newHeight)

			if got.Min.Y != tt.wantMinY {
				t.Errorf("Min.Y = %d; want %d", got.Min.Y, tt.wantMinY)
			}
			if got.Max.Y != tt.wantMaxY {
				t.Errorf("Max.Y = %d; want %d", got.Max.Y, tt.wantMaxY)
			}
			if base != image.Rect(0, 0, 32, 64) {
				t.Error("input rect was mutated")
			}
		})
	}
}

func TestResizeFixedTop(t *testing.T) {
	base := image.Rect(0, 0, 32, 64)

	tests := []struct {
		name      string
		newHeight int
		wantMinY  int
		wantMaxY  int
	}{
		{"normal shrink", 32, 0, 32},
		{"zero height", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResizeFixedTop(base, tt.newHeight)

			if got.Min.Y != tt.wantMinY {
				t.Errorf("Min.Y = %d; want %d", got.Min.Y, tt.wantMinY)
			}
			if got.Max.Y != tt.wantMaxY {
				t.Errorf("Max.Y = %d; want %d", got.Max.Y, tt.wantMaxY)
			}
			if base != image.Rect(0, 0, 32, 64) {
				t.Error("input rect was mutated")
			}
		})
	}
}
