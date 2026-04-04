package scene

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestTilemapScene_Basics(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	am := &audio.AudioManager{}
	ctx := &app.AppContext{AudioManager: am}

	s := NewTilemapScene(ctx)

	if s.AppContext() != ctx {
		t.Error("AppContext not set correctly")
	}

	if s.Camera() == nil {
		t.Error("Camera not initialized")
	}

	if s.Audiomanager() != am {
		t.Error("Audiomanager() returned wrong manager")
	}

	// Test Camera Config
	s.SetCameraConfig(CameraConfig{Mode: CameraModeFollow})
	if !s.Camera().IsFollowing() {
		t.Error("SetCameraConfig(Follow) failed to set camera following state")
	}

	s.SetCameraConfig(CameraConfig{Mode: CameraModeFixed})
	if s.Camera().IsFollowing() {
		t.Error("SetCameraConfig(Fixed) failed to unset camera following state")
	}

	// Test Tilemap Width/Height defaults (no tilemap loaded)
	if s.GetTilemapWidth() != 320 {
		t.Errorf("expected default width 320; got %d", s.GetTilemapWidth())
	}
	if s.GetTilemapHeight() != 240 {
		t.Errorf("expected default height 240; got %d", s.GetTilemapHeight())
	}

	// Test Camera Bounds
	bounds := image.Rect(0, 0, 100, 100)
	s.Camera().SetBounds(&bounds)
	b, ok := s.GetCameraBounds()
	if !ok || b != bounds {
		t.Errorf("GetCameraBounds returned %v, %v; want %v, true", b, ok, bounds)
	}

	// Test Update
	if err := s.Update(); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
}

func TestTilemapScene_SetPlayerStartPosition_NoTilemap(t *testing.T) {
	s := &TilemapScene{} // No tilemap
	p := &mocks.MockActor{}
	// Should not panic
	s.SetPlayerStartPosition(p)
}
