package movement

import (
	"math"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

type BeatEmUpMovementModel struct {
	playerMovementBlocker PlayerMovementBlocker
	isScripted            bool
}

func NewBeatEmUpMovementModel(playerMovementBlocker PlayerMovementBlocker) *BeatEmUpMovementModel {
	return &BeatEmUpMovementModel{
		playerMovementBlocker: playerMovementBlocker,
	}
}

func (m *BeatEmUpMovementModel) SetIsScripted(isScripted bool) {
	m.isScripted = isScripted
}

func (m *BeatEmUpMovementModel) Update(b body.MovableCollidable, space body.BodiesSpace) error {
	if b.Freeze() {
		return nil
	}

	vx16, vy16 := b.Velocity()

	// Apply previous-frame velocity with collision resolution
	_, _, _ = b.ApplyValidPosition(vx16, true, space)  // X axis
	_, _, _ = b.ApplyValidPosition(vy16, false, space) // Y axis
	vx16, vy16 = b.Velocity()

	// Prevents leaving the play area
	clampToPlayArea(b, space)

	// Integrate acceleration (set externally by skill — passive model)
	accX, accY := b.Acceleration()
	scaledAccX, scaledAccY := smoothDiagonalMovement(accX, accY)
	vx16 = increaseVelocity(vx16, scaledAccX)
	vy16 = increaseVelocity(vy16, scaledAccY)

	// 2D speed cap
	speedMax16 := fp16.To16(b.MaxSpeed())
	if mult := config.Get().Physics.SpeedMultiplier; mult != 0 {
		speedMax16 = int(float64(speedMax16) * mult)
	}
	velSq := int64(vx16)*int64(vx16) + int64(vy16)*int64(vy16)
	maxSq := int64(speedMax16) * int64(speedMax16)
	if velSq > maxSq {
		scale := float64(speedMax16) / math.Sqrt(float64(velSq))
		vx16 = int(float64(vx16) * scale)
		vy16 = int(float64(vy16) * scale)
	}

	b.CheckMovementDirectionX()
	b.SetAcceleration(0, 0)

	// Friction both axes
	vx16 = reduceVelocity(vx16)
	vy16 = reduceVelocity(vy16)
	b.SetVelocity(vx16, vy16)
	return nil
}
