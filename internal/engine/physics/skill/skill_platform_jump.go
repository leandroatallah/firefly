package skill

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	spacephysics "github.com/leandroatallah/firefly/internal/engine/physics/space"
)

type JumpSkill struct {
	SkillBase
	activationKey ebiten.Key

	coyoteTimeCounter int
	jumpBufferCounter int
}

func NewJumpSkill() *JumpSkill {
	return &JumpSkill{
		SkillBase: SkillBase{
			state: StateReady,
		},
		activationKey: ebiten.KeySpace,
	}
}

func (s *JumpSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

// HandleInput checks for the dash activation key.
func (s *JumpSkill) HandleInput(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if inpututil.IsKeyJustPressed(s.activationKey) {
		s.tryActivate(body, model, space)
	}
}

func (s *JumpSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	s.SkillBase.Update(b, model)

	s.handleCoyoteAndJumpBuffering(b, model, model.OnGround())
}

func (s *JumpSkill) tryActivate(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	cfg := config.Get()
	if model.OnGround() || s.coyoteTimeCounter > 0 {
		body.TryJump(cfg.Physics.JumpForce)

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
		body.TryJump(cfg.Physics.JumpForce)
		model.SetOnGround(false)
		s.jumpBufferCounter = 0
		s.coyoteTimeCounter = 0
	}
}
