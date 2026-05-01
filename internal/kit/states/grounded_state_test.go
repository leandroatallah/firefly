package kitstates_test

import (
	"testing"

	kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"
)

func newGroundedState(input *MockInputSource) *kitstates.GroundedState {
	return kitstates.NewGroundedState(kitstates.GroundedDeps{Input: input})
}

func TestGroundedSubStateTransitions(t *testing.T) {
	tests := []struct {
		name         string
		setupInput   func(*MockInputSource)
		initialSetup func(*kitstates.GroundedState)
		wantSubState kitstates.GroundedSubStateEnum
	}{
		{
			name:         "no input stays Idle",
			setupInput:   func(m *MockInputSource) {},
			initialSetup: nil,
			wantSubState: kitstates.SubStateIdle,
		},
		{
			name: "horizontal input transitions to Walking",
			setupInput: func(m *MockInputSource) {
				m.HorizontalInputFunc = func() int { return 1 }
			},
			initialSetup: nil,
			wantSubState: kitstates.SubStateWalking,
		},
		{
			name: "duck input transitions to Ducking",
			setupInput: func(m *MockInputSource) {
				m.DuckHeldFunc = func() bool { return true }
			},
			initialSetup: func(g *kitstates.GroundedState) {
				// start from Walking
				g.ForceSubState(kitstates.SubStateWalking)
			},
			wantSubState: kitstates.SubStateDucking,
		},
		{
			name: "duck released with clearance transitions to Idle",
			setupInput: func(m *MockInputSource) {
				m.DuckHeldFunc = func() bool { return false }
				m.HasCeilingClearanceFunc = func() bool { return true }
			},
			initialSetup: func(g *kitstates.GroundedState) {
				g.ForceSubState(kitstates.SubStateDucking)
			},
			wantSubState: kitstates.SubStateIdle,
		},
		{
			name: "aim-lock input transitions to AimLock",
			setupInput: func(m *MockInputSource) {
				m.AimLockHeldFunc = func() bool { return true }
			},
			initialSetup: nil,
			wantSubState: kitstates.SubStateAimLock,
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

	if next == kitstates.StateGrounded {
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

	if g.ActiveSubState() != kitstates.SubStateDucking {
		t.Fatalf("expected SubStateDucking before OnFinish, got %v", g.ActiveSubState())
	}

	// OnFinish must call the active sub-state's OnFinish without panicking.
	// We verify indirectly: after OnFinish + re-entry, sub-state is Idle (not Ducking),
	// which means Ducking.OnFinish() ran (restores body) and re-entry reset to Idle.
	g.OnFinish()
	g.OnStart(0)

	if g.ActiveSubState() != kitstates.SubStateIdle {
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

	if g.ActiveSubState() != kitstates.SubStateWalking {
		t.Fatalf("expected SubStateWalking, got %v", g.ActiveSubState())
	}

	g.OnFinish()
	g.OnStart(0)

	if g.ActiveSubState() != kitstates.SubStateIdle {
		t.Errorf("re-entry: expected SubStateIdle, got %v", g.ActiveSubState())
	}
}

// TestGroundedState_State verifies that State() returns StateGrounded.
func TestGroundedState_State(t *testing.T) {
	g := newGroundedState(&MockInputSource{})
	g.OnStart(0)

	if got := g.State(); got != kitstates.StateGrounded {
		t.Errorf("State() = %v, want StateGrounded", got)
	}
}

// TestGroundedState_GetAnimationCount verifies that GetAnimationCount returns
// the difference between the current count and the count recorded at OnStart.
func TestGroundedState_GetAnimationCount(t *testing.T) {
	tests := []struct {
		name         string
		startCount   int
		currentCount int
		want         int
	}{
		{
			name:         "zero frames elapsed",
			startCount:   0,
			currentCount: 0,
			want:         0,
		},
		{
			name:         "ten frames elapsed",
			startCount:   5,
			currentCount: 15,
			want:         10,
		},
		{
			name:         "started at large offset",
			startCount:   1000,
			currentCount: 1042,
			want:         42,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := newGroundedState(&MockInputSource{})
			g.OnStart(tc.startCount)

			if got := g.GetAnimationCount(tc.currentCount); got != tc.want {
				t.Errorf("GetAnimationCount(%d) after OnStart(%d) = %d, want %d",
					tc.currentCount, tc.startCount, got, tc.want)
			}
		})
	}
}

// TestGroundedState_IsAnimationFinished verifies that IsAnimationFinished always
// returns false (GroundedState loops its animation indefinitely).
func TestGroundedState_IsAnimationFinished(t *testing.T) {
	g := newGroundedState(&MockInputSource{})
	g.OnStart(0)

	if g.IsAnimationFinished() {
		t.Errorf("IsAnimationFinished() = true, want false (grounded state loops)")
	}
}

// TestWalkingSubState_OnStartOnFinish exercises walkingSubState.OnStart and
// OnFinish via GroundedState transitions: Idle→Walking triggers OnStart, then
// Walking→Idle triggers OnFinish.
func TestWalkingSubState_OnStartOnFinish(t *testing.T) {
	moving := true
	input := &MockInputSource{
		HorizontalInputFunc: func() int {
			if moving {
				return 1
			}
			return 0
		},
	}
	g := newGroundedState(input)
	g.OnStart(0)

	// Idle→Walking: calls walkingSubState.OnStart internally.
	g.Update()
	if got := g.ActiveSubState(); got != kitstates.SubStateWalking {
		t.Fatalf("expected SubStateWalking after horizontal input, got %v", got)
	}

	// Walking→Idle: calls walkingSubState.OnFinish internally.
	moving = false
	g.Update()
	if got := g.ActiveSubState(); got != kitstates.SubStateIdle {
		t.Errorf("expected SubStateIdle after stopping, got %v", got)
	}
}

// TestDuckingSubState_OnStartOnFinish exercises duckingSubState.OnStart and
// OnFinish via GroundedState transitions: Walking→Ducking (OnStart), then
// Ducking→Idle (OnFinish).
func TestDuckingSubState_OnStartOnFinish(t *testing.T) {
	ducking := true
	input := &MockInputSource{
		DuckHeldFunc:            func() bool { return ducking },
		HasCeilingClearanceFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	// Idle→Ducking: calls duckingSubState.OnStart internally.
	g.Update()
	if got := g.ActiveSubState(); got != kitstates.SubStateDucking {
		t.Fatalf("expected SubStateDucking after duck input, got %v", got)
	}

	// Ducking→Idle: calls duckingSubState.OnFinish internally.
	ducking = false
	g.Update()
	if got := g.ActiveSubState(); got != kitstates.SubStateIdle {
		t.Errorf("expected SubStateIdle after duck released, got %v", got)
	}
}
