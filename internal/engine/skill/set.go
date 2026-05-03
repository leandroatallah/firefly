package skill

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/hajimehoshi/ebiten/v2"
)

// Set is a registry of skills belonging to an actor.
type Set struct {
	skills []Skill
}

// NewSet creates an empty Set.
func NewSet() *Set {
	return &Set{}
}

// Add registers a skill in the set.
func (s *Set) Add(sk Skill) {
	s.skills = append(s.skills, sk)
}

// Get returns the ActiveSkill registered for the given key, if any.
func (s *Set) Get(key ebiten.Key) (ActiveSkill, bool) {
	for _, sk := range s.skills {
		if as, ok := sk.(ActiveSkill); ok {
			if as.ActivationKey() == key {
				return as, true
			}
		}
	}
	return nil, false
}

// Update calls Update on every registered skill.
func (s *Set) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	for _, sk := range s.skills {
		sk.Update(actor, model)
	}
}

// ActiveCount returns the number of currently active skills.
func (s *Set) ActiveCount() int {
	count := 0
	for _, sk := range s.skills {
		if sk.IsActive() {
			count++
		}
	}
	return count
}

// All returns the registered skills in insertion order.
func (s *Set) All() []Skill {
	return s.skills
}
