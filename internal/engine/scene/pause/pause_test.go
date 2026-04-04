package pause

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestPauseScreen_Toggle(t *testing.T) {
	onStartCalled := false
	onFinishCalled := false
	p := NewPauseScreen(ebiten.KeyEscape, 0)
	p.SetOnStart(func(ps *PauseScreen) { onStartCalled = true })
	p.SetOnFinish(func(ps *PauseScreen) { onFinishCalled = true })

	if p.IsPaused() {
		t.Error("expected initial isPaused to be false")
	}

	// Toggle ON
	p.Toggle()
	if !p.IsPaused() {
		t.Error("expected isPaused to be true after toggle")
	}
	if !onStartCalled {
		t.Error("onStart callback not called")
	}

	// Toggle OFF
	p.Toggle()
	if p.IsPaused() {
		t.Error("expected isPaused to be false after second toggle")
	}
	if !onFinishCalled {
		t.Error("onFinish callback not called")
	}
}

func TestPauseScreen_UpdateAndCount(t *testing.T) {
	p := NewPauseScreen(ebiten.KeyEscape, 0)

	p.Update() // Should not increment count while not paused
	if p.Count() != 0 {
		t.Errorf("expected count 0; got %d", p.Count())
	}

	p.Toggle()
	p.Update() // Increment count
	if p.Count() != 1 {
		t.Errorf("expected count 1; got %d", p.Count())
	}

	p.Toggle() // Should reset count
	if p.Count() != 0 {
		t.Errorf("expected count reset to 0; got %d", p.Count())
	}
}

func TestPauseScreen_DisableFor(t *testing.T) {
	// Pause for 10 frames (approx 166ms at 60fps)
	disableFor := 100 * time.Millisecond
	p := NewPauseScreen(ebiten.KeyEscape, disableFor)

	p.Toggle() // Should set disable = true
	if !p.disable {
		t.Error("expected disable to be true after toggle with disableFor")
	}

	// Try to toggle off immediately - should be ignored
	p.Toggle()
	if !p.IsPaused() {
		t.Error("toggle should have been ignored while disabled")
	}

	// Update many times to exceed disableFor duration
	// 100ms at 60TPS is 6 frames.
	for i := 0; i < 10; i++ {
		p.Update()
	}

	// Now it should be possible to toggle off
	p.Toggle()
	if p.IsPaused() {
		t.Error("expected toggle to work after disableFor duration")
	}
}
