package tween

import "math"

// InOutSineTween interpolates from `from` to `to` over `durationFrames` using an InOutSine curve.
type InOutSineTween struct {
	from, to       float64
	durationFrames int
	currentFrame   int
}

func NewInOutSineTween(from, to float64, durationFrames int) *InOutSineTween {
	return &InOutSineTween{from: from, to: to, durationFrames: durationFrames}
}

// Tick returns the interpolated value at the current frame, then advances the counter.
// First call returns a value at progress 0 (≈ from); after durationFrames calls Done() is true.
func (t *InOutSineTween) Tick() float64 {
	if t.currentFrame >= t.durationFrames {
		return t.to
	}
	if t.durationFrames <= 1 {
		t.currentFrame++
		return t.to
	}
	progress := float64(t.currentFrame) / float64(t.durationFrames-1)
	val := t.from + (t.to-t.from)*(1-math.Cos(math.Pi*progress))/2
	t.currentFrame++
	return val
}

// Done returns true when the tween has completed.
func (t *InOutSineTween) Done() bool {
	return t.currentFrame >= t.durationFrames
}

// Reset restarts the tween from the beginning.
func (t *InOutSineTween) Reset() {
	t.currentFrame = 0
}
