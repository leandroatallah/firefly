package body

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

// mockBodiesSpace implements body.BodiesSpace for testing
type mockBodiesSpace struct {
	bodies []body.Collidable
}

func newMockBodiesSpace() *mockBodiesSpace {
	return &mockBodiesSpace{bodies: make([]body.Collidable, 0)}
}

func (m *mockBodiesSpace) AddBody(b body.Collidable) {
	m.bodies = append(m.bodies, b)
}

func (m *mockBodiesSpace) Bodies() []body.Collidable {
	return m.bodies
}

func (m *mockBodiesSpace) RemoveBody(b body.Collidable) {
	for i, body := range m.bodies {
		if body == b {
			m.bodies = append(m.bodies[:i], m.bodies[i+1:]...)
			break
		}
	}
}

func (m *mockBodiesSpace) QueueForRemoval(b body.Collidable) {
	// No-op for testing
}

func (m *mockBodiesSpace) ProcessRemovals() {
	// No-op for testing
}

func (m *mockBodiesSpace) Clear() {
	m.bodies = make([]body.Collidable, 0)
}

func (m *mockBodiesSpace) ResolveCollisions(b body.Collidable) (touching bool, blocking bool) {
	pos := b.Position()
	for _, other := range m.bodies {
		if other == b {
			continue
		}
		otherPos := other.Position()
		if pos.Overlaps(otherPos) {
			if other.IsObstructive() {
				return true, true
			}
			return true, false
		}
	}
	return false, false
}

func (m *mockBodiesSpace) SetTilemapDimensionsProvider(p tilemaplayer.TilemapDimensionsProvider) {
	// No-op for testing
}

func (m *mockBodiesSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return nil
}

func (m *mockBodiesSpace) Find(id string) body.Collidable {
	for _, b := range m.bodies {
		if b.ID() == id {
			return b
		}
	}
	return nil
}

func (m *mockBodiesSpace) Query(rect image.Rectangle) []body.Collidable {
	var result []body.Collidable
	for _, b := range m.bodies {
		if b.Position().Overlaps(rect) {
			result = append(result, b)
		}
	}
	return result
}

func TestNewCollidableBody(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))
	cb := NewCollidableBody(b)

	if cb == nil {
		t.Fatal("NewCollidableBody returned nil")
	}
	if cb.Body != b {
		t.Errorf("expected Body to be set; got %v", cb.Body)
	}
}

func TestNewCollidableBody_NilBody(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewCollidableBody did not panic with nil body")
		}
	}()
	NewCollidableBody(nil)
}

func TestNewCollidableBodyFromRect(t *testing.T) {
	rect := NewRect(5, 10, 20, 30)
	cb := NewCollidableBodyFromRect(rect)

	if cb == nil {
		t.Fatal("NewCollidableBodyFromRect returned nil")
	}
	if cb.Body == nil {
		t.Error("Body was not initialized")
	}

	pos := cb.Position()
	if pos.Min.X != 0 || pos.Min.Y != 0 {
		t.Errorf("expected position (0,0); got (%d,%d)", pos.Min.X, pos.Min.Y)
	}
}

func TestCollidableBody_SetTouchable(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	trigger := NewTouchTrigger(func() {}, nil)

	cb.SetTouchable(trigger)
	got := cb.GetTouchable()
	if got != trigger {
		t.Errorf("expected touchable %v; got %v", trigger, got)
	}
}

func TestCollidableBody_GetTouchable(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))

	got := cb.GetTouchable()
	if got != nil {
		t.Errorf("expected nil touchable by default; got %v", got)
	}
}

func TestCollidableBody_AddCollision(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test-body")

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))

	cb.AddCollision(other)

	shapes := cb.CollisionShapes()
	if len(shapes) != 1 {
		t.Errorf("expected 1 collision shape; got %d", len(shapes))
	}

	expectedID := "test-body_COLLISION_0"
	if other.ID() != expectedID {
		t.Errorf("expected collision ID '%s'; got '%s'", expectedID, other.ID())
	}
}

func TestCollidableBody_AddCollision_Multiple(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("parent")

	for i := 0; i < 3; i++ {
		other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
		cb.AddCollision(other)
	}

	shapes := cb.CollisionShapes()
	if len(shapes) != 3 {
		t.Errorf("expected 3 collision shapes; got %d", len(shapes))
	}
}

