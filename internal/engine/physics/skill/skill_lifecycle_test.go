package skill

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	contractsbody "github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

// mockMovableCollidable for skill tests
type mockMovableCollidable struct {
	*bodyphysics.ObstacleRect
}

// mockPlayerMovementBlocker implements movement.PlayerMovementBlocker for testing
type mockPlayerMovementBlocker struct {
	blocked bool
}

func (m *mockPlayerMovementBlocker) IsMovementBlocked() bool {
	return m.blocked
}

func newMockMovableCollidable() *mockMovableCollidable {
	rect := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	rect.SetID("mock")
	rect.SetSpeed(2)
	rect.SetMaxSpeed(10)
	rect.SetJumpForceMultiplier(1.0)
	return &mockMovableCollidable{
		ObstacleRect: rect,
	}
}

// Test SkillBase
func TestSkillBase_InitialState(t *testing.T) {
	sb := &SkillBase{
		state:    StateReady,
		duration: 10,
		cooldown: 20,
	}

	if sb.state != StateReady {
		t.Errorf("expected state Ready; got %s", sb.state)
	}
}

func TestSkillBase_Update(t *testing.T) {
	sb := &SkillBase{}
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	// Should not panic
	sb.Update(actor, model)
}

func TestSkillBase_IsActive(t *testing.T) {
	sb := &SkillBase{state: StateReady}

	if sb.IsActive() {
		t.Error("expected IsActive=false for Ready state")
	}

	sb.state = StateActive
	if !sb.IsActive() {
		t.Error("expected IsActive=true for Active state")
	}

	sb.state = StateCooldown
	if sb.IsActive() {
		t.Error("expected IsActive=false for Cooldown state")
	}
}

// Test DashSkill
func TestDashSkill_New(t *testing.T) {
	d := NewDashSkill()

	if d == nil {
		t.Fatal("NewDashSkill returned nil")
	}
	if d.state != StateReady {
		t.Errorf("expected state Ready; got %s", d.state)
	}
	if d.duration <= 0 {
		t.Error("expected positive duration")
	}
	if d.cooldown <= 0 {
		t.Error("expected positive cooldown")
	}
	if d.activationKey != ebiten.KeyShift {
		t.Errorf("expected activationKey Shift; got %v", d.activationKey)
	}
	if !d.canAirDash {
		t.Error("expected canAirDash=true")
	}
}

func TestDashSkill_ActivationKey(t *testing.T) {
	d := NewDashSkill()
	if d.ActivationKey() != ebiten.KeyShift {
		t.Error("expected Shift key")
	}
}

func TestDashSkill_Update_StateTransitions(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	d := NewDashSkill()

	// Test Ready state - no changes
	d.state = StateReady
	d.Update(actor, model)
	if d.state != StateReady {
		t.Errorf("expected state to remain Ready; got %s", d.state)
	}

	// Test Active state
	d.state = StateActive
	d.timer = 2
	d.Update(actor, model)
	if d.state != StateActive {
		t.Errorf("expected state to remain Active; got %s", d.state)
	}
	// Note: We can't check model.dashActive directly as it's unexported
	// The integration test verifies the behavior

	// Active state expires
	d.timer = 0
	d.Update(actor, model)
	if d.state != StateCooldown {
		t.Errorf("expected state Cooldown; got %s", d.state)
	}

	// Cooldown state
	d.state = StateCooldown
	d.timer = 2
	d.Update(actor, model)
	if d.state != StateCooldown {
		t.Errorf("expected state to remain Cooldown; got %s", d.state)
	}

	// Cooldown expires
	d.timer = 0
	d.Update(actor, model)
	if d.state != StateReady {
		t.Errorf("expected state Ready after cooldown; got %s", d.state)
	}
}

func TestDashSkill_Update_AirDashReset(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	d := NewDashSkill()
	d.airDashUsed = true

	// Simulate landing
	model.SetOnGround(true)
	d.Update(actor, model)

	if d.airDashUsed {
		t.Error("expected airDashUsed to be reset on ground")
	}
}

func TestDashSkill_Update_FaceDirection(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	d := NewDashSkill()
	d.state = StateActive
	d.timer = 2

	// Test facing right
	actor.SetFaceDirection(animation.FaceDirectionRight)
	d.Update(actor, model)
	// Note: Can't check dashVelocityX directly as it's unexported
	// Integration test verifies behavior

	// Test facing left
	actor.SetFaceDirection(animation.FaceDirectionLeft)
	d.Update(actor, model)
	// Note: Can't check dashVelocityX directly as it's unexported
}

func TestDashSkill_tryActivate_Ready(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.state = StateReady
	model.SetOnGround(true)

	d.tryActivate(actor, model, sp)

	if d.state != StateActive {
		t.Errorf("expected state Active; got %s", d.state)
	}
	if d.timer != d.duration {
		t.Errorf("expected timer=duration; got %d", d.timer)
	}
}

