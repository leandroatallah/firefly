package app

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestFastForwardTPS validates the pure fast-forward TPS helper.
//
// The helper computes the effective Ebitengine TPS for fast-forward mode and
// signals via the second return whether ebiten.SetTPS should be called at all.
// Cases cover: disabled flag, default 4x, 2x, 1.0 no-op, clamping (below min,
// above max), and rounding for a non-integer factor.
//
// At least one case uses ebiten.DefaultTPS to lock the constant the production
// call site uses.
func TestFastForwardTPS(t *testing.T) {
	tests := []struct {
		name        string
		fastForward bool
		factor      float64
		defaultTPS  int
		wantTPS     int
		wantOK      bool
	}{
		{
			name:        "T-F1 fast-forward disabled",
			fastForward: false,
			factor:      4.0,
			defaultTPS:  60,
			wantTPS:     0,
			wantOK:      false,
		},
		{
			name:        "T-F2 quad speed default uses ebiten.DefaultTPS",
			fastForward: true,
			factor:      4.0,
			defaultTPS:  ebiten.DefaultTPS,
			wantTPS:     240,
			wantOK:      true,
		},
		{
			name:        "T-F3 double speed",
			fastForward: true,
			factor:      2.0,
			defaultTPS:  60,
			wantTPS:     120,
			wantOK:      true,
		},
		{
			name:        "T-F4 factor 1.0 is no-op",
			fastForward: true,
			factor:      1.0,
			defaultTPS:  60,
			wantTPS:     0,
			wantOK:      false,
		},
		{
			name:        "T-F5 factor below min clamped to no-op",
			fastForward: true,
			factor:      0.5,
			defaultTPS:  60,
			wantTPS:     0,
			wantOK:      false,
		},
		{
			name:        "T-F6 factor above max clamped to 16x",
			fastForward: true,
			factor:      32.0,
			defaultTPS:  60,
			wantTPS:     960,
			wantOK:      true,
		},
		{
			name:        "T-F7 rounding 2.5x",
			fastForward: true,
			factor:      2.5,
			defaultTPS:  60,
			wantTPS:     150,
			wantOK:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTPS, gotOK := FastForwardTPS(tc.fastForward, tc.factor, tc.defaultTPS)
			if gotTPS != tc.wantTPS || gotOK != tc.wantOK {
				t.Fatalf("FastForwardTPS(%v, %v, %d) = (%d, %v); want (%d, %v)",
					tc.fastForward, tc.factor, tc.defaultTPS, gotTPS, gotOK, tc.wantTPS, tc.wantOK)
			}
		})
	}
}
