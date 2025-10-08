package physics

import (
	"sync"
)

// Space centralizes physics bodies and collision resolution.
type Space struct {
	mu                        sync.RWMutex
	bodies                    map[string]Body
	tilemapDimensionsProvider TilemapDimensionsProvider
}

func NewSpace() *Space {
	return &Space{
		bodies: make(map[string]Body),
	}
}

func (s *Space) AddBody(body Body) {
	if body == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bodies == nil {
		s.bodies = make(map[string]Body)
	}

	s.bodies[body.ID()] = body
}

func (s *Space) RemoveBody(body Body) {
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

func (s *Space) Bodies() []Body {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Body, 0, len(s.bodies))
	for _, b := range s.bodies {
		if b == nil {
			continue
		}
		res = append(res, b)
	}

	return res
}

func (s *Space) ResolveCollisions(body Body) (touching bool, blocking bool) {
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

func hasCollision(a, b Body) bool {
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
