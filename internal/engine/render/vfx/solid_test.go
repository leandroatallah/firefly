package vfx

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestSolidColor_Activate(t *testing.T) {
	f := NewSolidColor()
	f.FadeOut(10)

	if !f.IsActive() {
		t.Error("expected overlay to be active after Activate()")
	}
	if f.duration != 10 {
		t.Errorf("expected duration=10 after Activate, got %v", f.duration)
	}
}

func TestSolidColor_Update(t *testing.T) {
	f := NewSolidColor()
	f.FadeOut(2)

	done := f.Update()
	if done {
		t.Error("expected not done after 1 frame")
	}
	if f.duration != 1 {
		t.Errorf("expected duration=1, got %v", f.duration)
	}

	done = f.Update()
	if !done {
		t.Error("expected done after 2 frames")
	}
	if f.IsActive() {
		t.Error("expected not active after completion")
	}
}

func TestSolidColor_Reset(t *testing.T) {
	f := NewSolidColor()
	f.FadeOut(10)
	f.Reset()

	if f.IsActive() {
		t.Error("expected overlay inactive after Reset()")
	}
}

func TestSolidColor_Draw(t *testing.T) {
	f := NewSolidColor()
	screen := ebiten.NewImage(320, 180)

	// Inactive
	f.Draw(screen)

	// Active
	f.FadeOut(10)
	f.Draw(screen)
}

func TestSolidColor_SetColor(t *testing.T) {
	f := NewSolidColor()
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	f.SetColor(c)

	if f.color != c {
		t.Errorf("expected color %v, got %v", c, f.color)
	}
}

func TestSolidColor_Update_Inactive(t *testing.T) {
	f := NewSolidColor()
	if !f.Update() {
		t.Error("Update() on inactive overlay should return true")
	}
}

func TestSolidColor_ActivateFadeIn(t *testing.T) {
	f := NewSolidColor()
	// Should not panic and should be a no-op
	f.FadeIn(10)
	if f.IsActive() {
		t.Error("ActivateFadeIn should be a no-op for SolidColor")
	}
}
