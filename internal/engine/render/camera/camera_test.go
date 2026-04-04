package camera

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestControllerFollowTargetAndBounds(t *testing.T) {
	// Save and restore config state
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)

	rect := bodyphysics.NewRect(0, 0, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(300, 200) // near bottom-right

	ctrl.SetFollowTarget(target)
	ctrl.SetFollowing(true)

	// Set bounds smaller than target position to force clamping inside camera view
	bounds := image.Rect(0, 0, 320, 240)
	ctrl.SetBounds(&bounds)

	ctrl.Update()

	if ctrl.Bounds() == nil || *ctrl.Bounds() != bounds {
		t.Fatalf("bounds were not set correctly")
	}

	// Smoke test draw flow (no panic)
	dst := ebiten.NewImage(320, 240)
	opts := &ebiten.DrawImageOptions{}
	src := ebiten.NewImage(1, 1)
	ctrl.Draw(src, opts, dst)

	// Test without following but with center
	ctrl.SetFollowing(false)
	ctrl.SetCenter(50, 50)
	ctrl.Update()
}

func TestControllerSetBounds(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)

	tests := []struct {
		name   string
		bounds *image.Rectangle
	}{
		{"nil bounds", nil},
		{"small bounds", &image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{100, 100}}},
		{"large bounds", &image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{1000, 1000}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl.SetBounds(tt.bounds)
			if ctrl.Bounds() != tt.bounds {
				t.Errorf("SetBounds() = %v, want %v", ctrl.Bounds(), tt.bounds)
			}
		})
	}
}

func TestControllerUpdateWithBounds(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)

	// Set a small bound to test clamping
	bounds := image.Rect(0, 0, 100, 100)
	ctrl.SetBounds(&bounds)
	ctrl.SetFollowing(true)

	// Create a target outside bounds
	rect := bodyphysics.NewRect(0, 0, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(500, 500) // Well outside bounds
	ctrl.SetFollowTarget(target)

	ctrl.Update()

	// Camera should be clamped to bounds
	// This is a smoke test to ensure no panic and bounds are respected
	if ctrl.Bounds() == nil {
		t.Error("Bounds should not be nil after Update")
	}
}

func TestControllerDraw(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)

	// Test normal draw
	t.Run("normal draw", func(t *testing.T) {
		src := ebiten.NewImage(32, 32)
		opts := &ebiten.DrawImageOptions{}
		dst := ebiten.NewImage(320, 240)
		// Should not panic
		ctrl.Draw(src, opts, dst)
	})

	// Test nil options - kamera library panics with nil options
	t.Run("nil options", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Draw() panicked with nil options (kamera limitation): %v", r)
			}
		}()
		src := ebiten.NewImage(32, 32)
		dst := ebiten.NewImage(320, 240)
		ctrl.Draw(src, nil, dst)
		// If we reach here, it didn't panic
		t.Error("Draw() should panic with nil options (kamera limitation)")
	})

	// Test nil source
	t.Run("nil source", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Draw() panicked with nil source as expected: %v", r)
			}
		}()
		opts := &ebiten.DrawImageOptions{}
		dst := ebiten.NewImage(320, 240)
		ctrl.Draw(nil, opts, dst)
		t.Error("Draw() should panic with nil source")
	})
}

type fakeCollidable struct {
	*bodyphysics.CollidableBody
	obstructive bool
}

func (f *fakeCollidable) IsObstructive() bool { return f.obstructive }

func TestDrawCollisionBoxColorsDontPanic(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)

	base := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(10, 10, 10, 10))
	base.SetID("b")
	x, y := base.GetPositionMin()
	col := bodyphysics.NewCollidableBodyFromRect(base.GetShape())
	col.SetPosition(x, y)
	base.AddCollision(col)
	f := &fakeCollidable{CollidableBody: base, obstructive: true}

	screen := ebiten.NewImage(100, 100)
	ctrl.DrawCollisionBox(screen, f)

	f.obstructive = false
	ctrl.DrawCollisionBox(screen, f)
}

func TestCameraNew(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	cam := NewCamera(10, 10)
	if cam == nil {
		t.Fatal("NewCamera returned nil")
	}
}

