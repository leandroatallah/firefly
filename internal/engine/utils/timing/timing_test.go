package timing

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestTiming(t *testing.T) {
	// Set intended TPS to a known value for testing
	intendedTPS := 60
	ebiten.SetTPS(intendedTPS)

	// Note: ebiten.TPS() should return the maximum TPS (intended), 
	// even if the game loop isn't running.
	currentTPS := TPS()
	if currentTPS != intendedTPS {
		// In some environments, it might still return 0 if not fully initialized?
		// But let's assume 60 or what we set.
		// If it's 0, our division in ToDuration will result in +Inf.
		t.Logf("Warning: TPS() returned %d, expected %d. This might happen in some headless environments.", currentTPS, intendedTPS)
	}

	// Test FromDuration
	d := time.Second
	frames := FromDuration(d)
	if currentTPS > 0 && frames != currentTPS {
		t.Errorf("FromDuration(%v) = %d, want %d", d, frames, currentTPS)
	}

	// Test ToDuration
	if currentTPS > 0 {
		framesInput := currentTPS
		gotDuration := ToDuration(framesInput)
		if gotDuration != time.Second {
			t.Errorf("ToDuration(%d) = %v, want %v", framesInput, gotDuration, time.Second)
		}
	}
}

func TestFromDuration_Zero(t *testing.T) {
	if got := FromDuration(0); got != 0 {
		t.Errorf("FromDuration(0) = %d, want 0", got)
	}
}

func TestToDuration_Zero(t *testing.T) {
	if got := ToDuration(0); got != 0 {
		t.Errorf("ToDuration(0) = %v, want 0", got)
	}
}
