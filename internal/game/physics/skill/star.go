package skill

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type StarSkill struct {
	state               engineskill.SkillState
	duration            int
	cooldown            int
	timer               int
	activationRequested bool

	// Callbacks for external systems (audio, vfx)
	OnActive func()

	// Respawn tracking
	itemToRespawn body.Collidable
	itemPosition  image.Point
	respawnSpace  body.BodiesSpace
}

func NewStarSkill() *StarSkill {
	return &StarSkill{
		state:    engineskill.StateReady,
		duration: timing.FromDuration(10 * time.Second),
		cooldown: timing.FromDuration(5 * time.Second),
	}
}

// ActivationKey returns 0 as this skill is item-activated
func (s *StarSkill) ActivationKey() ebiten.Key {
	return 0
}

func (s *StarSkill) IsActive() bool {
	return s.state == engineskill.StateActive
}

// RequestActivation flags the skill to be activated on the next update cycle
func (s *StarSkill) RequestActivation() {
	s.activationRequested = true
}

// RequestActivationWithItem flags the skill to be activated and registers the item for respawn
func (s *StarSkill) RequestActivationWithItem(item body.Collidable, space body.BodiesSpace) {
	s.activationRequested = true
	if item != nil {
		s.itemToRespawn = item
		pos := item.Position()
		s.itemPosition = image.Point{X: pos.Min.X, Y: pos.Min.Y}
		s.respawnSpace = space
	}
}

func (s *StarSkill) Reset() {
	if s.state == engineskill.StateActive {
		s.deactivate()
	}
	s.state = engineskill.StateReady
	s.timer = 0
	s.activationRequested = false
}

func (s *StarSkill) HandleInput(player body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if s.activationRequested {
		s.activationRequested = false
		if s.state == engineskill.StateReady {
			s.activate()
		} else if s.state == engineskill.StateCooldown {
			// Allow collecting power-up during cooldown to reset cooldown timer
			// This enables players to collect multiple power-ups in sequence
			s.state = engineskill.StateReady
			s.timer = 0
			s.activate()
		}
		// If state is Active, ignore the request (already active)
	}
}

func (s *StarSkill) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	switch s.state {
	case engineskill.StateActive:
		s.timer--
		if s.OnActive != nil {
			s.OnActive()
		}
		if s.timer <= 0 {
			s.deactivate()
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

func (s *StarSkill) activate() {
	s.state = engineskill.StateActive
	s.timer = s.duration
}

func (s *StarSkill) deactivate() {
	// Respawn the power-up item if registered
	s.respawnItem()
}

// respawnItem restores the power-up item to its original position
func (s *StarSkill) respawnItem() {
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
