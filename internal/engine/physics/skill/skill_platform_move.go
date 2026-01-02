package skill

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/input"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type HorizontalMovementSkill struct {
	SkillBase
	activationKey ebiten.Key
}

func NewHorizontalMovementSkill() *HorizontalMovementSkill {
	return &HorizontalMovementSkill{
		SkillBase: SkillBase{
			state: StateReady,
		},
	}
}

func (s *HorizontalMovementSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	s.SkillBase.Update(b, model)
}

func (s *HorizontalMovementSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

func (s *HorizontalMovementSkill) HandleInput(body body.MovableCollidable, _ *physicsmovement.PlatformMovementModel, _ body.BodiesSpace) {
	if body.Immobile() {
		_, vy16 := body.Velocity()
		_, accY := body.Acceleration()
		body.SetVelocity(0, vy16)
		body.SetAcceleration(0, accY)
		return
	}

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
