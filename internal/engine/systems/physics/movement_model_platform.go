package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
)

type PlatformMovementModel struct {
	playerMovementBlocker PlayerMovementBlocker
	onGround              bool
	maxFallSpeed          int
	isScripted            bool
	coyoteTimeCounter     int
	jumpBufferCounter     int
	// skills                []Skill
}

// NewPlatformMovementModel creates a new PlatformMovementModel with default values.
func NewPlatformMovementModel(playerMovementBlocker PlayerMovementBlocker) *PlatformMovementModel {
	m := &PlatformMovementModel{
		playerMovementBlocker: playerMovementBlocker,
		maxFallSpeed:          config.Get().Physics.MaxFallSpeed,
	}
	return m
}

// TODO: Maybe I should update the method named, because this handle the acceleration only
func (m *PlatformMovementModel) UpdateHorizontalVelocity(body body.MovableCollidable) (int, int) {
	cfg := config.Get()

	if cfg.Physics.HorizontalInertia > 0 {
		// Acceleration-based movement
		accX, _ := body.Acceleration()
		scaledAccX, _ := smoothDiagonalMovement(accX, 0)

		// Apply air control multiplier if the player is in the air
		if !m.onGround {
			scaledAccX = int(float64(scaledAccX) * cfg.Physics.AirControlMultiplier)
		}

		vx16, vy16 := body.Velocity()

		vx16 = increaseVelocity(vx16, scaledAccX)
		vx16 = clampAxisVelocity(vx16, body.MaxSpeed()*cfg.Unit)

		// Apply friction if the player is not actively moving
		if accX == 0 {
			baseFriction := int(float64(cfg.Unit/4) * cfg.Physics.HorizontalInertia)
			friction := baseFriction

			// Apply air friction multiplier if the player is in the air
			if !m.onGround {
				friction = int(float64(baseFriction) * cfg.Physics.AirFrictionMultiplier)
			}

			if vx16 > friction {
				vx16 -= friction
			} else if vx16 < -friction {
				vx16 += friction
			} else {
				vx16 = 0
			}
		}

		body.SetVelocity(vx16, vy16)
	}

	return body.Velocity()
}

// Coyote Time & Jump Buffering
func (m *PlatformMovementModel) handleCoyoteAndJumpBuffering(body body.MovableCollidable, wasOnGround bool) {
	cfg := config.Get()

	if m.onGround {
		m.coyoteTimeCounter = cfg.Physics.CoyoteTimeFrames
	} else {
		if m.coyoteTimeCounter > 0 {
			m.coyoteTimeCounter--
		}
	}

	if m.jumpBufferCounter > 0 {
		m.jumpBufferCounter--
	}

	if !wasOnGround && m.onGround && m.jumpBufferCounter > 0 {
		body.TryJump(cfg.Physics.JumpForce)
		m.onGround = false
		m.jumpBufferCounter = 0
		m.coyoteTimeCounter = 0
	}
}

func (m *PlatformMovementModel) handleGravity(b body.MovableCollidable) (int, int) {
	vx16, vy16 := b.Velocity()

	if m.onGround {
		return vx16, vy16
	}

	cfg := config.Get()

	// Apply gravity when in the air
	if vy16 < 0 {
		vy16 += cfg.Physics.UpwardGravity
	} else {
		vy16 += cfg.Physics.DownwardGravity
	}

	// Clamp fall speed
	if vy16 > m.maxFallSpeed {
		vy16 = m.maxFallSpeed
	}

	return vx16, vy16
}

