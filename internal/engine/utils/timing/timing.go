package timing

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// TPS returns the Ticks Per Second of the game engine.
// This delegates to Ebitengine's current TPS value.
func TPS() int {
	return ebiten.TPS()
}

// FromDuration converts a time.Duration to the number of frames (ticks) based on TPS.
func FromDuration(d time.Duration) int {
	return int(d.Seconds() * float64(TPS()))
}

// ToDuration converts a number of frames (ticks) to time.Duration based on TPS.
func ToDuration(frames int) time.Duration {
	return time.Duration(float64(frames) / float64(TPS()) * float64(time.Second))
}