func TestDashSkill_tryActivate_NotReady(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.state = StateActive // Already active

	d.tryActivate(actor, model, sp)

	if d.state != StateActive {
		t.Errorf("expected state to remain Active; got %s", d.state)
	}
}

func TestDashSkill_tryActivate_AirDash(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.state = StateReady
	model.SetOnGround(false)

	// First air dash should work
	d.tryActivate(actor, model, sp)
	if d.state != StateActive {
		t.Errorf("expected first air dash to work; got state %s", d.state)
	}
	if !d.airDashUsed {
		t.Error("expected airDashUsed to be set")
	}

	// Reset and try second air dash (should fail)
	d.state = StateReady
	d.tryActivate(actor, model, sp)
	if d.state == StateActive {
		t.Error("expected second air dash to fail")
	}
}

func TestDashSkill_tryActivate_NoAirDash(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.canAirDash = false
	d.state = StateReady
	model.SetOnGround(false)

	d.tryActivate(actor, model, sp)

	if d.state == StateActive {
		t.Error("expected air dash to fail when canAirDash=false")
	}
}

// Test JumpSkill
func TestJumpSkill_New(t *testing.T) {
	j := NewJumpSkill()

	if j == nil {
		t.Fatal("NewJumpSkill returned nil")
	}
	if j.state != StateReady {
		t.Errorf("expected state Ready; got %s", j.state)
	}
	if j.activationKey != ebiten.KeySpace {
		t.Errorf("expected activationKey Space; got %v", j.activationKey)
	}
}

func TestJumpSkill_ActivationKey(t *testing.T) {
	j := NewJumpSkill()
	if j.ActivationKey() != ebiten.KeySpace {
		t.Error("expected Space key")
	}
}

func TestJumpSkill_Update_CoyoteTime(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			CoyoteTimeFrames: 5,
			DownwardGravity:  4,
			MaxFallSpeed:     128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	j := NewJumpSkill()

	// Simulate being on ground
	model.SetOnGround(true)
	j.Update(actor, model)

	if j.coyoteTimeCounter != 5 {
		t.Errorf("expected coyoteTimeCounter=5; got %d", j.coyoteTimeCounter)
	}

	// Simulate leaving ground
	model.SetOnGround(false)
	j.Update(actor, model)

	if j.coyoteTimeCounter != 4 {
		t.Errorf("expected coyoteTimeCounter=4 after one frame; got %d", j.coyoteTimeCounter)
	}
}

func TestJumpSkill_Update_JumpBuffer(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpBufferFrames: 5,
			DownwardGravity:  4,
			MaxFallSpeed:     128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	j := NewJumpSkill()
	j.jumpBufferCounter = 3

	// Decrement buffer counter
	j.Update(actor, model)

	if j.jumpBufferCounter != 2 {
		t.Errorf("expected jumpBufferCounter=2; got %d", j.jumpBufferCounter)
	}
}

func TestJumpSkill_tryActivate_OnGround(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpForce:       8,
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	j := NewJumpSkill()
	model.SetOnGround(true)

	j.tryActivate(actor, model, sp)

	_, vy := actor.Velocity()
	if vy >= 0 {
		t.Errorf("expected negative (upward) velocity after jump; got %d", vy)
	}
	if model.OnGround() {
		t.Error("expected model to be set airborne")
	}
	if j.coyoteTimeCounter != 0 {
		t.Errorf("expected coyoteTimeCounter=0; got %d", j.coyoteTimeCounter)
	}
	if j.jumpBufferCounter != 0 {
		t.Errorf("expected jumpBufferCounter=0; got %d", j.jumpBufferCounter)
	}
}

func TestJumpSkill_tryActivate_CoyoteTime(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpForce:        8,
			CoyoteTimeFrames: 5,
			DownwardGravity:  4,
			MaxFallSpeed:     128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	j := NewJumpSkill()
	j.coyoteTimeCounter = 3 // Has coyote time
	model.SetOnGround(false)

	j.tryActivate(actor, model, sp)

	_, vy := actor.Velocity()
	if vy >= 0 {
		t.Errorf("expected negative velocity with coyote time; got %d", vy)
	}
}

func TestJumpSkill_tryActivate_Airborne(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpForce:        8,
			JumpBufferFrames: 5,
			DownwardGravity:  4,
			MaxFallSpeed:     128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	j := NewJumpSkill()
	model.SetOnGround(false)
	j.coyoteTimeCounter = 0 // No coyote time

	j.tryActivate(actor, model, sp)

	if j.jumpBufferCounter != 5 {
		t.Errorf("expected jumpBufferCounter=5; got %d", j.jumpBufferCounter)
	}
}

func TestJumpSkill_tryActivate_ZeroJumpForce(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpForce:       0,
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	j := NewJumpSkill()
	model.SetOnGround(true)

	j.tryActivate(actor, model, sp)

	_, vy := actor.Velocity()
	if vy != 0 {
		t.Errorf("expected zero velocity with zero jump force; got %d", vy)
	}
}