// Update handles the physics for a platformer-style character.
func (m *PlatformMovementModel) Update(body body.MovableCollidable, space body.BodiesSpace) error {
	cfg := config.Get()

	// Handle input for player movement. This needs to be done before physics calculations.
	m.InputHandler(body, space)

	vx16, vy16 := body.Velocity()

	// Apply horizontal movement to the body and check for collisions.
	vx16, _ = m.UpdateHorizontalVelocity(body)
	_, _, isBlockingX := body.ApplyValidPosition(vx16, true, space)
	if isBlockingX {
		vx16 = 0
	}

	// TODO: Check if it should get rid of wasOnGround
	wasOnGround := m.onGround
	// Apply vertical movement to the body and check for collisions.
	_, _, isBlockingY := body.ApplyValidPosition(vy16, false, space)
	vx16, vy16 = body.Velocity()
	if isBlockingY {
		if !m.onGround && vy16 > 0 { // Moving down and collided (i.e., landed)
			m.onGround = true
			// Set a small downward velocity to "stick" to the ground, ensuring it's less than the falling threshold.
			vy16 = cfg.Physics.DownwardGravity - 1
			body.SetVelocity(vx16, vy16)
		}
	} else {
		m.onGround = false
	}

	if clampToPlayArea(body, space.(*Space)) {
		vy16 = cfg.Physics.DownwardGravity - 1
		body.SetVelocity(vx16, vy16)
	}

	m.handleCoyoteAndJumpBuffering(body, wasOnGround)

	// --- Final State Updates ---
	body.CheckMovementDirectionX()
	body.SetAcceleration(0, 0)

	// Only apply gravity when airborne. The sticking force is handled above.
	_, vy16 = m.handleGravity(body)

	body.SetVelocity(vx16, vy16)

	return nil
}

// SetIsScripted sets the scripted mode for the movement model.
func (m *PlatformMovementModel) SetIsScripted(isScripted bool) {
	m.isScripted = isScripted
}

// InputHandler processes player input for movement.
// TODO: Move movement behavior to physics.MovableBody
func (m *PlatformMovementModel) InputHandler(body body.MovableCollidable, space body.BodiesSpace) {
	if m.isScripted || m.playerMovementBlocker.IsMovementBlocked() {
		return // Ignore player input when scripted or movement is blocked
	}

	// TODO: Should this be check here?
	_, vy16 := body.Velocity()
	if body.Immobile() {
		_, accY := body.Acceleration()
		body.SetVelocity(0, vy16)
		body.SetAcceleration(0, accY)
		return
	}

	m.InputHandlerHorizontal(body)
	m.InputHandlerJump(body, space)
}

func (m *PlatformMovementModel) InputHandlerHorizontal(body body.MovableCollidable) {
	cfg := config.Get()
	vx16, vy16 := body.Velocity()

	moveLeft := input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft)
	moveRight := input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight)

	if cfg.Physics.HorizontalInertia > 0 {
		if moveLeft {
			body.OnMoveLeft(body.Speed())
		}
		if moveRight {
			body.OnMoveRight(body.Speed())
		}
	} else {
		switch {
		case moveLeft:
			vx16 = -body.Speed() * cfg.Unit
		case moveRight:
			vx16 = body.Speed() * cfg.Unit
		default:
			vx16 = 0
		}
	}

	body.SetVelocity(vx16, vy16)
}

func (m *PlatformMovementModel) InputHandlerJump(body body.MovableCollidable, space body.BodiesSpace) (int, int) {
	cfg := config.Get()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if m.onGround || m.coyoteTimeCounter > 0 {
			// vx16, vy16 := body.Velocity()
			body.TryJump(cfg.Physics.JumpForce)
			// blocking := false
			// n_vx16, n_vy16 := body.Velocity()

			for _, other := range space.Bodies() {
				if other == nil || other.ID() == body.ID() {
					continue
				}

				// TODO: Evaluate this condition
				if other.ID() == "OBSTACLE_GROUND" {
					continue
				}

				if !hasCollision(body, other) {
					continue
				}

				if other.IsObstructive() {
					// blocking = true
					break
				}
			}

			m.onGround = false
			m.coyoteTimeCounter = 0
			m.jumpBufferCounter = 0

			// check collision
			// TODO: Should return?
			return body.Velocity()
		} else {
			m.jumpBufferCounter = cfg.Physics.JumpBufferFrames
		}
	}

	vx16, vy16 := body.Velocity()

	// Variable Jump Height when release jump key
	if !m.onGround && !input.IsSomeKeyPressed(ebiten.KeySpace) && vy16 < 0 {
		vy16 = int(float64(vy16) * cfg.Physics.JumpCutMultiplier)
	}

	body.SetVelocity(vx16, vy16)
	// TODO: Should return?
	return vx16, vy16
}
