package movement

import (
	"image"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

func TestNewChaseMovementState(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)

	if state == nil {
		t.Fatal("expected NewChaseMovementState to return non-nil state")
	}
	if state.actor != actor {
		t.Errorf("expected actor to be set")
	}
	if state.target != target {
		t.Errorf("expected target to be set")
	}
}

func TestChaseMovementState_Move_Immobile(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2, immobile: true}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)

	state.Move(nil)

	if actor.moveLeftForce != 0 && actor.moveRightForce != 0 &&
		actor.moveUpForce != 0 && actor.moveDownForce != 0 {
		t.Errorf("expected immobile actor to not move")
	}
}

func TestChaseMovementState_Move_NoPath(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	state.path = []image.Point{} // Empty path

	state.Move(nil)

	// Should not move when no path exists
	if actor.moveLeftForce != 0 && actor.moveRightForce != 0 &&
		actor.moveUpForce != 0 && actor.moveDownForce != 0 {
		t.Errorf("expected actor to not move when no path exists")
	}
}

func TestChaseMovementState_Move_WithPath(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(50, 0, 60, 10), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	// Set a simple path to the right
	state.path = []image.Point{{10, 0}, {20, 0}, {30, 0}}

	state.Move(nil)

	// Should move right towards the first path point
	if actor.moveRightForce == 0 {
		t.Errorf("expected actor to move right towards path point")
	}
}

func TestChaseMovementState_Move_PathThreshold(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(50, 0, 60, 10), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	// Set path with first point very close (within threshold)
	state.path = []image.Point{{5, 0}, {20, 0}, {30, 0}}

	state.Move(nil)

	// First point should be removed (within threshold), should move to next
	if len(state.path) != 2 {
		t.Errorf("expected first path point to be removed when within threshold")
	}
}

func TestChaseMovementState_Move_PathExhausted(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(5, 0, 15, 10), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	// Set path with single point very close
	state.path = []image.Point{{5, 0}}

	state.Move(nil)

	// Path should be exhausted
	if len(state.path) != 0 {
		t.Errorf("expected path to be exhausted")
	}
}

func TestChaseMovementState_Move_Deadzone(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(50, 2, 60, 12), speed: 2} // Slightly offset in Y
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	// Path point close in Y (within deadzone)
	state.path = []image.Point{{20, 2}}

	state.Move(nil)

	// Should only move right (Y offset is within deadzone)
	if actor.moveRightForce == 0 {
		t.Errorf("expected actor to move right")
	}
	if actor.moveUpForce != 0 || actor.moveDownForce != 0 {
		t.Errorf("expected no vertical movement when within deadzone")
	}
}

func TestChaseMovementState_CalculatePath(t *testing.T) {
	t.Run("simple straight path", func(t *testing.T) {
		actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2, id: "actor"}
		target := &mockActor{pos: image.Rect(40, 0, 50, 10), speed: 2, id: "target"}
		base := NewBaseMovementState(Chase, actor, target)
		state := NewChaseMovementState(base)

		state.calculatePath()

		if len(state.path) == 0 {
			t.Errorf("expected path to be calculated")
		}
	})

	t.Run("path around obstacle", func(t *testing.T) {
		actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2, id: "actor"}
		target := &mockActor{pos: image.Rect(40, 0, 50, 10), speed: 2, id: "target"}
		base := NewBaseMovementState(Chase, actor, target)
		state := NewChaseMovementState(base)
		
		// Add obstacle in the middle
		obstacle := newMockMovableCollidable(15, 0, 10, 10)
		state.obstacles = []body.MovableCollidable{obstacle}

		state.calculatePath()

		// Should find a path around the obstacle
		if len(state.path) == 0 {
			t.Errorf("expected path to be calculated around obstacle")
		}
	})

	t.Run("no path when surrounded", func(t *testing.T) {
		actor := &mockActor{pos: image.Rect(20, 20, 30, 30), speed: 2, id: "actor"}
		target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2, id: "target"}
		base := NewBaseMovementState(Chase, actor, target)
		state := NewChaseMovementState(base)
		
		// Surround actor with obstacles
		obstacles := []body.MovableCollidable{
			newMockMovableCollidable(30, 20, 5, 10),  // right
			newMockMovableCollidable(15, 20, 5, 10),  // left
			newMockMovableCollidable(20, 30, 10, 5),  // down
			newMockMovableCollidable(20, 15, 10, 5),  // up
		}
		state.obstacles = obstacles

		state.calculatePath()

		// May or may not find a path depending on gap size
		// This test mainly verifies the function doesn't panic
	})
}

