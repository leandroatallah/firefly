package movement

import (
	"testing"

	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
)

func TestPatrolMovementState(t *testing.T) {
	actor := &mockActor{speed: 2}
	actor.SetPosition(0, 0)

	wp1 := bodyphysics.NewRect(20, 0, 10, 10)
	config := NewPredefinedWaypointConfig([]*bodyphysics.Rect{wp1}, 5)

	base := NewBaseMovementState(Patrol, actor, nil)
	state := NewPatrolMovementState(base)
	
	state.SetWaypointConfig(config)
	state.OnStart()

	// Manually ensure we are in chase state and moving right for the test
	state.patrolState = patrolChase
	state.movementDirections = MovementDirections{Right: true}

	state.Move(nil)
	if actor.moveRightForce != 2 {
		t.Errorf("expected moveRightForce 2, got %d", actor.moveRightForce)
	}
}
