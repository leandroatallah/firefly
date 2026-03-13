package movement

import (
	"image"
	"testing"
)

func TestNewBaseMovementState(t *testing.T) {
	actor := &mockActor{pos: image.Rect(0, 0, 10, 10), speed: 2}
	target := &mockActor{pos: image.Rect(50, 50, 60, 60), speed: 2}

	state := NewBaseMovementState(Chase, actor, target)

	if state.state != Chase {
		t.Errorf("expected state Chase, got %v", state.state)
	}
	if state.actor != actor {
		t.Errorf("expected actor to be set")
	}
	if state.target != target {
		t.Errorf("expected target to be set")
	}
}

func TestBaseMovementState_State(t *testing.T) {
	actor := &mockActor{}
	state := NewBaseMovementState(Patrol, actor, nil)

	if state.State() != Patrol {
		t.Errorf("expected State to return Patrol, got %v", state.State())
	}
}

func TestBaseMovementState_Target(t *testing.T) {
	actor := &mockActor{}
	target := &mockActor{pos: image.Rect(100, 100, 110, 110)}
	state := NewBaseMovementState(Idle, actor, target)

	if state.Target() != target {
		t.Errorf("expected Target to return the target actor")
	}
}

func TestBaseMovementState_SetTarget(t *testing.T) {
	actor := &mockActor{}
	initialTarget := &mockActor{pos: image.Rect(10, 10, 20, 20)}
	newTarget := &mockActor{pos: image.Rect(50, 50, 60, 60)}
	state := NewBaseMovementState(Idle, actor, initialTarget)

	state.SetTarget(newTarget)

	if state.Target() != newTarget {
		t.Errorf("expected SetTarget to update the target")
	}
}

func TestBaseMovementState_Actor(t *testing.T) {
	actor := &mockActor{}
	state := NewBaseMovementState(Idle, actor, nil)

	if state.Actor() != actor {
		t.Errorf("expected Actor to return the actor")
	}
}

func TestBaseMovementState_OnStart(t *testing.T) {
	actor := &mockActor{}
	state := NewBaseMovementState(Idle, actor, nil)

	// Should not panic - base implementation is a no-op
	state.OnStart()
}

func TestCalculateMovementDirections(t *testing.T) {
	tests := []struct {
		name      string
		actorPos  image.Rectangle
		targetPos image.Rectangle
		isAvoid   bool
		want      MovementDirections
	}{
		{
			name:      "target to the right",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(20, 0, 30, 10),
			isAvoid:   false,
			want:      MovementDirections{Right: true},
		},
		{
			name:      "target to the left",
			actorPos:  image.Rect(20, 0, 30, 10),
			targetPos: image.Rect(0, 0, 10, 10),
			isAvoid:   false,
			want:      MovementDirections{Left: true},
		},
		{
			name:      "target below",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(0, 20, 10, 30),
			isAvoid:   false,
			want:      MovementDirections{Down: true},
		},
		{
			name:      "target above",
			actorPos:  image.Rect(0, 20, 10, 30),
			targetPos: image.Rect(0, 0, 10, 10),
			isAvoid:   false,
			want:      MovementDirections{Up: true},
		},
		{
			name:      "target diagonal (down-right)",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(20, 20, 30, 30),
			isAvoid:   false,
			want:      MovementDirections{Right: true, Down: true},
		},
		{
			name:      "target diagonal (up-left)",
			actorPos:  image.Rect(20, 20, 30, 30),
			targetPos: image.Rect(0, 0, 10, 10),
			isAvoid:   false,
			want:      MovementDirections{Left: true, Up: true},
		},
		{
			name:      "avoid mode - target to right",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(20, 0, 30, 10),
			isAvoid:   true,
			want:      MovementDirections{Up: true, Down: true, Left: true},
		},
		{
			name:      "avoid mode - target diagonal",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(20, 20, 30, 30),
			isAvoid:   true,
			want:      MovementDirections{Left: true, Up: true},
		},
		{
			name:      "overlapping positions",
			actorPos:  image.Rect(0, 0, 10, 10),
			targetPos: image.Rect(5, 5, 15, 15),
			isAvoid:   false,
			want:      MovementDirections{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorBody := &mockActor{pos: tt.actorPos}
			targetBody := &mockActor{pos: tt.targetPos}

			got := calculateMovementDirections(actorBody, targetBody, tt.isAvoid)

			if got != tt.want {
				t.Errorf("calculateMovementDirections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecuteMovement(t *testing.T) {
	tests := []struct {
		name       string
		directions MovementDirections
		speed      int
		wantLeft   int
		wantRight  int
		wantUp     int
		wantDown   int
	}{
		{
			name:       "move right",
			directions: MovementDirections{Right: true},
			speed:      2,
			wantRight:  2,
		},
		{
			name:       "move left",
			directions: MovementDirections{Left: true},
			speed:      3,
			wantLeft:   3,
		},
		{
			name:       "move up",
			directions: MovementDirections{Up: true},
			speed:      2,
			wantUp:     2,
		},
		{
			name:       "move down",
			directions: MovementDirections{Down: true},
			speed:      2,
			wantDown:   2,
		},
		{
			name:       "move up-right",
			directions: MovementDirections{Up: true, Right: true},
			speed:      2,
			wantUp:     2,
			wantRight:  2,
		},
		{
			name:       "move up-left",
			directions: MovementDirections{Up: true, Left: true},
			speed:      2,
			wantUp:     2,
			wantLeft:   2,
		},
		{
			name:       "move down-right",
			directions: MovementDirections{Down: true, Right: true},
			speed:      2,
			wantDown:   2,
			wantRight:  2,
		},
		{
			name:       "move down-left",
			directions: MovementDirections{Down: true, Left: true},
			speed:      2,
			wantDown:   2,
			wantLeft:   2,
		},
		{
			name:       "no movement",
			directions: MovementDirections{},
			speed:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := &mockActor{speed: tt.speed}

			executeMovement(actor, tt.directions)

			if actor.moveLeftForce != tt.wantLeft {
				t.Errorf("moveLeftForce = %d, want %d", actor.moveLeftForce, tt.wantLeft)
			}
			if actor.moveRightForce != tt.wantRight {
				t.Errorf("moveRightForce = %d, want %d", actor.moveRightForce, tt.wantRight)
			}
			if actor.moveUpForce != tt.wantUp {
				t.Errorf("moveUpForce = %d, want %d", actor.moveUpForce, tt.wantUp)
			}
			if actor.moveDownForce != tt.wantDown {
				t.Errorf("moveDownForce = %d, want %d", actor.moveDownForce, tt.wantDown)
			}
		})
	}
}

func TestExecuteMovement_NoDirections(t *testing.T) {
	actor := &mockActor{speed: 2}

	executeMovement(actor, MovementDirections{})

	if actor.moveLeftForce != 0 || actor.moveRightForce != 0 ||
		actor.moveUpForce != 0 || actor.moveDownForce != 0 {
		t.Errorf("expected no movement when all directions are false")
	}
}
