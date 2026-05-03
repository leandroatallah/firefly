package kitskills

import (
	"testing"

	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/hajimehoshi/ebiten/v2"
)

// --- OffsetToggler ---

func TestOffsetToggler_NewAndNext(t *testing.T) {
	o := NewOffsetToggler(5)
	if o == nil {
		t.Fatal("NewOffsetToggler returned nil")
	}
	// starts at -5, Next toggles to +5
	got := o.Next()
	if got != 5 {
		t.Errorf("first Next: got %d, want 5", got)
	}
	got = o.Next()
	if got != -5 {
		t.Errorf("second Next: got %d, want -5", got)
	}
}

// --- ShootingSkill.ActivationKey ---

func TestShootingSkill_ActivationKey(t *testing.T) {
	s := NewShootingSkill(nil)
	if s.ActivationKey() != ebiten.KeyX {
		t.Errorf("expected KeyX, got %v", s.ActivationKey())
	}
}

// --- DashSkill.HandleInput ---

func TestDashSkill_HandleInput_TriggersActivation(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()
	model.SetOnGround(true)

	d := NewDashSkill()

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	// First call: dash pressed (edge: not pressed before)
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Dash: true}
	}
	d.HandleInput(actor, model, sp)
	if d.State() != skill.StateActive {
		t.Errorf("expected Active after dash press; got %s", d.State())
	}

	// Second call: dash still held (no re-trigger)
	d.HandleInput(actor, model, sp)
	// state stays Active (tryActivate skipped because dashPressed already true)
}

func TestDashSkill_HandleInput_NoTriggerWhenHeld(t *testing.T) {
	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()
	model.SetOnGround(true)

	d := NewDashSkill()
	d.dashPressed = true // simulate already held

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Dash: true}
	}

	d.HandleInput(actor, model, sp)
	// Should not activate because dashPressed was already true
	if d.State() != skill.StateReady {
		t.Errorf("expected Ready (no re-trigger); got %s", d.State())
	}
}

func TestDashSkill_HandleInput_ReleaseClearsDashPressed(t *testing.T) {
	d := NewDashSkill()
	d.dashPressed = true

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Dash: false}
	}

	d.HandleInput(newMockMovableCollidable(), movement.NewPlatformMovementModel(nil), space.NewSpace())
	if d.dashPressed {
		t.Error("expected dashPressed=false after release")
	}
}

// --- HorizontalMovementSkill.Update and ActivationKey ---

func TestHorizontalMovementSkill_Update(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	s := NewHorizontalMovementSkill()
	// Update should not panic and delegates to SkillBase
	s.Update(actor, model)
}

func TestHorizontalMovementSkill_ActivationKey(t *testing.T) {
	s := NewHorizontalMovementSkill()
	// activationKey is zero-value; just ensure it's callable
	_ = s.ActivationKey()
}

// --- HorizontalMovementSkill.HandleInput: inertia=0 branch and axis-release paths ---

func TestHorizontalMovementSkill_HandleInput_InertiaZero_MoveLeft(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{HorizontalInertia: 0}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetSpeed(3)
	actor.SetHorizontalInertia(-1) // use config value

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	s := NewHorizontalMovementSkill()

	// Press left
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Left: true} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	vx, _ := actor.Velocity()
	if vx >= 0 {
		t.Errorf("expected negative vx when moving left with inertia=0; got %d", vx)
	}
}

func TestHorizontalMovementSkill_HandleInput_InertiaZero_MoveRight(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{HorizontalInertia: 0}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetSpeed(3)
	actor.SetHorizontalInertia(-1)

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	s := NewHorizontalMovementSkill()

	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Right: true} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	vx, _ := actor.Velocity()
	if vx <= 0 {
		t.Errorf("expected positive vx when moving right with inertia=0; got %d", vx)
	}
}

func TestHorizontalMovementSkill_HandleInput_InertiaZero_NoInput(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{HorizontalInertia: 0}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetSpeed(3)
	actor.SetHorizontalInertia(-1)
	actor.SetVelocity(fp16.To16(5), 0)

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	s := NewHorizontalMovementSkill()

	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	vx, _ := actor.Velocity()
	if vx != 0 {
		t.Errorf("expected vx=0 with no input and inertia=0; got %d", vx)
	}
}

func TestHorizontalMovementSkill_HandleInput_AxisRelease(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{HorizontalInertia: 0}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetSpeed(3)
	actor.SetHorizontalInertia(-1)

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	s := NewHorizontalMovementSkill()

	// Press left, then release — exercises the axis.Release(-1) path
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Left: true} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Left: false} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	// Press right, then release — exercises the axis.Release(1) path
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Right: true} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Right: false} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)
}

