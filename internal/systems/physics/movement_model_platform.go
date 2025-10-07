package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/input"
)

type PlatformMovementModel struct {
	onGround     bool
	maxFallSpeed int
}

// NewPlatformMovementModel creates a new PlatformMovementModel with default values.
func NewPlatformMovementModel() *PlatformMovementModel {
	return &PlatformMovementModel{
		maxFallSpeed: 12 * config.Unit,
	}
}

// Update handles the physics for a platformer-style character.
// It processes input, applies movement and collisions, handles gravity.
func (m *PlatformMovementModel) Update(body *PhysicsBody, space *Space) error {
	// Handle input for player movement
	m.InputHandler(body)

	// Apply horizontal movement and check for collisions.
	if applyAxisMovement(body, body.vx16, true, space) {
		body.vx16 = 0
	}

	// Apply vertical movement and check for collisions, detecting if the body has landed.
	verticalBlocked := applyAxisMovement(body, body.vy16, false, space)
	landed := false
	if verticalBlocked {
		if body.vy16 > 0 {
			landed = true
		}
		body.vy16 = 0
	}

	// Clamp the body's position to the play area boundaries.
	if clampToPlayArea(body, space) {
		landed = true
		body.vy16 = 0
	}

	// Update the 'onGround' state based on landing or vertical movement.
	if landed {
		m.onGround = true
	} else if body.vy16 != 0 {
		m.onGround = false
	}

	// Apply acceleration and clamp horizontal velocity.
	scaledAccX, _ := smoothDiagonalMovement(body.accelerationX, 0)
	body.vx16 = increaseVelocity(body.vx16, scaledAccX)
	body.vx16 = clampAxisVelocity(body.vx16, body.maxSpeed*config.Unit)

	body.CheckMovementDirectionX()

	// Apply friction to slow down horizontal movement.
	body.accelerationX, body.accelerationY = 0, 0
	body.vx16 = reduceVelocity(body.vx16)

	if m.onGround {
		// By setting vy16 to a small positive value, we ensure that the collision
		// detection for the ground is triggered on the next frame. This allows
		// the system to detect when the player walks off a platform.
		body.vy16 = 1
	} else {
		// Apply gravity if the body is in the air.
		applyGravity(body, m.maxFallSpeed)
	}

	return nil
}

// InputHandler processes player input for movement.
// Platform player can move horinzontally acceleration based
// and triggers a jump if the jump key is pressed while on the ground.
func (m *PlatformMovementModel) InputHandler(body *PhysicsBody) {
	if body.Immobile() {
		return
	}

	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		body.OnMoveLeft(body.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		body.OnMoveRight(body.Speed())
	}
	if m.onGround && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		body.TryJump(8) // Replace magic number
		m.onGround = false
	}
}
