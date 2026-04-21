package vfx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestManager(t *testing.T) {
	// Create a dummy vfx.json
	vfxData := []VFXConfig{
		{
			Type: "jump",
			ParticleData: schemas.ParticleData{
				Image: "jump.png",
			},
		},
	}

	jsonData, _ := json.Marshal(vfxData)
	_ = os.WriteFile("vfx_test.json", jsonData, 0644)
	defer os.Remove("vfx_test.json")

	m := NewManagerFromPath("vfx_test.json")
	if m == nil {
		t.Fatal("NewManager returned nil")
	}

	m.SetDefaultFont(nil)

	m.SpawnJumpPuff(10, 10, 5)
	m.SpawnLandingPuff(20, 20, 5)
	m.SpawnPuff("non_existent", 0, 0, 1, 0)

	m.SpawnFloatingText("hello", 10, 10, 10)

	m.Update()

	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	m.Draw(screen, cam)
}

func TestSpawnDirectionalPuff_AnchorAndFlip(t *testing.T) {
	vfxData := []VFXConfig{
		{
			Type: "muzzle",
			ParticleData: schemas.ParticleData{
				Image:     "muzzle.png",
				FrameRate: 1,
			},
		},
	}
	jsonData, _ := json.Marshal(vfxData)
	_ = os.WriteFile("vfx_dir_test.json", jsonData, 0644)
	defer os.Remove("vfx_dir_test.json")

	tests := []struct {
		name       string
		faceRight  bool
		wantAnchor float64
		wantFlipX  bool
	}{
		{"facing right anchors left edge, no flip", true, 0.0, false},
		{"facing left anchors right edge, flips sprite", false, 1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManagerFromPath("vfx_dir_test.json")
			if m == nil {
				t.Fatal("NewManager returned nil")
			}
			m.SpawnDirectionalPuff("muzzle", 50, 60, tt.faceRight, 1, 0)

			parts := m.system.Particles()
			if len(parts) != 1 {
				t.Fatalf("expected 1 particle, got %d", len(parts))
			}
			p := parts[0]
			if p.AnchorX != tt.wantAnchor {
				t.Errorf("AnchorX: got %v, want %v", p.AnchorX, tt.wantAnchor)
			}
			if p.FlipX != tt.wantFlipX {
				t.Errorf("FlipX: got %v, want %v", p.FlipX, tt.wantFlipX)
			}
		})
	}
}
