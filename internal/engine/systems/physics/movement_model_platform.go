package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
)

const (
	// horizontalInertia controls the smoothness of horizontal movement.
	// 0.0 means instant movement (no inertia).
	// 1.0 is a good default, providing a smoother, acceleration-based movement.
	// Higher values will make the character slide more.
	horizontalInertia = 2.0

	// airFrictionMultiplier controls how much friction is applied in the air.
	// 1.0 means air friction is the same as ground friction.
	// 0.0 means no air friction (full momentum).
	// 0.5 means air friction is half of the ground friction.
	airFrictionMultiplier = 0.5

	// airControlMultiplier controls how much acceleration is applied in the air.
	// 1.0 means air control is the same as on the ground.
	// Values > 1.0 provide stronger air control, making it easier to change direction.
	// Values < 1.0 provide weaker air control, making it harder to change direction.
	airControlMultiplier = 0.25

	// coyoteTimeFrames is the number of frames the player can still jump after leaving a ledge.
	coyoteTimeFrames = 6

	// jumpBufferFrames is the number of frames a jump input is remembered before landing.
	jumpBufferFrames = 6

	jumpHeight        = 6
	jumpCutMultiplier = 0.5
	upwardGravity     = 6 // Gravity when going up
	downwardGravity   = 8 // Gravity when falling
)

type PlatformMovementModel struct {
	onGround     bool
	maxFallSpeed int
	isScripted   bool

	// coyoteTimeCounter allows the player to jump for a few frames after leaving a ledge.
	coyoteTimeCounter int
	// jumpBufferCounter remembers a jump input for a few frames, executing it upon landing.
	jumpBufferCounter int

	skills []Skill
}

// NewPlatformMovementModel creates a new PlatformMovementModel with default values.
func NewPlatformMovementModel() *PlatformMovementModel {
	m := &PlatformMovementModel{
		maxFallSpeed: 12 * config.Get().Unit,
	}
	m.skills = append(m.skills, NewDashSkill())
	return m
}

// Update handles the physics for a platformer-style character.
// It processes input, applies movement and collisions, handles gravity.
func (m *PlatformMovementModel) Update(body *PhysicsBody, space body.BodiesSpace) error {
	// Handle input for player movement. This needs to be done before physics calculations.
	m.InputHandler(body)

	// Update all skills
	var skillIsActive bool
	for _, skill := range m.skills {
		skill.Update(body, m)
		if skill.IsActive() {
			skillIsActive = true
		}
	}

	// --- Horizontal Movement ---
	if !skillIsActive && horizontalInertia > 0 {
		// Acceleration-based movement
		scaledAccX, _ := smoothDiagonalMovement(body.accelerationX, 0)

		// Apply air control multiplier if the player is in the air
		if !m.onGround {
			scaledAccX = int(float64(scaledAccX) * airControlMultiplier)
		}

		body.vx16 = increaseVelocity(body.vx16, scaledAccX)
		body.vx16 = clampAxisVelocity(body.vx16, body.maxSpeed*config.Get().Unit)

		// Apply friction if the player is not actively moving
		if body.accelerationX == 0 {
			baseFriction := int(float64(config.Get().Unit/4) * horizontalInertia)
			friction := baseFriction

			// Apply air friction multiplier if the player is in the air
			if !m.onGround {
				friction = int(float64(baseFriction) * airFrictionMultiplier)
			}

			if body.vx16 > friction {
				body.vx16 -= friction
			} else if body.vx16 < -friction {
				body.vx16 += friction
			} else {
				body.vx16 = 0
			}
		}
	}

	// Apply horizontal movement to the body and check for collisions.
	if applyAxisMovement(body, body.vx16, true, space) {
		body.vx16 = 0
	}

	// --- Vertical Movement & State ---
	// Store previous onGround status to detect landing
	wasOnGround := m.onGround

	// Apply vertical movement and check for collisions to determine if the body has landed.
	verticalBlocked := applyAxisMovement(body, body.vy16, false, space)
	landed := false
	if verticalBlocked {
		if body.vy16 > 0 {
			landed = true
		}
		body.vy16 = 0
	}

	// Also consider landing if the body hits the bottom of the play area.
	if clampToPlayArea(body, space.(*Space)) {
		landed = true
		body.vy16 = 0
	}

	// Update the final 'onGround' state.
	m.onGround = landed

	// --- Coyote Time & Jump Buffering ---
	if m.onGround {
		// If on the ground, reset the coyote time counter.
		m.coyoteTimeCounter = coyoteTimeFrames
	} else {
		// If in the air, decrement the coyote time counter.
		if m.coyoteTimeCounter > 0 {
			m.coyoteTimeCounter--
		}
	}

	// Decrement the jump buffer counter each frame.
	if m.jumpBufferCounter > 0 {
		m.jumpBufferCounter--
	}

	// Check for and execute a buffered jump if we just landed.
	if !wasOnGround && m.onGround && m.jumpBufferCounter > 0 {
		body.TryJump(jumpHeight)
		m.onGround = false      // We are jumping, so we are no longer on the ground.
		m.jumpBufferCounter = 0 // Consume the buffer.
		m.coyoteTimeCounter = 0 // Ensure coyote time isn't used as well.
	}

	// --- Final State Updates ---
	body.CheckMovementDirectionX()
	body.accelerationX, body.accelerationY = 0, 0

	if m.onGround {
		// By setting vy16 to a small positive value, we ensure that ground collision
		// is checked on the next frame. This allows the system to detect when the
		// player walks off a platform.
		body.vy16 = 1
	} else if !skillIsActive {
		// Apply custom gravity based on vertical velocity
		if body.vy16 < 0 {
			// Player is moving up (jumping)
			body.vy16 += upwardGravity
		} else {
			// Player is moving down (falling)
			body.vy16 += downwardGravity
		}

		// Clamp fall speed to terminal velocity
		if body.vy16 > m.maxFallSpeed {
			body.vy16 = m.maxFallSpeed
		}
	}

	return nil
}

