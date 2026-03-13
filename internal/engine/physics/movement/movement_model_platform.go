package movement

import (
	"image"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type PlatformMovementModel struct {
	playerMovementBlocker PlayerMovementBlocker
	onGround              bool
	maxFallSpeed          int
	isScripted            bool
	dashActive            bool
	dashVelocityX         int
	gravityEnabled        bool
}

// NewPlatformMovementModel creates a new PlatformMovementModel with default values.
func NewPlatformMovementModel(playerMovementBlocker PlayerMovementBlocker) *PlatformMovementModel {
	m := &PlatformMovementModel{
		playerMovementBlocker: playerMovementBlocker,
		maxFallSpeed:          config.Get().Physics.MaxFallSpeed,
		gravityEnabled:        true,
	}
	return m
}

func (m *PlatformMovementModel) UpdateHorizontalVelocity(body body.MovableCollidable) (int, int) {
	cfg := config.Get()
	horizontalInertia := cfg.Physics.HorizontalInertia
	if val := body.HorizontalInertia(); val >= 0 {
		horizontalInertia = val
	}

	if horizontalInertia > 0 {
		// Acceleration-based movement
		accX, _ := body.Acceleration()
		scaledAccX, _ := smoothDiagonalMovement(accX, 0)

		// Apply air control multiplier if the player is in the air
		if !m.onGround && m.gravityEnabled {
			scaledAccX = int(float64(scaledAccX) * cfg.Physics.AirControlMultiplier)
		}

		vx16, vy16 := body.Velocity()

		vx16 = increaseVelocity(vx16, scaledAccX)
		limit := fp16.To16(body.MaxSpeed())
		if cfg.Physics.SpeedMultiplier != 0 {
			limit = int(float64(limit) * cfg.Physics.SpeedMultiplier)
		}
		vx16 = clampAxisVelocity(vx16, limit)

		// Apply friction if the player is not actively moving
		if accX == 0 {
			// Increase friction coefficient to reduce sliding distance
			baseFriction := int(float64(fp16.To16(1)/2) * horizontalInertia)
			friction := baseFriction

			// Apply air friction multiplier if the player is in the air
			if !m.onGround && m.gravityEnabled {
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
	} else {
		// If inertia is disabled, check if we have acceleration pending from a state/command
		accX, _ := body.Acceleration()
		if accX != 0 {
			vx16, vy16 := body.Velocity()
			vx16 = accX
			body.SetVelocity(vx16, vy16)
		}
	}

	return body.Velocity()
}

func (m *PlatformMovementModel) UpdateVerticalVelocity(body body.MovableCollidable) (int, int) {
	if m.gravityEnabled {
		return body.Velocity()
	}

	_, accY := body.Acceleration()
	if accY != 0 {
		vx16, vy16 := body.Velocity()
		vy16 = accY
		body.SetVelocity(vx16, vy16)
	}
	return body.Velocity()
}

func (m *PlatformMovementModel) handleGravity(b body.MovableCollidable) (int, int) {
	vx16, vy16 := b.Velocity()

	if !m.gravityEnabled {
		return vx16, vy16
	}

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
	if body.Freeze() {
		return nil
	}

	cfg := config.Get()

	vx16, vy16 := body.Velocity()

	// Handle horizontal movement based on dash state or normal acceleration/friction
	if m.dashActive {
		vx16 = m.dashVelocityX
	} else {
		vx16, _ = m.UpdateHorizontalVelocity(body)
	}

	// Handle vertical movement if gravity is disabled
	if !m.gravityEnabled {
		_, vy16 = m.UpdateVerticalVelocity(body)
	}

	// Apply horizontal movement to the body and check for collisions.
	_, _, isBlockingX := body.ApplyValidPosition(vx16, true, space)
	if isBlockingX {
		vx16 = 0
	}

	// Apply vertical movement to the body and check for collisions.
	_, _, isBlockingY := body.ApplyValidPosition(vy16, false, space)
	vx16, vy16 = body.Velocity()

	isGrounded := false
	if isBlockingY {
		if vy16 > 0 {
			isGrounded = true
		}
	} else {
		if vy16 >= 0 && m.CheckGround(body, space) {
			isGrounded = true
		}
	}

	if isGrounded {
		m.onGround = true
		// Set a small downward velocity to "stick" to the ground, ensuring it's less than the falling threshold.
		if m.gravityEnabled && vy16 >= 0 {
			vy16 = cfg.Physics.DownwardGravity - 1
			body.SetVelocity(vx16, vy16)
		}
	} else {
		m.onGround = false
	}

	if clampToPlayArea(body, space) {
		if m.gravityEnabled {
			vy16 = cfg.Physics.DownwardGravity - 1
			body.SetVelocity(vx16, vy16)
		}
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

// IsInputBlocked reports whether player input should be ignored.
func (m *PlatformMovementModel) IsInputBlocked() bool {
	return m.playerMovementBlocker != nil && m.playerMovementBlocker.IsMovementBlocked()
}

func (m *PlatformMovementModel) OnGround() bool {
	return m.onGround
}

func (m *PlatformMovementModel) SetOnGround(value bool) {
	m.onGround = value
}

// SetDashActive sets the dash state and velocity for the movement model.
func (m *PlatformMovementModel) SetDashActive(active bool, vx int) {
	m.dashActive = active
	m.dashVelocityX = vx
}

func (m *PlatformMovementModel) SetGravityEnabled(enabled bool) {
	m.gravityEnabled = enabled
}

func (m *PlatformMovementModel) CheckGround(b body.MovableCollidable, space body.BodiesSpace) bool {
	collisionRects := b.CollisionPosition()
	// If no specific collision shapes are defined, fall back to the main body shape.
	if len(collisionRects) == 0 {
		collisionRects = []image.Rectangle{b.Position()}
	}

	for _, pos := range collisionRects {
		// Create a rect shifted 1 pixel down
		checkRect := pos.Add(image.Point{0, 1})

		// Shrink horizontally to avoid wall friction
		width := checkRect.Dx()
		if width > 4 {
			checkRect.Min.X += 2
			checkRect.Max.X -= 2
		} else if width > 2 {
			checkRect.Min.X += 1
			checkRect.Max.X -= 1
		}

		collidables := space.Query(checkRect)
		for _, c := range collidables {
			if c.ID() == b.ID() {
				continue
			}
			if c.IsObstructive() {
				return true
			}
		}
	}
	return false
}
