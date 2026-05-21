package space

import (
	"image"
	"log"
	"sort"
	"sync"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
)

// Space centralizes physics bodies and collision resolution.
type Space struct {
	mu                        sync.RWMutex
	bodies                    map[string]body.Collidable
	bodiesCache               []body.Collidable
	cacheDirty                bool
	toBeRemoved               []body.Collidable
	tilemapDimensionsProvider tilemaplayer.TilemapDimensionsProvider
}

func NewSpace() body.BodiesSpace {
	return &Space{
		bodies:     make(map[string]body.Collidable),
		cacheDirty: true,
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
	s.cacheDirty = true
}

func (s *Space) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bodies = make(map[string]body.Collidable)
	s.cacheDirty = true
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
	s.cacheDirty = true
}

func (s *Space) QueueForRemoval(body body.Collidable) {
	s.toBeRemoved = append(s.toBeRemoved, body)
}

func (s *Space) ProcessRemovals() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.toBeRemoved) == 0 {
		return
	}

	for _, b := range s.toBeRemoved {
		if b == nil {
			continue
		}
		delete(s.bodies, b.ID())
	}
	s.toBeRemoved = nil
	s.cacheDirty = true
}

// Bodies returns a slice of all collidable bodies in the space.
// This method uses a cache for performance. The returned slice is a direct
// reference to the cache and MUST NOT be modified by the caller. If modifications
// are needed, the caller should create a copy. The slice is sorted by body ID.
func (s *Space) Bodies() []body.Collidable {
	s.mu.RLock()
	if !s.cacheDirty {
		defer s.mu.RUnlock()
		return s.bodiesCache
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	// Re-check condition, as another goroutine could have updated the cache
	// between the RUnlock and Lock.
	if s.cacheDirty {
		s.bodiesCache = make([]body.Collidable, 0, len(s.bodies))
		for _, b := range s.bodies {
			if b == nil {
				continue
			}
			s.bodiesCache = append(s.bodiesCache, b)
		}

		// Sort the bodies by ID
		sort.Slice(s.bodiesCache, func(i, j int) bool {
			return s.bodiesCache[i].ID() < s.bodiesCache[j].ID()
		})

		s.cacheDirty = false
	}
	return s.bodiesCache
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

		if !HasCollision(body, other) {
			continue
		}

		body = s.Find(body.ID())
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

// Find return a body with the given ID.
func (s *Space) Find(id string) body.Collidable {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.bodies[id]
	if !ok {
		return nil
	}
	return b
}

// Query returns all bodies that overlap with the given rectangle.
func (s *Space) Query(rect image.Rectangle) []body.Collidable {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []body.Collidable

	for _, b := range s.bodies {
		if b == nil {
			continue
		}

		isOverlapping := false
		// Check the main body position
		if b.Position().Overlaps(rect) {
			isOverlapping = true
		} else {
			// Check all collision shapes of the body
			for _, collisionShape := range b.CollisionPosition() {
				if collisionShape.Overlaps(rect) {
					isOverlapping = true
					break
				}
			}
		}

		if isOverlapping {
			result = append(result, b)
		}
	}

	return result
}

// collisionRects returns the collision rectangles for a body.
// If no specific collision shapes are defined, it falls back to the body's own position.
// This is consistent with the pattern used in CheckGround.
func collisionRects(b body.Collidable) []image.Rectangle {
	rects := b.CollisionPosition()
	if len(rects) == 0 {
		return []image.Rectangle{b.Position()}
	}
	return rects
}

// HasCollision reports whether two bodies are currently colliding.
//
// The check runs in two phases:
//  1. 2D AABB overlap on the bodies' collision rectangles (screen-space).
//     This is the only check that runs for plain 2D bodies.
//  2. Depth-lane gate (opt-in, 2.5D): if BOTH a and b implement
//     DepthLaneBody, an additional constraint requires their ground-Y
//     coordinates to lie within max(a.LaneHalfWidth, b.LaneHalfWidth) of
//     each other. This filters bbox overlaps that LOOK on-screen but
//     happen at different floor depths in a 2.5D beat-em-up plane.
//
// Returns false for empty IDs or self-pair.
func HasCollision(a, b body.Collidable) bool {
	// Every body must have an ID
	if a.ID() == "" || b.ID() == "" {
		return false
	}

	// Prevent to check the same body
	if a.ID() == b.ID() {
		return false
	}

	rectsA := collisionRects(a)
	rectsB := collisionRects(b)

	// Phase 1 — 2D bbox overlap (screen-space AABB).
	// A miss here short-circuits; bodies that do not visually overlap can
	// never collide, regardless of depth.
	bboxOverlap := false
	for _, r := range rectsA {
		for _, s := range rectsB {
			if r.Overlaps(s) {
				bboxOverlap = true
				break
			}
		}
		if bboxOverlap {
			break
		}
	}

	if !bboxOverlap {
		return false
	}

	// Phase 2 — depth-lane gate (opt-in, 2.5D).
	// When BOTH bodies implement DepthLaneBody, the bbox overlap alone is
	// not enough: their ground projections must also lie within the same
	// depth lane. Bodies that do not opt in keep the legacy 2D behavior
	// (bbox overlap == collision), so plain 2D scenes are unaffected.
	da, okA := a.(DepthLaneBody)
	db, okB := b.(DepthLaneBody)
	if !okA || !okB {
		return true
	}

	// Pair tolerance is the LARGER of the two half-widths so a "wide" body
	// (e.g. a hitbox tuned to forgive depth misalignment) governs the match.
	tol := da.LaneHalfWidth()
	if db.LaneHalfWidth() > tol {
		tol = db.LaneHalfWidth()
	}
	// Compare ground-Y (depth axis, pre-altitude) — NOT screen-Y. This is
	// why DepthLaneBody exposes GroundY() explicitly rather than reusing
	// GetPosition16(), whose Y already has altitude subtracted.
	diff := da.GroundY() - db.GroundY()
	if diff < 0 {
		diff = -diff
	}
	return diff <= tol
}

func (s *Space) SetTilemapDimensionsProvider(provider tilemaplayer.TilemapDimensionsProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tilemapDimensionsProvider = provider
}

func (s *Space) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tilemapDimensionsProvider
}
