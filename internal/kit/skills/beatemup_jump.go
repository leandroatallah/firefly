package kitskills

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/hajimehoshi/ebiten/v2"
)

// BeatEmUpJumpSkill drives altitude-axis jumps on *BeatEmUpMovementModel.
type BeatEmUpJumpSkill struct {
	skill.SkillBase
	activationKey     ebiten.Key
	coyoteTimeCounter int
	jumpBufferCounter int
	jumpCutMultiplier float64
	jumpCutPending    bool
	jumpPressed       bool

	// OnJump is invoked whenever the actor begins an altitude-axis jump.
	OnJump func(b body.MovableCollidable)
}

// NewBeatEmUpJumpSkill returns a ready skill with default jump-cut multiplier (1.0).
func NewBeatEmUpJumpSkill() *BeatEmUpJumpSkill {
	s := &BeatEmUpJumpSkill{
		activationKey:     ebiten.KeySpace,
		jumpCutMultiplier: 1.0,
	}
	s.SetState(skill.StateReady)
	return s
}

// SetJumpCutMultiplier clamps m to (0.1, 1.0]; 1.0 disables the cut.
func (s *BeatEmUpJumpSkill) SetJumpCutMultiplier(m float64) {
	if m <= 0 {
		s.jumpCutMultiplier = 0.1
	} else if m > 1 {
		s.jumpCutMultiplier = 1.0
	} else {
		s.jumpCutMultiplier = m
	}
}

// ActivationKey returns the key that triggers the jump.
func (s *BeatEmUpJumpSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

// HandleInput processes leading/trailing-edge jump input on the altitude axis.
func (s *BeatEmUpJumpSkill) HandleInput(b body.MovableCollidable, model physicsmovement.MovementModel, _ body.BodiesSpace) {
	bm, _ := model.(*physicsmovement.BeatEmUpMovementModel)
	if bm == nil {
		return
	}
	if bm.IsInputBlocked() {
		return
	}

	pressed := input.CommandsReader().Jump

	if pressed && !s.jumpPressed {
		s.tryActivate(b)
	}
	if !pressed && s.jumpPressed && s.jumpCutPending {
		s.applyJumpCut(b)
	}

	s.jumpPressed = pressed
}

func (s *BeatEmUpJumpSkill) tryActivate(b body.MovableCollidable) {
	cfg := config.Get()
	grounded := b.Altitude() <= 0
	if grounded || s.coyoteTimeCounter > 0 {
		force := int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())
		if force <= 0 {
			return
		}
		b.SetVAltitude16(-force)
		s.jumpCutPending = true
		if s.OnJump != nil {
			s.OnJump(b)
		}
		s.coyoteTimeCounter = 0
		s.jumpBufferCounter = 0
	} else {
		s.jumpBufferCounter = cfg.Physics.JumpBufferFrames
	}
}

func (s *BeatEmUpJumpSkill) applyJumpCut(b body.MovableCollidable) {
	v := b.VAltitude16()
	if v < 0 {
		b.SetVAltitude16(int(float64(v) * s.jumpCutMultiplier))
	}
	s.jumpCutPending = false
}

// Update advances coyote/buffer counters and fires buffered jumps on landing.
func (s *BeatEmUpJumpSkill) Update(b body.MovableCollidable, model physicsmovement.MovementModel) {
	s.SkillBase.Update(b, model)
	bm, _ := model.(*physicsmovement.BeatEmUpMovementModel)
	if bm == nil {
		return
	}
	if b.Freeze() {
		return
	}

	// Clear jumpCutPending once apex passes
	if s.jumpCutPending && b.VAltitude16() >= 0 {
		s.jumpCutPending = false
	}

	grounded := b.Altitude() <= 0
	cfg := config.Get()

	// Coyote
	if grounded {
		s.coyoteTimeCounter = cfg.Physics.CoyoteTimeFrames
	} else if s.coyoteTimeCounter > 0 {
		s.coyoteTimeCounter--
	}

	// Buffer decay
	if s.jumpBufferCounter > 0 {
		s.jumpBufferCounter--
	}

	// Buffered jump fires on landing
	if grounded && s.jumpBufferCounter > 0 {
		force := int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())
		if force <= 0 {
			return
		}
		b.SetVAltitude16(-force)
		s.jumpCutPending = true
		if s.OnJump != nil {
			s.OnJump(b)
		}
		s.jumpBufferCounter = 0
		s.coyoteTimeCounter = 0
	}
}
