package skill

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
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
	OnActive     func()
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
		if s.OnActive != nil {
			s.OnActive()
		}
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

	// Capture current bottom-center position to maintain it after size change
	x16, y16 := player.GetPosition16()
	centerX16 := x16 + (s.originalWidth*16)/2
	bottomY16 := y16 + (s.originalHeight*16)

	// Double size (Physical)
	newWidth := s.originalWidth * 2
	newHeight := s.originalHeight * 2
	player.SetSize(newWidth, newHeight)

	// Re-position to maintain bottom-center
	newX16 := centerX16 - (newWidth*16)/2
	newY16 := bottomY16 - (newHeight*16)
	player.SetPosition16(newX16, newY16)

	// Double scale (Visual) - This is the target scale, 
	// but the Growing state might override it to 1.0 for the animation.
	if c, ok := player.(interface{ SetScale(float64) }); ok {
		c.SetScale(2.0)
	}

	// Set State Growing
	if stateSetter, ok := player.(interface{ SetNewState(actors.ActorStateEnum) error }); ok {
		_ = stateSetter.SetNewState(gamestates.Growing)
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
	// Capture current bottom-center position to maintain it after size change
	shape := player.GetShape()
	currentWidth := shape.Width()
	currentHeight := shape.Height()

	x16, y16 := player.GetPosition16()
	centerX16 := x16 + (currentWidth*16)/2
	bottomY16 := y16 + (currentHeight*16)

	// Restore size
	player.SetSize(s.originalWidth, s.originalHeight)

	// Re-position to maintain bottom-center
	newX16 := centerX16 - (s.originalWidth*16)/2
	newY16 := bottomY16 - (s.originalHeight*16)
	player.SetPosition16(newX16, newY16)

	// Restore scale (Visual)
	if c, ok := player.(interface{ SetScale(float64) }); ok {
		c.SetScale(1.0)
	}

	// Set State Shrinking
	if stateSetter, ok := player.(interface{ SetNewState(actors.ActorStateEnum) error }); ok {
		_ = stateSetter.SetNewState(gamestates.Shrinking)
	}

	// Refresh collision bodies if possible
	if r, ok := player.(interface{ RefreshCollisions() }); ok {
		r.RefreshCollisions()
	}

	if s.OnDeactivate != nil {
		s.OnDeactivate()
	}
}
