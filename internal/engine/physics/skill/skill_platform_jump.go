package skill

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	spacephysics "github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

type JumpSkill struct {
	SkillBase
	activationKey ebiten.Key

	coyoteTimeCounter int
	jumpBufferCounter int

	jumpCutMultiplier float64
	jumpCutPending    bool

	OnJump func(body body.MovableCollidable)
}

func NewJumpSkill() *JumpSkill {
	return &JumpSkill{
		SkillBase: SkillBase{
			state: StateReady,
		},
		activationKey:     ebiten.KeySpace,
		jumpCutMultiplier: 1.0,
	}
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

func (s *JumpSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

// HandleInput checks for the jump activation key.
func (s *JumpSkill) HandleInput(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if model != nil && model.IsInputBlocked() {
		return
	}
	if inpututil.IsKeyJustPressed(s.activationKey) {
		s.tryActivate(body, model, space)
	}
	if inpututil.IsKeyJustReleased(s.activationKey) && s.jumpCutPending {
		s.applyJumpCut(body)
	}
}

func (s *JumpSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	s.SkillBase.Update(b, model)

	if s.jumpCutPending && !b.IsGoingUp() {
		s.jumpCutPending = false
	}

	s.handleCoyoteAndJumpBuffering(b, model, model.OnGround())
}

func (s *JumpSkill) tryActivate(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	cfg := config.Get()
	if model.OnGround() || s.coyoteTimeCounter > 0 {
		force := int(float64(cfg.Physics.JumpForce) * body.JumpForceMultiplier())
		if force <= 0 {
			return
		}

		body.TryJump(force)
		s.jumpCutPending = true

		if s.OnJump != nil {
			s.OnJump(body)
		}

		// Check against map boundaries if the actor has a physics space.
		for _, other := range space.Bodies() {
			if other == nil || other.ID() == body.ID() {
				continue
			}

			if !spacephysics.HasCollision(body, other) {
				continue
			}

			if other.IsObstructive() {
				// blocking = true
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

// Coyote Time & Jump Buffering
func (s *JumpSkill) handleCoyoteAndJumpBuffering(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, wasOnGround bool) {
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
		force := int(float64(cfg.Physics.JumpForce) * body.JumpForceMultiplier())
		if force <= 0 {
			return
		}

		body.TryJump(force)
		s.jumpCutPending = true
		if s.OnJump != nil {
			s.OnJump(body)
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
