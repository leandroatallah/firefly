package movement

import (
	"image"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

type dimsProvider struct{ w, h int }

func (d dimsProvider) GetTilemapWidth() int                     { return d.w }
func (d dimsProvider) GetTilemapHeight() int                    { return d.h }
func (d dimsProvider) GetCameraBounds() (image.Rectangle, bool) { return image.Rectangle{}, false }

func TestClampToPlayAreaWithTilemapDimensions(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	sp := space.NewSpace()
	sp.SetTilemapDimensionsProvider(dimsProvider{w: 200, h: 100})

	rect := bodyphysics.NewRect(0, 0, 20, 20)
	actor := bodyphysics.NewObstacleRect(rect)
	actor.SetID("actor")
	actor.SetPosition(-10, -10) // outside

	onGround := clampToPlayArea(actor, sp)
	if actor.Position().Min.X != 0 || actor.Position().Min.Y != 0 {
		t.Fatalf("expected position clamped to (0,0), got %v", actor.Position().Min)
	}
	if onGround {
		t.Fatalf("should not be ground when at top-left corner")
	}

	actor.SetPosition(195, 95) // beyond bounds; should clamp inside
	_ = clampToPlayArea(actor, sp)
	if actor.Position().Max.X > 200 || actor.Position().Max.Y > 100 {
		t.Fatalf("expected max inside bounds, got %v", actor.Position().Max)
	}
}

func TestPlatformMovementGroundDetection(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			DownwardGravity:   4,
			UpwardGravity:     2,
			MaxFallSpeed:      128,
			HorizontalInertia: 1.0,
		},
	}
	config.Set(cfg)

	sp := space.NewSpace()

	// Actor at (10,10) size 10x10
	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(10, 10)

	// Ground tile directly beneath (y=20..30)
	ground := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 100, 10))
	ground.SetID("ground")
	ground.SetIsObstructive(true)
	ground.SetPosition(0, 20)
	ground.AddCollisionBodies()
	sp.AddBody(actor)
	sp.AddBody(ground)

	model := NewPlatformMovementModel(nil)

	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("update error: %v", err)
	}

	if !model.OnGround() {
		t.Fatalf("expected actor to be grounded")
	}
}
