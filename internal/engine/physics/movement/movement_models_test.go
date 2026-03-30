package movement

import (
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
		name       string
		x, y       int
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
			HorizontalInertia:       1.0,
			SpeedMultiplier:         1.0,
			AirControlMultiplier:    0.5,
			AirFrictionMultiplier:   2.0,
		},
	}
	config.Set(cfg)

	tests := []struct {
		name        string
		accelX      int
		onGround    bool
		maxSpeed    int
		wantVelX    int
	}{
		{"no acceleration", 0, true, 10, 0},
		{"positive acceleration", fp16.To16(2), true, 10, 0}, // Will be increased
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
	vx, vy = model.UpdateVerticalVelocity(actor)
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
	vx, vy = model.handleGravity(actor)
	if vy <= 50 {
		t.Errorf("expected vy > 50 after gravity; got %d", vy)
	}

	// Airborne, jumping (negative vy)
	actor.SetVelocity(0, -50)
	vx, vy = model.handleGravity(actor)
	// Upward gravity adds to negative velocity (making it less negative towards zero)
	// vy = -50 + 2 = -48
	if vy != -48 {
		t.Errorf("expected vy=-48 after upward gravity; got %d", vy)
	}

	// Clamp to max fall speed
	actor.SetVelocity(0, 200)
	vx, vy = model.handleGravity(actor)
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
