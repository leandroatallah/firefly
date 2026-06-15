package app

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestGameUpdateAndDrawIntegration(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
	}

	ctx := &AppContext{
		Config:          cfg,
		DialogueManager: &mocks.MockDialogueManager{},
	}

	sm := &mocks.MockSceneManager{}
	ctx.SceneManager = sm

	game := NewGame(ctx)

	if ctx.FrameCount != 0 {
		t.Fatalf("expected initial FrameCount 0, got %d", ctx.FrameCount)
	}

	if err := game.Update(); err != nil {
		t.Fatalf("unexpected error on Update: %v", err)
	}

	if ctx.FrameCount != 1 {
		t.Fatalf("expected FrameCount 1 after Update, got %d", ctx.FrameCount)
	}

	if !sm.UpdateCalled {
		t.Fatalf("expected SceneManager.Update to be called")
	}

	screen := ebiten.NewImage(cfg.ScreenWidth, cfg.ScreenHeight)
	game.Draw(screen)

	if !sm.DrawCalled {
		t.Fatalf("expected SceneManager.Draw to be called")
	}

	w, h := game.Layout(0, 0)
	if w != cfg.ScreenWidth || h != cfg.ScreenHeight {
		t.Fatalf("Layout() = (%d,%d), want (%d,%d)", w, h, cfg.ScreenWidth, cfg.ScreenHeight)
	}
}

// TestGameUpdateSlowMoAppliedGuard verifies the one-time guard that applies
// ebiten.SetTPS at the first Game.Update tick when slow-motion is configured.
//
// We intentionally do NOT assert on ebiten.CurrentTPS() because Ebitengine
// does not guarantee that CurrentTPS reflects a SetTPS call outside an active
// RunGame loop — asserting it would be flaky. The pure EffectiveTPS helper
// (see slowmo_test.go) exercises the clamp/rounding/no-op math; here we only
// verify the observable wiring on the Game struct:
//   - slowMoApplied flips to true on the first Update regardless of cfg.SlowMo
//   - the SceneManager.Update path still runs when the slow-mo branch fires
func TestGameUpdateSlowMoAppliedGuard(t *testing.T) {
	t.Run("T-G1 slowMoApplied flips on first Update when SlowMo disabled", func(t *testing.T) {
		cfg := &config.AppConfig{
			ScreenWidth:  320,
			ScreenHeight: 240,
			SlowMo:       false,
		}
		ctx := &AppContext{
			Config:          cfg,
			DialogueManager: &mocks.MockDialogueManager{},
			SceneManager:    &mocks.MockSceneManager{},
		}
		game := NewGame(ctx)

		if game.slowMoApplied {
			t.Fatalf("expected slowMoApplied=false before any Update, got true")
		}

		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error on first Update: %v", err)
		}
		if !game.slowMoApplied {
			t.Fatalf("expected slowMoApplied=true after first Update (regardless of cfg.SlowMo)")
		}

		// Second update must remain safe (idempotent guard).
		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error on second Update: %v", err)
		}
		if !game.slowMoApplied {
			t.Fatalf("expected slowMoApplied to remain true after second Update")
		}
	})

	t.Run("T-G2 SceneManager.Update still called when slow-mo branch runs", func(t *testing.T) {
		cfg := &config.AppConfig{
			ScreenWidth:  320,
			ScreenHeight: 240,
			SlowMo:       true,
			SlowMoFactor: 0.25,
		}
		sm := &mocks.MockSceneManager{}
		ctx := &AppContext{
			Config:          cfg,
			DialogueManager: &mocks.MockDialogueManager{},
			SceneManager:    sm,
		}
		game := NewGame(ctx)

		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error on Update: %v", err)
		}
		if !sm.UpdateCalled {
			t.Fatalf("expected SceneManager.Update to be called when slow-mo branch runs")
		}
		if ctx.FrameCount != 1 {
			t.Fatalf("expected FrameCount=1 after one Update, got %d", ctx.FrameCount)
		}
		if !game.slowMoApplied {
			t.Fatalf("expected slowMoApplied=true after first Update with SlowMo enabled")
		}
	})
}

