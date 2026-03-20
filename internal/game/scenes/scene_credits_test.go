package gamescene

import (
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

func TestCreditsScene_Structure(t *testing.T) {
	moduleRoot := getModuleRoot()
	fontMain, err := font.NewFontText(os.DirFS(moduleRoot), config.Get().MainFontFace)
	if err != nil {
		t.Fatalf("failed to load font: %v", err)
	}

	mockNav := &mocks.MockSceneManager{}
	ctx := &app.AppContext{
		SceneManager: mockNav,
		Font:         fontMain,
	}

	s := NewCreditsScene(ctx)
	if s == nil {
		t.Fatal("NewCreditsScene returned nil")
	}

	if s.fontTitle == nil {
		t.Error("fontTitle is nil")
	}
	if s.fontText == nil {
		t.Error("fontText is nil")
	}
	if s.credits == nil || len(s.credits) == 0 {
		t.Error("credits should not be empty")
	}
}

func TestCreditsScene_Lifecycle(t *testing.T) {
	moduleRoot := getModuleRoot()
	fontMain, err := font.NewFontText(os.DirFS(moduleRoot), config.Get().MainFontFace)
	if err != nil {
		t.Fatalf("failed to load font: %v", err)
	}

	mockNav := &mocks.MockSceneManager{}
	ctx := &app.AppContext{
		SceneManager: mockNav,
		Space:        space.NewSpace(),
		Font:         fontMain,
	}
	s := NewCreditsScene(ctx)

	// Test OnStart
	s.OnStart()

	// Test Update doesn't panic
	err = s.Update()
	if err != nil {
		t.Errorf("Update returned error: %v", err)
	}

	// Verify scroll offset decreases (auto-scroll)
	initialOffset := s.scrollOffset
	s.Update()
	if s.scrollOffset >= initialOffset {
		t.Error("scrollOffset should decrease after Update (auto-scroll)")
	}
}

func TestCreditsScene_Draw(t *testing.T) {
	moduleRoot := getModuleRoot()
	fontMain, err := font.NewFontText(os.DirFS(moduleRoot), config.Get().MainFontFace)
	if err != nil {
		t.Fatalf("failed to load font: %v", err)
	}

	mockNav := &mocks.MockSceneManager{}
	ctx := &app.AppContext{
		SceneManager: mockNav,
		Font:         fontMain,
		Space:        space.NewSpace(),
	}
	s := NewCreditsScene(ctx)
	s.OnStart()

	screen := ebiten.NewImage(320, 240)

	// Test Draw doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked: %v", r)
		}
	}()

	s.Draw(screen)
}

func TestCreditsScene_BuildCredits(t *testing.T) {
	credits := buildCredits()

	if len(credits) == 0 {
		t.Fatal("buildCredits returned empty list")
	}

	// Verify structure: should have sections and entries
	hasSection := false
	hasEntry := false
	for _, c := range credits {
		if c.name == "" && c.title != "" {
			hasSection = true
		}
		if c.name != "" {
			hasEntry = true
		}
	}

	if !hasSection {
		t.Error("credits should have section headers")
	}
	if !hasEntry {
		t.Error("credits should have entries")
	}
}
