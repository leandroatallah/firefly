package physics

import (
	"sync"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// Space centralizes physics bodies and collision resolution.
type Space struct {
	mu                        sync.RWMutex
	bodies                    map[string]body.Body
	tilemapDimensionsProvider TilemapDimensionsProvider
}

func NewSpace() *Space {
	return &Space{
		bodies: make(map[string]body.Body),
	}
}

func (s *Space) AddBody(b body.Body) {
	if b == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bodies == nil {
		s.bodies = make(map[string]body.Body)
	}

	s.bodies[b.ID()] = b
}

func (s *Space) RemoveBody(body body.Body) {
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

func (s *Space) Bodies() []body.Body {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]body.Body, 0, len(s.bodies))
	for _, b := range s.bodies {
		if b == nil {
			continue
		}
		res = append(res, b)
	}

	return res
}

func (s *Space) ResolveCollisions(body body.Body) (touching bool, blocking bool) {
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

func hasCollision(a, b body.Body) bool {
	rectsA := a.CollisionPosition()
	rectsB := b.CollisionPosition()

	for _, rectA := range rectsA {
		for _, rectB := range rectsB {
			if rectA.Overlaps(rectB) {
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
