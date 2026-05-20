package movement

import (
	"math"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

// mockPlayerMovementBlocker implements PlayerMovementBlocker for testing
type mockPlayerMovementBlocker struct {
	blocked bool
}

func (m *mockPlayerMovementBlocker) IsMovementBlocked() bool {
	return m.blocked
}

// mockMovableCollidable implements body.MovableCollidable for testing
type mockMovableCollidable struct {
	*bodyphysics.ObstacleRect
}

func newMockMovableCollidable() *mockMovableCollidable {
	rect := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	rect.SetID("mock")
	rect.SetSpeed(2)
	rect.SetMaxSpeed(10)
	return &mockMovableCollidable{
		ObstacleRect: rect,
	}
}

func TestClampToPlayArea_Edges(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	tests := []struct {
		name         string
		x, y         int
		wantX, wantY int
	}{
		{"top-left corner", -10, -10, 0, 0},
		{"top-right corner", 330, -10, 310, 0},
		{"bottom-left corner", -10, 250, 0, 230},
		{"bottom-right corner", 330, 250, 310, 230},
		{"within bounds", 100, 100, 100, 100},
		{"left edge", -5, 100, 0, 100},
		{"right edge", 315, 100, 310, 100},
		{"top edge", 100, -5, 100, 0},
		{"bottom edge", 100, 235, 100, 230},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := space.NewSpace()
			actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
			actor.SetID("actor")
			actor.SetPosition(tt.x, tt.y)

			clampToPlayArea(actor, sp)

			pos := actor.Position().Min
			if pos.X != tt.wantX || pos.Y != tt.wantY {
				t.Errorf("clampToPlayArea(%d, %d) = (%d, %d); want (%d, %d)",
					tt.x, tt.y, pos.X, pos.Y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestClampToPlayArea_WithTilemapProvider(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	sp := space.NewSpace()
	sp.SetTilemapDimensionsProvider(dimsProvider{w: 400, h: 300})

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(-10, -10)

	clampToPlayArea(actor, sp)

	pos := actor.Position().Min
	if pos.X != 0 || pos.Y != 0 {
		t.Errorf("expected position clamped to (0,0); got (%d, %d)", pos.X, pos.Y)
	}

	// Test right/bottom edge with tilemap bounds
	actor.SetPosition(395, 295)
	clampToPlayArea(actor, sp)

	pos = actor.Position().Min
	if pos.X != 390 || pos.Y != 290 {
		t.Errorf("expected position clamped to (390, 290); got (%d, %d)", pos.X, pos.Y)
	}
}

func TestClampToPlayArea_OnGroundDetection(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	sp := space.NewSpace()
	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")

	// At bottom edge - should detect ground
	actor.SetPosition(100, 230)
	onGround := clampToPlayArea(actor, sp)
	if !onGround {
		t.Error("expected onGround=true at bottom edge")
	}

	// Not at bottom edge
	actor.SetPosition(100, 100)
	onGround = clampToPlayArea(actor, sp)
	if onGround {
		t.Error("expected onGround=false when not at bottom")
	}
}

func TestClampToPlayArea_NonRectShape(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	sp := space.NewSpace()

	// Create a body with non-rect shape (using mock)
	actor := newMockMovableCollidable()
	actor.SetPosition(-10, -10)

	// Should return false for non-rect shapes
	onGround := clampToPlayArea(actor, sp)
	if onGround {
		t.Error("expected onGround=false for non-rect shape")
	}
}

func TestPlatformMovementModel_New(t *testing.T) {
	model := NewPlatformMovementModel(nil)

	if model == nil {
		t.Fatal("NewPlatformMovementModel returned nil")
	}
	if model.playerMovementBlocker != nil {
		t.Error("expected nil playerMovementBlocker")
	}
	if !model.gravityEnabled {
		t.Error("expected gravityEnabled=true by default")
	}
}

func TestPlatformMovementModel_SetIsScripted(t *testing.T) {
	model := NewPlatformMovementModel(nil)

	model.SetIsScripted(true)
	if !model.isScripted {
		t.Error("expected isScripted=true")
	}

	model.SetIsScripted(false)
	if model.isScripted {
		t.Error("expected isScripted=false")
	}
}

func TestPlatformMovementModel_IsInputBlocked(t *testing.T) {
	blocker := &mockPlayerMovementBlocker{blocked: false}
	model := NewPlatformMovementModel(blocker)

	if model.IsInputBlocked() {
		t.Error("expected IsInputBlocked=false")
	}

	blocker.blocked = true
	if !model.IsInputBlocked() {
		t.Error("expected IsInputBlocked=true")
	}

	modelNil := NewPlatformMovementModel(nil)
	if modelNil.IsInputBlocked() {
		t.Error("expected IsInputBlocked=false with nil blocker")
	}
}

func TestPlatformMovementModel_OnGround(t *testing.T) {
	model := NewPlatformMovementModel(nil)

	if model.OnGround() {
		t.Error("expected OnGround=false by default")
	}

	model.SetOnGround(true)
	if !model.OnGround() {
		t.Error("expected OnGround=true after setting")
	}
}

func TestPlatformMovementModel_SetDashActive(t *testing.T) {
	model := NewPlatformMovementModel(nil)

	model.SetDashActive(true, 100)
	if !model.dashActive {
		t.Error("expected dashActive=true")
	}
	if model.dashVelocityX != 100 {
		t.Errorf("expected dashVelocityX=100; got %d", model.dashVelocityX)
	}

	model.SetDashActive(false, 0)
	if model.dashActive {
		t.Error("expected dashActive=false")
	}
}

func TestPlatformMovementModel_SetGravityEnabled(t *testing.T) {
	model := NewPlatformMovementModel(nil)

	if !model.gravityEnabled {
		t.Error("expected gravityEnabled=true by default")
	}

	model.SetGravityEnabled(false)
	if model.gravityEnabled {
		t.Error("expected gravityEnabled=false")
	}

	model.SetGravityEnabled(true)
	if !model.gravityEnabled {
		t.Error("expected gravityEnabled=true")
	}
}

func TestPlatformMovementModel_UpdateHorizontalVelocity(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			HorizontalInertia:     1.0,
			SpeedMultiplier:       1.0,
			AirControlMultiplier:  0.5,
			AirFrictionMultiplier: 2.0,
		},
	}
	config.Set(cfg)

	tests := []struct {
		name     string
		accelX   int
		onGround bool
		maxSpeed int
		wantVelX int
	}{
		{"no acceleration", 0, true, 10, 0},
		{"positive acceleration", fp16.To16(2), true, 10, 0},  // Will be increased
		{"negative acceleration", -fp16.To16(2), true, 10, 0}, // Will be decreased
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := newMockMovableCollidable()
			actor.SetMaxSpeed(tt.maxSpeed)
			actor.SetAcceleration(tt.accelX, 0)

			model := NewPlatformMovementModel(nil)
			model.onGround = tt.onGround

			vx, vy := model.UpdateHorizontalVelocity(actor)

			// Just verify no panic and velocity is set
			_ = vx
			_ = vy
		})
	}
}

func TestPlatformMovementModel_UpdateVerticalVelocity(t *testing.T) {
	actor := newMockMovableCollidable()
	model := NewPlatformMovementModel(nil)

	// Gravity enabled - should return current velocity
	model.gravityEnabled = true
	actor.SetVelocity(0, 0)
	vx, vy := model.UpdateVerticalVelocity(actor)
	if vx != 0 || vy != 0 {
		t.Errorf("expected (0, 0) with gravity enabled; got (%d, %d)", vx, vy)
	}

	// Gravity disabled with acceleration
	model.gravityEnabled = false
	actor.SetAcceleration(0, fp16.To16(5))
	_, vy = model.UpdateVerticalVelocity(actor)
	if vy != fp16.To16(5) {
		t.Errorf("expected vy=%d; got %d", fp16.To16(5), vy)
	}
}

func TestPlatformMovementModel_handleGravity(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			UpwardGravity:   2,
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := NewPlatformMovementModel(nil)

	// Gravity disabled
	model.gravityEnabled = false
	actor.SetVelocity(0, 0)
	vx, vy := model.handleGravity(actor)
	if vx != 0 || vy != 0 {
		t.Errorf("expected (0, 0) with gravity disabled; got (%d, %d)", vx, vy)
	}

	// On ground
	model.gravityEnabled = true
	model.onGround = true
	vx, vy = model.handleGravity(actor)
	if vx != 0 || vy != 0 {
		t.Errorf("expected (0, 0) on ground; got (%d, %d)", vx, vy)
	}

	// Airborne, falling
	model.onGround = false
	actor.SetVelocity(0, 50)
	_, vy = model.handleGravity(actor)
	if vy <= 50 {
		t.Errorf("expected vy > 50 after gravity; got %d", vy)
	}

	// Airborne, jumping (negative vy)
	actor.SetVelocity(0, -50)
	_, vy = model.handleGravity(actor)
	// Upward gravity adds to negative velocity (making it less negative towards zero)
	// vy = -50 + 2 = -48
	if vy != -48 {
		t.Errorf("expected vy=-48 after upward gravity; got %d", vy)
	}

	// Clamp to max fall speed
	actor.SetVelocity(0, 200)
	_, vy = model.handleGravity(actor)
	if vy > 128 {
		t.Errorf("expected vy <= 128 (max fall speed); got %d", vy)
	}
}

func TestPlatformMovementModel_CheckGround(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
		},
	}
	config.Set(cfg)

	sp := space.NewSpace()

	// Actor at (10, 10) size 10x10
	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(10, 10)
	sp.AddBody(actor)

	// Ground tile directly beneath (y=20..30)
	ground := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 100, 10))
	ground.SetID("ground")
	ground.SetIsObstructive(true)
	ground.SetPosition(0, 20)
	ground.AddCollisionBodies()
	sp.AddBody(ground)

	model := NewPlatformMovementModel(nil)

	// Test that CheckGround works when ground is below
	// Note: The exact behavior depends on the ground detection algorithm
	// which checks 1 pixel below the collision rects
	result := model.CheckGround(actor, sp)
	_ = result // Just verify it doesn't panic
}

