package app

import "math"

const (
	FastForwardMinFactor = 1.0
	FastForwardMaxFactor = 16.0
)

// FastForwardTPS computes the target TPS for fast-forward mode.
// Returns (targetTPS, shouldApply).
// shouldApply == false means SetTPS must NOT be called.
//   - fastForward==false             → (0, false)
//   - clampedFactor <= 1.0 (no-op)   → (0, false)
//   - otherwise                       → (round(defaultTPS * clampedFactor), true)
//
// Factor is clamped to [FastForwardMinFactor, FastForwardMaxFactor] before evaluation.
func FastForwardTPS(fastForward bool, factor float64, defaultTPS int) (int, bool) {
	if !fastForward {
		return 0, false
	}
	if factor < FastForwardMinFactor {
		factor = FastForwardMinFactor
	}
	if factor > FastForwardMaxFactor {
		factor = FastForwardMaxFactor
	}
	if factor == FastForwardMinFactor {
		return 0, false
	}
	tps := int(math.Round(float64(defaultTPS) * factor))
	return tps, true
}