func TestJumpSkill_tryActivate_OnJumpCallback(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpForce:       8,
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	j := NewJumpSkill()
	called := false
	j.OnJump = func(_ contractsbody.MovableCollidable) {
		called = true
	}
	model.SetOnGround(true)

	j.tryActivate(actor, model, sp)

	if !called {
		t.Error("expected OnJump callback to be called")
	}
}

func TestJumpSkill_handleCoyoteAndJumpBuffering_Landing(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			JumpBufferFrames: 5,
			DownwardGravity:  4,
			MaxFallSpeed:     128,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	j := NewJumpSkill()
	j.jumpBufferCounter = 2

	// Simulate landing with buffered jump
	wasOnGround := false
	model.SetOnGround(true)

	// Just verify it doesn't panic
	j.handleCoyoteAndJumpBuffering(actor, model, wasOnGround)
}

func TestJumpSkill_HandleInput_Blocked(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(&mockPlayerMovementBlocker{blocked: true})
	sp := space.NewSpace()

	j := NewJumpSkill()

	// Should not activate when blocked
	j.HandleInput(actor, model, sp)

	if j.state != StateReady {
		t.Errorf("expected state to remain Ready; got %s", j.state)
	}
}

// Test HorizontalMovementSkill
func TestHorizontalMovementSkill_New(t *testing.T) {
	s := NewHorizontalMovementSkill()

	if s == nil {
		t.Fatal("NewHorizontalMovementSkill returned nil")
	}
	if s.state != StateReady {
		t.Errorf("expected state Ready; got %s", s.state)
	}
}

func TestHorizontalMovementSkill_ActivationKey(t *testing.T) {
	s := NewHorizontalMovementSkill()
	// Returns zero value since not set
	if s.ActivationKey() != ebiten.Key(0) {
		t.Errorf("expected zero key; got %v", s.ActivationKey())
	}
}

func TestHorizontalMovementSkill_Update(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	s := NewHorizontalMovementSkill()
	s.Update(actor, model)

	// Should not panic
}

func TestHorizontalMovementSkill_HandleInput_Blocked(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(&mockPlayerMovementBlocker{blocked: true})

	s := NewHorizontalMovementSkill()
	s.HandleInput(actor, model, nil)

	// Should not modify velocity when blocked
	vx, vy := actor.Velocity()
	if vx != 0 || vy != 0 {
		t.Errorf("expected zero velocity when blocked; got (%d, %d)", vx, vy)
	}
}

func TestHorizontalMovementSkill_HandleInput_Immobile(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			HorizontalInertia: 0,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	actor.SetImmobile(true)
	actor.SetVelocity(fp16.To16(5), fp16.To16(3))
	actor.SetAcceleration(fp16.To16(2), fp16.To16(1))

	model := movement.NewPlatformMovementModel(nil)

	s := NewHorizontalMovementSkill()
	s.HandleInput(actor, model, nil)

	vx, vy := actor.Velocity()
	if vx != 0 {
		t.Errorf("expected vx=0 when immobile; got %d", vx)
	}
	if vy != fp16.To16(3) {
		t.Errorf("expected vy unchanged; got %d", vy)
	}

	accX, accY := actor.Acceleration()
	if accX != 0 {
		t.Errorf("expected accX=0 when immobile; got %d", accX)
	}
	if accY != fp16.To16(1) {
		t.Errorf("expected accY unchanged; got %d", accY)
	}
}

func TestHorizontalMovementSkill_HandleInput_WithInertia(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			HorizontalInertia: 1.0,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	s := NewHorizontalMovementSkill()

	// Simulate key press (we can't actually press keys, so test the logic path)
	// When inertia > 0, it uses OnMoveLeft/OnMoveRight
	// This is tested via the input package which we can't easily mock
	// So we just verify no panic
	s.HandleInput(actor, model, nil)
}

func TestHorizontalMovementSkill_HandleInput_WithoutInertia(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			HorizontalInertia: 0,
		},
	}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	s := NewHorizontalMovementSkill()

	// Without inertia, it sets velocity directly based on input
	// Again, we can't simulate key presses easily
	// Just verify no panic
	s.HandleInput(actor, model, nil)
}

// Test timing integration
func TestDashSkill_TimingConstants(t *testing.T) {
	d := NewDashSkill()

	// Verify timing constants are reasonable
	expectedDuration := timing.FromDuration(133 * time.Millisecond)
	expectedCooldown := timing.FromDuration(750 * time.Millisecond)

	if d.duration != expectedDuration {
		t.Errorf("expected duration %d; got %d", expectedDuration, d.duration)
	}
	if d.cooldown != expectedCooldown {
		t.Errorf("expected cooldown %d; got %d", expectedCooldown, d.cooldown)
	}
}