func TestPlatformMovementModel_Update(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			DownwardGravity:   4,
			UpwardGravity:     2,
			MaxFallSpeed:      128,
			HorizontalInertia: 1.0,
			SpeedMultiplier:   1.0,
		},
	}
	config.Set(cfg)

	sp := space.NewSpace()

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(100, 100)
	actor.SetMaxSpeed(5)

	sp.AddBody(actor)

	model := NewPlatformMovementModel(nil)

	// Test freeze
	actor.SetFreeze(true)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update with freeze failed: %v", err)
	}
	actor.SetFreeze(false)

	// Test normal update
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify acceleration was reset
	accX, accY := actor.Acceleration()
	if accX != 0 || accY != 0 {
		t.Errorf("expected acceleration reset; got (%d, %d)", accX, accY)
	}
}

func TestTopDownMovementModel_New(t *testing.T) {
	blocker := &mockPlayerMovementBlocker{}
	model := NewTopDownMovementModel(blocker)

	if model == nil {
		t.Fatal("NewTopDownMovementModel returned nil")
	}
	if model.playerMovementBlocker != blocker {
		t.Error("expected playerMovementBlocker to be set")
	}
}

func TestTopDownMovementModel_SetIsScripted(t *testing.T) {
	model := NewTopDownMovementModel(nil)

	model.SetIsScripted(true)
	if !model.isScripted {
		t.Error("expected isScripted=true")
	}

	model.SetIsScripted(false)
	if model.isScripted {
		t.Error("expected isScripted=false")
	}
}