func TestCollidableBody_AddCollision_NoID(t *testing.T) {
	// AddCollision calls log.Fatal when body has no ID
	// We skip this test as log.Fatal calls os.Exit which is hard to test
	// The behavior is documented: AddCollision requires an ID
	t.Skip("AddCollision calls log.Fatal which exits the process")
}

func TestCollidableBody_ClearCollisions(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	cb.AddCollision(NewCollidableBodyFromRect(NewRect(0, 0, 5, 5)))
	cb.AddCollision(NewCollidableBodyFromRect(NewRect(0, 0, 5, 5)))

	cb.ClearCollisions()

	shapes := cb.CollisionShapes()
	if len(shapes) != 0 {
		t.Errorf("expected 0 collision shapes after clear; got %d", len(shapes))
	}
}

func TestCollidableBody_SetIsObstructive(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))

	if cb.IsObstructive() {
		t.Error("expected obstructive to be false by default")
	}

	cb.SetIsObstructive(true)
	if !cb.IsObstructive() {
		t.Error("expected obstructive to be true after setting")
	}
}

func TestCollidableBody_CollisionPosition(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetPosition(100, 200)
	cb.AddCollision(other)

	positions := cb.CollisionPosition()
	if len(positions) != 1 {
		t.Fatalf("expected 1 collision position; got %d", len(positions))
	}

	pos := positions[0]
	if pos.Min.X != 100 || pos.Min.Y != 200 {
		t.Errorf("expected position (100,200); got (%d,%d)", pos.Min.X, pos.Min.Y)
	}
}

func TestCollidableBody_SetPosition(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetPosition(0, 0)
	cb.AddCollision(other)

	cb.SetPosition(50, 60)

	x, y := cb.GetPositionMin()
	if x != 50 || y != 60 {
		t.Errorf("expected body position (50,60); got (%d,%d)", x, y)
	}

	// Collision should move with body
	otherX, otherY := other.GetPositionMin()
	if otherX != 50 || otherY != 60 {
		t.Errorf("expected collision position (50,60); got (%d,%d)", otherX, otherY)
	}
}

func TestCollidableBody_SetPosition16(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	cb.AddCollision(other)

	x16 := fp16.To16(100)
	y16 := fp16.To16(200)
	cb.SetPosition16(x16, y16)

	gotX16, gotY16 := cb.GetPosition16()
	if gotX16 != x16 || gotY16 != y16 {
		t.Errorf("expected (%d,%d); got (%d,%d)", x16, y16, gotX16, gotY16)
	}
}

func TestCollidableBody_OnTouch(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	touched := false
	trigger := NewTouchTrigger(func() { touched = true }, cb)
	cb.SetTouchable(trigger)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("test") // Must match subject ID
	cb.OnTouch(other)

	if !touched {
		t.Error("expected OnTouch to trigger callback")
	}
}

func TestCollidableBody_OnBlock(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	blocked := false
	trigger := NewTouchTrigger(func() { blocked = true }, cb)
	cb.SetTouchable(trigger)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("test")
	
	// OnBlock is intentionally a no-op for TouchTrigger
	// This test verifies that behavior
	cb.OnBlock(other)

	if blocked {
		t.Error("expected OnBlock to NOT trigger callback (no-op for TouchTrigger)")
	}
}

func TestCollidableBody_DrawCollisionBox(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("test")

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	cb.AddCollision(other)

	screen := ebiten.NewImage(100, 100)
	position := image.Rect(0, 0, 100, 100)

	// Should not panic
	cb.DrawCollisionBox(screen, position)
}

func TestCollidableBody_ApplyValidPosition_NoMovement(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	space := newMockBodiesSpace()

	x, y, blocked := cb.ApplyValidPosition(0, true, space)

	if x != 0 || y != 0 {
		t.Errorf("expected position (0,0); got (%d,%d)", x, y)
	}
	if blocked {
		t.Error("expected not blocked")
	}
}

func TestCollidableBody_ApplyValidPosition_NilSpace(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))

	x, y, blocked := cb.ApplyValidPosition(10, true, nil)

	if x != 0 || y != 0 {
		t.Errorf("expected position (0,0); got (%d,%d)", x, y)
	}
	if blocked {
		t.Error("expected not blocked")
	}
}

func TestCollidableBody_ApplyValidPosition_XAxis(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	space := newMockBodiesSpace()

	x, y, blocked := cb.ApplyValidPosition(fp16.To16(5), true, space)

	if x != 5 || y != 0 {
		t.Errorf("expected position (5,0); got (%d,%d)", x, y)
	}
	if blocked {
		t.Error("expected not blocked")
	}
}

