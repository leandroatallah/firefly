package camera

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
)

func setupConfig() {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
	})
}

func TestNewController(t *testing.T) {
	setupConfig()
	base := enginecamera.NewController(0, 0)
	ctrl := NewController(base)
	if ctrl.Base() != base {
		t.Error("expected Base() to return the base controller")
	}
}

func TestControllerVerticalOnlyUpward(t *testing.T) {
	setupConfig()
	base := enginecamera.NewController(160, 120)
	ctrl := NewController(base)

	rect := bodyphysics.NewRect(100, 100, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(100, 100) // target center Y is 108

	ctrl.SetFollowTarget(target)
	ctrl.SetFollowing(true)

	// Initial center should be (116, 108)
	ctrl.Update()
	cx, cy := base.Kamera().Center()
	if cx != 116 || cy != 108 {
		t.Errorf("expected initial center (116, 108), got (%f, %f)", cx, cy)
	}

	// Move target UP (Y decreases) - camera SHOULD follow
	target.SetPosition(100, 50) // target center Y is 58
	ctrl.Update()
	cx, cy = base.Kamera().Center()
	if cx != 116 || cy != 58 {
		t.Errorf("expected center after moving UP (116, 58), got (%f, %f)", cx, cy)
	}

	// Move target DOWN (Y increases) - camera SHOULD NOT follow downward
	target.SetPosition(100, 100) // target center Y is 108
	ctrl.Update()
	cx, cy = base.Kamera().Center()
	if cx != 116 || cy != 58 {
		t.Errorf("expected center after moving DOWN (116, 58), got (%f, %f) - downward movement should be blocked", cx, cy)
	}

	// Move target SIDEWAYS while at same height - camera SHOULD follow X
	target.SetPosition(200, 50) // target center X is 216
	ctrl.Update()
	cx, cy = base.Kamera().Center()
	if cx != 216 || cy != 58 {
		t.Errorf("expected center after moving sideways (216, 58), got (%f, %f)", cx, cy)
	}
}

func TestControllerBoundsClamping(t *testing.T) {
	setupConfig()
	base := enginecamera.NewController(160, 120)
	ctrl := NewController(base)

	// Set bounds [0, 0] to [1000, 1000]
	// Viewport is [320, 240], half [160, 120]
	// Min center: [160, 120], Max center: [840, 880]
	bounds := image.Rect(0, 0, 1000, 1000)
	ctrl.SetBounds(&bounds)

	rect := bodyphysics.NewRect(0, 0, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	ctrl.SetFollowTarget(target)
	ctrl.SetFollowing(true)

	// Move target way UP-LEFT (outside bounds)
	target.SetPosition(-100, -100)
	ctrl.Update()
	cx, cy := base.Kamera().Center()
	if cx != 160 || cy != 120 {
		t.Errorf("expected clamped center at min bounds (160, 120), got (%f, %f)", cx, cy)
	}

	// Move target way DOWN-RIGHT (outside bounds)
	// targetY > lastCameraY so it will block, let's reset lastCameraY first
	ctrl.SetCenter(160, 120)
	target.SetPosition(2000, 2000)
	// Since 2000 > 120, it will block downward movement unless we force it or move up first.
	// Let's just move it right but at same height
	target.SetPosition(2000, 112) // target center Y is 120
	ctrl.Update()
	cx, cy = base.Kamera().Center()
	if cx != 840 || cy != 120 {
		t.Errorf("expected clamped center at max X bounds (840, 120), got (%f, %f)", cx, cy)
	}
}

func TestControllerSettersUpdateLastCameraY(t *testing.T) {
	setupConfig()
	base := enginecamera.NewController(0, 0)
	ctrl := NewController(base)

	// SetCenter should update lastCameraY
	ctrl.SetCenter(100, 200)
	
	// If we now follow a target at 300, it should NOT move down
	rect := bodyphysics.NewRect(100, 300, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(100, 300) // center Y is 308
	ctrl.SetFollowTarget(target) // SetFollowTarget also updates lastCameraY to target center
	
	ctrl.Update()
	_, cy := base.Kamera().Center()
	if cy != 308 {
		t.Errorf("expected center Y 308 after SetFollowTarget, got %f", cy)
	}
	
	// Move target down
	target.SetPosition(100, 400) // center Y is 408
	ctrl.Update()
	_, cy = base.Kamera().Center()
	if cy != 308 {
		t.Errorf("expected center Y to stay 308, got %f", cy)
	}

	// SetPositionTopLeft should update lastCameraY
	// screen height is 240, half is 120
	ctrl.SetPositionTopLeft(0, 0) // center Y becomes 120
	ctrl.Update()
	_, cy = base.Kamera().Center()
	if cy != 120 {
		t.Errorf("expected center Y 120 after SetPositionTopLeft, got %f", cy)
	}
	
	// Move target to 200 center Y (down from 120)
	target.SetPosition(100, 192) // center Y 200
	ctrl.Update()
	_, cy = base.Kamera().Center()
	if cy != 120 {
		t.Errorf("expected center Y to stay 120, got %f", cy)
	}
}

func TestControllerDelegations(t *testing.T) {
	setupConfig()
	base := enginecamera.NewController(0, 0)
	ctrl := NewController(base)

	if ctrl.IsFollowing() {
		t.Error("expected IsFollowing to be false initially")
	}
	ctrl.SetFollowing(true)
	if !ctrl.IsFollowing() {
		t.Error("expected IsFollowing to be true after SetFollowing")
	}

	if ctrl.FollowTarget() != nil {
		t.Error("expected FollowTarget to be nil initially")
	}
	
	rect := bodyphysics.NewRect(0, 0, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	ctrl.SetFollowTarget(target)
	if ctrl.FollowTarget() != target {
		t.Error("expected FollowTarget to return the target")
	}

	if ctrl.Target() != target {
		t.Error("expected Target() to return the target")
	}

	if ctrl.Width() != 320 {
		t.Errorf("expected Width 320, got %f", ctrl.Width())
	}
	if ctrl.Height() != 240 {
		t.Errorf("expected Height 240, got %f", ctrl.Height())
	}

	// Smoke tests for other delegations
	ctrl.SetBounds(&image.Rectangle{})
	if ctrl.Bounds() == nil {
		t.Error("expected Bounds to be set")
	}

	ctrl.AddTrauma(0.5)
	ctrl.CamDebug()
	
	if ctrl.Kamera() != nil {
		t.Error("expected Kamera() to return nil as per placeholder")
	}
	
	pos := ctrl.Position()
	if pos.Dx() != 32 || pos.Dy() != 16 {
		t.Errorf("expected Position to return target rect, got %v", pos)
	}

	// Test Draw
	src := ebiten.NewImage(32, 32)
	dst := ebiten.NewImage(320, 240)
	opts := &ebiten.DrawImageOptions{}
	ctrl.Draw(src, opts, dst)

	// Test DrawCollisionBox
	ctrl.DrawCollisionBox(dst, target)
}
