package app

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestEffectiveTPS validates the pure slow-motion TPS helper.
//
// The helper computes the effective Ebitengine TPS for slow-motion mode and
// signals via the second return whether ebiten.SetTPS should be called at all.
// Cases cover: disabled flag, quarter speed, half speed, 1.0 no-op, clamping
// (zero, negative, > max), and rounding for a repeating-decimal factor.
//
// At least one case uses ebiten.DefaultTPS to lock the constant the production
// call site uses.
func TestEffectiveTPS(t *testing.T) {
	tests := []struct {
		name       string
		slowMo     bool
		factor     float64
		defaultTPS int
		wantTPS    int
		wantOK     bool
	}{
		{
			name:       "T-S1 slow-mo disabled",
			slowMo:     false,
			factor:     0.25,
			defaultTPS: 60,
			wantTPS:    0,
			wantOK:     false,
		},
		{
			name:       "T-S2 quarter speed default uses ebiten.DefaultTPS",
			slowMo:     true,
			factor:     0.25,
			defaultTPS: ebiten.DefaultTPS,
			wantTPS:    15,
			wantOK:     true,
		},
		{
			name:       "T-S3 half speed",
			slowMo:     true,
			factor:     0.5,
			defaultTPS: 60,
			wantTPS:    30,
			wantOK:     true,
		},
		{
			name:       "T-S4 factor 1.0 is no-op",
			slowMo:     true,
			factor:     1.0,
			defaultTPS: 60,
			wantTPS:    0,
			wantOK:     false,
		},
		{
			name:       "T-S5 factor 0.0 clamped to min",
			slowMo:     true,
			factor:     0.0,
			defaultTPS: 60,
			wantTPS:    3,
			wantOK:     true,
		},
		{
			name:       "T-S6 factor 2.0 clamped to max no-op",
			slowMo:     true,
			factor:     2.0,
			defaultTPS: 60,
			wantTPS:    0,
			wantOK:     false,
		},
		{
			name:       "T-S7 negative factor clamped to min",
			slowMo:     true,
			factor:     -0.5,
			defaultTPS: 60,
			wantTPS:    3,
			wantOK:     true,
		},
		{
			name:       "T-S8 rounding one third",
			slowMo:     true,
			factor:     1.0 / 3.0,
			defaultTPS: 60,
			wantTPS:    20,
			wantOK:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTPS, gotOK := EffectiveTPS(tc.slowMo, tc.factor, tc.defaultTPS)
			if gotTPS != tc.wantTPS || gotOK != tc.wantOK {
				t.Fatalf("EffectiveTPS(%v, %v, %d) = (%d, %v); want (%d, %v)",
					tc.slowMo, tc.factor, tc.defaultTPS, gotTPS, gotOK, tc.wantTPS, tc.wantOK)
			}
		})
	}
}