func TestCollidableBody_ApplyValidPosition_YAxis(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	space := newMockBodiesSpace()

	x, y, blocked := cb.ApplyValidPosition(fp16.To16(10), false, space)

	if x != 0 || y != 10 {
		t.Errorf("expected position (0,10); got (%d,%d)", x, y)
	}
	if blocked {
		t.Error("expected not blocked")
	}
}

func TestCollidableBody_ApplyValidPosition_Blocked(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("moving")

	// Create an obstacle
	obstacle := NewCollidableBodyFromRect(NewRect(5, 0, 10, 10))
	obstacle.SetID("obstacle")
	obstacle.SetIsObstructive(true)

	space := newMockBodiesSpace()
	space.AddBody(obstacle)
	space.AddBody(cb)

	x, _, blocked := cb.ApplyValidPosition(fp16.To16(10), true, space)

	// Should stop before or at collision
	if x >= 5 {
		t.Errorf("expected x < 5 (blocked); got %d", x)
	}
	if !blocked {
		t.Error("expected to be blocked")
	}
}

func TestCollidableBody_ApplyValidPosition_Negative(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("moving")
	// Set initial position to (20, 0)
	cb.SetPosition(20, 0)
	
	space := newMockBodiesSpace()
	space.AddBody(cb)

	// Move left by 5 pixels (negative direction)
	// fp16 scale is 16, so 5 pixels = 5 * 16 = 80
	x, y, blocked := cb.ApplyValidPosition(-fp16.To16(5), true, space)

	if x != 15 || y != 0 {
		t.Errorf("expected position (15,0); got (%d,%d)", x, y)
	}
	if blocked {
		t.Error("expected not blocked")
	}
}

func TestCollidableBody_ApplyValidPosition_SubPixel(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	space := newMockBodiesSpace()
	space.AddBody(cb)

	// Move by less than 1 pixel (0.5 pixels = 8 in fp16 with scale 16)
	x, y, _ := cb.ApplyValidPosition(8, true, space)

	// Should not move yet (accumulates)
	if x != 0 || y != 0 {
		t.Errorf("expected position (0,0) for sub-pixel; got (%d,%d)", x, y)
	}

	// Move again to accumulate (another 0.5 pixels)
	x, y, _ = cb.ApplyValidPosition(8, true, space)

	// Now should have moved 1 pixel (0.5 + 0.5 = 1)
	if x != 1 || y != 0 {
		t.Errorf("expected position (1,0) after accumulation; got (%d,%d)", x, y)
	}
}

func TestCollidableBody_ApplyValidPosition_ClearsAccumulatorOnBlock(t *testing.T) {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	cb.SetID("moving")

	obstacle := NewCollidableBodyFromRect(NewRect(3, 0, 10, 10))
	obstacle.SetID("obstacle")
	obstacle.SetIsObstructive(true)

	space := newMockBodiesSpace()
	space.AddBody(obstacle)
	space.AddBody(cb)

	// Try to move 5 pixels but blocked at 3
	cb.ApplyValidPosition(fp16.To16(5), true, space)

	// Accumulator should be cleared
	if cb.accumulatorX16 != 0 {
		t.Errorf("expected accumulatorX16 to be 0 after block; got %d", cb.accumulatorX16)
	}
}

func TestCollidableBody_SetPositionRequiresID(t *testing.T) {
	// SetPosition calls log.Fatal when body has no ID
	// We skip this test as log.Fatal calls os.Exit which is hard to test
	// The behavior is documented: SetPosition requires an ID
	t.Skip("SetPosition calls log.Fatal which exits the process")
}

func TestAbs(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for _, tt := range tests {
		if got := abs(tt.n); got != tt.want {
			t.Errorf("abs(%d) = %d; want %d", tt.n, got, tt.want)
		}
	}
}

func TestCollidableBody_AddCollision_NilShape(t *testing.T) {
	// AddCollision calls log.Fatal when body has no shape
	// We skip this test as log.Fatal calls os.Exit which is hard to test
	t.Skip("AddCollision calls log.Fatal when shape is nil")
}

func TestCollidableBody_AddCollision_EmptyID(t *testing.T) {
	// AddCollision calls log.Fatal when body has no ID
	// We skip this test as log.Fatal calls os.Exit which is hard to test
	t.Skip("AddCollision calls log.Fatal when ID is empty")
}
