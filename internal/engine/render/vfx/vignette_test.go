package vfx

import (
	"image"
	"os"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
)

func TestMain(m *testing.M) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 224,
	}
	config.Set(cfg)
	os.Exit(m.Run())
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

	pixels := v.BuildMaskPixels(baseCam, actor, cfg.ScreenWidth, cfg.ScreenHeight)
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
	pixels := v.BuildMaskPixels(baseCam, actor, cfg.ScreenWidth, cfg.ScreenHeight)

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
