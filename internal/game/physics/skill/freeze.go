package skill

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type FreezeSkill struct {
	state               engineskill.SkillState
	duration            int
	cooldown            int
	timer               int
	activationRequested bool

	frozenBodies []body.Movable

	// Callbacks for external systems (audio, vfx)
	OnActivate   func()
	OnActive     func()
	OnDeactivate func()
}

func NewFreezeSkill() *FreezeSkill {
	return &FreezeSkill{
		state:        engineskill.StateReady,
		duration:     timing.FromDuration(3 * time.Second),  // Default 3s duration
		cooldown:     timing.FromDuration(10 * time.Second), // Default 10s cooldown
		frozenBodies: make([]body.Movable, 0),
	}
}

// ActivationKey returns 0 as this skill is item-activated
func (s *FreezeSkill) ActivationKey() ebiten.Key {
	return 0
}

func (s *FreezeSkill) IsActive() bool {
	return s.state == engineskill.StateActive
}

// RequestActivation flags the skill to be activated on the next update cycle
func (s *FreezeSkill) RequestActivation() {
	s.activationRequested = true
}

func (s *FreezeSkill) Reset() {
	s.deactivate()
	s.state = engineskill.StateReady
	s.timer = 0
	s.activationRequested = false
}

func (s *FreezeSkill) HandleInput(player body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	if s.activationRequested {
		s.activationRequested = false
		if s.state == engineskill.StateReady {
			s.activate(player, space)
		} else if s.state == engineskill.StateCooldown {
			// Allow collecting power-up during cooldown to reset cooldown timer
			// This enables players to collect multiple power-ups in sequence
			s.state = engineskill.StateReady
			s.timer = 0
			s.activate(player, space)
		}
		// If state is Active, ignore the request (already active)
	}
}

func (s *FreezeSkill) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
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

func (s *FreezeSkill) activate(player body.MovableCollidable, space body.BodiesSpace) {
	s.state = engineskill.StateActive
	s.timer = s.duration

	// Freeze all other movable bodies
	s.frozenBodies = s.frozenBodies[:0] // Clear previous list
	
	if space == nil {
		return
	}

	for _, b := range space.Bodies() {
		if b == nil {
			continue
		}
		
		// Skip player using ID check for robustness
		if b.ID() == player.ID() {
			continue
		}

		if movable, ok := b.(body.Movable); ok {
			// Only freeze if not already frozen and not obstructive (e.g. walls shouldn't be frozen? wait, walls are usually immovable)
			// But movable obstacles might be frozen.
			// Let's assume we freeze everything movable except player.
			if !movable.Freeze() {
				movable.SetFreeze(true)
				s.frozenBodies = append(s.frozenBodies, movable)
			}
		}
	}

	if s.OnActivate != nil {
		s.OnActivate()
	}
}

func (s *FreezeSkill) deactivate() {
	// Unfreeze bodies we froze
	for _, b := range s.frozenBodies {
		// Only unfreeze if we were the ones who froze it (implied by being in list)
		// Check if it's still frozen? Yes, SetFreeze(false) is idempotent usually.
		b.SetFreeze(false)
	}
	s.frozenBodies = s.frozenBodies[:0]

	if s.OnDeactivate != nil {
		s.OnDeactivate()
	}
}
