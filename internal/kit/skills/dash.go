package kitskills

import (
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	"github.com/hajimehoshi/ebiten/v2"
)

// DashSkill implements a dash and air dash ability.
type DashSkill struct {
	skill.SkillBase

	canAirDash    bool
	airDashUsed   bool
	activationKey ebiten.Key
	dashPressed   bool
}

// NewDashSkill creates a new DashSkill with default values.
func NewDashSkill() *DashSkill {
	d := &DashSkill{
		canAirDash:    true,
		airDashUsed:   false,
		activationKey: ebiten.KeyShift,
	}
	d.SetState(skill.StateReady)
	d.SetDuration(timing.FromDuration(200 * time.Millisecond))
	d.SetCooldown(timing.FromDuration(750 * time.Millisecond))
	d.SetSpeed(fp16.To16(6))
	return d
}

// ActivationKey returns the activation key for the dash skill.
func (d *DashSkill) ActivationKey() ebiten.Key {
	return d.activationKey
}

// HandleInput checks for the dash activation key.
func (d *DashSkill) HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	cmds := input.CommandsReader()
	dashPressed := cmds.Dash
	if dashPressed && !d.dashPressed {
		d.tryActivate(b, model, space)
	}
	d.dashPressed = dashPressed
}

// Update manages the skill's state, timers, and applies its effects.
func (d *DashSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	d.SkillBase.Update(b, model)

	// Reset air dash capability when the player lands.
	if model.OnGround() {
		d.airDashUsed = false
	}

	switch d.State() {
	case skill.StateActive:
		d.SetTimer(d.Timer() - 1)
		if d.Timer() <= 0 {
			d.SetState(skill.StateCooldown)
			d.SetTimer(d.Cooldown())
			model.SetDashActive(false, 0)
			model.SetGravityEnabled(true)
		} else {
			// Apply dash movement by setting it in the movement model
			dirX := 1
			if b.FaceDirection() == animation.FaceDirectionLeft {
				dirX = -1
			}
			model.SetDashActive(true, d.Speed()*dirX)
			// Override gravity during air dash to maintain height (Cuphead-style)
			if !model.OnGround() {
				model.SetGravityEnabled(false)
				vx, _ := b.Velocity()
				b.SetVelocity(vx, 0)
			}
		}
	case skill.StateCooldown:
		d.SetTimer(d.Timer() - 1)
		if d.Timer() <= 0 {
			d.SetState(skill.StateReady)
		}
	}
}

func (d *DashSkill) tryActivate(_ body.MovableCollidable, model *physicsmovement.PlatformMovementModel, _ body.BodiesSpace) {
	if d.State() != skill.StateReady {
		return
	}

	// Check for air dash conditions
	if !model.OnGround() {
		if !d.canAirDash || d.airDashUsed {
			return
		}
		d.airDashUsed = true
	}

	d.SetState(skill.StateActive)
	d.SetTimer(d.Duration())
}
