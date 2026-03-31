package input

import "testing"

func TestHorizontalAxisValue(t *testing.T) {
	tests := []struct {
		name    string
		actions func(h *HorizontalAxis)
		want    int
	}{
		{
			name:    "only left held",
			actions: func(h *HorizontalAxis) { h.Press(-1) },
			want:    -1,
		},
		{
			name:    "only right held",
			actions: func(h *HorizontalAxis) { h.Press(1) },
			want:    1,
		},
		{
			name:    "both held, right last",
			actions: func(h *HorizontalAxis) { h.Press(-1); h.Press(1) },
			want:    1,
		},
		{
			name:    "both held, left last",
			actions: func(h *HorizontalAxis) { h.Press(1); h.Press(-1) },
			want:    -1,
		},
		{
			name:    "release last-pressed, other still held",
			actions: func(h *HorizontalAxis) { h.Press(-1); h.Press(1); h.Release(1) },
			want:    -1,
		},
		{
			name:    "release all",
			actions: func(h *HorizontalAxis) { h.Press(1); h.Release(1) },
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHorizontalAxis()
			tt.actions(h)
			if got := h.Value(); got != tt.want {
				t.Errorf("Value() = %d, want %d", got, tt.want)
			}
		})
	}
}
