package kitskills

import (
	"strings"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/debug"
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
	wasMoving     bool
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

	debug.Watch("eight_dir_move", b.ID(), s.DebugCmd(cmds))

	moving := cmds.Left || cmds.Right || cmds.Up || cmds.Down
	if moving && !s.wasMoving {
		debug.Log("skill_activated", "skill=EightDirectionalMovementSkill player=%s", b.ID())
	}
	s.wasMoving = moving

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

func (s *EightDirectionalMovementSkill) DebugCmd(cmds input.PlayerCommands) string {
	out := []string{}

	if cmds.Left {
		out = append(out, "Left")
	}
	if cmds.Right {
		out = append(out, "Right")
	}
	if cmds.Up {
		out = append(out, "Up")
	}
	if cmds.Down {
		out = append(out, "Down")
	}

	return strings.Join(out, "+")
}
