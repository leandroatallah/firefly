package kitskills

import (
	"testing"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	"github.com/hajimehoshi/ebiten/v2"
)

// mockMovableCollidable for skill tests
type mockMovableCollidable struct {
	*bodyphysics.ObstacleRect
	isDuckingFunc func() bool
}

func (m *mockMovableCollidable) IsDucking() bool {
	if m.isDuckingFunc != nil {
		return m.isDuckingFunc()
	}
	return false
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

// Test DashSkill
func TestDashSkill_New(t *testing.T) {
	d := NewDashSkill()

	if d == nil {
		t.Fatal("NewDashSkill returned nil")
	}
	if d.State() != skill.StateReady {
		t.Errorf("expected state Ready; got %s", d.State())
	}
	if d.Duration() <= 0 {
		t.Error("expected positive duration")
	}
	if d.Cooldown() <= 0 {
		t.Error("expected positive cooldown")
	}
	if d.ActivationKey() != ebiten.KeyShift {
		t.Errorf("expected activationKey Shift; got %v", d.ActivationKey())
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
	d.SetState(skill.StateReady)
	d.Update(actor, model)
	if d.State() != skill.StateReady {
		t.Errorf("expected state to remain Ready; got %s", d.State())
	}

	// Test Active state
	d.SetState(skill.StateActive)
	d.SetTimer(2)
	d.Update(actor, model)
	if d.State() != skill.StateActive {
		t.Errorf("expected state to remain Active; got %s", d.State())
	}

	// Active state expires
	d.SetTimer(0)
	d.Update(actor, model)
	if d.State() != skill.StateCooldown {
		t.Errorf("expected state Cooldown; got %s", d.State())
	}

	// Cooldown state
	d.SetState(skill.StateCooldown)
	d.SetTimer(2)
	d.Update(actor, model)
	if d.State() != skill.StateCooldown {
		t.Errorf("expected state to remain Cooldown; got %s", d.State())
	}

	// Cooldown expires
	d.SetTimer(0)
	d.Update(actor, model)
	if d.State() != skill.StateReady {
		t.Errorf("expected state Ready after cooldown; got %s", d.State())
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

func TestDashSkill_tryActivate_Ready(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.SetState(skill.StateReady)
	model.SetOnGround(true)

	d.tryActivate(actor, model, sp)

	if d.State() != skill.StateActive {
		t.Errorf("expected state Active; got %s", d.State())
	}
	if d.Timer() != d.Duration() {
		t.Errorf("expected timer=duration; got %d", d.Timer())
	}
}

func TestDashSkill_tryActivate_NotReady(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.SetState(skill.StateActive) // Already active

	d.tryActivate(actor, model, sp)

	if d.State() != skill.StateActive {
		t.Errorf("expected state to remain Active; got %s", d.State())
	}
}

func TestDashSkill_tryActivate_AirDash(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.SetState(skill.StateReady)
	model.SetOnGround(false)

	// First air dash should work
	d.tryActivate(actor, model, sp)
	if d.State() != skill.StateActive {
		t.Errorf("expected first air dash to work; got state %s", d.State())
	}
	if !d.airDashUsed {
		t.Error("expected airDashUsed to be set")
	}

	// Reset and try second air dash (should fail)
	d.SetState(skill.StateReady)
	d.tryActivate(actor, model, sp)
	if d.State() == skill.StateActive {
		t.Error("expected second air dash to fail")
	}
}

func TestDashSkill_tryActivate_NoAirDash(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	d := NewDashSkill()
	d.canAirDash = false
	d.SetState(skill.StateReady)
	model.SetOnGround(false)

	d.tryActivate(actor, model, sp)

	if d.State() == skill.StateActive {
		t.Error("expected air dash to fail when canAirDash=false")
	}
}

// Test JumpSkill
func TestJumpSkill_New(t *testing.T) {
	j := NewJumpSkill()

	if j == nil {
		t.Fatal("NewJumpSkill returned nil")
	}
	if j.State() != skill.StateReady {
		t.Errorf("expected state Ready; got %s", j.State())
	}
	if j.ActivationKey() != ebiten.KeySpace {
		t.Errorf("expected activationKey Space; got %v", j.ActivationKey())
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

func TestJumpSkill_HandleInput_Blocked(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(&mockPlayerMovementBlocker{blocked: true})
	sp := space.NewSpace()

	j := NewJumpSkill()

	// Should not activate when blocked
	j.HandleInput(actor, model, sp)

	if j.State() != skill.StateReady {
		t.Errorf("expected state to remain Ready; got %s", j.State())
	}
}

// Test HorizontalMovementSkill
func TestHorizontalMovementSkill_New(t *testing.T) {
	s := NewHorizontalMovementSkill()

	if s == nil {
		t.Fatal("NewHorizontalMovementSkill returned nil")
	}
	if s.State() != skill.StateReady {
		t.Errorf("expected state Ready; got %s", s.State())
	}
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

// Test timing integration
func TestDashSkill_TimingConstants(t *testing.T) {
	d := NewDashSkill()

	// Verify timing constants are reasonable
	expectedDuration := timing.FromDuration(200 * time.Millisecond)
	expectedCooldown := timing.FromDuration(750 * time.Millisecond)

	if d.Duration() != expectedDuration {
		t.Errorf("expected duration %d; got %d", expectedDuration, d.Duration())
	}
	if d.Cooldown() != expectedCooldown {
		t.Errorf("expected cooldown %d; got %d", expectedCooldown, d.Cooldown())
	}
}

func TestHorizontalMovementSkillImmobileBehavior(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{HorizontalInertia: 0}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetVelocity(fp16.To16(5), fp16.To16(3))
	actor.SetAcceleration(fp16.To16(2), fp16.To16(1))
	actor.SetImmobile(true)

	s := NewHorizontalMovementSkill()
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	vx, vy := actor.Velocity()
	if vx != 0 {
		t.Fatalf("expected vx=0 when immobile, got %d", vx)
	}
	if vy != fp16.To16(3) {
		t.Fatalf("vy should remain unchanged, got %d", vy)
	}
}

func TestDashSkillIntegratesWithMovementModel(t *testing.T) {
	cfg := &config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			MaxFallSpeed:    128,
		},
	}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetFaceDirection(animation.FaceDirectionRight)

	sp := space.NewSpace()
	model := movement.NewPlatformMovementModel(nil)

	d := NewDashSkill()
	// Force active state for a few frames
	d.SetState(skill.StateActive)
	d.SetTimer(2)

	// First update should set dash active and move horizontally
	prevX, _ := actor.GetPositionMin()
	d.Update(actor, model)
	_ = model.Update(actor, sp)

	newX, _ := actor.GetPositionMin()
	if newX == prevX {
		t.Fatalf("expected actor to move horizontally due to dash, prevX=%d newX=%d", prevX, newX)
	}

	// Advance to end active -> cooldown
	d.Update(actor, model)
	_ = model.Update(actor, sp) // should keep dash while timer > 0

	// Next tick transitions to cooldown and clears dash on model
	d.Update(actor, model)
	_ = model.Update(actor, sp)
}
