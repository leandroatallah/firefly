package movement

import (
	"image"
	"testing"
)

func TestNewWanderMovementState(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base)

	wanderState, ok := state.(*WanderMovementState)
	if !ok {
		t.Fatal("expected NewWanderMovementState to return *WanderMovementState")
	}

	if wanderState.maxDistance != 50 {
		t.Errorf("expected default maxDistance to be 50, got %d", wanderState.maxDistance)
	}
}

func TestWanderMovementState_OnStart(t *testing.T) {
	actor := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)

	state.OnStart()

	if state.anchorX != 100 {
		t.Errorf("expected anchorX to be set to actor's X position (100), got %d", state.anchorX)
	}
	if state.state != wanderIdle {
		t.Errorf("expected initial state to be wanderIdle, got %v", state.state)
	}
	if state.timer != 0 {
		t.Errorf("expected timer to be 0, got %d", state.timer)
	}
}

func TestWanderMovementState_Move_Immobile(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2, immobile: true}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.state = wanderMove
	state.movingRight = true

	state.Move(nil)

	if actor.moveLeftForce != 0 && actor.moveRightForce != 0 {
		t.Errorf("expected immobile actor to not move")
	}
}

func TestWanderMovementState_Move_Idle(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.state = wanderIdle
	state.timer = 0
	state.idleTime = 5 // Short idle time for testing

	// Should stay idle until timer expires
	for i := 0; i < 5; i++ {
		state.Move(nil)
		if state.state != wanderIdle {
			t.Errorf("expected to stay idle until timer expires, got state %v at iteration %d", state.state, i)
		}
	}

	// After timer expires, should transition to move state
	state.Move(nil)
	if state.state != wanderMove {
		t.Errorf("expected to transition to wanderMove after idle timer, got %v", state.state)
	}
}

func TestWanderMovementState_Move_Wander(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.state = wanderMove
	state.movingRight = true
	state.moveTime = 5
	state.timer = 0

	// Should move right while timer hasn't expired
	state.Move(nil)
	if actor.moveRightForce != 2 {
		t.Errorf("expected actor to move right with speed 2, got %d", actor.moveRightForce)
	}

	// Continue until timer expires
	for i := 0; i < 5; i++ {
		state.Move(nil)
	}

	// After timer expires, should transition to idle
	if state.state != wanderIdle {
		t.Errorf("expected to transition to wanderIdle after move timer, got %v", state.state)
	}
}

func TestWanderMovementState_Move_WanderLeft(t *testing.T) {
	actor := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 3}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.state = wanderMove
	state.movingRight = false
	state.moveTime = 10
	state.timer = 0

	state.Move(nil)

	if actor.moveLeftForce != 3 {
		t.Errorf("expected actor to move left with speed 3, got %d", actor.moveLeftForce)
	}
}

func TestWanderMovementState_ShouldStop_Ledge(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.movingRight = true

	// Create space with no ground ahead (ledge)
	space := newMockSpaceWithGround([]image.Point{})

	shouldStop := state.shouldStop(space)
	if !shouldStop {
		t.Errorf("expected shouldStop to return true at ledge")
	}
}

func TestWanderMovementState_ShouldStop_Ground(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.movingRight = true

	// Create space with ground ahead
	groundPositions := []image.Point{
		{60, 61}, // Ground at bottom-right + 1 pixel down
	}
	space := newMockSpaceWithGround(groundPositions)

	shouldStop := state.shouldStop(space)
	if shouldStop {
		t.Errorf("expected shouldStop to return false when ground exists")
	}
}

func TestWanderMovementState_ShouldStop_Wall(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.movingRight = true

	// Create space with wall ahead
	wallRect := image.Rect(60, 50, 61, 60) // Wall directly to the right
	space := newMockSpaceWithObstacles([]image.Rectangle{wallRect})

	shouldStop := state.shouldStop(space)
	if !shouldStop {
		t.Errorf("expected shouldStop to return true when wall is ahead")
	}
}

func TestWanderMovementState_ShouldStop_NilSpace(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.movingRight = true

	shouldStop := state.shouldStop(nil)
	if shouldStop {
		t.Errorf("expected shouldStop to return false when space is nil")
	}
}

func TestWanderMovementState_PickNextMove_BeyondMaxDistanceRight(t *testing.T) {
	actor := &mockActor{pos: image.Rect(160, 100, 170, 110), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.anchorX = 100
	state.maxDistance = 50

	state.pickNextMove()

	if state.movingRight {
		t.Errorf("expected to move left when beyond maxDistance to the right")
	}
}

func TestWanderMovementState_PickNextMove_BeyondMaxDistanceLeft(t *testing.T) {
	actor := &mockActor{pos: image.Rect(40, 100, 50, 110), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.anchorX = 100
	state.maxDistance = 50

	state.pickNextMove()

	if !state.movingRight {
		t.Errorf("expected to move right when beyond maxDistance to the left")
	}
}

func TestWanderMovementState_StartIdle(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.state = wanderMove
	state.timer = 100

	state.startIdle()

	if state.state != wanderIdle {
		t.Errorf("expected state to be wanderIdle, got %v", state.state)
	}
	if state.timer != 0 {
		t.Errorf("expected timer to be reset to 0, got %d", state.timer)
	}
	if state.idleTime < 60 || state.idleTime > 180 {
		t.Errorf("expected idleTime to be between 60 and 180, got %d", state.idleTime)
	}
}

func TestWanderMovementState_Move_ShouldStopTransition(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	base := NewBaseMovementState(Wander, actor, nil)
	state := NewWanderMovementState(base).(*WanderMovementState)
	state.state = wanderMove
	state.movingRight = true

	// Create space with ledge ahead
	space := newMockSpaceWithGround([]image.Point{})

	state.Move(space)

	// Should transition to idle when hitting ledge
	if state.state != wanderIdle {
		t.Errorf("expected to transition to idle when shouldStop returns true")
	}
}
