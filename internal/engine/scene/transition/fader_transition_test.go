package transition

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
)

func setupConfig() {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
	})
}

func TestNewFader(t *testing.T) {
	hold := 100 * time.Millisecond
	visible := 200 * time.Millisecond
	f := NewFader(hold, visible)

	if f == nil {
		t.Fatal("NewFader returned nil")
	}
	if f.holdDuration != hold {
		t.Errorf("expected holdDuration %v; got %v", hold, f.holdDuration)
	}
	if f.visibleDuration != visible {
		t.Errorf("expected visibleDuration %v; got %v", visible, f.visibleDuration)
	}
	if f.fadeSpeed != 15 {
		t.Errorf("expected default fadeSpeed 15; got %f", f.fadeSpeed)
	}
}

func TestFader_FadeOutFadeInSequence(t *testing.T) {
	setupConfig()
	// Fader with no hold/visible wait for quick testing
	f := NewFader(0, 0)
	f.fadeSpeed = 100 // fast fade

	called := false
	cb := func() { called = true }

	// Start FadeOut
	f.fadeOut(cb)
	if !f.active || !f.exiting || f.alpha != 0 {
		t.Error("fadeOut did not initialize correctly")
	}

	// Update until alpha reaches 255
	for f.exiting {
		f.Update()
	}

	if f.alpha != 255 {
		t.Errorf("expected alpha 255 after exit; got %f", f.alpha)
	}
	if !called {
		t.Error("expected callback to be called when alpha reaches 255 (no hold)")
	}

	// Next update will transition to fade in (starting)
	f.Update()

	if !f.starting {
		t.Error("expected fade in (starting) to begin immediately (no visible wait)")
	}

	// Update until alpha reaches 0
	for f.starting {
		f.Update()
	}

	if f.alpha != 0 {
		t.Errorf("expected alpha 0 after start; got %f", f.alpha)
	}
	if f.active {
		t.Error("expected fader to be inactive after fade in")
	}
}

func TestFader_WithHoldDuration(t *testing.T) {
	setupConfig()
	// Fader with 1 second hold (60 frames at 60 TPS)
	// New behavior: callback is called immediately when fade-out completes
	// Hold duration is the wait AFTER callback before fade-in starts
	f := NewFader(time.Second, 0)
	f.fadeSpeed = 255 // instant fade out

	called := false
	cb := func() { called = true }

	f.fadeOut(cb)
	f.Update() // Alpha 0 -> 255, exiting becomes false, callback called immediately

	if f.alpha != 255 || f.exiting {
		t.Fatalf("expected instant fade out; alpha=%f, exiting=%v", f.alpha, f.exiting)
	}
	if !called {
		t.Error("callback should be called immediately when fade-out completes")
	}
	if f.starting {
		t.Error("fade in should not start yet due to holdDuration")
	}

	// Wait for holdDuration (approx 60 frames) - fade in should not start yet
	for i := 0; i < 59; i++ {
		f.Update()
		if f.starting {
			t.Errorf("fade in started too early at frame %d", i)
		}
	}

	f.Update() // Frame 60 - hold complete, fade in should start
	if !f.starting {
		t.Error("fade in should start after holdDuration")
	}
}

func TestFader_WithVisibleDuration(t *testing.T) {
	setupConfig()
	// Fader with no hold, but 1 second visible wait after callback
	f := NewFader(0, time.Second)
	f.fadeSpeed = 255 // instant fade out

	called := false
	cb := func() { called = true }

	f.fadeOut(cb)
	f.Update() // Instant fade out + callback (no hold)

	if !called {
		t.Error("callback should be called immediately with no hold")
	}
	if f.starting {
		t.Error("fade in should not start yet due to visibleDuration")
	}

	// Wait for visibleDuration
	for i := 0; i < 60; i++ {
		if f.starting {
			t.Errorf("fade in started too early at frame %d", i)
		}
		f.Update()
	}

	if !f.starting {
		t.Error("fade in should start after visibleDuration")
	}
}

func TestFader_Draw(t *testing.T) {
	setupConfig()
	f := NewFader(0, 0)
	screen := ebiten.NewImage(320, 240)

	// Inactive - should not draw
	f.Draw(screen)

	// Active
	f.active = true
	f.alpha = 128
	// Should not panic
	f.Draw(screen)
}

func TestFader_StartTransition(t *testing.T) {
	setupConfig()
	f := NewFader(0, 0)
	called := false
	f.StartTransition(func() { called = true })

	if !f.active || !f.exiting {
		t.Error("StartTransition did not initialize fade out")
	}
	
	// StartTransition calls fadeOut with a callback that calls fadeIn
	// fadeIn sets starting = true and calls the user callback immediately
	
	// Simulate reaching end of fade out
	f.alpha = 255
	f.Update()
	
	if !called {
		t.Error("callback should be called via fadeIn when fadeOut completes")
	}
}

func TestFader_MultipleFadeOutCalls(t *testing.T) {
	f := NewFader(0, 0)
	f.fadeOut(nil)
	f.fadeOut(func() { t.Error("second call should be ignored") })
}

func TestFader_EndTransition(t *testing.T) {
	f := NewFader(0, 0)
	f.EndTransition(func() {})
}
