package skill

import (
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
)

// SkillDeps contains all dependencies required to instantiate skills from config.
type SkillDeps struct {
	Inventory    combat.Inventory
	OnJump       func(interface{})
	EventManager interface{ Publish(interface{}) }
}

// FromConfig instantiates skills from a SkillsConfig.
// Returns an empty slice if cfg is nil.
// Skips skills with nil sub-config or Enabled == false.
// Skips shooting skill if Inventory is nil.
func FromConfig(cfg *schemas.SkillsConfig, deps SkillDeps) []Skill {
	if cfg == nil {
		return []Skill{}
	}

	var skills []Skill

	// Movement skill
	if cfg.Movement != nil && isEnabled(cfg.Movement.Enabled) {
		skills = append(skills, NewHorizontalMovementSkill())
	}

	// Jump skill
	if cfg.Jump != nil && isEnabled(cfg.Jump.Enabled) {
		jumpSkill := NewJumpSkill()
		if cfg.Jump.JumpCutMultiplier > 0 {
			jumpSkill.SetJumpCutMultiplier(cfg.Jump.JumpCutMultiplier)
		}
		if deps.OnJump != nil {
			jumpSkill.OnJump = func(b body.MovableCollidable) {
				deps.OnJump(b)
			}
		}
		skills = append(skills, jumpSkill)
	}

	// Dash skill
	if cfg.Dash != nil && isEnabled(cfg.Dash.Enabled) {
		dashSkill := NewDashSkill()
		if cfg.Dash.DurationMs > 0 {
			dashSkill.duration = timing.FromDuration(time.Duration(cfg.Dash.DurationMs) * time.Millisecond)
		}
		if cfg.Dash.CooldownMs > 0 {
			dashSkill.cooldown = timing.FromDuration(time.Duration(cfg.Dash.CooldownMs) * time.Millisecond)
		}
		if cfg.Dash.Speed > 0 {
			dashSkill.speed = fp16.To16(cfg.Dash.Speed)
		}
		if cfg.Dash.CanAirDash != nil {
			dashSkill.canAirDash = *cfg.Dash.CanAirDash
		}
		skills = append(skills, dashSkill)
	}

	// Shooting skill
	if cfg.Shooting != nil && isEnabled(cfg.Shooting.Enabled) && deps.Inventory != nil {
		shootingSkill := NewShootingSkill(deps.Inventory)
		skills = append(skills, shootingSkill)
	}

	return skills
}

// isEnabled returns true if the pointer is nil or points to true.
func isEnabled(enabled *bool) bool {
	return enabled == nil || *enabled
}
