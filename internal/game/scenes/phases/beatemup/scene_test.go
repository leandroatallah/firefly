package gamebeatemupphase

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
)

// OnStart() is not unit-tested directly — it requires a full AppContext with GPU resources.

func TestReachEndpointGoal_Completion(t *testing.T) {
	cases := []struct {
		name            string
		reachedEndpoint bool
		wantCompleted   bool
	}{
		{"not completed when endpoint not reached", false, false},
		{"completed when endpoint reached", true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &BeatemupPhaseScene{reachedEndpoint: tc.reachedEndpoint}
			goal := &ReachEndpointGoal{scene: s}
			if got := goal.IsCompleted(); got != tc.wantCompleted {
				t.Errorf("IsCompleted() = %v, want %v", got, tc.wantCompleted)
			}
		})
	}
}

func TestReachEndpointGoal_OnCompletion_NilAudioManager(t *testing.T) {
	s := &BeatemupPhaseScene{reachedEndpoint: true}
	goal := &ReachEndpointGoal{scene: s}
	// Must not panic when AudioManager is nil.
	goal.OnCompletion()
}

func TestBodyCounter_CountsWolfBodies(t *testing.T) {
	space := &mockBodiesSpace{}
	space.AddBody(&mockCollidable{id: "other"})
	bc := &BodyCounter{}
	bc.setBodyCounter(space)
	if bc.wolf != 0 {
		t.Errorf("wolf = %d, want 0", bc.wolf)
	}
}

func TestDefaultCompletion_EnablesTrigger(t *testing.T) {
	s := &BeatemupPhaseScene{}
	s.defaultCompletion()
	if !s.completionTrigger.IsEnabled() {
		t.Error("completionTrigger not enabled")
	}
}

func TestEndpointTrigger_SetsReachedEndpoint(t *testing.T) {
	s := &BeatemupPhaseScene{
		TilemapScene: &scene.TilemapScene{},
		hasPlayer:    true,
	}
	ctx := &app.AppContext{Space: &mockBodiesSpace{}}
	s.SetAppContext(ctx)
	s.endpointTrigger("test_endpoint")
	if !s.reachedEndpoint {
		t.Error("reachedEndpoint not set to true")
	}
}

func TestTriggerScreenFlash_SetsFlashCounter(t *testing.T) {
	s := &BeatemupPhaseScene{}
	s.TriggerScreenFlash()
	if s.ShowDrawScreenFlash == 0 {
		t.Error("ShowDrawScreenFlash not set")
	}
}

func TestCanPause_RequiresAllowPauseAndNoSequence(t *testing.T) {
	cases := []struct {
		name       string
		allowPause bool
		playing    bool
		want       bool
	}{
		{"cannot pause when not allowed", false, false, false},
		{"cannot pause during sequence", true, true, false},
		{"can pause when allowed and no sequence", true, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &BeatemupPhaseScene{
				allowPause:     tc.allowPause,
				sequencePlayer: &mockSequencePlayer{playing: tc.playing},
			}
			if got := s.canPause(); got != tc.want {
				t.Errorf("canPause() = %v, want %v", got, tc.want)
			}
		})
	}
}

func newUpdateHarness(goal phases.Goal, sm *mockSceneManager) *BeatemupPhaseScene {
	space := &mockBodiesSpace{}
	ctx := &app.AppContext{
		Space:        space,
		SceneManager: sm,
	}
	s := &BeatemupPhaseScene{
		TilemapScene: scene.NewTilemapScene(ctx),
		goal:         goal,
		hasPlayer:    false,
	}
	s.SetAppContext(ctx)
	return s
}

func TestUpdate_GoalCompletion_CallsOnCompletion(t *testing.T) {
	goal := &mockGoal{completed: true}
	s := newUpdateHarness(goal, &mockSceneManager{})
	if err := s.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if !goal.onCompletionCalled {
		t.Error("OnCompletion() not called when goal is completed")
	}
}

func TestUpdate_GoalPartial_DoesNotCallOnCompletion(t *testing.T) {
	goal := &mockGoal{completed: false}
	s := newUpdateHarness(goal, &mockSceneManager{})
	if err := s.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if goal.onCompletionCalled {
		t.Error("OnCompletion() must not be called when goal is not completed")
	}
}

func TestUpdate_SequencePlayerUpdated(t *testing.T) {
	sp := &trackingSequencePlayer{}
	space := &mockBodiesSpace{}
	ctx := &app.AppContext{Space: space, SceneManager: &mockSceneManager{}}
	s := &BeatemupPhaseScene{
		TilemapScene:   scene.NewTilemapScene(ctx),
		sequencePlayer: sp,
		hasPlayer:      false,
	}
	s.SetAppContext(ctx)
	if err := s.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if !sp.updateCalled {
		t.Error("sequencePlayer.Update() not called during Update()")
	}
}

func TestUpdate_DeathTrigger_NavigatesToReboot(t *testing.T) {
	sm := &mockSceneManager{}
	s := newUpdateHarness(nil, sm)
	s.deathTrigger.Enable(0)
	if err := s.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if !sm.navigateToCalled {
		t.Error("SceneManager.NavigateTo not called when deathTrigger fires")
	}
}

// trackingSequencePlayer wraps mockSequencePlayer and records Update calls.
type trackingSequencePlayer struct {
	mockSequencePlayer
	updateCalled bool
}

func (p *trackingSequencePlayer) Update() { p.updateCalled = true }

func TestOnFinish_NoPlayer_DoesNotPanic(t *testing.T) {
	ctx := &app.AppContext{Space: &mockBodiesSpace{}}
	s := &BeatemupPhaseScene{
		TilemapScene: scene.NewTilemapScene(ctx),
		hasPlayer:    false,
	}
	s.SetAppContext(ctx)
	s.OnFinish() // must not panic
}
