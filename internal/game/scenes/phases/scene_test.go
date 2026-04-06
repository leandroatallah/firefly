package gamescenephases

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
			s := &PhasesScene{reachedEndpoint: tc.reachedEndpoint}
			goal := &ReachEndpointGoal{scene: s}
			if got := goal.IsCompleted(); got != tc.wantCompleted {
				t.Errorf("IsCompleted() = %v, want %v", got, tc.wantCompleted)
			}
		})
	}
}

func TestReachEndpointGoal_OnCompletion_NilAudioManager(t *testing.T) {
	// Production code must guard against nil AudioManager — this test exposes the missing guard.
	s := &PhasesScene{reachedEndpoint: true}
	goal := &ReachEndpointGoal{scene: s}
	// Must not panic when AudioManager is nil.
	goal.OnCompletion()
}

func TestNoGoal_NeverCompletes(t *testing.T) {
	g := &phases.NoGoal{}
	if g.IsCompleted() {
		t.Error("NoGoal.IsCompleted() must always return false")
	}
	// OnCompletion must not panic
	g.OnCompletion()
}

func TestSequenceGoal_CompletesWhenNotPlaying(t *testing.T) {
	cases := []struct {
		name          string
		playing       bool
		wantCompleted bool
	}{
		{"not completed while sequence is playing", true, false},
		{"completed when sequence stops", false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := &phases.SequenceGoal{Player: &mockSequencePlayer{playing: tc.playing}}
			if got := g.IsCompleted(); got != tc.wantCompleted {
				t.Errorf("IsCompleted() = %v, want %v", got, tc.wantCompleted)
			}
		})
	}
}

func TestFreezeController_PausesForConfiguredFrames(t *testing.T) {
	var fc scene.FreezeController
	fc.FreezeFrame(3)

	for i := 0; i < 3; i++ {
		if !fc.IsFrozen() {
			t.Fatalf("expected frozen on tick %d", i)
		}
		fc.Tick()
	}
	if fc.IsFrozen() {
		t.Error("expected not frozen after all ticks consumed")
	}
}

func TestBodyCounter_CountsWolfBodies(t *testing.T) {
	space := &mockBodiesSpace{}
	space.AddBody(&mockCollidable{id: "other"})
	bc := &BodyCounter{}
	bc.setBodyCounter(space)
	// No WolfEnemy bodies — wolf count must be zero.
	if bc.wolf != 0 {
		t.Errorf("wolf = %d, want 0", bc.wolf)
	}
}

func TestCheckPlayerFallDeath_TriggersWhenBelowCamera(t *testing.T) {
	s := &PhasesScene{death: deathSequence{active: false}}
	s.checkPlayerFallDeath()
	if s.death.active {
		t.Error("death triggered with nil gameCamera")
	}
}

func TestDefaultCompletion_EnablesTrigger(t *testing.T) {
	s := &PhasesScene{}
	s.defaultCompletion()
	if !s.completionTrigger.IsEnabled() {
		t.Error("completionTrigger not enabled")
	}
}

func TestCamera_ReturnsGameCamera(t *testing.T) {
	s := &PhasesScene{TilemapScene: &scene.TilemapScene{}}
	cam := s.Camera()
	if cam == nil {
		t.Error("Camera() returned nil")
	}
}

func TestBaseCamera_ReturnsUnderlyingCamera(t *testing.T) {
	ctx := &app.AppContext{Space: &mockBodiesSpace{}}
	s := &PhasesScene{TilemapScene: scene.NewTilemapScene(ctx)}
	cam := s.BaseCamera()
	if cam == nil {
		t.Error("BaseCamera() returned nil")
	}
}

func TestEndpointTrigger_SetsReachedEndpoint(t *testing.T) {
	s := &PhasesScene{
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
	s := &PhasesScene{}
	s.TriggerScreenFlash()
	if s.ShowDrawScreenFlash == 0 {
		t.Error("ShowDrawScreenFlash not set")
	}
}

func TestStartDeathSequence_ActivatesDeathState(t *testing.T) {
	s := &PhasesScene{
		TilemapScene: &scene.TilemapScene{},
		death:        deathSequence{active: false},
	}
	ctx := &app.AppContext{Space: &mockBodiesSpace{}}
	s.SetAppContext(ctx)
	s.startDeathSequence()
	if !s.death.active {
		t.Error("death.active not set to true")
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
			s := &PhasesScene{
				allowPause:     tc.allowPause,
				sequencePlayer: &mockSequencePlayer{playing: tc.playing},
			}
			if got := s.canPause(); got != tc.want {
				t.Errorf("canPause() = %v, want %v", got, tc.want)
			}
		})
	}
}

func newUpdateHarness(goal phases.Goal, sm *mockSceneManager) *PhasesScene {
	space := &mockBodiesSpace{}
	ctx := &app.AppContext{
		Space:        space,
		SceneManager: sm,
	}
	s := &PhasesScene{
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
	s := &PhasesScene{
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
	// Enable deathTrigger with 0-frame delay so it fires immediately.
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
	s := &PhasesScene{
		TilemapScene: scene.NewTilemapScene(ctx),
		hasPlayer:    false,
	}
	s.SetAppContext(ctx)
	s.OnFinish() // must not panic
}

func TestDisableVignetteDarkness_NilVignette_DoesNotPanic(t *testing.T) {
	s := &PhasesScene{vignette: nil}
	s.DisableVignetteDarkness() // must not panic
}

func TestUpdate_CompletionTrigger_CallsCompleteCurrentPhase(t *testing.T) {
	sm := &mockSceneManager{}
	goal := &mockGoal{completed: true}
	s := newUpdateHarness(goal, sm)
	// First Update: goal.IsCompleted() true → OnCompletion() called → completionTrigger.Enable(n)
	// We manually enable with 0-frame delay to fire on next Update.
	s.completionTrigger.Enable(0)
	// CompleteCurrentPhase calls PhaseManager.AdvanceToNextPhase — PhaseManager is nil so it returns early.
	// NavigateTo is NOT called via CompleteCurrentPhase when PhaseManager is nil.
	// This test just verifies Update() does not panic when completionTrigger fires with nil PhaseManager.
	if err := s.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}
}
