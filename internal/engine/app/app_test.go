package app

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestGameUpdateAndDrawIntegration(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
	}

	ctx := &AppContext{
		Config:          cfg,
		DialogueManager: speech.NewManager(),
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
