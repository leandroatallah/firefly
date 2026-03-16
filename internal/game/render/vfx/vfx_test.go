package vfx

import (
	"image"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
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
		Id: "test-actor",
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

func TestVignetteDisabledNoChanges(t *testing.T) {
	cfg := config.Get()
	baseCam := enginecamera.NewController(0, 0)
	baseCam.SetCenter(float64(cfg.ScreenWidth)/2, float64(cfg.ScreenHeight)/2)

	actor := &mocks.MockActor{
		Id:  "player",
		Pos: image.Rect(cfg.ScreenWidth/2-5, cfg.ScreenHeight/2-5, cfg.ScreenWidth/2+5, cfg.ScreenHeight/2+5),
	}

	v := NewVignette()
	// Ensure disabled even if toggled before.
	v.Enable(12)
	v.Disable()

	pixels := v.buildMaskPixels(baseCam, actor, cfg.ScreenWidth, cfg.ScreenHeight)
	for i := 3; i < len(pixels); i += 4 {
		if pixels[i] != 0 {
			t.Fatalf("expected alpha 0 everywhere when disabled; got %d at alpha index %d", pixels[i], i)
		}
	}
}

func TestVignetteEnabledBlacksOutOutsideRadius(t *testing.T) {
	cfg := config.Get()
	baseCam := enginecamera.NewController(0, 0)
	baseCam.SetCenter(float64(cfg.ScreenWidth)/2, float64(cfg.ScreenHeight)/2)

	actor := &mocks.MockActor{
		Id:  "player",
		Pos: image.Rect(cfg.ScreenWidth/2-5, cfg.ScreenHeight/2-5, cfg.ScreenWidth/2+5, cfg.ScreenHeight/2+5),
	}

	v := NewVignette()
	v.Enable(12)
	pixels := v.buildMaskPixels(baseCam, actor, cfg.ScreenWidth, cfg.ScreenHeight)

	centerX, centerY := cfg.ScreenWidth/2, cfg.ScreenHeight/2
	// Inside circle: should be transparent in the overlay.
	insideIdx := (centerY*cfg.ScreenWidth + centerX) * 4
	if pixels[insideIdx+3] != 0 {
		t.Fatalf("expected center alpha 0 (transparent), got %d", pixels[insideIdx+3])
	}

	// Far outside circle: should be black (0,0,0,255).
	outX, outY := 0, 0
	outIdx := (outY*cfg.ScreenWidth + outX) * 4
	if pixels[outIdx+3] != 255 || pixels[outIdx+0] != 0 || pixels[outIdx+1] != 0 || pixels[outIdx+2] != 0 {
		t.Fatalf("expected outside pixel to be black opaque, got (%d,%d,%d,%d)",
			pixels[outIdx+0], pixels[outIdx+1], pixels[outIdx+2], pixels[outIdx+3],
		)
	}
}
