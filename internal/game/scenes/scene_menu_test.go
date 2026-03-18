package gamescene

import (
	"os"
	"testing"
	"time"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/data/i18n"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

func createMenuSceneContext() *app.AppContext {
	mockNav := &mocks.MockSceneManager{}
	i18nManager := i18n.NewI18nManager(os.DirFS("."))
	i18nManager.Load("en")
	return &app.AppContext{
		SceneManager: mockNav,
		I18n:         i18nManager,
	}
}

func TestMenuScene_Structure(t *testing.T) {
	ctx := createMenuSceneContext()
	s := NewMenuScene(ctx)
	if s == nil {
		t.Fatal("NewMenuScene returned nil")
	}

	if s.mainMenu == nil {
		t.Error("mainMenu is nil")
	}
	if s.optionsMenu == nil {
		t.Error("optionsMenu is nil")
	}

	s.OnStart()

	if !s.mainMenu.Visible() {
		t.Error("mainMenu should be visible after OnStart")
	}
	if s.optionsMenu.Visible() {
		t.Error("optionsMenu should be invisible after OnStart")
	}
}

func TestMenuScene_Interaction(t *testing.T) {
	ctx := createMenuSceneContext()
	s := NewMenuScene(ctx)
	s.OnStart()

	// Simulate delay to enable interaction
	ticks := timing.FromDuration(time.Second) + 1
	for i := 0; i < ticks; i++ {
		s.Update()
	}

	// Test 1: Select "Game Start" (Index 0)
	s.mainMenu.Select()

	if !s.isNavigating {
		t.Error("Expected isNavigating to be true after selecting Game Start")
	}

	// Reset for next test
	s.OnStart()

	// Test 2: Navigate to "Options" (Index 1)
	s.mainMenu.NavigateDown() // Selection 0 -> 1 ("Options")
	s.mainMenu.Select()

	if s.mainMenu.Visible() {
		t.Error("mainMenu should be hidden after selecting Options")
	}
	if !s.optionsMenu.Visible() {
		t.Error("optionsMenu should be visible after selecting Options")
	}

	// Test 3: Back from Options (Index 0)
	s.optionsMenu.Select()

	if !s.mainMenu.Visible() {
		t.Error("mainMenu should be visible after Back")
	}
	if s.optionsMenu.Visible() {
		t.Error("optionsMenu should be hidden after Back")
	}
}

func TestMenuScene_FullscreenToggle(t *testing.T) {
	ctx := createMenuSceneContext()
	// Initial config state - modify only Fullscreen
	cfg := config.Get()
	cfg.Fullscreen = false
	config.Set(cfg)

	s := NewMenuScene(ctx)
	s.OnStart()

	// Navigate to Options
	s.mainMenu.NavigateDown()
	s.mainMenu.Select()

	if !s.optionsMenu.Visible() {
		t.Fatal("optionsMenu should be visible")
	}

	// Navigate to Fullscreen (Index 2)
	s.optionsMenu.NavigateDown() // Index 0 -> 1 ("Language")
	s.optionsMenu.NavigateDown() // Index 1 -> 2 ("Fullscreen")

	// Select Fullscreen
	s.optionsMenu.Select()

	if !config.Get().Fullscreen {
		t.Error("Expected Fullscreen to be true after toggle")
	}

	// Toggle back
	s.optionsMenu.Select()
	if config.Get().Fullscreen {
		t.Error("Expected Fullscreen to be false after second toggle")
	}
}

func TestMenuScene_LanguageToggle(t *testing.T) {
	ctx := createMenuSceneContext()
	cfg := config.Get()
	cfg.Language = "en"
	config.Set(cfg)

	s := NewMenuScene(ctx)
	s.OnStart()

	// Navigate to Options
	s.mainMenu.NavigateDown()
	s.mainMenu.Select()

	// Navigate to Language (Index 1)
	s.optionsMenu.NavigateDown() // Index 0 -> 1 ("Language")

	// Select Language
	s.optionsMenu.Select()

	if config.Get().Language != "pt-br" {
		t.Errorf("Expected Language to be pt-br, got %s", config.Get().Language)
	}

	// Toggle back
	s.optionsMenu.Select()
	if config.Get().Language != "en" {
		t.Errorf("Expected Language to be en, got %s", config.Get().Language)
	}
}