func TestHorizontalMovementSkill_HandleInput_WithInertia_MoveLeftRight(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{HorizontalInertia: 5}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetSpeed(3)
	actor.SetHorizontalInertia(-1) // use config value (5)

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	s := NewHorizontalMovementSkill()

	// Move left with inertia — exercises OnMoveLeft branch
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Left: true} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)

	// Move right with inertia — exercises OnMoveRight branch
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Right: true} }
	s.HandleInput(actor, movement.NewPlatformMovementModel(nil), nil)
}

// --- JumpSkill.HandleInput: jump-cut path ---

func TestJumpSkill_HandleInput_JumpCutOnRelease(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		DownwardGravity:  4,
		MaxFallSpeed:     128,
	}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	actor.SetVelocity(0, -320) // going up
	model := movement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	j := NewJumpSkill()
	j.SetJumpCutMultiplier(0.5)
	j.jumpCutPending = true
	j.jumpPressed = true // was held

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	// Release jump — exercises the !jumpPressed && s.jumpPressed && s.jumpCutPending path
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Jump: false} }
	j.HandleInput(actor, model, sp)

	if j.jumpCutPending {
		t.Error("expected jumpCutPending=false after release")
	}
}

func TestJumpSkill_HandleInput_PressActivates(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		DownwardGravity:  4,
		MaxFallSpeed:     128,
	}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	model.SetOnGround(true)
	sp := space.NewSpace()

	j := NewJumpSkill()
	j.jumpPressed = false

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	// Press jump — exercises the jumpPressed && !s.jumpPressed path → tryActivate
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Jump: true} }
	j.HandleInput(actor, model, sp)

	_, vy := actor.Velocity()
	if vy >= 0 {
		t.Errorf("expected upward velocity after jump press; got %d", vy)
	}
}

// --- JumpSkill.handleCoyoteAndJumpBuffering: buffered jump on landing ---

func TestJumpSkill_HandleCoyoteAndJumpBuffering_BufferedJumpOnLanding(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		DownwardGravity:  4,
		MaxFallSpeed:     128,
	}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)

	j := NewJumpSkill()
	j.jumpBufferCounter = 3

	// Simulate: was NOT on ground, now IS on ground → buffered jump fires
	model.SetOnGround(false)
	// Call handleCoyoteAndJumpBuffering with wasOnGround=false, model.OnGround()=true
	// We do this via Update which calls handleCoyoteAndJumpBuffering(b, model, model.OnGround())
	// but Update reads model.OnGround() twice. Set it to true before Update.
	model.SetOnGround(true)

	// Manually call to exercise the buffered-jump branch
	j.handleCoyoteAndJumpBuffering(actor, model, false /* wasOnGround */)

	_, vy := actor.Velocity()
	if vy >= 0 {
		t.Errorf("expected upward velocity from buffered jump; got %d", vy)
	}
	if j.jumpBufferCounter != 0 {
		t.Errorf("expected jumpBufferCounter=0 after buffered jump; got %d", j.jumpBufferCounter)
	}
}

func TestJumpSkill_HandleCoyoteAndJumpBuffering_OnJumpCallback(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		DownwardGravity:  4,
		MaxFallSpeed:     128,
	}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	model := movement.NewPlatformMovementModel(nil)
	model.SetOnGround(true)

	called := false
	j := NewJumpSkill()
	j.jumpBufferCounter = 3
	j.OnJump = func(_ contractsbody.MovableCollidable) { called = true }

	j.handleCoyoteAndJumpBuffering(actor, model, false)

	if !called {
		t.Error("expected OnJump callback to be called on buffered jump")
	}
}

// --- JumpSkill.tryActivate: zero force guard ---

func TestJumpSkill_tryActivate_ZeroForce(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:       0, // zero force → early return
		DownwardGravity: 4,
		MaxFallSpeed:    128,
	}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	actor.SetJumpForceMultiplier(1.0)
	model := movement.NewPlatformMovementModel(nil)
	model.SetOnGround(true)

	j := NewJumpSkill()
	j.tryActivate(actor, model, space.NewSpace())

	_, vy := actor.Velocity()
	if vy != 0 {
		t.Errorf("expected no jump with zero force; got vy=%d", vy)
	}
}

// --- JumpSkill.Update: jumpCutPending cleared when not going up ---

func TestJumpSkill_Update_ClearsPendingWhenFalling(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		CoyoteTimeFrames: 5,
		DownwardGravity:  4,
		MaxFallSpeed:     128,
	}}
	config.Set(cfg)

	actor := newMockMovableCollidable()
	actor.SetVelocity(0, fp16.To16(5)) // falling (positive vy)
	model := movement.NewPlatformMovementModel(nil)

	j := NewJumpSkill()
	j.jumpCutPending = true

	j.Update(actor, model)

	if j.jumpCutPending {
		t.Error("expected jumpCutPending=false when not going up")
	}
}
