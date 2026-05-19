// Red-Phase tests for story 059-thin-game-phase-scenes [AC-11].
// SPEC.md §2 moves the goal-type constants and ReachEndpointGoal
// implementation from internal/game/scenes/phases into the engine
// package so both kit genres can resolve them.
//
// These tests intentionally fail until the engine phases package
// exposes ReactEndpointType, SequenceGoalType, NoGoalType, and a
// ReachEndpointGoal struct with Reach() and an optional OnCompletion_
// callback.
package phases

import "testing"

func TestEngineGoalTypeConstants_AreDefined(t *testing.T) {
	cases := []struct {
		name string
		got  GoalType
		want string
	}{
		{name: "reach_endpoint", got: ReactEndpointType, want: "reach_endpoint"},
		{name: "sequence", got: SequenceGoalType, want: "sequence"},
		{name: "no_goal", got: NoGoalType, want: "no_goal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if string(tc.got) != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, string(tc.got))
			}
		})
	}
}

func TestReachEndpointGoal_NotCompletedBeforeReach(t *testing.T) {
	g := &ReachEndpointGoal{}
	if g.IsCompleted() {
		t.Fatal("expected IsCompleted() == false before Reach()")
	}
}

func TestReachEndpointGoal_CompletedAfterReach(t *testing.T) {
	g := &ReachEndpointGoal{}
	g.Reach()
	if !g.IsCompleted() {
		t.Fatal("expected IsCompleted() == true after Reach()")
	}
}

func TestReachEndpointGoal_OnCompletionInvokesCallback(t *testing.T) {
	calls := 0
	g := &ReachEndpointGoal{
		OnCompletion_: func() { calls++ },
	}
	g.OnCompletion()
	if calls != 1 {
		t.Fatalf("expected callback invoked exactly once, got %d", calls)
	}
}

func TestReachEndpointGoal_OnCompletionNoOpWhenCallbackNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic when callback is nil, got: %v", r)
		}
	}()
	g := &ReachEndpointGoal{}
	g.OnCompletion()
}