func TestAppContextPhaseNavigationIntegration(t *testing.T) {
	pm := phases.NewManager()
	sceneType1 := navigation.SceneType(1)
	sceneType2 := navigation.SceneType(2)

	pm.AddPhase(phases.Phase{ID: 1, SceneType: sceneType1, NextPhaseID: 2})
	pm.AddPhase(phases.Phase{ID: 2, SceneType: sceneType2})

	if err := pm.SetCurrentPhase(1); err != nil {
		t.Fatalf("SetCurrentPhase: %v", err)
	}

	sm := &mocks.MockSceneManager{}
	ctx := &AppContext{
		PhaseManager: pm,
		SceneManager: sm,
	}

	ctx.GoToCurrentPhaseScene(nil, true)

	if sm.NavigateCalls != 1 {
		t.Fatalf("expected 1 NavigateTo call, got %d", sm.NavigateCalls)
	}
	if sm.LastSceneType != sceneType1 {
		t.Fatalf("GoToCurrentPhaseScene navigated to scene %d, want %d", sm.LastSceneType, sceneType1)
	}

	ctx.CompleteCurrentPhase(nil, true)

	if pm.CurrentPhase != 2 {
		t.Fatalf("expected CurrentPhase to be 2 after completion, got %d", pm.CurrentPhase)
	}
	if sm.NavigateCalls != 2 {
		t.Fatalf("expected 2 NavigateTo calls after completing phase, got %d", sm.NavigateCalls)
	}
	if sm.LastSceneType != sceneType2 {
		t.Fatalf("CompleteCurrentPhase navigated to scene %d, want %d", sm.LastSceneType, sceneType2)
	}
}

// TestGame_F3F4_ToggleSlowMoFastForward verifies that flipping cfg.SlowMo /
// cfg.FastForward between Update calls is reflected in the game's last-seen
// values on the next tick (the TPS dirty-check wiring).
func TestGame_F3F4_ToggleSlowMoFastForward(t *testing.T) {
	t.Run("T-FF1 lastSlowMo tracks cfg.SlowMo across ticks", func(t *testing.T) {
		cfg := &config.AppConfig{
			ScreenWidth:  320,
			ScreenHeight: 240,
			SlowMo:       false,
			SlowMoFactor: 0.5,
		}
		ctx := &AppContext{Config: cfg, SceneManager: &mocks.MockSceneManager{}, DialogueManager: &mocks.MockDialogueManager{}}
		game := NewGame(ctx)

		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if game.lastSlowMo != false {
			t.Fatalf("lastSlowMo = %v, want false", game.lastSlowMo)
		}

		cfg.SlowMo = true
		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if game.lastSlowMo != true {
			t.Fatalf("lastSlowMo = %v, want true after toggle", game.lastSlowMo)
		}
	})

	t.Run("T-FF2 lastFastForward tracks cfg.FastForward across ticks", func(t *testing.T) {
		cfg := &config.AppConfig{
			ScreenWidth:       320,
			ScreenHeight:      240,
			FastForward:       false,
			FastForwardFactor: 2.0,
		}
		ctx := &AppContext{Config: cfg, SceneManager: &mocks.MockSceneManager{}, DialogueManager: &mocks.MockDialogueManager{}}
		game := NewGame(ctx)

		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if game.lastFastForward != false {
			t.Fatalf("lastFastForward = %v, want false", game.lastFastForward)
		}

		cfg.FastForward = true
		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if game.lastFastForward != true {
			t.Fatalf("lastFastForward = %v, want true after toggle", game.lastFastForward)
		}
	})
}

// overlay is open, Game.Update skips both SceneManager.Update and
// DialogueManager.Update, but still advances FrameCount.
//
// AC-8: overlay open → SceneManager.Update and DialogueManager.Update not called.
func TestGame_OverlayOpenSuppressesSceneUpdate(t *testing.T) {
	t.Run("T-G1 overlay open suppresses scene and dialogue update", func(t *testing.T) {
		cfg := &config.AppConfig{
			ScreenWidth:  320,
			ScreenHeight: 240,
		}

		sm := &mocks.MockSceneManager{}
		dm := &mocks.MockDialogueManager{}
		ctx := &AppContext{
			Config:          cfg,
			SceneManager:    sm,
			DialogueManager: dm,
		}

		game := NewGame(ctx)

		// Open the overlay via the test-friendly accessor.
		overlay := game.DebugOverlay()
		if overlay == nil {
			t.Fatalf("DebugOverlay() returned nil, want non-nil overlay")
		}
		overlay.Open()

		startFrame := ctx.FrameCount
		if err := game.Update(); err != nil {
			t.Fatalf("unexpected error on Update: %v", err)
		}

		if sm.UpdateCalled {
			t.Fatalf("SceneManager.Update was called while overlay open, want suppressed")
		}
		if dm.UpdateCalls != 0 {
			t.Fatalf("DialogueManager.Update was called %d times while overlay open, want 0", dm.UpdateCalls)
		}
		if ctx.FrameCount != startFrame+1 {
			t.Fatalf("FrameCount = %d, want %d (frame should still advance)", ctx.FrameCount, startFrame+1)
		}
	})
}
