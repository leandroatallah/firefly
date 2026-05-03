package kitskills

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	spacephysics "github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/hajimehoshi/ebiten/v2"
)

// JumpSkill implements a platformer jump with coyote time and jump buffering.
type JumpSkill struct {
	skill.SkillBase
	activationKey ebiten.Key

	coyoteTimeCounter int
	jumpBufferCounter int

	jumpCutMultiplier float64
	jumpCutPending    bool
	jumpPressed       bool

	// OnJump is called when the actor jumps.
	OnJump func(body body.MovableCollidable)
}

// NewJumpSkill creates a new JumpSkill with default values.
func NewJumpSkill() *JumpSkill {
	s := &JumpSkill{
		activationKey:     ebiten.KeySpace,
		jumpCutMultiplier: 1.0,
	}
	s.SetState(skill.StateReady)
	return s
}

// SetJumpCutMultiplier sets the velocity multiplier applied on early jump release.
// Clamped to (0.1, 1.0]; 1.0 disables the cut.
func (s *JumpSkill) SetJumpCutMultiplier(m float64) {
	if m <= 0 {
		m = 0.1
	} else if m > 1 {
		m = 1.0
	}
	s.jumpCutMultiplier = m
}

// ActivationKey returns the key that triggers the jump.
func (s *JumpSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

// HandleInput checks for the jump activation key.
func (s *JumpSkill) HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if model != nil && model.IsInputBlocked() {
		return
	}
	cmds := input.CommandsReader()
	jumpPressed := cmds.Jump
	if jumpPressed && !s.jumpPressed {
		s.tryActivate(b, model, space)
	}
	if !jumpPressed && s.jumpPressed && s.jumpCutPending {
		s.applyJumpCut(b)
	}
	s.jumpPressed = jumpPressed
}

// Update advances jump state (coyote time, jump buffering).
func (s *JumpSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	s.SkillBase.Update(b, model)

	if s.jumpCutPending && !b.IsGoingUp() {
		s.jumpCutPending = false
	}

	s.handleCoyoteAndJumpBuffering(b, model, model.OnGround())
}

func (s *JumpSkill) tryActivate(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	cfg := config.Get()
	if model.OnGround() || s.coyoteTimeCounter > 0 {
		force := int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())
		if force <= 0 {
			return
		}

		b.TryJump(force)
		s.jumpCutPending = true

		if s.OnJump != nil {
			s.OnJump(b)
		}

		// Check against map boundaries if the actor has a physics space.
		for _, other := range space.Bodies() {
			if other == nil || other.ID() == b.ID() {
				continue
			}

			if !spacephysics.HasCollision(b, other) {
				continue
			}

			if other.IsObstructive() {
				break
			}
		}

		model.SetOnGround(false)
		s.coyoteTimeCounter = 0
		s.jumpBufferCounter = 0
	} else {
		s.jumpBufferCounter = cfg.Physics.JumpBufferFrames
	}
}

// handleCoyoteAndJumpBuffering manages coyote time and jump buffering logic.
func (s *JumpSkill) handleCoyoteAndJumpBuffering(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, wasOnGround bool) {
	cfg := config.Get()

	if model.OnGround() {
		s.coyoteTimeCounter = cfg.Physics.CoyoteTimeFrames
	} else {
		if s.coyoteTimeCounter > 0 {
			s.coyoteTimeCounter--
		}
	}

	if s.jumpBufferCounter > 0 {
		s.jumpBufferCounter--
	}

	if !wasOnGround && model.OnGround() && s.jumpBufferCounter > 0 {
		force := int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())
		if force <= 0 {
			return
		}

		b.TryJump(force)
		s.jumpCutPending = true
		if s.OnJump != nil {
			s.OnJump(b)
		}
		model.SetOnGround(false)
		s.jumpBufferCounter = 0
		s.coyoteTimeCounter = 0
	}
}

func (s *JumpSkill) applyJumpCut(b body.MovableCollidable) {
	if b.IsGoingUp() {
		vx, vy := b.Velocity()
		b.SetVelocity(vx, int(float64(vy)*s.jumpCutMultiplier))
	}
	s.jumpCutPending = false
}
