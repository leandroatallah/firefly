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
	skills                []Skill
}

// NewPlatformMovementModel creates a new PlatformMovementModel with default values.
func NewPlatformMovementModel(playerMovementBlocker PlayerMovementBlocker) *PlatformMovementModel {
	m := &PlatformMovementModel{
		playerMovementBlocker: playerMovementBlocker,
		maxFallSpeed:          config.Get().Physics.MaxFallSpeed,
	}
	m.skills = append(m.skills, NewDashSkill())
	return m
}

// Update handles the physics for a platformer-style character.
func (m *PlatformMovementModel) Update(body *PhysicsBody, space body.BodiesSpace) error {
	cfg := config.Get().Physics

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
	if !skillIsActive && cfg.HorizontalInertia > 0 {
		// Acceleration-based movement
		scaledAccX, _ := smoothDiagonalMovement(body.accelerationX, 0)

		// Apply air control multiplier if the player is in the air
		if !m.onGround {
			scaledAccX = int(float64(scaledAccX) * cfg.AirControlMultiplier)
		}

		body.vx16 = increaseVelocity(body.vx16, scaledAccX)
		body.vx16 = clampAxisVelocity(body.vx16, body.maxSpeed*config.Get().Unit)

		// Apply friction if the player is not actively moving
		if body.accelerationX == 0 {
			baseFriction := int(float64(config.Get().Unit/4) * cfg.HorizontalInertia)
			friction := baseFriction

			// Apply air friction multiplier if the player is in the air
			if !m.onGround {
				friction = int(float64(baseFriction) * cfg.AirFrictionMultiplier)
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
	wasOnGround := m.onGround

	verticalBlocked := applyAxisMovement(body, body.vy16, false, space)
	landed := false
	if verticalBlocked {
		if body.vy16 > 0 {
			landed = true
		}
		body.vy16 = 0
	}

	if clampToPlayArea(body, space.(*Space)) {
		landed = true
		body.vy16 = 0
	}

	m.onGround = landed

	// --- Coyote Time & Jump Buffering ---
	if m.onGround {
		m.coyoteTimeCounter = cfg.CoyoteTimeFrames
	} else {
		if m.coyoteTimeCounter > 0 {
			m.coyoteTimeCounter--
		}
	}

	if m.jumpBufferCounter > 0 {
		m.jumpBufferCounter--
	}

	if !wasOnGround && m.onGround && m.jumpBufferCounter > 0 {
		body.TryJump(cfg.JumpForce)
		m.onGround = false
		m.jumpBufferCounter = 0
		m.coyoteTimeCounter = 0
	}

	// --- Final State Updates ---
	body.CheckMovementDirectionX()
	body.accelerationX, body.accelerationY = 0, 0

	if m.onGround {
		body.vy16 = 1
	} else if !skillIsActive {
		if body.vy16 < 0 {
			body.vy16 += cfg.UpwardGravity
		} else {
			body.vy16 += cfg.DownwardGravity
		}

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
// TODO: Move movement behavior to physics.MovableBody
func (m *PlatformMovementModel) InputHandler(body *PhysicsBody) {
	cfg := config.Get().Physics

	if m.isScripted || m.playerMovementBlocker.IsMovementBlocked() {
		return // Ignore player input when scripted or movement is blocked
	}

	// Let skills handle their input first.
	for _, skill := range m.skills {
		if activeSkill, ok := skill.(ActiveSkill); ok {
			activeSkill.HandleInput(body, m)
		}
	}

	var skillIsActive bool
	for _, skill := range m.skills {
		if skill.IsActive() {
			skillIsActive = true
			break
		}
	}

	if skillIsActive {
		body.accelerationX = 0
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

	if cfg.HorizontalInertia > 0 {
		if moveLeft {
			body.OnMoveLeft(body.Speed())
		}
		if moveRight {
			body.OnMoveRight(body.Speed())
		}
	} else {
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
		if m.onGround || m.coyoteTimeCounter > 0 {
			body.TryJump(cfg.JumpForce)
			m.onGround = false
			m.coyoteTimeCounter = 0
			m.jumpBufferCounter = 0
		} else {
			m.jumpBufferCounter = cfg.JumpBufferFrames
		}
	}

	// --- Variable Jump Height ---
	if !m.onGround && !input.IsSomeKeyPressed(ebiten.KeySpace) && body.vy16 < 0 {
		body.vy16 = int(float64(body.vy16) * cfg.JumpCutMultiplier)
	}
}
