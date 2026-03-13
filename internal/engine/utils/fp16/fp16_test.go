package fp16

import "testing"

func TestTo16(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 16},
		{2, 32},
		{-1, -16},
		{10, 160},
	}

	for _, tt := range tests {
		if got := To16(tt.input); got != tt.want {
			t.Errorf("To16(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestFrom16(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{16, 1},
		{32, 2},
		{-16, -1},
		{160, 10},
		{15, 0}, // Integer division behavior
		{17, 1},
	}

	for _, tt := range tests {
		if got := From16(tt.input); got != tt.want {
			t.Errorf("From16(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
