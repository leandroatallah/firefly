package body

import (
	"fmt"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type accVect struct {
	accelerationX int
	accelerationY int
}

func NewAccVect(x16, y16 int) accVect {
	return accVect{x16, y16}
}

func (v accVect) String() string {
	return fmt.Sprintf("{ x16: %d, y16: %d }", v.accelerationX, v.accelerationY)
}

func TestNewMovableBody_NilBody(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewMovableBody did not panic with nil body")
		}
	}()
	NewMovableBody(nil)
}

func TestMovableBody_Movement(t *testing.T) {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 1.0,
		},
	})

	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))

	distance := 5
	distancex16 := fp16.To16(distance)

	tests := []struct {
		name string
		fn   func(int)
		want accVect
	}{
		{"Move Left", b.OnMoveLeft, NewAccVect(-distancex16, 0)},
		{"Move Right", b.OnMoveRight, NewAccVect(distancex16, 0)},
		{"Move Up", b.OnMoveUp, NewAccVect(0, -distancex16)},
		{"Move Down", b.OnMoveDown, NewAccVect(0, distancex16)},
		{"Move Up Left", b.OnMoveUpLeft, NewAccVect(-distancex16, -distancex16)},
		{"Move Up Right", b.OnMoveUpRight, NewAccVect(distancex16, -distancex16)},
		{"Move Down Left", b.OnMoveDownLeft, NewAccVect(-distancex16, distancex16)},
		{"Move Down Right", b.OnMoveDownRight, NewAccVect(distancex16, distancex16)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset
			b.accelerationX, b.accelerationY = 0, 0

			tt.fn(distance)
			if b.accelerationX != tt.want.accelerationX || b.accelerationY != tt.want.accelerationY {
				t.Errorf(
					"expected %v; got %v",
					tt.want, fmt.Sprintf(" { x16: %d , y16: %d }", b.accelerationX, b.accelerationY))
			}
		})
	}
}

func TestMovableBody_SpeedMultiplier(t *testing.T) {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 2.0,
		},
	})

	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))
	b.MoveX(10)

	expected := int(float64(fp16.To16(10)) * 2.0)
	if b.accelerationX != expected {
		t.Errorf("expected accelerationX %d; got %d", expected, b.accelerationX)
	}
}

func TestMovableBody_Properties(t *testing.T) {
	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))

	// Speed
	if err := b.SetSpeed(10); err != nil {
		t.Errorf("unexpected error setting speed: %v", err)
	}
	if b.Speed() != 10 {
		t.Errorf("expected speed 10; got %d", b.Speed())
	}
	if err := b.SetSpeed(-1); err == nil {
		t.Error("expected error setting negative speed")
	}

	// MaxSpeed
	if err := b.SetMaxSpeed(20); err != nil {
		t.Errorf("unexpected error setting maxSpeed: %v", err)
	}
	if b.MaxSpeed() != 20 {
		t.Errorf("expected maxSpeed 20; got %d", b.MaxSpeed())
	}
	if err := b.SetMaxSpeed(-1); err == nil {
		t.Error("expected error setting negative maxSpeed")
	}

	// Velocity
	b.SetVelocity(100, 200)
	vx, vy := b.Velocity()
	if vx != 100 || vy != 200 {
		t.Errorf("expected velocity (100, 200); got (%d, %d)", vx, vy)
	}

	// Acceleration
	b.SetAcceleration(5, 10)
	ax, ay := b.Acceleration()
	if ax != 5 || ay != 10 {
		t.Errorf("expected acceleration (5, 10); got (%d, %d)", ax, ay)
	}

	// Immobile
	b.SetImmobile(true)
	if !b.Immobile() {
		t.Error("expected immobile to be true")
	}

	// Freeze
	b.SetFreeze(true)
	if !b.Freeze() {
		t.Error("expected freeze to be true")
	}

	// FaceDirection
	b.SetFaceDirection(animation.FaceDirectionLeft)
	if b.FaceDirection() != animation.FaceDirectionLeft {
		t.Errorf("expected FaceDirectionLeft; got %v", b.FaceDirection())
	}

	// JumpForceMultiplier
	b.SetJumpForceMultiplier(1.5)
	if b.JumpForceMultiplier() != 1.5 {
		t.Errorf("expected JumpForceMultiplier 1.5; got %f", b.JumpForceMultiplier())
	}

	// HorizontalInertia
	b.SetHorizontalInertia(0.8)
	if b.HorizontalInertia() != 0.8 {
		t.Errorf("expected HorizontalInertia 0.8; got %f", b.HorizontalInertia())
	}
}

