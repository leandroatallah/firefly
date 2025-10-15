package physics

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
)

type TopDownMovementModel struct{
	isScripted bool
}

func NewTopDownMovementModel() *TopDownMovementModel {
	return &TopDownMovementModel{}
}

func (m *TopDownMovementModel) Update(body *PhysicsBody, space body.BodiesSpace) error {
	// Handle input for player movement
	m.InputHandler(body)

	// Apply physics to player's position based on the velocity from previous frame.
	// This is a simple Euler integration step: position += velocity * deltaTime (where deltaTime=1 frame).
	body.ApplyValidMovement(body.vx16, true, space)
	body.ApplyValidMovement(body.vy16, false, space)

	// Prevents leaving the play area`
	clampToPlayArea(body, space.(*Space))

	// Convert the raw input acceleration into a scaled and normalized vector.
	scaledAccX, scaledAccY := smoothDiagonalMovement(body.accelerationX, body.accelerationY)

	body.vx16 = increaseVelocity(body.vx16, scaledAccX)
	body.vy16 = increaseVelocity(body.vy16, scaledAccY)

	// Cap the magnitude of the velocity vector to enforce a maximum speed.
	// This is crucial for preventing faster movement on diagonals.
	// We need to check if the velocity magnitude `sqrt(vx² + vy²)` exceeds `speedMax16²`.
	// To avoid a costly square root, we can compare the squared values:
	speedMax16 := body.maxSpeed * config.Get().Unit
	// Use int64 for squared values to prevent potential overflow.
	velSq := int64(body.vx16)*int64(body.vx16) + int64(body.vy16)*int64(body.vy16)
	maxSq := int64(speedMax16) * int64(speedMax16)

	if velSq > maxSq {
		// If the speed is too high, we need to scale the velocity vector down.
		// The scaling factor is `scale = speedMax16 / current_speed`.
		// `current_speed` is `sqrt(velSq)`.
		// So, `scale = speedMax16 / sqrt(velSq)`.
		scale := float64(speedMax16) / math.Sqrt(float64(velSq))
		body.vx16 = int(float64(body.vx16) * scale)
		body.vy16 = int(float64(body.vy16) * scale)
	}

	body.CheckMovementDirectionX()

	// Reset frame-specific acceleration.
	// It will be recalculated on the next frame from input.
	body.accelerationX, body.accelerationY = 0, 0

	// Apply friction to slow the player down when there is no input.
	body.vx16 = reduceVelocity(body.vx16)
	body.vy16 = reduceVelocity(body.vy16)

	return nil
}

func (m *TopDownMovementModel) SetIsScripted(isScripted bool) {
	m.isScripted = isScripted
}

// InputHandler processes player input for movement.
// Top-Down player can move for all directions and diagonals.
func (m *TopDownMovementModel) InputHandler(body *PhysicsBody) {
	if m.isScripted {
		return // Ignore player input when scripted
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