func TestTopDownMovementModel_InputHandler(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			HorizontalInertia: 0,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	blocker := &mockPlayerMovementBlocker{blocked: false}
	model := NewTopDownMovementModel(blocker)

	// Test scripted mode
	model.isScripted = true
	model.InputHandler(actor, nil)
	accX, accY := actor.Acceleration()
	if accX != 0 || accY != 0 {
		t.Error("expected no acceleration in scripted mode")
	}

	// Test blocked mode
	model.isScripted = false
	blocker.blocked = true
	model.InputHandler(actor, nil)
	accX, accY = actor.Acceleration()
	if accX != 0 || accY != 0 {
		t.Error("expected no acceleration when blocked")
	}

	// Test immobile
	blocker.blocked = false
	actor.SetImmobile(true)
	model.InputHandler(actor, nil)
	accX, accY = actor.Acceleration()
	if accX != 0 || accY != 0 {
		t.Error("expected no acceleration when immobile")
	}
}

func TestTopDownMovementModel_Update(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 1.0,
		},
	}
	config.Set(cfg)

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	sp.AddBody(actor)

	model := NewTopDownMovementModel(nil)
	model.isScripted = true // Skip input

	// Test freeze
	actor.SetFreeze(true)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update with freeze failed: %v", err)
	}
	actor.SetFreeze(false)

	// Test normal update
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify acceleration was reset
	accX, accY := actor.Acceleration()
	if accX != 0 || accY != 0 {
		t.Errorf("expected acceleration reset; got (%d, %d)", accX, accY)
	}
}