// SetIsScripted sets the scripted mode for the movement model.
func (m *PlatformMovementModel) SetIsScripted(isScripted bool) {
	m.isScripted = isScripted
}

// InputHandler processes player input for movement.
func (m *PlatformMovementModel) InputHandler(body *PhysicsBody) {
	if m.isScripted {
		return // Ignore player input when scripted
	}
	// Let skills handle their input first.
	for _, skill := range m.skills {
		if activeSkill, ok := skill.(ActiveSkill); ok {
			activeSkill.HandleInput(body, m)
		}
	}

	// Check if any skill is active, which might block normal input.
	var skillIsActive bool
	for _, skill := range m.skills {
		if skill.IsActive() {
			skillIsActive = true
			break
		}
	}

	// If a skill is active, skip normal movement input.
	if skillIsActive {
		body.accelerationX = 0 // Prevent acceleration from previous frame carrying over
		return
	}

	if body.Immobile() {
		body.vx16 = 0
		body.accelerationX = 0
		return
	}

	// --- Horizontal Input ---
	moveLeft := input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft)
	moveRight := input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight)

	if horizontalInertia > 0 {
		// Acceleration-based movement
		if moveLeft {
			body.OnMoveLeft(body.Speed())
		}
		if moveRight {
			body.OnMoveRight(body.Speed())
		}
	} else {
		// Instant movement
		switch {
		case moveLeft:
			body.vx16 = -body.Speed() * config.Get().Unit
		case moveRight:
			body.vx16 = body.Speed() * config.Get().Unit
		default:
			body.vx16 = 0
		}
	}

	// --- Jump Input ---
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// A jump is triggered if the player is on the ground OR if coyote time is active.
		if m.onGround || m.coyoteTimeCounter > 0 {
			body.TryJump(jumpHeight)
			m.onGround = false      // Immediately leave the ground state.
			m.coyoteTimeCounter = 0 // Consume coyote time to prevent double jumps.
			m.jumpBufferCounter = 0 // Clear any buffered jump.
		} else {
			// If in the air and unable to jump, buffer the jump input.
			m.jumpBufferCounter = jumpBufferFrames
		}
	}

	// --- Variable Jump Height ---
	// If the player releases the jump button mid-air, reduce the upward velocity.
	if !m.onGround && !input.IsSomeKeyPressed(ebiten.KeySpace) && body.vy16 < 0 {
		body.vy16 = int(float64(body.vy16) * jumpCutMultiplier)
	}
}