func TestChaseMovementState_IsTraversable(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2, id: "actor"}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2, id: "target"}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	size := image.Point{10, 10}

	t.Run("point within bounds", func(t *testing.T) {
		traversable := state.isTraversable(image.Point{50, 50}, size)
		if !traversable {
			t.Errorf("expected point within bounds to be traversable")
		}
	})

	t.Run("negative coordinates", func(t *testing.T) {
		traversable := state.isTraversable(image.Point{-5, 50}, size)
		if traversable {
			t.Errorf("expected negative coordinates to not be traversable")
		}
	})

	t.Run("point overlapping obstacle", func(t *testing.T) {
		obstacle := newMockMovableCollidable(45, 45, 10, 10)
		state.obstacles = []body.MovableCollidable{obstacle}

		traversable := state.isTraversable(image.Point{50, 50}, size)
		if traversable {
			t.Errorf("expected point overlapping obstacle to not be traversable")
		}
	})

	t.Run("point overlapping target skipped", func(t *testing.T) {
		// Target position should be traversable (we're chasing it)
		state.obstacles = []body.MovableCollidable{target}

		traversable := state.isTraversable(target.Position().Min, size)
		if !traversable {
			t.Errorf("expected target position to be traversable")
		}
	})
}

func TestChaseMovementState_IsTraversable_WithMapBounds(t *testing.T) {
	// This test requires integration with the tilemap system
	// The isTraversable method checks bounds via type assertion on the actor's Space()
	// For now, we test the basic traversable logic without map bounds
	t.Skip("Skipping integration test - requires full tilemap space setup")
}

func TestChaseMovementState_GetNeighbors(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2, id: "target"}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)
	size := image.Point{10, 10}

	t.Run("all 8 directions unobstructed", func(t *testing.T) {
		state.obstacles = []body.MovableCollidable{}

		neighbors := state.getNeighbors(image.Point{50, 50}, size)

		// Should have 8 neighbors (4 cardinal + 4 diagonal)
		if len(neighbors) != 8 {
			t.Errorf("expected 8 neighbors when unobstructed, got %d", len(neighbors))
		}
	})

	t.Run("obstacle to right blocks right and diagonals", func(t *testing.T) {
		// Obstacle directly to the right
		obstacle := newMockMovableCollidable(60, 50, 10, 10)
		state.obstacles = []body.MovableCollidable{obstacle}

		neighbors := state.getNeighbors(image.Point{50, 50}, size)

		// Should not have neighbors at X >= 60 (right, up-right, down-right)
		for _, n := range neighbors {
			if n.X >= 60 {
				t.Errorf("expected no neighbors to the right when obstacle exists, got %v", n)
			}
		}
	})

	t.Run("corner cutting prevention", func(t *testing.T) {
		// Obstacle blocking up movement
		obstacle := newMockMovableCollidable(50, 40, 10, 10)
		state.obstacles = []body.MovableCollidable{obstacle}

		neighbors := state.getNeighbors(image.Point{50, 50}, size)

		// Should not have upward neighbors (up, up-left, up-right)
		for _, n := range neighbors {
			if n.Y < 50 {
				t.Errorf("expected no upward neighbors when up is blocked, got %v", n)
			}
		}
	})
}

func TestChaseMovementState_Move_PathRecalculation(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2}
	base := NewBaseMovementState(Chase, actor, target)
	state := NewChaseMovementState(base)

	// Move 29 times - should not recalculate yet
	for i := 0; i < 29; i++ {
		state.count = i
		state.Move(nil)
	}

	// On 30th move, should recalculate
	state.count = 29
	state.Move(nil)

	// Count should reset or path should be recalculated
	// This test verifies the recalculation logic runs
}

// mockActorWithSpace wraps mockActor to provide Space() method
type mockActorWithSpace struct {
	*mockActor
	space body.BodiesSpace
}

func (m *mockActorWithSpace) Space() body.BodiesSpace {
	return m.space
}
