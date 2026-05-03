package skill

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
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
	Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel)
	IsActive() bool
}

// ActiveSkill defines the interface for a skill that requires user input.
type ActiveSkill interface {
	Skill
	HandleInput(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace)
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

// State returns the current skill state.
func (s *SkillBase) State() SkillState { return s.state }

// SetState sets the current skill state.
func (s *SkillBase) SetState(st SkillState) { s.state = st }

// Duration returns the active duration in frames.
func (s *SkillBase) Duration() int { return s.duration }

// SetDuration sets the active duration in frames.
func (s *SkillBase) SetDuration(d int) { s.duration = d }

// Cooldown returns the cooldown duration in frames.
func (s *SkillBase) Cooldown() int { return s.cooldown }

// SetCooldown sets the cooldown duration in frames.
func (s *SkillBase) SetCooldown(c int) { s.cooldown = c }

// Speed returns the skill speed value.
func (s *SkillBase) Speed() int { return s.speed }

// SetSpeed sets the skill speed value.
func (s *SkillBase) SetSpeed(sp int) { s.speed = sp }

// Timer returns the current timer value.
func (s *SkillBase) Timer() int { return s.timer }

// SetTimer sets the timer value.
func (s *SkillBase) SetTimer(t int) { s.timer = t }

// IncTimer increments the timer by one.
func (s *SkillBase) IncTimer() { s.timer++ }

// Update is a no-op base implementation.
func (s *SkillBase) Update(_ body.MovableCollidable, _ *physicsmovement.PlatformMovementModel) {}

// IsActive returns true if the skill is currently in its active phase.
func (s *SkillBase) IsActive() bool {
	return s.state == StateActive
}
