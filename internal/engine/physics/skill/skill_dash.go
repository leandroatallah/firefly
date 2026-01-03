package skill

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
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
			speed:    fp16.To16(10),
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
	if inpututil.IsKeyJustPressed(d.activationKey) {
		d.tryActivate(body, model, space)
	}
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
			model.SetDashActive(false, 0) // Signal to the movement model that dash is no longer active
		} else {
			// Apply dash movement by setting it in the movement model
			var dirX int = 1
			if b.FaceDirection() == body.FaceDirectionLeft {
				dirX = -1
			}
			model.SetDashActive(true, d.speed*dirX)
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
