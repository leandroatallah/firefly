package app

import "math"

const (
	SlowMoMinFactor = 0.05
	SlowMoMaxFactor = 1.0
)

// EffectiveTPS computes the target TPS for slow-motion mode.
// Returns (targetTPS, shouldApply).
// shouldApply == false means SetTPS must NOT be called.
//   - slowMo==false                  → (0, false)
//   - clampedFactor >= 1.0 (no-op)   → (0, false)
//   - otherwise                       → (round(defaultTPS * clampedFactor), true)
//
// Factor is clamped to [SlowMoMinFactor, SlowMoMaxFactor] before evaluation.
func EffectiveTPS(slowMo bool, factor float64, defaultTPS int) (int, bool) {
	if !slowMo {
		return 0, false
	}
	if factor < SlowMoMinFactor {
		factor = SlowMoMinFactor
	}
	if factor > SlowMoMaxFactor {
		factor = SlowMoMaxFactor
	}
	if factor == SlowMoMaxFactor {
		return 0, false
	}
	tps := int(math.Round(float64(defaultTPS) * factor))
	return tps, true
}