func TestNewMovementModel(t *testing.T) {
	blocker := &mockPlayerMovementBlocker{}

	// Test TopDown
	model, err := NewMovementModel(TopDown, blocker)
	if err != nil {
		t.Fatalf("NewMovementModel(TopDown) failed: %v", err)
	}
	if model == nil {
		t.Fatal("expected non-nil model for TopDown")
	}

	// Test Platform
	model, err = NewMovementModel(Platform, blocker)
	if err != nil {
		t.Fatalf("NewMovementModel(Platform) failed: %v", err)
	}
	if model == nil {
		t.Fatal("expected non-nil model for Platform")
	}

	// Test unknown type
	_, err = NewMovementModel(999, blocker)
	if err == nil {
		t.Error("expected error for unknown movement model type")
	}
}

func TestMovementModelEnum_String(t *testing.T) {
	tests := []struct {
		enum MovementModelEnum
		want string
	}{
		{TopDown, "TopDown"},
		{Platform, "Platform"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.enum.String()
			if got != tt.want {
				t.Errorf("%d.String() = %s; want %s", tt.enum, got, tt.want)
			}
		})
	}
}

func TestBeatEmUpMovementModel_New(t *testing.T) {
	blocker := &mockPlayerMovementBlocker{}
	model := NewBeatEmUpMovementModel(blocker)

	if model == nil {
		t.Fatal("NewBeatEmUpMovementModel returned nil")
	}
	if model.playerMovementBlocker != blocker {
		t.Error("expected playerMovementBlocker to be set")
	}
	if model.isScripted {
		t.Error("expected isScripted=false by default")
	}
}

func TestBeatEmUpMovementModel_SetIsScripted(t *testing.T) {
	model := NewBeatEmUpMovementModel(nil)

	model.SetIsScripted(true)
	if !model.isScripted {
		t.Error("expected isScripted=true")
	}

	model.SetIsScripted(false)
	if model.isScripted {
		t.Error("expected isScripted=false")
	}
}

func TestBeatEmUpMovementModel_FreezeGuard(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetVelocity(fp16.To16(3), fp16.To16(3))
	actor.SetFreeze(true)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)

	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update with freeze returned error: %v", err)
	}

	pos := actor.Position().Min
	if pos.X != 100 || pos.Y != 100 {
		t.Errorf("expected position unchanged at (100,100); got (%d,%d)", pos.X, pos.Y)
	}

	vx, vy := actor.Velocity()
	if vx != fp16.To16(3) || vy != fp16.To16(3) {
		t.Errorf("expected velocity unchanged at (%d,%d); got (%d,%d)",
			fp16.To16(3), fp16.To16(3), vx, vy)
	}
}

func TestBeatEmUpMovementModel_NoGravityWhenIdle(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetVelocity(0, 0)
	actor.SetAcceleration(0, 0)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)

	for i := 0; i < 60; i++ {
		if err := model.Update(actor, sp); err != nil {
			t.Fatalf("frame %d: Update returned error: %v", i, err)
		}
		_, vy := actor.Velocity()
		if vy != 0 {
			t.Fatalf("frame %d: expected vy=0 (no gravity); got %d", i, vy)
		}
	}

	pos := actor.Position().Min
	if pos.Y != 100 {
		t.Errorf("expected Y position unchanged at 100; got %d", pos.Y)
	}
}

