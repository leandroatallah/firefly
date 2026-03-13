package movement

import (
	"testing"
)

func TestDumbChaseMovementState(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)
	
	target := &mockActor{}
	target.SetPosition(100, 100)

	base := NewBaseMovementState(DumbChase, actor, target)
	state := NewDumbChaseMovementState(base)

	state.Move(nil)

	if actor.moveLeftForce != 0 || actor.moveRightForce != 5 {
		t.Errorf("expected moveRightForce 5, got Left:%d Right:%d", actor.moveLeftForce, actor.moveRightForce)
	}
}

func TestAvoidMovementState(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(50, 50)
	
	target := &mockActor{}
	target.SetPosition(100, 100)

	base := NewBaseMovementState(Avoid, actor, target)
	state := NewAvoidMovementState(base)

	state.Move(nil)

	if actor.moveLeftForce != 5 {
		t.Errorf("expected moveLeftForce 5 when avoiding target to the right, got %d", actor.moveLeftForce)
	}
}
