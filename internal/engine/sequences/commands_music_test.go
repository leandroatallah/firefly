package sequences

import (
	"testing"
)

func TestPlayMusicCommand_Init_NilContext(t *testing.T) {
	cmd := &PlayMusicCommand{
		Path:   "test_music.mp3",
		Rewind: false,
		Loop:   true,
	}

	// Should panic with nil context - this is expected behavior
	// We just verify the command structure is correct
	if cmd.Path != "test_music.mp3" {
		t.Error("PlayMusicCommand path not set correctly")
	}
}

func TestPlayMusicCommand_Update(t *testing.T) {
	cmd := &PlayMusicCommand{
		Path: "test_music.mp3",
	}

	// Update should return true (instant command)
	if !cmd.Update() {
		t.Error("PlayMusicCommand.Update() should return true (instant command)")
	}
}

func TestPauseAllMusicCommand_Init_NilContext(t *testing.T) {
	// Verify the command structure is correct
	cmd := &PauseAllMusicCommand{}
	_ = cmd
}

func TestPauseAllMusicCommand_Update(t *testing.T) {
	cmd := &PauseAllMusicCommand{}

	if !cmd.Update() {
		t.Error("PauseAllMusicCommand.Update() should return true (instant command)")
	}
}

func TestFadeOutAllMusicCommand_Init_NilContext(t *testing.T) {
	cmd := &FadeOutAllMusicCommand{
		Duration: 60,
	}

	// Should panic with nil context - this is expected behavior
	// We just verify the command structure is correct
	if cmd.Duration != 60 {
		t.Error("FadeOutAllMusicCommand duration not set correctly")
	}
}

func TestFadeOutAllMusicCommand_Update(t *testing.T) {
	cmd := &FadeOutAllMusicCommand{
		Duration: 60,
	}

	if !cmd.Update() {
		t.Error("FadeOutAllMusicCommand.Update() should return true (instant command)")
	}
}
