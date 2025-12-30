package movement

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/input"
	spacephysics "github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type TopDownMovementModel struct {
	playerMovementBlocker PlayerMovementBlocker
	isScripted            bool
}

func NewTopDownMovementModel(playerMovementBlocker PlayerMovementBlocker) *TopDownMovementModel {
	return &TopDownMovementModel{
		playerMovementBlocker: playerMovementBlocker,
	}
}

func (m *TopDownMovementModel) Update(body body.MovableCollidable, space body.BodiesSpace) error {
	// Handle input for player movement
	if m.playerMovementBlocker != nil {
		m.InputHandler(body, space)
	}

	vx16, vy16 := body.Velocity()

	// Apply physics to player's position based on the velocity from previous frame.
	// This is a simple Euler integration step: position += velocity * deltaTime (where deltaTime=1 frame).
	_, _, _ = body.ApplyValidPosition(vx16, true, space)
	_, _, _ = body.ApplyValidPosition(vy16, false, space)

	vx16, vy16 = body.Velocity()

	// Prevents leaving the play area`
	clampToPlayArea(body, space.(*spacephysics.Space))

	// Convert the raw input acceleration into a scaled and normalized vector.
	accX, accY := body.Acceleration()
	scaledAccX, scaledAccY := smoothDiagonalMovement(accX, accY)

	vx16 = increaseVelocity(vx16, scaledAccX)
	vy16 = increaseVelocity(vy16, scaledAccY)

	// Cap the magnitude of the velocity vector to enforce a maximum speed.
	// This is crucial for preventing faster movement on diagonals.
	// We need to check if the velocity magnitude `sqrt(vx² + vy²)` exceeds `speedMax16²`.
	// To avoid a costly square root, we can compare the squared values:
	speedMax16 := fp16.To16(body.MaxSpeed())
	// Use int64 for squared values to prevent potential overflow.
	velSq := int64(vx16) + int64(vy16)*int64(vy16)
	maxSq := int64(speedMax16) * int64(speedMax16)

	if velSq > maxSq {
		// If the speed is too high, we need to scale the velocity vector down.
		// The scaling factor is `scale = speedMax16 / current_speed`.
		// `current_speed` is `sqrt(velSq)`.
		// So, `scale = speedMax16 / sqrt(velSq)`.
		scale := float64(speedMax16) / math.Sqrt(float64(velSq))
		vx16 = int(float64(vx16) * scale)
		vy16 = int(float64(vy16) * scale)
	}

	body.CheckMovementDirectionX()

	// Reset frame-specific acceleration.
	// It will be recalculated on the next frame from input.
	body.SetAcceleration(0, 0)

	// Apply friction to slow the player down when there is no input.
	vx16 = reduceVelocity(vx16)
	vy16 = reduceVelocity(vy16)
	body.SetVelocity(vx16, vy16)

	return nil
}

func (m *TopDownMovementModel) SetIsScripted(isScripted bool) {
	m.isScripted = isScripted
}

// InputHandler processes player input for movement.
// Top-Down player can move for all directions and diagonals.
func (m *TopDownMovementModel) InputHandler(body body.MovableCollidable, space body.BodiesSpace) {
	if m.isScripted {
		return // Ignore player input when scripted
	}
	if m.playerMovementBlocker != nil && m.playerMovementBlocker.IsMovementBlocked() {
		return // Ignore player input when movement is blocked
	}
	if body.Immobile() {
		return
	}

	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		body.OnMoveLeft(body.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		body.OnMoveRight(body.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		body.OnMoveUp(body.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		body.OnMoveDown(body.Speed())
	}
}
