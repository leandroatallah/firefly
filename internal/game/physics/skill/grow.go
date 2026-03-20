package skill

import (
	"image"
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

	// Respawn tracking
	itemToRespawn body.Collidable
	itemPosition  image.Point
	respawnSpace  body.BodiesSpace
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

// RequestActivationWithItem flags the skill to be activated and registers the item for respawn
func (s *GrowSkill) RequestActivationWithItem(item body.Collidable, space body.BodiesSpace) {
	s.activationRequested = true
	if item != nil {
		s.itemToRespawn = item
		pos := item.Position()
		s.itemPosition = image.Point{X: pos.Min.X, Y: pos.Min.Y}
		s.respawnSpace = space
	}
}

func (s *GrowSkill) Reset(player body.MovableCollidable) {
	if s.state == engineskill.StateActive {
		s.RestoreNormalSize(player)
	}
	s.state = engineskill.StateReady
	s.timer = 0
	s.activationRequested = false
}

func (s *GrowSkill) RestoreNormalSize(player body.MovableCollidable) {
	// Restore size
	player.SetSize(s.originalWidth, s.originalHeight)

	// Restore scale (Visual)
	if c, ok := player.(interface{ SetScale(float64) }); ok {
		c.SetScale(1.0)
	}

	// Refresh collision bodies if possible
	if r, ok := player.(interface{ RefreshCollisions() }); ok {
		r.RefreshCollisions()
	}
}

func (s *GrowSkill) HandleInput(player body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if s.activationRequested {
		s.activationRequested = false
		if s.state == engineskill.StateReady {
			s.activate(player)
		} else if s.state == engineskill.StateCooldown {
			// Allow collecting power-up during cooldown to reset cooldown timer
			// This enables players to collect multiple power-ups in sequence
			s.state = engineskill.StateReady
			s.timer = 0
			s.activate(player)
		}
		// If state is Active, ignore the request (already active)
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
	// Transition to Shrinking state.
	// We DO NOT change the physical body size yet; the actor's state transition logic
	// will handle the final size restoration when the animation finishes.
	// This keeps the 32x32 body aligned with the 32x32 animation frame during transition.

	// Set State Shrinking
	if stateSetter, ok := player.(interface{ SetNewState(actors.ActorStateEnum) error }); ok {
		_ = stateSetter.SetNewState(gamestates.Shrinking)
	}

	// Refresh collision bodies if possible
	if r, ok := player.(interface{ RefreshCollisions() }); ok {
		r.RefreshCollisions()
	}

	// Respawn the power-up item if registered
	s.respawnItem()

	if s.OnDeactivate != nil {
		s.OnDeactivate()
	}
}

// respawnItem restores the power-up item to its original position
func (s *GrowSkill) respawnItem() {
	if s.itemToRespawn == nil || s.respawnSpace == nil {
		return
	}

	// Check if item implements the removed interface
	type Removable interface {
		IsRemoved() bool
		SetRemoved(bool)
	}

	if item, ok := s.itemToRespawn.(Removable); ok {
		// Mark item as not removed
		item.SetRemoved(false)

		// Restore position
		s.itemToRespawn.SetPosition(s.itemPosition.X, s.itemPosition.Y)

		// Re-add to physics space
		s.respawnSpace.AddBody(s.itemToRespawn)
	}

	// Clear respawn tracking
	s.itemToRespawn = nil
	s.itemPosition = image.Point{}
	s.respawnSpace = nil
}
