package tween_test

import (
	"math"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/physics/tween"
)

func TestInOutSineTween(t *testing.T) {
	const (
		from     = 160.0
		to       = 0.0
		duration = 18
		epsilon  = 10.0
	)

	tests := []struct {
		name      string
		ticksTo   int
		checkDone bool
		wantNear  float64
	}{
		{
			name:     "first tick is close to from",
			ticksTo:  1,
			wantNear: from,
		},
		{
			name:     "middle tick is near midpoint",
			ticksTo:  duration / 2,
			wantNear: (from + to) / 2,
		},
		{
			name:      "after duration ticks Done is true and value is close to to",
			ticksTo:   duration,
			checkDone: true,
			wantNear:  to,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tw := tween.NewInOutSineTween(from, to, duration)

			var val float64
			for i := 0; i < tc.ticksTo; i++ {
				val = tw.Tick()
			}

			if math.Abs(val-tc.wantNear) > epsilon {
				t.Errorf("got %.4f, want near %.4f (±%.1f)", val, tc.wantNear, epsilon)
			}

			if tc.checkDone && !tw.Done() {
				t.Error("expected Done() == true after durationFrames ticks")
			}
		})
	}

	t.Run("Reset restarts tween", func(t *testing.T) {
		tw := tween.NewInOutSineTween(from, to, duration)
		for i := 0; i < duration; i++ {
			tw.Tick()
		}
		if !tw.Done() {
			t.Fatal("expected Done before Reset")
		}
		tw.Reset()
		if tw.Done() {
			t.Error("expected Done() == false after Reset")
		}
	})
}