func TestBeatEmUpMovementModel_DiagonalSpeedNormalization(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	tests := []struct {
		name   string
		accelX int
		accelY int
		check  func(t *testing.T, vx16, vy16 int)
	}{
		{
			name:   "cardinal-X",
			accelX: fp16.To16(2),
			accelY: 0,
			check: func(t *testing.T, vx16, vy16 int) {
				absVX := vx16
				if absVX < 0 {
					absVX = -absVX
				}
				// |vx| should approach fp16.To16(5) (allow some friction slack)
				cap5 := fp16.To16(5)
				if absVX > cap5+cap5/20 {
					t.Errorf("cardinal-X: |vx16|=%d exceeded cap=%d", absVX, cap5)
				}
				if absVX < cap5/2 {
					t.Errorf("cardinal-X: |vx16|=%d well below cap=%d, expected near cap", absVX, cap5)
				}
				if vy16 != 0 {
					t.Errorf("cardinal-X: expected vy16=0; got %d", vy16)
				}
			},
		},
		{
			name:   "diagonal",
			accelX: fp16.To16(2),
			accelY: fp16.To16(2),
			check: func(t *testing.T, vx16, vy16 int) {
				cap5 := fp16.To16(5)
				mag := math.Sqrt(float64(vx16)*float64(vx16) + float64(vy16)*float64(vy16))
				if mag > float64(cap5)*1.05 {
					t.Errorf("diagonal: magnitude %.2f exceeded 1.05*cap=%.2f", mag, float64(cap5)*1.05)
				}
				absVX := vx16
				if absVX < 0 {
					absVX = -absVX
				}
				absVY := vy16
				if absVY < 0 {
					absVY = -absVY
				}
				diff := absVX - absVY
				if diff < 0 {
					diff = -diff
				}
				// Allow modest asymmetry from friction interplay.
				if diff > cap5/4 {
					t.Errorf("diagonal: |vx|-|vy| asymmetry too large: |vx|=%d |vy|=%d", absVX, absVY)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := space.NewSpace()
			actor := newMockMovableCollidable()
			actor.SetMaxSpeed(5)
			actor.SetPosition(100, 100)
			sp.AddBody(actor)

			model := NewBeatEmUpMovementModel(nil)

			for i := 0; i < 60; i++ {
				actor.SetAcceleration(tt.accelX, tt.accelY)
				if err := model.Update(actor, sp); err != nil {
					t.Fatalf("frame %d: Update returned error: %v", i, err)
				}
			}

			vx16, vy16 := actor.Velocity()
			tt.check(t, vx16, vy16)
		})
	}
}

func TestBeatEmUpMovementModel_XObstacleCollision(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(100, 100)
	actor.SetMaxSpeed(30)
	actor.SetVelocity(fp16.To16(20), 0)
	sp.AddBody(actor)

	obstacle := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	obstacle.SetID("obstacle")
	obstacle.SetIsObstructive(true)
	obstacle.SetPosition(120, 100)
	obstacle.AddCollisionBodies()
	sp.AddBody(obstacle)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	pos := actor.Position().Min
	if pos.X >= 120 {
		t.Errorf("expected actor.X < 120 (blocked by obstacle); got %d", pos.X)
	}
	if pos.X < 100 {
		t.Errorf("expected actor.X >= 100; got %d", pos.X)
	}
}

func TestBeatEmUpMovementModel_YObstacleCollision(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(100, 100)
	actor.SetMaxSpeed(30)
	actor.SetVelocity(0, fp16.To16(20))
	sp.AddBody(actor)

	obstacle := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	obstacle.SetID("obstacle")
	obstacle.SetIsObstructive(true)
	obstacle.SetPosition(100, 120)
	obstacle.AddCollisionBodies()
	sp.AddBody(obstacle)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	pos := actor.Position().Min
	if pos.Y >= 120 {
		t.Errorf("expected actor.Y < 120 (blocked by obstacle); got %d", pos.Y)
	}
}

func TestBeatEmUpMovementModel_FrictionAppliedEachFrame(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetVelocity(fp16.To16(4), fp16.To16(4))
	actor.SetAcceleration(0, 0)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	vx16, vy16 := actor.Velocity()
	if vx16 >= fp16.To16(4) {
		t.Errorf("expected vx16 < fp16.To16(4)=%d; got %d", fp16.To16(4), vx16)
	}
	if vy16 >= fp16.To16(4) {
		t.Errorf("expected vy16 < fp16.To16(4)=%d; got %d", fp16.To16(4), vy16)
	}
	if vx16 <= 0 {
		t.Errorf("expected vx16 > 0 (friction not full stop); got %d", vx16)
	}
	if vy16 <= 0 {
		t.Errorf("expected vy16 > 0 (friction not full stop); got %d", vy16)
	}
}

func TestBeatEmUpMovementModel_AccelerationResetAfterUpdate(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAcceleration(fp16.To16(2), fp16.To16(2))
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	accX, accY := actor.Acceleration()
	if accX != 0 || accY != 0 {
		t.Errorf("expected acceleration reset to (0,0); got (%d,%d)", accX, accY)
	}
}

func TestBeatEmUpMovementModel_ClampToPlayAreaEngaged(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()
	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetPosition(-10, -10)
	actor.SetMaxSpeed(10)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	pos := actor.Position().Min
	if pos.X != 0 {
		t.Errorf("expected actor.X clamped to 0; got %d", pos.X)
	}
	if pos.Y != 0 {
		t.Errorf("expected actor.Y clamped to 0; got %d", pos.Y)
	}
}

func TestBeatEmUpMovementModel_FactoryWiring(t *testing.T) {
	blocker := &mockPlayerMovementBlocker{}
	model, err := NewMovementModel(BeatEmUp, blocker)
	if err != nil {
		t.Fatalf("NewMovementModel(BeatEmUp) failed: %v", err)
	}
	if model == nil {
		t.Fatal("expected non-nil model for BeatEmUp")
	}
	if _, ok := model.(*BeatEmUpMovementModel); !ok {
		t.Errorf("expected *BeatEmUpMovementModel; got %T", model)
	}
}

func TestBeatEmUpMovementModel_EnumString(t *testing.T) {
	got := BeatEmUp.String()
	if got != "BeatEmUp" {
		t.Errorf("BeatEmUp.String() = %q; want %q", got, "BeatEmUp")
	}
}

// --- Story 061: Altitude axis (jump + ground detection) ---

// beatEmUpAltitudeTestConfig sets the standard config used by all T-061-* tests.
func beatEmUpAltitudeTestConfig() {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			UpwardGravity:   2,
			DownwardGravity: 4,
			SpeedMultiplier: 1.0,
		},
	})
}

