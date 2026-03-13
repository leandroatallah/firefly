package movement

import (
	"image"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/contracts/tilemaplayer"
)

type mockSpace struct {
	bodies          []body.Collidable
	tilemapProvider tilemaplayer.TilemapDimensionsProvider
}

func (m *mockSpace) Query(rect image.Rectangle) []body.Collidable {
	var result []body.Collidable
	for _, b := range m.bodies {
		if b.Position().Overlaps(rect) {
			result = append(result, b)
		}
	}
	return result
}

func (m *mockSpace) AddBody(b body.Collidable) {}
func (m *mockSpace) RemoveBody(b body.Collidable) {}
func (m *mockSpace) Bodies() []body.Collidable { return m.bodies }
func (m *mockSpace) ResolveCollisions(b body.Collidable) (bool, bool) { return false, false }
func (m *mockSpace) SetTilemapDimensionsProvider(p tilemaplayer.TilemapDimensionsProvider) {
	m.tilemapProvider = p
}
func (m *mockSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return m.tilemapProvider
}
func (m *mockSpace) QueueForRemoval(body body.Collidable) {}
func (m *mockSpace) ProcessRemovals()                     {}
func (m *mockSpace) Clear()                               {}
func (m *mockSpace) Find(id string) body.Collidable       { return nil }

func TestSideToSideMovementState_WallDetection(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0) // 10x10 size by default in mock

	wall := &mockActor{}
	wall.pos = image.Rect(10, 0, 20, 10)

	space := &mockSpace{bodies: []body.Collidable{wall}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)

	state.Move(space)

	if actor.moveLeftForce != 5 {
		t.Errorf("expected moveLeftForce 5 after hitting wall, got %d", actor.moveLeftForce)
	}
}

func TestSideToSideMovementState_LedgeDetection(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)

	// Actor at edge of platform, no ground ahead
	space := &mockSpace{bodies: []body.Collidable{}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)

	state.Move(space)

	// Should turn around at ledge
	if actor.moveLeftForce != 5 {
		t.Errorf("expected moveLeftForce 5 at ledge, got %d", actor.moveLeftForce)
	}
}

func TestSideToSideMovementState_Immobile(t *testing.T) {
	actor := &mockActor{speed: 5, immobile: true}
	actor.SetPosition(0, 0)

	space := &mockSpace{bodies: []body.Collidable{}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)

	state.Move(space)

	if actor.moveRightForce != 0 || actor.moveLeftForce != 0 {
		t.Error("expected no movement when immobile")
	}
}

func TestSideToSideMovementState_WaitBeforeTurn(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)

	// Wall ahead
	wall := &mockActor{}
	wall.pos = image.Rect(10, 0, 20, 10)

	space := &mockSpace{bodies: []body.Collidable{wall}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)
	state.waitDuration = 10

	// First move - should hit wall and start waiting
	state.Move(space)

	if !state.isWaiting {
		t.Error("expected state to be waiting")
	}

	// During wait, should not change direction
	for i := 0; i < 9; i++ {
		state.Move(space)
		if !state.isWaiting {
			t.Errorf("expected state to still be waiting at iteration %d", i)
		}
	}

	// After wait completes, should turn (movingPositive flips from true to false)
	state.Move(space)

	if state.isWaiting {
		t.Error("expected waiting to be complete")
	}
	if state.movingPositive {
		t.Error("expected movingPositive to be false after wait (should have turned)")
	}
}

func TestSideToSideMovementState_Vertical(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)

	space := &mockSpace{bodies: []body.Collidable{}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)
	state.vertical = true

	state.Move(space)

	// Vertical movement uses SetVelocity, not OnMoveDown
	vx, vy := actor.Velocity()
	if vy <= 0 {
		t.Errorf("expected positive velocity Y for vertical movement, got %d", vy)
	}
	if vx != 0 {
		t.Errorf("expected velocity X to be 0, got %d", vx)
	}
}

func TestSideToSideMovementState_IgnoreLedges(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)

	// No ground ahead
	space := &mockSpace{bodies: []body.Collidable{}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)
	state.ignoreLedges = true

	state.Move(space)

	// Should continue moving right despite no ground
	if actor.moveRightForce != 5 {
		t.Errorf("expected moveRightForce 5 when ignoring ledges, got %d", actor.moveRightForce)
	}
}

func TestSideToSideMovementState_LimitToRoom(t *testing.T) {
	actor := &mockActor{speed: 5}
	// Position at right edge of room
	actor.SetPosition(310, 0)

	space := &mockSpace{bodies: []body.Collidable{}}
	space.tilemapProvider = &mockTilemapProvider{
		bounds: image.Rect(0, 0, 320, 240),
		hasBounds: true,
	}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)
	state.limitToRoom = true

	state.Move(space)

	// Should turn around at room limit
	if actor.moveLeftForce != 5 {
		t.Errorf("expected moveLeftForce 5 at room limit, got %d", actor.moveLeftForce)
	}
}

func TestSideToSideMovementState_LimitToRoomVertical(t *testing.T) {
	actor := &mockActor{speed: 5}
	// Position at bottom edge of room
	actor.SetPosition(0, 230)

	space := &mockSpace{bodies: []body.Collidable{}}
	space.tilemapProvider = &mockTilemapProvider{
		bounds: image.Rect(0, 0, 320, 240),
		hasBounds: true,
	}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)
	state.vertical = true
	state.limitToRoom = true

	state.Move(space)

	// After hitting room limit, should reverse direction
	// When moving negative vertically, it uses TryJump
	if actor.moveRightForce <= 0 {
		t.Errorf("expected TryJump to be called after hitting room limit, got %d", actor.moveRightForce)
	}
}

func TestSideToSideMovementState_NilSpace(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)

	// Should not panic with nil space
	state.Move(nil)

	// Should continue moving right
	if actor.moveRightForce != 5 {
		t.Errorf("expected moveRightForce 5 with nil space, got %d", actor.moveRightForce)
	}
}

func TestSideToSideMovementState_TurnAfterWall(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)

	wall := &mockActor{}
	wall.pos = image.Rect(10, 0, 20, 10)

	space := &mockSpace{bodies: []body.Collidable{wall}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)

	// First move - hit wall, turn around
	state.Move(space)

	// Now moving left, remove wall
	space.bodies = []body.Collidable{}

	// Move left until past starting position
	actor.SetPosition(-20, 0)

	state.Move(space)

	// Should continue moving left
	if actor.moveLeftForce != 5 {
		t.Errorf("expected moveLeftForce 5 after turning, got %d", actor.moveLeftForce)
	}
}

