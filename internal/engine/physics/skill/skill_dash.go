package skill

import (
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	"github.com/hajimehoshi/ebiten/v2"
)

// DashSkill implements a dash and air dash ability.
type DashSkill struct {
	SkillBase

	canAirDash    bool
	airDashUsed   bool
	activationKey ebiten.Key
	dashPressed   bool
}

// NewDashSkill creates a new DashSkill with default values.
func NewDashSkill() *DashSkill {
	return &DashSkill{
		SkillBase: SkillBase{
			state:    StateReady,
			duration: timing.FromDuration(200 * time.Millisecond), // 8 frames (short burst)
			cooldown: timing.FromDuration(750 * time.Millisecond), // 45 frames
			speed:    fp16.To16(6),
		},
		canAirDash:    true,
		airDashUsed:   false,
		activationKey: ebiten.KeyShift,
	}
}

// ActivationKey returns the activation key for the dash skill.
func (d *DashSkill) ActivationKey() ebiten.Key {
	return d.activationKey
}

// HandleInput checks for the dash activation key.
func (d *DashSkill) HandleInput(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	cmds := input.CommandsReader()
	dashPressed := cmds.Dash
	if dashPressed && !d.dashPressed {
		d.tryActivate(body, model, space)
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

	switch d.state {
	case StateActive:
		d.timer--
		if d.timer <= 0 {
			d.state = StateCooldown
			d.timer = d.cooldown
			model.SetDashActive(false, 0)
			model.SetGravityEnabled(true)
		} else {
			// Apply dash movement by setting it in the movement model
			dirX := 1
			if b.FaceDirection() == animation.FaceDirectionLeft {
				dirX = -1
			}
			model.SetDashActive(true, d.speed*dirX)
			// Override gravity during air dash to maintain height (Cuphead-style)
			if !model.OnGround() {
				model.SetGravityEnabled(false)
				vx, _ := b.Velocity()
				b.SetVelocity(vx, 0)
			}
		}
	case StateCooldown:
		d.timer--
		if d.timer <= 0 {
			d.state = StateReady
		}
	}
}

func (d *DashSkill) tryActivate(_ body.MovableCollidable, model *physicsmovement.PlatformMovementModel, _ body.BodiesSpace) {
	if d.state != StateReady {
		return
	}

	// Check for air dash conditions
	if !model.OnGround() {
		if !d.canAirDash || d.airDashUsed {
			return
		}
		d.airDashUsed = true
	}

	d.state = StateActive
	d.timer = d.duration
	// Optional: trigger a sound or visual effect here
}