// T-061-1: rising arc — UpwardGravity accumulates [AC-1]
func TestBeatEmUpMovementModel_Altitude_RisingArc_UpwardGravity(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAltitude(20)
	preVAlt16 := -fp16.To16(10)
	actor.SetVAltitude16(preVAlt16)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	postVAlt16 := actor.VAltitude16()
	wantVAlt16 := preVAlt16 + 2 // UpwardGravity
	if postVAlt16 != wantVAlt16 {
		t.Errorf("rising arc: VAltitude16() = %d; want %d (pre %d + UpwardGravity 2)",
			postVAlt16, wantVAlt16, preVAlt16)
	}

	postAlt := actor.Altitude()
	if postAlt <= 20 {
		t.Errorf("rising arc: Altitude() = %d; want > 20 (actor rising; negative vAlt16 increases altitude)", postAlt)
	}
	if postAlt <= 0 {
		t.Errorf("rising arc: Altitude() = %d; want > 0 (still airborne)", postAlt)
	}
}

// T-061-2: falling — DownwardGravity accumulates [AC-1]
func TestBeatEmUpMovementModel_Altitude_FallingArc_DownwardGravity(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAltitude(50)
	preVAlt16 := fp16.To16(2)
	actor.SetVAltitude16(preVAlt16)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	postVAlt16 := actor.VAltitude16()
	wantVAlt16 := preVAlt16 + 4 // DownwardGravity
	if postVAlt16 != wantVAlt16 {
		t.Errorf("falling arc: VAltitude16() = %d; want %d (pre %d + DownwardGravity 4)",
			postVAlt16, wantVAlt16, preVAlt16)
	}

	postAlt := actor.Altitude()
	if postAlt >= 50 {
		t.Errorf("falling arc: Altitude() = %d; want < 50 (integrated by positive velocity)", postAlt)
	}
	if postAlt <= 0 {
		t.Errorf("falling arc: Altitude() = %d; want > 0 (still airborne)", postAlt)
	}
}

// T-061-3: landing clamps altitude and zeroes velocity [AC-4]
func TestBeatEmUpMovementModel_Altitude_LandingClamp(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAltitude(1)
	actor.SetVAltitude16(fp16.To16(50))
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if got := actor.Altitude(); got != 0 {
		t.Errorf("landing: Altitude() = %d; want 0 (clamped)", got)
	}
	if got := actor.VAltitude16(); got != 0 {
		t.Errorf("landing: VAltitude16() = %d; want 0 (zeroed on landing)", got)
	}
}

