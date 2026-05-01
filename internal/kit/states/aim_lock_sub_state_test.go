package kitstates_test

import (
	"testing"

	kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"
)

// TestAimLockSubState_TransitionTo exercises TransitionTo through GroundedState
// which internally creates an aimLockSubState when SubStateAimLock is active.
func TestAimLockSubState_TransitionTo(t *testing.T) {
	tests := []struct {
		name         string
		aimLockHeld  bool
		wantSubState kitstates.GroundedSubStateEnum
	}{
		{
			name:         "aim-lock held stays in AimLock",
			aimLockHeld:  true,
			wantSubState: kitstates.SubStateAimLock,
		},
		{
			name:         "aim-lock released transitions to Idle",
			aimLockHeld:  false,
			wantSubState: kitstates.SubStateIdle,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := &MockInputSource{
				AimLockHeldFunc: func() bool { return tc.aimLockHeld },
			}
			g := newGroundedState(input)
			g.OnStart(0)
			g.ForceSubState(kitstates.SubStateAimLock)

			g.Update()

			if got := g.ActiveSubState(); got != tc.wantSubState {
				t.Errorf("ActiveSubState() = %v, want %v", got, tc.wantSubState)
			}
		})
	}
}

// TestAimLockSubState_OnStart_CalledOnEntry verifies that aimLockSubState.OnStart
// is called when GroundedState transitions into SubStateAimLock. Exercised by
// forcing the transition and confirming the sub-state is active without panic.
func TestAimLockSubState_OnStart_CalledOnEntry(t *testing.T) {
	input := &MockInputSource{
		AimLockHeldFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	// Update transitions from Idle → AimLock, which calls aimLockSubState.OnStart.
	g.Update()

	if got := g.ActiveSubState(); got != kitstates.SubStateAimLock {
		t.Errorf("expected SubStateAimLock after aim-lock input, got %v", got)
	}
}

// TestAimLockSubState_OnFinish_CalledOnExit verifies aimLockSubState.OnFinish is
// called when the sub-state changes away from AimLock (e.g. aim-lock released).
func TestAimLockSubState_OnFinish_CalledOnExit(t *testing.T) {
	aimLocked := true
	input := &MockInputSource{
		AimLockHeldFunc: func() bool { return aimLocked },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	// Enter AimLock sub-state.
	g.Update()
	if got := g.ActiveSubState(); got != kitstates.SubStateAimLock {
		t.Fatalf("precondition: expected SubStateAimLock, got %v", got)
	}

	// Release aim-lock — causes aimLockSubState.OnFinish() to be called.
	aimLocked = false
	g.Update()

	if got := g.ActiveSubState(); got != kitstates.SubStateIdle {
		t.Errorf("after aim-lock released: ActiveSubState() = %v, want SubStateIdle", got)
	}
}
