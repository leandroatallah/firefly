package movement

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/input"
	spacephysics "github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type PlatformMovementModel struct {
	playerMovementBlocker PlayerMovementBlocker
	onGround              bool
	maxFallSpeed          int
	isScripted            bool
}

// NewPlatformMovementModel creates a new PlatformMovementModel with default values.
func NewPlatformMovementModel(playerMovementBlocker PlayerMovementBlocker) *PlatformMovementModel {
	m := &PlatformMovementModel{
		playerMovementBlocker: playerMovementBlocker,
		maxFallSpeed:          config.Get().Physics.MaxFallSpeed,
	}
	return m
}

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
		vx16 = clampAxisVelocity(vx16, fp16.To16(body.MaxSpeed()))

		// Apply friction if the player is not actively moving
		if accX == 0 {
			baseFriction := int(float64(fp16.To16(1)/4) * cfg.Physics.HorizontalInertia)
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
	if m.playerMovementBlocker != nil {
		m.InputHandler(body, space)
	}

	vx16, vy16 := body.Velocity()

	// Apply horizontal movement to the body and check for collisions.
	vx16, _ = m.UpdateHorizontalVelocity(body)
	_, _, isBlockingX := body.ApplyValidPosition(vx16, true, space)
	if isBlockingX {
		vx16 = 0
	}

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

	if clampToPlayArea(body, space.(*spacephysics.Space)) {
		vy16 = cfg.Physics.DownwardGravity - 1
		body.SetVelocity(vx16, vy16)
	}

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
func (m *PlatformMovementModel) InputHandler(body body.MovableCollidable, space body.BodiesSpace) {
	if m.isScripted {
		return // Ignore player input when scripted
	}
	if m.playerMovementBlocker != nil && m.playerMovementBlocker.IsMovementBlocked() {
		return // Ignore player input when movement is blocked
	}

	_, vy16 := body.Velocity()
	if body.Immobile() {
		_, accY := body.Acceleration()
		body.SetVelocity(0, vy16)
		body.SetAcceleration(0, accY)
		return
	}

	m.InputHandlerHorizontal(body)
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
			vx16 = -fp16.To16(body.Speed())
		case moveRight:
			vx16 = fp16.To16(body.Speed())
		default:
			vx16 = 0
		}
	}

	body.SetVelocity(vx16, vy16)
}

func (m *PlatformMovementModel) OnGround() bool {
	return m.onGround
}

func (m *PlatformMovementModel) SetOnGround(value bool) {
	m.onGround = value
}