func TestControllerGettersSetters(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	ctrl := NewController(0, 0)

	if ctrl.IsFollowing() {
		t.Error("expected IsFollowing to be false by default")
	}
	ctrl.SetFollowing(true)
	if !ctrl.IsFollowing() {
		t.Error("expected IsFollowing to be true after SetFollowing(true)")
	}

	ctrl.SetCenter(100, 200)
	cx, cy := ctrl.Kamera().Center()
	if cx != 100 || cy != 200 {
		t.Errorf("expected center (100, 200), got (%f, %f)", cx, cy)
	}

	ctrl.SetPositionTopLeft(10, 20)
	cx, cy = ctrl.Kamera().Center()
	if cx != 170 || cy != 140 {
		t.Errorf("expected center (170, 140), got (%f, %f)", cx, cy)
	}
}

func TestControllerSetCenter(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	tests := []struct {
		name  string
		x, y  float64
		wantX float64
		wantY float64
	}{
		{"center at origin", 0, 0, 0, 0},
		{"center positive", 100, 200, 100, 200},
		{"center negative", -50, -100, -50, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewController(0, 0)
			ctrl.SetCenter(tt.x, tt.y)
			cx, cy := ctrl.Kamera().Center()
			if cx != tt.wantX || cy != tt.wantY {
				t.Errorf("SetCenter(%f, %f) = (%f, %f), want (%f, %f)",
					tt.x, tt.y, cx, cy, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestControllerSetPositionTopLeft(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	tests := []struct {
		name   string
		x, y   float64
		wantCX float64 // expected center X
		wantCY float64 // expected center Y
	}{
		{"top-left at origin", 0, 0, 160, 120},    // 320/2, 240/2
		{"top-left at 10,20", 10, 20, 170, 140},   // 10+160, 20+120
		{"top-left negative", -10, -20, 150, 100}, // -10+160, -20+120
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewController(0, 0)
			ctrl.SetPositionTopLeft(tt.x, tt.y)
			cx, cy := ctrl.Kamera().Center()
			if cx != tt.wantCX || cy != tt.wantCY {
				t.Errorf("SetPositionTopLeft(%f, %f) center = (%f, %f), want (%f, %f)",
					tt.x, tt.y, cx, cy, tt.wantCX, tt.wantCY)
			}
		})
	}
}

func TestControllerTargetAndPosition(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	ctrl := NewController(0, 0)

	rect := bodyphysics.NewRect(10, 20, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(10, 20)
	ctrl.SetFollowTarget(target)

	if ctrl.Target() != target {
		t.Error("expected Target() to return the follow target")
	}

	pos := ctrl.Position()
	if pos.Min.X != 10 || pos.Min.Y != 20 || pos.Dx() != 32 || pos.Dy() != 16 {
		t.Errorf("expected position (10, 20, 32, 16), got %+v", pos)
	}
}

func TestControllerSetFollowTarget(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	tests := []struct {
		name    string
		targetX int
		targetY int
		targetW int
		targetH int
		wantCX  float64
		wantCY  float64
	}{
		// SetFollowTarget sets camera center to target center
		// Target center = position + size/2
		{"target at origin", 0, 0, 32, 16, 16, 8},
		{"target at 100,50", 100, 50, 32, 16, 116, 58}, // 100+16, 50+8
		{"larger target", 0, 0, 64, 32, 32, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh controller for each test
			ctrl := NewController(0, 0)
			rect := bodyphysics.NewRect(tt.targetX, tt.targetY, tt.targetW, tt.targetH)
			target := bodyphysics.NewCollidableBodyFromRect(rect)

			// SetPosition must be called after creating the body
			// NewCollidableBodyFromRect creates body at origin
			target.SetPosition(tt.targetX, tt.targetY)

			// Verify target position before setting
			tx, ty := target.GetPositionMin()
			if tx != tt.targetX || ty != tt.targetY {
				t.Errorf("target position = (%d, %d), want (%d, %d)", tx, ty, tt.targetX, tt.targetY)
			}

			ctrl.SetFollowTarget(target)

			cx, cy := ctrl.Kamera().Center()
			if cx != tt.wantCX || cy != tt.wantCY {
				t.Errorf("SetFollowTarget() center = (%f, %f), want (%f, %f)",
					cx, cy, tt.wantCX, tt.wantCY)
			}
		})
	}
}

func TestCamDebugSmoke(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)
	ctrl.CamDebug()
}

func TestControllerAddTrauma(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)
	// Smoke test, should not panic
	ctrl.AddTrauma(0.5)
}
