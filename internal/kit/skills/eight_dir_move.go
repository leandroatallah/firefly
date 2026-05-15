package kitskills

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/hajimehoshi/ebiten/v2"
)

// EightDirectionalMovementSkill drives a body in 8 directions on the X/Y
// ground plane. Genre-agnostic: imports neither beat-em-up, platformer, nor
// top-down packages. Y is ground-plane depth; altitude is never written.
type EightDirectionalMovementSkill struct {
	skill.SkillBase
	activationKey ebiten.Key
}

func NewEightDirectionalMovementSkill() *EightDirectionalMovementSkill {
	s := &EightDirectionalMovementSkill{}
	s.SetState(skill.StateReady)
	return s
}

func (s *EightDirectionalMovementSkill) Update(b body.MovableCollidable, model physicsmovement.MovementModel) {
	s.SkillBase.Update(b, model)
}

func (s *EightDirectionalMovementSkill) ActivationKey() ebiten.Key {
	return s.activationKey
}

func (s *EightDirectionalMovementSkill) HandleInput(b body.MovableCollidable, model physicsmovement.MovementModel, _ body.BodiesSpace) {
	if blocker, ok := model.(physicsmovement.InputBlocker); ok && blocker.IsInputBlocked() {
		return
	}
	if b.Immobile() {
		b.SetVelocity(0, 0)
		b.SetAcceleration(0, 0)
		return
	}

	cmds := input.CommandsReader()
	speed := b.Speed()

	if cmds.Left {
		b.OnMoveLeft(speed)
	}
	if cmds.Right {
		b.OnMoveRight(speed)
	}
	if cmds.Up {
		b.OnMoveUp(speed)
	}
	if cmds.Down {
		b.OnMoveDown(speed)
	}
}
