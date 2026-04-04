package skill

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/hajimehoshi/ebiten/v2"
)

type HorizontalMovementSkill struct {
	SkillBase
	activationKey ebiten.Key
	axis          *input.HorizontalAxis
	prevLeft      bool
	prevRight     bool
}

func NewHorizontalMovementSkill() *HorizontalMovementSkill {
	return &HorizontalMovementSkill{
		SkillBase: SkillBase{
			state: StateReady,
		},
		axis: input.NewHorizontalAxis(),
	}
}

func (s *HorizontalMovementSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	s.SkillBase.Update(b, model)
}

func (s *HorizontalMovementSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

func (s *HorizontalMovementSkill) HandleInput(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, _ body.BodiesSpace) {
	if model != nil && model.IsInputBlocked() {
		return
	}
	if body.Immobile() {
		_, vy16 := body.Velocity()
		_, accY := body.Acceleration()
		body.SetVelocity(0, vy16)
		body.SetAcceleration(0, accY)
		return
	}

	cfg := config.Get()
	vx16, vy16 := body.Velocity()

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
	if val := body.HorizontalInertia(); val >= 0 {
		horizontalInertia = val
	}

	dir := s.axis.Value()
	if horizontalInertia > 0 {
		if dir < 0 {
			body.OnMoveLeft(body.Speed())
		} else if dir > 0 {
			body.OnMoveRight(body.Speed())
		}
	} else {
		vx16 = fp16.To16(dir * body.Speed())
	}

	body.SetVelocity(vx16, vy16)
}
