package gamescene

import (
	"testing"
	"time"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

func TestMenuScene_Structure(t *testing.T) {
	mockNav := &mocks.MockSceneManager{}
	ctx := &app.AppContext{
		SceneManager: mockNav,
	}

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

	// Verify initial state (before OnStart)
	// Menus are created but visibility logic is in OnStart
	// Default visibility for Menu is false
	if s.mainMenu.Visible() {
		t.Log("mainMenu visible before OnStart (unexpected but allowed if init changed)")
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
	mockNav := &mocks.MockSceneManager{}
	ctx := &app.AppContext{
		SceneManager: mockNav,
	}
	s := NewMenuScene(ctx)
	s.OnStart()

	// Simulate delay to enable interaction
	ticks := timing.FromDuration(time.Second) + 1
	for i := 0; i < ticks; i++ {
		s.Update()
	}

	// Test 1: Select "Game Start" (Index 0)
	// We manually trigger Select() as we can't mock input easily
	s.mainMenu.Select()

	if !s.isNavigating {
		t.Error("Expected isNavigating to be true after selecting Game Start")
	}

	// Reset for next test
	s.OnStart()
	// Need to simulate delay again?
	// s.count is reset to 0 in OnStart.
	// But for manual Select() calls, we don't strictly need s.count unless callbacks check it.
	// The callback checks: if !s.isNavigating { ... }
	// It does NOT check s.count. s.count is checked in Update before calling Update/Select.
	// So manual Select() works regardless of count, but logic inside callback is what matters.

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
	// Options menu items: Back, Language, Fullscreen.
	// Selection defaults to 0 ("Back").
	s.optionsMenu.Select()

	if !s.mainMenu.Visible() {
		t.Error("mainMenu should be visible after Back")
	}
	if s.optionsMenu.Visible() {
		t.Error("optionsMenu should be hidden after Back")
	}
}