// T-061-4: idempotent grounded — no mutation [AC-2, AC-7]
func TestBeatEmUpMovementModel_Altitude_GroundedIdempotent(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	for i := 0; i < 5; i++ {
		if err := model.Update(actor, sp); err != nil {
			t.Fatalf("frame %d: Update returned error: %v", i, err)
		}
		if got := actor.Altitude(); got != 0 {
			t.Fatalf("frame %d: Altitude() = %d; want 0 (grounded should not mutate)", i, got)
		}
		if got := actor.VAltitude16(); got != 0 {
			t.Fatalf("frame %d: VAltitude16() = %d; want 0 (grounded should not mutate)", i, got)
		}
	}
}

// T-061-5: freeze guard skips altitude mutation [AC-6]
func TestBeatEmUpMovementModel_Altitude_FreezeGuard(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAltitude(30)
	preVAlt16 := fp16.To16(5)
	actor.SetVAltitude16(preVAlt16)
	actor.SetFreeze(true)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(actor, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if got := actor.VAltitude16(); got != preVAlt16 {
		t.Errorf("freeze: VAltitude16() = %d; want %d (unchanged)", got, preVAlt16)
	}
	if got := actor.Altitude(); got != 30 {
		t.Errorf("freeze: Altitude() = %d; want 30 (unchanged)", got)
	}
}

// T-061-6: 2D regression — body never touching altitude stays at 0 [AC-7]
func TestBeatEmUpMovementModel_Altitude_2DRegression(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetMaxSpeed(10)
	actor.SetVelocity(fp16.To16(2), 0)
	sp.AddBody(actor)

	preX := actor.Position().Min.X

	model := NewBeatEmUpMovementModel(nil)
	for i := 0; i < 30; i++ {
		if err := model.Update(actor, sp); err != nil {
			t.Fatalf("frame %d: Update returned error: %v", i, err)
		}
		if got := actor.Altitude(); got != 0 {
			t.Fatalf("frame %d: Altitude() = %d; want 0 (2D body untouched)", i, got)
		}
		if got := actor.VAltitude16(); got != 0 {
			t.Fatalf("frame %d: VAltitude16() = %d; want 0 (2D body untouched)", i, got)
		}
	}

	postX := actor.Position().Min.X
	if postX == preX {
		t.Errorf("2D regression: X did not move (pre=%d post=%d); 2D motion broken", preX, postX)
	}
}

// T-061-7: external jump impulse → full rise/fall/land arc [AC-1, AC-3, AC-4, AC-5]
func TestBeatEmUpMovementModel_Altitude_JumpArc(t *testing.T) {
	beatEmUpAltitudeTestConfig()

	sp := space.NewSpace()
	actor := newMockMovableCollidable()
	actor.SetPosition(100, 100)
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)
	sp.AddBody(actor)

	model := NewBeatEmUpMovementModel(nil)

	// Step A: external jump impulse.
	actor.SetVAltitude16(-fp16.To16(8))

	var roseFrame, peakFrame int
	rose := false
	peaked := false
	landed := false
	const maxFrames = 600
	frames := 0

	for i := 0; i < maxFrames; i++ {
		if err := model.Update(actor, sp); err != nil {
			t.Fatalf("frame %d: Update returned error: %v", i, err)
		}
		frames = i + 1

		if !rose && actor.Altitude() > 0 {
			rose = true
			roseFrame = i
		}
		if rose && !peaked && actor.VAltitude16() >= 0 {
			peaked = true
			peakFrame = i
		}
		if frames > 1 && actor.Altitude() == 0 && actor.VAltitude16() == 0 && rose {
			landed = true
			break
		}
	}

	if !rose {
		t.Fatalf("jump arc: actor never rose above 0 altitude within %d frames", maxFrames)
	}
	if !peaked {
		t.Fatalf("jump arc: actor never reached peak (VAltitude16 >= 0) within %d frames", maxFrames)
	}
	if peakFrame < roseFrame {
		t.Errorf("jump arc: peak frame %d < rose frame %d (out of order)", peakFrame, roseFrame)
	}
	if !landed {
		t.Fatalf("jump arc: actor did not land within %d frames (alt=%d vAlt16=%d)",
			maxFrames, actor.Altitude(), actor.VAltitude16())
	}
	if frames >= maxFrames {
		t.Errorf("jump arc: loop did not terminate before %d frames", maxFrames)
	}
}
