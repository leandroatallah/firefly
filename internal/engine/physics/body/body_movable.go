package body

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type MovableBody struct {
	*Body

	Ownership

	vx16                int
	vy16                int
	accelerationX       int
	accelerationY       int
	speed               int
	maxSpeed            int
	jumpForceMultiplier float64
	horizontalInertia   float64
	immobile            bool
	freeze              bool
	faceDirection       animation.FacingDirectionEnum
}

func NewMovableBody(body *Body) *MovableBody {
	if body == nil {
		panic("NewMovableBody: body must not be nil")
	}
	return &MovableBody{
		Body:                body,
		jumpForceMultiplier: 1.0,
		horizontalInertia:   -1.0,
	}
}

func (b *MovableBody) MoveX(distance int) {
	multiplier := config.Get().Physics.SpeedMultiplier
	if multiplier == 0 {
		multiplier = 1.0
	}
	dist16 := fp16.To16(distance)
	b.accelerationX = int(float64(dist16) * multiplier)
}

func (b *MovableBody) MoveY(distance int) {
	multiplier := config.Get().Physics.SpeedMultiplier
	if multiplier == 0 {
		multiplier = 1.0
	}
	dist16 := fp16.To16(distance)
	b.accelerationY = int(float64(dist16) * multiplier)
}

func (b *MovableBody) OnMoveLeft(distance int) {
	b.MoveX(-distance)
}
func (b *MovableBody) OnMoveUpLeft(distance int) {
	b.MoveX(-distance)
	b.MoveY(-distance)
}
func (b *MovableBody) OnMoveDownLeft(distance int) {
	b.MoveX(-distance)
	b.MoveY(distance)
}
func (b *MovableBody) OnMoveRight(distance int) {
	b.MoveX(distance)
}
func (b *MovableBody) OnMoveUpRight(distance int) {
	b.MoveX(distance)
	b.MoveY(-distance)
}
func (b *MovableBody) OnMoveDownRight(distance int) {
	b.MoveX(distance)
	b.MoveY(distance)
}
func (b *MovableBody) OnMoveUp(distance int) {
	b.MoveY(-distance)
}
func (b *MovableBody) OnMoveDown(distance int) {
	b.MoveY(distance)
}

func (b *MovableBody) SetSpeed(speed int) error {
	if b == nil {
		return fmt.Errorf("MovableBody is nil")
	}
	if speed < 0 {
		return fmt.Errorf("speed must be >= 0; got %d", speed)
	}
	b.speed = speed
	return nil
}

func (b *MovableBody) SetMaxSpeed(maxSpeed int) error {
	if b == nil {
		return fmt.Errorf("MovableBody is nil")
	}
	if maxSpeed < 0 {
		return fmt.Errorf("maxSpeed must be >= 0; got %d", maxSpeed)
	}
	b.maxSpeed = maxSpeed
	return nil
}

func (b *MovableBody) Speed() int {
	return b.speed
}

func (b *MovableBody) MaxSpeed() int {
	return b.maxSpeed
}

func (b *MovableBody) Velocity() (int, int) {
	return b.vx16, b.vy16
}

func (b *MovableBody) SetVelocity(vx16, vy16 int) {
	b.vx16, b.vy16 = vx16, vy16
}

func (b *MovableBody) Acceleration() (accX, accY int) {
	return b.accelerationX, b.accelerationY
}

func (b *MovableBody) SetAcceleration(accX, accY int) {
	b.accelerationX, b.accelerationY = accX, accY
}

func (b *MovableBody) Immobile() bool {
	return b.immobile
}

func (b *MovableBody) SetImmobile(immobile bool) {
	b.immobile = immobile
}
func (b *MovableBody) SetFreeze(freeze bool) {
	b.freeze = freeze
}

func (b *MovableBody) Freeze() bool {
	return b.freeze
}

func (b *MovableBody) FaceDirection() animation.FacingDirectionEnum {
	return b.faceDirection
}

func (b *MovableBody) SetFaceDirection(value animation.FacingDirectionEnum) {
	b.faceDirection = value
}

func (b *MovableBody) IsIdle() bool {
	return !b.IsWalking() && !b.IsFalling() && !b.IsGoingUp()
}

func (b *MovableBody) IsWalking() bool {
	// A body cannot be walking if it is airborne.
	if b.IsFalling() || b.IsGoingUp() {
		return false
	}

	threshold := config.Get().Physics.DownwardGravity
	if b.vx16 > threshold || b.vx16 < -threshold {
		return true
	}

	return false
}

func (b *MovableBody) IsGoingUp() bool {
	threshold := config.Get().Physics.DownwardGravity
	if b.vy16 <= -threshold {
		return true
	}

	return false
}

func (b *MovableBody) IsFalling() bool {
	threshold := config.Get().Physics.DownwardGravity
	if b.vy16 >= threshold {
		return true
	}

	return false
}

// Platform methods
func (b *MovableBody) TryJump(force int) {
	b.vy16 = -fp16.To16(force)
}

// CheckMovementDirectionX set face direction based on accelerationX
func (b *MovableBody) CheckMovementDirectionX() {
	if b.accelerationX > 0 {
		b.SetFaceDirection(animation.FaceDirectionRight)
	} else if b.accelerationX < 0 {
		b.SetFaceDirection(animation.FaceDirectionLeft)
	}
}

func (b *MovableBody) SetJumpForceMultiplier(multiplier float64) {
	b.jumpForceMultiplier = multiplier
}

func (b *MovableBody) JumpForceMultiplier() float64 {
	return b.jumpForceMultiplier
}

func (b *MovableBody) SetHorizontalInertia(inertia float64) {
	b.horizontalInertia = inertia
}

func (b *MovableBody) HorizontalInertia() float64 {
	return b.horizontalInertia
}
