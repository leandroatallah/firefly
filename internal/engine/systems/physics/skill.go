package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SkillState represents the possible states of a skill.
type SkillState string

const (
	StateReady    SkillState = "ready"
	StateActive   SkillState = "active"
	StateCooldown SkillState = "cooldown"
)

// Skill defines the interface for a passive player ability.
type Skill interface {
	Update(body *PhysicsBody, model *PlatformMovementModel)
	IsActive() bool
}

// ActiveSkill defines the interface for a skill that requires user input.
type ActiveSkill interface {
	Skill
	HandleInput(body *PhysicsBody, model *PlatformMovementModel)
	ActivationKey() ebiten.Key
}

// SkillBase provides a base implementation for common skill attributes.
type SkillBase struct {
	state    SkillState
	duration int // frames
	cooldown int // frames
	speed    int
	timer    int
}

func (s *SkillBase) Update(body *PhysicsBody, model *PlatformMovementModel) {}

// IsActive returns true if the skill is currently in its active phase.
func (s *SkillBase) IsActive() bool {
	return s.state == StateActive
}
