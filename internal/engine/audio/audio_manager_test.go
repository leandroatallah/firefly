package audio

import (
	"encoding/binary"
	"testing"
	"time"
)

func createMinimalWAV() []byte {
	buf := make([]byte, 44+2)
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], 36+2)
	copy(buf[8:12], "WAVE")
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16)
	binary.LittleEndian.PutUint16(buf[20:22], 1) // PCM
	binary.LittleEndian.PutUint16(buf[22:24], 1) // Mono
	binary.LittleEndian.PutUint32(buf[24:28], 44100)
	binary.LittleEndian.PutUint32(buf[28:32], 44100*2)
	binary.LittleEndian.PutUint16(buf[32:34], 2)
	binary.LittleEndian.PutUint16(buf[34:36], 16)
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], 2)
	return buf
}

func TestAudioManagerAdd(t *testing.T) {
	am := getTestAudioManager()

	// Test unsupported format
	am.Add("test.txt", []byte("dummy"))
	if _, ok := am.audioPlayers["test.txt"]; ok {
		t.Error("expected test.txt to not be added")
	}

	// Test invalid wav
	am.Add("invalid.wav", []byte("dummy"))
	if _, ok := am.audioPlayers["invalid.wav"]; ok {
		t.Error("expected invalid.wav to not be added")
	}

	// Test valid wav
	wavData := createMinimalWAV()
	am.Add("test.wav", wavData)
	if _, ok := am.audioPlayers["test.wav"]; !ok {
		t.Error("expected test.wav to be added")
	}
}

func TestAudioManagerPlay(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_play.wav", wavData)

	// Test PlayMusic
	p := am.PlayMusic("test_play.wav", false)
	if p == nil {
		t.Fatal("expected player for test_play.wav")
	}
	if !am.IsPlaying("test_play.wav") {
		t.Error("expected test_play.wav to be playing")
	}
	if !am.IsPlayingSomething() {
		t.Error("expected IsPlayingSomething to be true")
	}

	// Test PauseMusic
	am.PauseMusic("test_play.wav")
	if am.IsPlaying("test_play.wav") {
		t.Error("expected test_play.wav to be paused")
	}

	// Test PlaySound
	p = am.PlaySound("test_play.wav")
	if p == nil {
		t.Fatal("expected player for test_play.wav")
	}

	// Test PauseAll
	am.PauseAll()
	if am.IsPlayingSomething() {
		t.Error("expected IsPlayingSomething to be false after PauseAll")
	}

	// Test missing
	if am.IsPlaying("missing.wav") {
		t.Error("expected missing.wav to not be playing")
	}
	am.PauseMusic("missing.wav") // should not panic
}

func TestAudioManagerPlayMusic_Loop(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_loop.wav", wavData)

	// Test PlayMusic with loop
	p := am.PlayMusic("test_loop.wav", true)
	if p == nil {
		t.Fatal("expected player for test_loop.wav")
	}

	// Wait a bit to ensure goroutine starts
	time.Sleep(50 * time.Millisecond)

	if !am.IsPlaying("test_loop.wav") {
		t.Error("expected test_loop.wav to be playing")
	}

	am.PauseMusic("test_loop.wav")
}

func TestAudioManagerFadeOut_AlreadyZero(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_zero.wav", wavData)

	p := am.PlayMusic("test_zero.wav", false)
	p.SetVolume(0)

	am.FadeOut("test_zero.wav", 50*time.Millisecond)
	// Should return early
}

func TestAudioManagerFadeOut_Complete(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_fade_complete.wav", wavData)

	am.PlayMusic("test_fade_complete.wav", false)
	am.FadeOut("test_fade_complete.wav", 50*time.Millisecond)

	// Wait enough for the ticker to trigger and duration to pass
	time.Sleep(200 * time.Millisecond)

	if am.IsPlaying("test_fade_complete.wav") {
		t.Error("expected test_fade_complete.wav to be stopped after fade out completion")
	}
}

func TestAudioManagerFadeOut_CancelExisting(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_cancel.wav", wavData)

	am.PlayMusic("test_cancel.wav", false)
	am.FadeOut("test_cancel.wav", 100*time.Millisecond)

	// FadeOut again should cancel the first one
	am.FadeOut("test_cancel.wav", 100*time.Millisecond)

	// FadeOutAll should cancel it too
	am.FadeOutAll(100 * time.Millisecond)

	// PlayMusic should cancel it too
	am.PlayMusic("test_cancel.wav", false)
}

func TestAudioManagerAdd_Formats(t *testing.T) {
	am := getTestAudioManager()

	// Test invalid mp3
	am.Add("test.mp3", []byte("invalid"))

	// Test invalid ogg
	am.Add("test.ogg", []byte("invalid"))
}

func TestAudioManagerFadeOutAll_AlreadyZero(t *testing.T) {
	am := getTestAudioManager()
	oldVolume := am.Volume()
	am.SetVolume(0)
	defer am.SetVolume(oldVolume)

	am.FadeOutAll(50 * time.Millisecond)
	// Should return early
}

func TestAudioManagerPlaySound_Missing(t *testing.T) {
	am := getTestAudioManager()
	if am.PlaySound("missing_sound.wav") != nil {
		t.Error("expected nil player for missing sound")
	}
}

func TestAudioManagerFadeOutAll_CancelExisting(t *testing.T) {
	am := getTestAudioManager()
	am.FadeOutAll(100 * time.Millisecond)
	am.FadeOutAll(100 * time.Millisecond) // Should cancel the first one
}

func TestAudioManagerNoSound(t *testing.T) {
	am := getTestAudioManager()
	oldNoSound := am.noSound
	am.noSound = true
	defer func() { am.noSound = oldNoSound }()

	if am.PlayMusic("any", false) != nil {
		t.Error("expected nil player when noSound is true")
	}
	if am.PlaySound("any") != nil {
		t.Error("expected nil player when noSound is true")
	}

	am.SetVolume(0.5)
	am.FadeOutAll(time.Second)
	am.FadeOut("any", time.Second)
}

func TestAudioManagerPauseResume(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_pause.wav", wavData)

	p := am.PlayMusic("test_pause.wav", true)
	if p == nil {
		t.Fatal("expected player for test_pause.wav")
	}

	if !am.IsPlaying("test_pause.wav") {
		t.Error("expected test_pause.wav to be playing")
	}

	am.PauseMusic("test_pause.wav")
	if am.IsPlaying("test_pause.wav") {
		t.Error("expected test_pause.wav to be paused")
	}

	am.ResumeMusic("test_pause.wav")
	if !am.IsPlaying("test_pause.wav") {
		t.Error("expected test_pause.wav to be playing after resume")
	}
}

func TestAudioManagerCurrentTrack(t *testing.T) {
	am := getTestAudioManager()
	wavData := createMinimalWAV()
	am.Add("test_track.wav", wavData)

	am.PlayMusic("test_track.wav", true)

	am.PauseCurrentMusic()
	if am.IsPlaying("test_track.wav") {
		t.Error("expected test_track.wav to be paused after PauseCurrentMusic")
	}

	am.ResumeCurrentMusic()
	if !am.IsPlaying("test_track.wav") {
		t.Error("expected test_track.wav to be playing after ResumeCurrentMusic")
	}
}
