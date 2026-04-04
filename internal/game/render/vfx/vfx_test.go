package vfx

import (
	"image"
	"os"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestMain(m *testing.M) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 224,
	}
	config.Set(cfg)
	os.Exit(m.Run())
}

func TestOverheadText(t *testing.T) {
	mockActor := &mocks.MockActor{
		Id:  "test-actor",
		Pos: image.Rect(100, 100, 110, 110),
	}

	ot := NewOverheadText(mockActor, "Hello", 60)
	if ot == nil {
		t.Fatal("NewOverheadText returned nil")
	}

	screen := ebiten.NewImage(320, 240)
	cam := camera.NewController(0, 0)

	// Draw should update position based on actor
	ot.Draw(screen, cam)

	if ot.X != 105 { // 100 + 10/2
		t.Errorf("expected X 105, got %f", ot.X)
	}
	if ot.Y != 90 { // 100 - 10
		t.Errorf("expected Y 90, got %f", ot.Y)
	}
}

func TestScreenText(t *testing.T) {
	st := NewScreenText(100, 100, "Screen Message", 60)
	if st == nil {
		t.Fatal("NewScreenText returned nil")
	}

	screen := ebiten.NewImage(320, 240)
	cam := camera.NewController(0, 0)

	st.Draw(screen, cam)
}
