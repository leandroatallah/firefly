package kitskills

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/hajimehoshi/ebiten/v2"
)

// HorizontalMovementSkill handles left/right movement input.
type HorizontalMovementSkill struct {
	skill.SkillBase
	activationKey ebiten.Key
	axis          *input.HorizontalAxis
	prevLeft      bool
	prevRight     bool
}

// NewHorizontalMovementSkill creates a new HorizontalMovementSkill.
func NewHorizontalMovementSkill() *HorizontalMovementSkill {
	s := &HorizontalMovementSkill{
		axis: input.NewHorizontalAxis(),
	}
	s.SetState(skill.StateReady)
	return s
}

// Update is a no-op for horizontal movement (input handled in HandleInput).
func (s *HorizontalMovementSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	s.SkillBase.Update(b, model)
}

// ActivationKey returns the activation key (unused for movement).
func (s *HorizontalMovementSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

// HandleInput processes left/right movement commands.
func (s *HorizontalMovementSkill) HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, _ body.BodiesSpace) {
	if model != nil && model.IsInputBlocked() {
		return
	}
	if b.Immobile() {
		_, vy16 := b.Velocity()
		_, accY := b.Acceleration()
		b.SetVelocity(0, vy16)
		b.SetAcceleration(0, accY)
		return
	}

	cfg := config.Get()
	vx16, vy16 := b.Velocity()

	cmds := input.CommandsReader()
	moveLeft := cmds.Left
	moveRight := cmds.Right

	if moveLeft && !s.prevLeft {
		s.axis.Press(-1)
	} else if !moveLeft && s.prevLeft {
		s.axis.Release(-1)
	}
	if moveRight && !s.prevRight {
		s.axis.Press(1)
	} else if !moveRight && s.prevRight {
		s.axis.Release(1)
	}

	s.prevLeft = moveLeft
	s.prevRight = moveRight

	horizontalInertia := cfg.Physics.HorizontalInertia
	if val := b.HorizontalInertia(); val >= 0 {
		horizontalInertia = val
	}

	dir := s.axis.Value()
	if horizontalInertia > 0 {
		if dir < 0 {
			b.OnMoveLeft(b.Speed())
		} else if dir > 0 {
			b.OnMoveRight(b.Speed())
		}
	} else {
		vx16 = fp16.To16(dir * b.Speed())
	}

	b.SetVelocity(vx16, vy16)
}
