package physics

import (
	"log"
	"sync"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// Space centralizes physics bodies and collision resolution.
type Space struct {
	mu                        sync.RWMutex
	bodies                    map[string]body.Collidable
	tilemapDimensionsProvider TilemapDimensionsProvider
}

func NewSpace() *Space {
	return &Space{
		bodies: make(map[string]body.Collidable),
	}
}

func (s *Space) AddBody(b body.Collidable) {
	if b == nil {
		return
	}

	if b.ID() == "" {
		log.Fatal("(*Space).AddBody: A body must have an ID")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bodies == nil {
		s.bodies = make(map[string]body.Collidable)
	}

	s.bodies[b.ID()] = b
}

func (s *Space) RemoveBody(body body.Collidable) {
	if body == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bodies == nil {
		return
	}

	delete(s.bodies, body.ID())
}

func (s *Space) Bodies() []body.Collidable {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]body.Collidable, 0, len(s.bodies))
	for _, b := range s.bodies {
		if b == nil {
			continue
		}
		res = append(res, b)
	}

	return res
}

// ResolveCollisions compare a body parameter with all bodies in space.
// Returns boolean values if is touching or blocking.
func (s *Space) ResolveCollisions(body body.Collidable) (touching bool, blocking bool) {
	if body == nil {
		return false, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, other := range s.bodies {
		if other == nil || other.ID() == body.ID() {
			continue
		}

		if !hasCollision(body, other) {
			continue
		}

		body.OnTouch(other)
		other.OnTouch(body)
		touching = true

		if other.IsObstructive() {
			body.OnBlock(other)
			other.OnBlock(body)
			blocking = true
			break
		}
	}

	return touching, blocking
}

func hasCollision(a, b body.Collidable) bool {
	// Every body must have an ID
	if a.ID() == "" || b.ID() == "" {
		return false
	}

	// Prevent to check the same body
	if a.ID() == b.ID() {
		return false
	}

	rectsA := a.CollisionPosition()
	rectsB := b.CollisionPosition()

	for _, r := range rectsA {
		for _, s := range rectsB {
			if r.Overlaps(s) {
				return true
			}
		}
	}

	return false
}

func (s *Space) SetTilemapDimensionsProvider(provider TilemapDimensionsProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tilemapDimensionsProvider = provider
}

func (s *Space) GetTilemapDimensionsProvider() TilemapDimensionsProvider {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tilemapDimensionsProvider
}
