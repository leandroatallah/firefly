package vfx

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestFadeOverlay_Activate(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)

	if !f.IsActive() {
		t.Error("expected overlay to be active after Activate()")
	}
	if f.alpha != 0 {
		t.Errorf("expected alpha=0 after Activate, got %v", f.alpha)
	}
}

func TestFadeOverlay_Update_IncrementAlpha(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)

	done := f.Update()
	if done {
		t.Error("expected fade not done after 1 frame")
	}
	if f.alpha <= 0 || f.alpha > 25.5 {
		t.Errorf("expected alpha in range (0, 25.5], got %v", f.alpha)
	}
}

func TestFadeOverlay_Update_CompleteAtMax(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)

	var done bool
	for i := 0; i < 11; i++ {
		done = f.Update()
	}

	if f.Alpha() < 255 {
		t.Errorf("expected alpha at 255, got %v", f.Alpha())
	}
	if done && f.IsActive() {
		t.Error("should not be animating after completion")
	}
	if !f.IsPersisting() {
		t.Error("expected fade to persist")
	}
}

func TestFadeOverlay_IsActive_False_WhenInactive(t *testing.T) {
	f := NewFadeOverlay()
	if f.IsActive() {
		t.Error("expected overlay not animating by default")
	}
}

func TestFadeOverlay_IsPersisting_True_WhenAlphaGreater(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)
	f.Update()

	if !f.IsPersisting() {
		t.Error("expected IsPersisting to return true when alpha > 0")
	}
}

func TestFadeOverlay_IsActive_False_AfterComplete(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)

	for i := 0; i < 15; i++ {
		f.Update()
	}

	if f.IsActive() {
		t.Error("expected IsActive to return false after animation completes")
	}
	if !f.IsPersisting() {
		t.Error("expected IsPersisting to return true (fade stays on screen)")
	}
}

func TestFadeOverlay_Reset_ClearsState(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)
	f.Update()
	f.Reset()

	if f.IsActive() {
		t.Error("expected overlay inactive after Reset()")
	}
	if f.alpha != 0 {
		t.Errorf("expected alpha=0 after Reset, got %v", f.alpha)
	}
}

func TestFadeOverlay_Draw_NoOpWhenInactive(t *testing.T) {
	f := NewFadeOverlay()
	screen := ebiten.NewImage(320, 180)

	// Should not panic when inactive
	f.Draw(screen)
}

func TestFadeOverlay_Draw_CreatesImage(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeOut(10)
	screen := ebiten.NewImage(320, 180)

	// Should not panic when active
	f.Draw(screen)
}

func TestFadeOverlay_FadeIn_SetsAlpha255(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeIn(10)

	if !f.IsActive() {
		t.Error("expected overlay to be active after ActivateFadeIn()")
	}
	if f.Alpha() != 255 {
		t.Errorf("expected alpha=255 after ActivateFadeIn, got %v", f.Alpha())
	}
}

func TestFadeOverlay_FadeIn_DecrementsAlpha(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeIn(10)

	done := f.Update()
	if done {
		t.Error("expected fade not done after 1 frame")
	}
	if f.Alpha() >= 255 {
		t.Errorf("expected alpha < 255 after one update, got %v", f.Alpha())
	}
}

func TestFadeOverlay_FadeIn_CompleteAtZero(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeIn(10)

	var done bool
	for i := 0; i < 15; i++ {
		done = f.Update()
	}

	if f.Alpha() != 0 {
		t.Errorf("expected alpha=0 after fade-in completes, got %v", f.Alpha())
	}
	if f.IsActive() {
		t.Error("expected IsActive=false after fade-in completes")
	}
	if done && f.IsPersisting() {
		t.Error("expected IsPersisting=false after fade-in completes")
	}
}

func TestFadeOverlay_FadeIn_DefaultFrames(t *testing.T) {
	f := NewFadeOverlay()
	f.FadeIn(0)

	if !f.IsActive() {
		t.Error("expected overlay active with default frames")
	}
}