func TestMovableBody_States(t *testing.T) {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 10,
		},
	})

	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))

	// Default state: Idle
	if !b.IsIdle() {
		t.Error("expected body to be idle by default")
	}

	// Walking
	b.SetVelocity(11, 0)
	if !b.IsWalking() {
		t.Error("expected walking to be true when vx > threshold")
	}
	if b.IsIdle() {
		t.Error("expected idle to be false when walking")
	}

	b.SetVelocity(-11, 0)
	if !b.IsWalking() {
		t.Error("expected walking to be true when vx < -threshold")
	}

	b.SetVelocity(5, 0)
	if b.IsWalking() {
		t.Error("expected walking to be false when vx within threshold")
	}

	// Falling
	b.SetVelocity(0, 10)
	if !b.IsFalling() {
		t.Error("expected falling to be true when vy >= threshold")
	}
	if b.IsWalking() {
		t.Error("expected walking to be false when falling")
	}

	// Going Up
	b.SetVelocity(0, -10)
	if !b.IsGoingUp() {
		t.Error("expected going up to be true when vy <= -threshold")
	}
	if b.IsWalking() {
		t.Error("expected walking to be false when going up")
	}
}

func TestMovableBody_TryJump(t *testing.T) {
	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))
	b.TryJump(10)

	expectedVY := -fp16.To16(10)
	_, vy := b.Velocity()
	if vy != expectedVY {
		t.Errorf("expected velocity Y %d; got %d", expectedVY, vy)
	}
}

func TestMovableBody_MoveX_DefaultMultiplier(t *testing.T) {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 0,
		},
	})

	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))
	b.MoveX(10)

	expected := fp16.To16(10) // multiplier 1.0 fallback
	if b.accelerationX != expected {
		t.Errorf("expected accelerationX %d; got %d", expected, b.accelerationX)
	}
}

func TestMovableBody_MoveY_DefaultMultiplier(t *testing.T) {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 0,
		},
	})

	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))
	b.MoveY(10)

	expected := fp16.To16(10) // multiplier 1.0 fallback
	if b.accelerationY != expected {
		t.Errorf("expected accelerationY %d; got %d", expected, b.accelerationY)
	}
}

func TestMovableBody_NilReceiver(t *testing.T) {
	var b *MovableBody
	if err := b.SetSpeed(10); err == nil {
		t.Error("expected error for nil receiver in SetSpeed")
	}
	if err := b.SetMaxSpeed(10); err == nil {
		t.Error("expected error for nil receiver in SetMaxSpeed")
	}
}

func TestMovableBody_CheckMovementDirectionX(t *testing.T) {
	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))

	b.SetAcceleration(10, 0)
	b.CheckMovementDirectionX()
	if b.FaceDirection() != animation.FaceDirectionRight {
		t.Error("expected FaceDirectionRight")
	}

	b.SetAcceleration(-10, 0)
	b.CheckMovementDirectionX()
	if b.FaceDirection() != animation.FaceDirectionLeft {
		t.Error("expected FaceDirectionLeft")
	}

	// Should not change if acceleration is 0
	b.SetFaceDirection(animation.FaceDirectionRight)
	b.SetAcceleration(0, 0)
	b.CheckMovementDirectionX()
	if b.FaceDirection() != animation.FaceDirectionRight {
		t.Error("expected FaceDirectionRight to be preserved")
	}
}
