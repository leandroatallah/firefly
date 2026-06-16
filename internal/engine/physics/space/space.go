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
		// Check all collision shapes of the body
		for _, collisionShape := range b.CollisionPosition() {
			if collisionShape.Overlaps(rect) {
				isOverlapping = true
				break
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
// resolveDepthLane returns the DepthLaneBody for b, checking b itself first
// and then walking the owner chain via LastOwner. This is necessary because
// ApplyValidPosition operates on the inner *CollidableBody rather than the
// outer game-layer struct (e.g. BeatEmUpCharacter) that implements
// DepthLaneBody.
func resolveDepthLane(b body.Collidable) (DepthLaneBody, bool) {
	if dl, ok := b.(DepthLaneBody); ok {
		return dl, true
	}
	if owner := b.LastOwner(); owner != nil {
		if dl, ok := owner.(DepthLaneBody); ok {
			return dl, true
		}
	}
	return nil, false
}

// altitudeOf returns the body's current altitude in pixels, checking the body
// itself and then its owner chain. Returns 0 if the body carries no altitude.
func altitudeOf(b body.Collidable) int {
	type altGetter interface{ Altitude() int }
	if ag, ok := b.(altGetter); ok {
		return ag.Altitude()
	}
	if owner := b.LastOwner(); owner != nil {
		if ag, ok := owner.(altGetter); ok {
			return ag.Altitude()
		}
	}
	return 0
}

// floorProjectedRects returns the body's collision rects shifted down by its
// altitude so they represent floor-plane positions rather than screen positions.
// For grounded bodies (altitude == 0) this is a no-op.
func floorProjectedRects(b body.Collidable) []image.Rectangle {
	rects := collisionRects(b)
	alt := altitudeOf(b)
	if alt == 0 {
		return rects
	}
	shifted := make([]image.Rectangle, len(rects))
	for i, r := range rects {
		shifted[i] = r.Add(image.Pt(0, alt))
	}
	return shifted
}

func collisionRects(b body.Collidable) []image.Rectangle {
	rects := b.CollisionPosition()
	if len(rects) == 0 {
		return []image.Rectangle{b.Position()}
	}
	return rects
}

// rectSlicesOverlap reports whether any rect in ra overlaps any rect in rb.
func rectSlicesOverlap(ra, rb []image.Rectangle) bool {
	for _, r := range ra {
		for _, s := range rb {
			if r.Overlaps(s) {
				return true
			}
		}
	}
	return false
}

// HasCollision reports whether two bodies are currently colliding.
//
// For plain 2D bodies (neither implements DepthLaneBody): screen-space AABB
// overlap only.
//
// For 2.5D beat-em-up bodies (both implement DepthLaneBody, resolved through
// the owner chain):
//  1. Depth-lane gate: different floor lanes → no collision (early exit).
//     This prevents airborne bodies from colliding with background obstacles
//     at different floor depths.
//  2. Floor-projected AABB: rects are shifted down by each body's altitude
//     before the overlap test. This ensures same-depth obstacles still block
//     an airborne body even when it has visually cleared them on screen.
//
// Returns false for empty IDs or self-pair.
func HasCollision(a, b body.Collidable) bool {
	if a.ID() == "" || b.ID() == "" {
		return false
	}
	if a.ID() == b.ID() {
		return false
	}

	da, okA := resolveDepthLane(a)
	db, okB := resolveDepthLane(b)

	if okA && okB {
		// Phase 1 — depth-lane gate.
		tol := da.LaneHalfWidth()
		if db.LaneHalfWidth() > tol {
			tol = db.LaneHalfWidth()
		}
		diff := da.GroundY() - db.GroundY()
		if diff < 0 {
			diff = -diff
		}
		if diff > tol {
			return false
		}
		// Phase 2 — floor-projected AABB.
		return rectSlicesOverlap(floorProjectedRects(a), floorProjectedRects(b))
	}

	// Legacy 2D: screen-space AABB only.
	return rectSlicesOverlap(collisionRects(a), collisionRects(b))
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
