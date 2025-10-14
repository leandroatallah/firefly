package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/config"
)

// DashSkill implements a dash and air dash ability.
type DashSkill struct {
	SkillBase

	canAirDash    bool
	airDashUsed   bool
	activationKey ebiten.Key
}

// NewDashSkill creates a new DashSkill with default values.
func NewDashSkill() *DashSkill {
	return &DashSkill{
		SkillBase: SkillBase{
			state:    StateReady,
			duration: 8,  // 8 frames (short burst)
			cooldown: 45, // 45 frames cooldown
			speed:    10 * config.Get().Unit,
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
func (d *DashSkill) HandleInput(body *PhysicsBody, model *PlatformMovementModel) {
	if inpututil.IsKeyJustPressed(d.activationKey) {
		d.tryActivate(body, model)
	}
}

// Update manages the skill's state, timers, and applies its effects.
func (d *DashSkill) Update(body *PhysicsBody, model *PlatformMovementModel) {
	d.SkillBase.Update(body, model)

	// Reset air dash capability when the player lands.
	if model.onGround {
		d.airDashUsed = false
	}

	switch d.state {
	case StateActive:
		d.timer--
		if d.timer <= 0 {
			d.state = StateCooldown
			d.timer = d.cooldown
			body.vx16 = 0 // Stop horizontal movement after dash
		} else {
			// Apply dash movement
			var dirX int = 1
			if body.FaceDirection() == FaceDirectionLeft {
				dirX = -1
			}
			body.vx16 = d.speed * dirX
			body.vy16 = 0 // Maintain horizontal trajectory
		}
	case StateCooldown:
		d.timer--
		if d.timer <= 0 {
			d.state = StateReady
		}
	}
}

func (d *DashSkill) tryActivate(body *PhysicsBody, model *PlatformMovementModel) {
	if d.state != StateReady {
		return
	}

	// Check for air dash conditions
	if !model.onGround {
		if !d.canAirDash || d.airDashUsed {
			return
		}
		d.airDashUsed = true
	}

	d.state = StateActive
	d.timer = d.duration
	// Optional: trigger a sound or visual effect here
}
