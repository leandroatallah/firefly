package skill

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type GrowSkill struct {
	state               engineskill.SkillState
	duration            int
	cooldown            int
	timer               int
	activationRequested bool

	originalWidth  int
	originalHeight int

	// Callbacks for external systems (audio, vfx)
	OnActivate   func()
	OnDeactivate func()
}

func NewGrowSkill() *GrowSkill {
	return &GrowSkill{
		state:    engineskill.StateReady,
		duration: timing.FromDuration(10 * time.Second), // Default 10s duration
		cooldown: timing.FromDuration(5 * time.Second),  // Default 5s cooldown
	}
}

// ActivationKey returns 0 as this skill is item-activated
func (s *GrowSkill) ActivationKey() ebiten.Key {
	return 0
}

func (s *GrowSkill) IsActive() bool {
	return s.state == engineskill.StateActive
}

// RequestActivation flags the skill to be activated on the next update cycle
func (s *GrowSkill) RequestActivation() {
	s.activationRequested = true
}

func (s *GrowSkill) Reset(player body.MovableCollidable) {
	if s.state == engineskill.StateActive {
		s.deactivate(player)
	}
	s.state = engineskill.StateReady
	s.timer = 0
	s.activationRequested = false
}

func (s *GrowSkill) HandleInput(player body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if s.activationRequested {
		s.activationRequested = false
		if s.state == engineskill.StateReady {
			s.activate(player)
		}
	}
}

func (s *GrowSkill) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	switch s.state {
	case engineskill.StateActive:
		s.timer--
		if s.timer <= 0 {
			s.deactivate(actor)
			s.state = engineskill.StateCooldown
			s.timer = s.cooldown
		}
	case engineskill.StateCooldown:
		s.timer--
		if s.timer <= 0 {
			s.state = engineskill.StateReady
		}
	}
}

func (s *GrowSkill) activate(player body.MovableCollidable) {
	s.state = engineskill.StateActive
	s.timer = s.duration

	shape := player.GetShape()
	s.originalWidth = shape.Width()
	s.originalHeight = shape.Height()

	// Double size (Physical)
	player.SetSize(s.originalWidth*2, s.originalHeight*2)

	// Double scale (Visual)
	if c, ok := player.(interface{ SetScale(float64) }); ok {
		c.SetScale(2.0)
	}

	// Refresh collision bodies if possible
	if r, ok := player.(interface{ RefreshCollisions() }); ok {
		r.RefreshCollisions()
	}

	if s.OnActivate != nil {
		s.OnActivate()
	}
}

func (s *GrowSkill) deactivate(player body.MovableCollidable) {
	// Restore size
	player.SetSize(s.originalWidth, s.originalHeight)

	// Restore scale
	if c, ok := player.(interface{ SetScale(float64) }); ok {
		c.SetScale(1.0)
	}

	// Refresh collision bodies if possible
	if r, ok := player.(interface{ RefreshCollisions() }); ok {
		r.RefreshCollisions()
	}

	if s.OnDeactivate != nil {
		s.OnDeactivate()
	}
}
