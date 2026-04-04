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
