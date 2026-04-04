package movement

import (
	"image"
	"testing"
)

func TestNewFollowMovementState(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Follow, actor, nil)
	state := NewFollowMovementState(base)

	if state.startDistance != 50 {
		t.Errorf("expected default startDistance to be 50, got %d", state.startDistance)
	}
	if state.stopDistance != 20 {
		t.Errorf("expected default stopDistance to be 20, got %d", state.stopDistance)
	}
}

func TestWithFollowDistances(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Follow, actor, nil)
	state := NewFollowMovementState(base)

	option := WithFollowDistances(100, 30)
	option(state)

	if state.startDistance != 100 {
		t.Errorf("expected startDistance to be 100, got %d", state.startDistance)
	}
	if state.stopDistance != 30 {
		t.Errorf("expected stopDistance to be 30, got %d", state.stopDistance)
	}
}

func TestWithFollowDistances_NotFollowState(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Patrol, actor, nil)
	state := NewPatrolMovementState(base)

	option := WithFollowDistances(100, 30)
	// Should not panic when applied to non-FollowMovementState
	option(state)
}

func TestWithPlatformFollow(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Follow, actor, nil)
	state := NewFollowMovementState(base)

	option := WithPlatformFollow()
	option(state)

	if !state.stayOnPlatform {
		t.Errorf("expected stayOnPlatform to be true")
	}
}

func TestWithPlatformFollow_NotFollowState(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Patrol, actor, nil)
	state := NewPatrolMovementState(base)

	option := WithPlatformFollow()
	// Should not panic when applied to non-FollowMovementState
	option(state)
}

func TestFollowMovementState_Move_Immobile(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2, immobile: true}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110), speed: 2}
	base := NewBaseMovementState(Follow, actor, target)
	state := NewFollowMovementState(base)

	state.Move(nil)

	if actor.moveLeftForce != 0 && actor.moveRightForce != 0 &&
		actor.moveUpForce != 0 && actor.moveDownForce != 0 {
		t.Errorf("expected immobile actor to not move")
	}
}

func TestFollowMovementState_Move_NoTarget(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	base := NewBaseMovementState(Follow, actor, nil)
	state := NewFollowMovementState(base)

	state.Move(nil)

	if actor.moveLeftForce != 0 && actor.moveRightForce != 0 &&
		actor.moveUpForce != 0 && actor.moveDownForce != 0 {
		t.Errorf("expected actor with no target to not move")
	}
}

func TestFollowMovementState_Move_StartFollowing(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(100, 0, 110, 10), speed: 2} // 100 pixels away
	base := NewBaseMovementState(Follow, actor, target)
	state := NewFollowMovementState(base)
	state.startDistance = 50
	state.stopDistance = 20

	state.Move(nil)

	// Should start moving since distance (100) > startDistance (50)
	if !state.isMoving {
		t.Errorf("expected to start moving when distance > startDistance")
	}
}

func TestFollowMovementState_Move_StopFollowing(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(10, 0, 20, 10), speed: 2} // 10 pixels away
	base := NewBaseMovementState(Follow, actor, target)
	state := NewFollowMovementState(base)
	state.startDistance = 50
	state.stopDistance = 20
	state.isMoving = true

	state.Move(nil)

	// Should stop moving since distance (10) < stopDistance (20)
	if state.isMoving {
		t.Errorf("expected to stop moving when distance < stopDistance")
	}
}

func TestFollowMovementState_Move_Hysteresis(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(30, 0, 40, 10), speed: 2} // 30 pixels away
	base := NewBaseMovementState(Follow, actor, target)
	state := NewFollowMovementState(base)
	state.startDistance = 50
	state.stopDistance = 20

	// Start not moving
	state.isMoving = false
	state.Move(nil)

	// Distance (30) is between stopDistance (20) and startDistance (50)
	// Should remain not moving (hysteresis)
	if state.isMoving {
		t.Errorf("expected to remain not moving when distance is between stop and start")
	}

	// Now set to moving
	state.isMoving = true
	state.Move(nil)

	// Should remain moving (hysteresis)
	if !state.isMoving {
		t.Errorf("expected to remain moving when distance is between stop and start")
	}
}

