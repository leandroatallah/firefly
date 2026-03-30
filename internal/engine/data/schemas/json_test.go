package schemas

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
)

func TestShapeRect_Rect(t *testing.T) {
	rect := ShapeRect{X: 10, Y: 20, Width: 100, Height: 200}

	x, y, width, height := rect.Rect()
	if x != 10 || y != 20 || width != 100 || height != 200 {
		t.Errorf("Rect() returned (%d, %d, %d, %d), expected (10, 20, 100, 200)", x, y, width, height)
	}
}

func TestShapeRect_ZeroValues(t *testing.T) {
	rect := ShapeRect{}

	x, y, width, height := rect.Rect()
	if x != 0 || y != 0 || width != 0 || height != 0 {
		t.Errorf("Rect() on zero value returned (%d, %d, %d, %d), expected (0, 0, 0, 0)", x, y, width, height)
	}
}

func TestAssetData_WithLoop(t *testing.T) {
	loop := true
	asset := AssetData{
		CollisionRects: []ShapeRect{
			{X: 0, Y: 0, Width: 10, Height: 10},
			{X: 10, Y: 10, Width: 20, Height: 20},
		},
		Loop: &loop,
	}

	if len(asset.CollisionRects) != 2 {
		t.Errorf("expected 2 CollisionRects, got %d", len(asset.CollisionRects))
	}
	if asset.Loop == nil || *asset.Loop != true {
		t.Error("expected Loop to be true")
	}
}

func TestAssetData_NoLoop(t *testing.T) {
	asset := AssetData{
		CollisionRects: []ShapeRect{},
		Loop:           nil,
	}

	if asset.Loop != nil {
		t.Error("expected Loop to be nil")
	}
}

func TestSpriteData(t *testing.T) {
	loop := false
	sprite := SpriteData{
		BodyRect:        ShapeRect{X: 0, Y: 0, Width: 32, Height: 32},
		Assets:          map[string]AssetData{"idle": {Loop: &loop}},
		FrameRate:       8,
		FacingDirection: animation.FaceDirectionRight,
	}

	if len(sprite.Assets) != 1 {
		t.Errorf("expected 1 Asset, got %d", len(sprite.Assets))
	}
	if sprite.FacingDirection != animation.FaceDirectionRight {
		t.Errorf("expected FacingDirection Right, got %d", sprite.FacingDirection)
	}
}

func TestParticleData(t *testing.T) {
	particle := ParticleData{
		Image:       "particles/fire.png",
		FrameWidth:  64,
		FrameHeight: 64,
		FrameRate:   12,
		Scale:       1.5,
	}

	// Spot-check key fields
	if particle.Image != "particles/fire.png" {
		t.Errorf("expected Image 'particles/fire.png', got %s", particle.Image)
	}
	if particle.Scale != 1.5 {
		t.Errorf("expected Scale 1.5, got %f", particle.Scale)
	}
}
