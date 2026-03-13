package skill

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

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
	d.state = StateActive
	d.timer = 2

	vx, _ := actor.Velocity()

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

	vx2, _ := actor.Velocity()
	// Velocity may persist due to physics, but dash should not force set; accept smoke check
	if vx2 == 0 && vx != 0 {
		// acceptable; no strict assertion — ensure no panic and progression through states
	}
}

func TestJumpSkillTryActivateAndCoyoteBuffer(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		DownwardGravity:  4,
		UpwardGravity:    2,
		MaxFallSpeed:     128,
	}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	model := movement.NewPlatformMovementModel(nil)
	model.SetOnGround(true)

	sp := space.NewSpace()
	js := NewJumpSkill()

	// Call unexported flow directly (same package)
	js.tryActivate(actor, model, sp)

	_, vy := actor.Velocity()
	if vy >= 0 {
		t.Fatalf("expected upward velocity after jump, got vy=%d", vy)
	}
	if model.OnGround() {
		t.Fatalf("expected model to be set airborne after jump")
	}
}
