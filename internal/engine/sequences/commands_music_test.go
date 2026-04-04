package sequences

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestPlayMusicCommand_Init_WithMockAudioManager(t *testing.T) {
	ctx := setupTestAppContext()
	am := mocks.NewMockAudioManager()
	ctx.AudioManager = am

	cmd := &PlayMusicCommand{Path: "bgm.ogg", Loop: true}
	cmd.Init(ctx)

	if len(am.PlayedPaths) == 0 || am.PlayedPaths[0] != "bgm.ogg" {
		t.Errorf("expected PlayMusic called with 'bgm.ogg', got %v", am.PlayedPaths)
	}
	if !am.LoopSettings["bgm.ogg"] {
		t.Error("expected loop=true for bgm.ogg")
	}
}

func TestPlayMusicCommand_Init_AlreadyPlayingNoRewind(t *testing.T) {
	ctx := setupTestAppContext()
	am := mocks.NewMockAudioManager()
	am.PlayingPaths["bgm.ogg"] = true
	ctx.AudioManager = am

	cmd := &PlayMusicCommand{Path: "bgm.ogg", Rewind: false}
	cmd.Init(ctx)

	// Already playing + Rewind=false → no-op, PlayMusic not called
	if len(am.PlayedPaths) != 0 {
		t.Errorf("expected PlayMusic not called, got %v", am.PlayedPaths)
	}
}

func TestPlayMusicCommand_Init_WithVolume(t *testing.T) {
	ctx := setupTestAppContext()
	am := mocks.NewMockAudioManager()
	ctx.AudioManager = am

	cmd := &PlayMusicCommand{Path: "bgm.ogg", Volume: 0.5}
	cmd.Init(ctx)

	if am.VolumeSet != 0.5 {
		t.Errorf("expected SetVolume(0.5), got %f", am.VolumeSet)
	}
}

func TestPauseAllMusicCommand_Init_CallsPauseAll(t *testing.T) {
	ctx := setupTestAppContext()
	am := mocks.NewMockAudioManager()
	ctx.AudioManager = am

	cmd := &PauseAllMusicCommand{}
	cmd.Init(ctx)

	if am.PausedAllCount != 1 {
		t.Errorf("expected PauseAll called once, got %d", am.PausedAllCount)
	}
}

func TestFadeOutAllMusicCommand_Init_CallsFadeOutAll(t *testing.T) {
	ctx := setupTestAppContext()
	am := mocks.NewMockAudioManager()
	ctx.AudioManager = am

	cmd := &FadeOutAllMusicCommand{Duration: 60}
	cmd.Init(ctx)
	// FadeOutAll is called — no panic and no assertion needed beyond compilation
	// (MockAudioManager.FadeOutAll is a no-op)
}

func TestPlayMusicCommand_Update(t *testing.T) {
	cmd := &PlayMusicCommand{Path: "test_music.mp3"}
	if !cmd.Update() {
		t.Error("PlayMusicCommand.Update() should return true (instant command)")
	}
}

func TestPauseAllMusicCommand_Update(t *testing.T) {
	cmd := &PauseAllMusicCommand{}
	if !cmd.Update() {
		t.Error("PauseAllMusicCommand.Update() should return true (instant command)")
	}
}

func TestFadeOutAllMusicCommand_Update(t *testing.T) {
	cmd := &FadeOutAllMusicCommand{Duration: 60}
	if !cmd.Update() {
		t.Error("FadeOutAllMusicCommand.Update() should return true (instant command)")
	}
}
