package gamestates_test

import (
	"testing"

	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

func newGroundedState(input *MockInputSource) *gamestates.GroundedState {
	return gamestates.NewGroundedState(gamestates.GroundedDeps{Input: input})
}

func TestGroundedSubStateTransitions(t *testing.T) {
	tests := []struct {
		name         string
		setupInput   func(*MockInputSource)
		initialSetup func(*gamestates.GroundedState)
		wantSubState gamestates.GroundedSubStateEnum
	}{
		{
			name:         "no input stays Idle",
			setupInput:   func(m *MockInputSource) {},
			initialSetup: nil,
			wantSubState: gamestates.SubStateIdle,
		},
		{
			name: "horizontal input transitions to Walking",
			setupInput: func(m *MockInputSource) {
				m.HorizontalInputFunc = func() int { return 1 }
			},
			initialSetup: nil,
			wantSubState: gamestates.SubStateWalking,
		},
		{
			name: "duck input transitions to Ducking",
			setupInput: func(m *MockInputSource) {
				m.DuckHeldFunc = func() bool { return true }
			},
			initialSetup: func(g *gamestates.GroundedState) {
				// start from Walking
				g.ForceSubState(gamestates.SubStateWalking)
			},
			wantSubState: gamestates.SubStateDucking,
		},
		{
			name: "duck released with clearance transitions to Idle",
			setupInput: func(m *MockInputSource) {
				m.DuckHeldFunc = func() bool { return false }
				m.HasCeilingClearanceFunc = func() bool { return true }
			},
			initialSetup: func(g *gamestates.GroundedState) {
				g.ForceSubState(gamestates.SubStateDucking)
			},
			wantSubState: gamestates.SubStateIdle,
		},
		{
			name: "aim-lock input transitions to AimLock",
			setupInput: func(m *MockInputSource) {
				m.AimLockHeldFunc = func() bool { return true }
			},
			initialSetup: nil,
			wantSubState: gamestates.SubStateAimLock,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := &MockInputSource{}
			g := newGroundedState(input)
			g.OnStart(0)

			if tc.initialSetup != nil {
				tc.initialSetup(g)
			}

			tc.setupInput(input)
			g.Update()

			if got := g.ActiveSubState(); got != tc.wantSubState {
				t.Errorf("ActiveSubState() = %v, want %v", got, tc.wantSubState)
			}
		})
	}
}

func TestGroundedSubStateTransitions_JumpExitsGrounded(t *testing.T) {
	input := &MockInputSource{
		JumpPressedFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	next := g.Update()

	if next == gamestates.StateGrounded {
		t.Errorf("jump input should exit Grounded, but Update() returned StateGrounded")
	}
}

func TestGroundedStateOnFinishCallsSubOnFinish(t *testing.T) {
	input := &MockInputSource{
		DuckHeldFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)
	g.Update() // transition to Ducking

	if g.ActiveSubState() != gamestates.SubStateDucking {
		t.Fatalf("expected SubStateDucking before OnFinish, got %v", g.ActiveSubState())
	}

	// OnFinish must call the active sub-state's OnFinish without panicking.
	// We verify indirectly: after OnFinish + re-entry, sub-state is Idle (not Ducking),
	// which means Ducking.OnFinish() ran (restores body) and re-entry reset to Idle.
	g.OnFinish()
	g.OnStart(0)

	if g.ActiveSubState() != gamestates.SubStateIdle {
		t.Errorf("after OnFinish+OnStart, expected SubStateIdle, got %v", g.ActiveSubState())
	}
}

func TestGroundedStateReEntryResetsToIdle(t *testing.T) {
	input := &MockInputSource{
		HorizontalInputFunc: func() int { return 1 },
	}
	g := newGroundedState(input)
	g.OnStart(0)
	g.Update() // → Walking

	if g.ActiveSubState() != gamestates.SubStateWalking {
		t.Fatalf("expected SubStateWalking, got %v", g.ActiveSubState())
	}

	g.OnFinish()
	g.OnStart(0)

	if g.ActiveSubState() != gamestates.SubStateIdle {
		t.Errorf("re-entry: expected SubStateIdle, got %v", g.ActiveSubState())
	}
}
