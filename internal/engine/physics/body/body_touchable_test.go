package body

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

func TestNewTouchTrigger(t *testing.T) {
	execute := func() {}
	subject := NewBody(NewRect(0, 0, 10, 10))

	tt := NewTouchTrigger(execute, subject)

	if tt == nil {
		t.Fatal("NewTouchTrigger returned nil")
	}
	if tt.subject != subject {
		t.Error("expected subject to be set")
	}
}

func TestTouchTrigger_OnTouch_MatchingSubject(t *testing.T) {
	called := false
	execute := func() { called = true }
	subject := NewBody(NewRect(0, 0, 10, 10))
	subject.SetID("subject-1")

	tt := NewTouchTrigger(execute, subject)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("subject-1")

	tt.OnTouch(other)

	if !called {
		t.Error("expected OnTouch to execute callback when subject matches")
	}
}

func TestTouchTrigger_OnTouch_NonMatchingSubject(t *testing.T) {
	called := false
	execute := func() { called = true }
	subject := NewBody(NewRect(0, 0, 10, 10))
	subject.SetID("subject-1")

	tt := NewTouchTrigger(execute, subject)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("different-subject")

	tt.OnTouch(other)

	if called {
		t.Error("expected OnTouch to NOT execute callback when subject doesn't match")
	}
}

func TestTouchTrigger_OnTouch_NilSubject(t *testing.T) {
	called := false
	execute := func() { called = true }

	tt := NewTouchTrigger(execute, nil)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("any-id")

	// Should not panic
	tt.OnTouch(other)

	if called {
		t.Error("expected OnTouch to NOT execute callback when subject is nil")
	}
}

func TestTouchTrigger_OnTouch_NilExecute(t *testing.T) {
	subject := NewBody(NewRect(0, 0, 10, 10))
	subject.SetID("subject-1")

	tt := NewTouchTrigger(nil, subject)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("subject-1")

	// Should not panic with nil execute function
	tt.OnTouch(other)
}

func TestTouchTrigger_OnBlock(t *testing.T) {
	called := false
	execute := func() { called = true }
	subject := NewBody(NewRect(0, 0, 10, 10))

	tt := NewTouchTrigger(execute, subject)

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))

	// OnBlock should be a no-op for TouchTrigger
	tt.OnBlock(other)

	if called {
		t.Error("expected OnBlock to NOT execute callback (no-op)")
	}
}

func TestTouchTrigger_OnTouch_WithTouchableInterface(t *testing.T) {
	touchCount := 0
	execute := func() { touchCount++ }
	subject := NewBody(NewRect(0, 0, 10, 10))
	subject.SetID("player")

	tt := NewTouchTrigger(execute, subject)

	// Create a collidable that will touch
	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("player")

	// Verify TouchTrigger implements body.Touchable
	var _ body.Touchable = tt

	tt.OnTouch(other)

	if touchCount != 1 {
		t.Errorf("expected touchCount to be 1; got %d", touchCount)
	}
}

func TestTouchTrigger_MultipleTouches(t *testing.T) {
	touchCount := 0
	execute := func() { touchCount++ }
	subject := NewBody(NewRect(0, 0, 10, 10))
	subject.SetID("subject")

	tt := NewTouchTrigger(execute, subject)

	for i := 0; i < 5; i++ {
		other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
		other.SetID("subject")
		tt.OnTouch(other)
	}

	if touchCount != 5 {
		t.Errorf("expected touchCount to be 5; got %d", touchCount)
	}
}

func TestTouchTrigger_DifferentSubjects(t *testing.T) {
	touchCount := 0
	execute := func() { touchCount++ }
	subject := NewBody(NewRect(0, 0, 10, 10))
	subject.SetID("target")

	tt := NewTouchTrigger(execute, subject)

	// Touch with matching subject
	matching := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	matching.SetID("target")
	tt.OnTouch(matching)

	// Touch with non-matching subject
	nonMatching := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	nonMatching.SetID("other")
	tt.OnTouch(nonMatching)

	// Touch with matching subject again
	matching2 := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	matching2.SetID("target")
	tt.OnTouch(matching2)

	if touchCount != 2 {
		t.Errorf("expected touchCount to be 2 (only matching); got %d", touchCount)
	}
}