func TestFollowMovementState_Move_Directions(t *testing.T) {
	tests := []struct {
		name      string
		actorPos  image.Rectangle
		targetPos image.Rectangle
		wantLeft  bool
		wantRight bool
		wantUp    bool
		wantDown  bool
	}{
		{
			name:      "target to right",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(100, 0, 110, 10),
			wantRight: true,
		},
		{
			name:      "target to left",
			actorPos:  image.Rect(100, 0, 110, 10),
			targetPos: image.Rect(0, 0, 10, 10),
			wantLeft:  true,
		},
		{
			name:      "target below",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(0, 100, 10, 110),
			wantDown:  true,
		},
		{
			name:      "target above",
			actorPos:  image.Rect(0, 100, 10, 110),
			targetPos: image.Rect(0, 0, 10, 10),
			wantUp:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := &mockActor{pos: tt.actorPos, speed: 2}
			target := &mockActor{pos: tt.targetPos, speed: 2}
			base := NewBaseMovementState(Follow, actor, target)
			state := NewFollowMovementState(base)
			state.startDistance = 10
			state.isMoving = true

			state.Move(nil)

			if tt.wantLeft && actor.moveLeftForce == 0 {
				t.Errorf("expected actor to move left")
			}
			if tt.wantRight && actor.moveRightForce == 0 {
				t.Errorf("expected actor to move right")
			}
			if tt.wantUp && actor.moveUpForce == 0 {
				t.Errorf("expected actor to move up")
			}
			if tt.wantDown && actor.moveDownForce == 0 {
				t.Errorf("expected actor to move down")
			}
		})
	}
}

func TestFollowMovementState_FilterSafeDirections(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	target := &mockActor{pos: image.Rect(100, 50, 110, 60), speed: 2}
	base := NewBaseMovementState(Follow, actor, target)
	state := NewFollowMovementState(base)

	t.Run("no ground left - remove left", func(t *testing.T) {
		space := newMockSpaceWithGround([]image.Point{}) // No ground
		directions := MovementDirections{Left: true, Right: false}

		safe := state.filterSafeDirections(directions, space)

		if safe.Left {
			t.Errorf("expected Left to be filtered out when no ground")
		}
	})

	t.Run("no ground right - remove right", func(t *testing.T) {
		space := newMockSpaceWithGround([]image.Point{}) // No ground
		directions := MovementDirections{Left: false, Right: true}

		safe := state.filterSafeDirections(directions, space)

		if safe.Right {
			t.Errorf("expected Right to be filtered out when no ground")
		}
	})

	t.Run("ground exists - keep directions", func(t *testing.T) {
		// Ground at both left and right check points
		groundPositions := []image.Point{
			{49, 61}, // Left ground check
			{60, 61}, // Right ground check
		}
		space := newMockSpaceWithGround(groundPositions)
		directions := MovementDirections{Left: true, Right: true}

		safe := state.filterSafeDirections(directions, space)

		if !safe.Left || !safe.Right {
			t.Errorf("expected both directions to be kept when ground exists")
		}
	})
}

func TestFollowMovementState_HasGroundInDirection(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	base := NewBaseMovementState(Follow, actor, nil)
	state := NewFollowMovementState(base)

	t.Run("ground to right", func(t *testing.T) {
		groundPositions := []image.Point{{60, 61}}
		space := newMockSpaceWithGround(groundPositions)

		hasGround := state.hasGroundInDirection(space, true)

		if !hasGround {
			t.Errorf("expected ground to exist to the right")
		}
	})

	t.Run("no ground to right", func(t *testing.T) {
		space := newMockSpaceWithGround([]image.Point{})

		hasGround := state.hasGroundInDirection(space, true)

		if hasGround {
			t.Errorf("expected no ground to the right")
		}
	})

	t.Run("ground to left", func(t *testing.T) {
		groundPositions := []image.Point{{49, 61}}
		space := newMockSpaceWithGround(groundPositions)

		hasGround := state.hasGroundInDirection(space, false)

		if !hasGround {
			t.Errorf("expected ground to exist to the left")
		}
	})

	t.Run("no ground to left", func(t *testing.T) {
		space := newMockSpaceWithGround([]image.Point{})

		hasGround := state.hasGroundInDirection(space, false)

		if hasGround {
			t.Errorf("expected no ground to the left")
		}
	})
}

func TestFollowMovementState_Move_WithPlatformFollow(t *testing.T) {
	actor := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2, id: "actor"}
	target := &mockActor{pos: image.Rect(100, 50, 110, 60), speed: 2}
	base := NewBaseMovementState(Follow, actor, target)
	state := NewFollowMovementState(base)
	state.startDistance = 10
	state.stopDistance = 5
	state.isMoving = true
	state.stayOnPlatform = true

	// No ground to the right (ledge)
	space := newMockSpaceWithGround([]image.Point{})

	state.Move(space)

	// Should not move right because there's no ground
	if actor.moveRightForce != 0 {
		t.Errorf("expected actor to not move right when no ground ahead with platform follow")
	}
}
