package speech

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
)

func TestManager(t *testing.T) {
	s1 := &mockSpeech{id: "speech1"}
	s2 := &mockSpeech{id: "speech2"}

	m := NewManager(s1, s2)

	if len(m.speeches) != 2 {
		t.Errorf("Expected 2 speeches, got %d", len(m.speeches))
	}

	m.SetActiveSpeech("speech1")
	if m.activeSpeech != "speech1" {
		t.Errorf("Expected active speech to be speech1, got %s", m.activeSpeech)
	}

	m.SetSpeech("speech2")
	if m.activeSpeech != "speech2" {
		t.Errorf("Expected active speech to be speech2, got %s", m.activeSpeech)
	}

	m.SetSpeech("non-existent")
	if m.activeSpeech != "speech2" {
		t.Errorf("Expected active speech to remain speech2, got %s", m.activeSpeech)
	}
}

func TestManager_ShowMessages(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")

	// Empty lines
	m.ShowMessages([]string{}, "top", 10)
	if m.IsSpeaking() {
		t.Error("Expected IsSpeaking to be false for empty lines")
	}

	lines := []string{"Line 1", "Line 2"}
	m.ShowMessages(lines, "top", 10)

	if !m.IsSpeaking() {
		t.Error("Expected IsSpeaking to be true")
	}

	if s.setPositionCalled != "top" {
		t.Errorf("Expected position top, got %s", s.setPositionCalled)
	}

	if s.setSpeedCalled != 10 {
		t.Errorf("Expected speed 10, got %d", s.setSpeedCalled)
	}

	// Default speed
	m.ShowMessages(lines, "bottom", 0)
	if s.setSpeedCalled != 4 {
		t.Errorf("Expected default speed 4, got %d", s.setSpeedCalled)
	}

	if s.resetTextCalled != 2 { // 2 from non-empty ShowMessages calls
		t.Errorf("Expected ResetText to be called twice, got %d", s.resetTextCalled)
	}

	if !s.visible {
		t.Error("Expected speech to be visible")
	}
}

func TestManager_Update(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")

	// Not speaking, Update should do nothing
	err := m.Update()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if s.updateCalled != 0 {
		t.Error("Expected Update not to be called on speech when not speaking")
	}

	lines := []string{"Line 1"}
	m.ShowMessages(lines, "bottom", 0)

	// Update while speaking
	err = m.Update()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if s.updateCalled != 1 {
		t.Errorf("Expected Update to be called on speech, got %d", s.updateCalled)
	}

	// Complete spelling
	s.spellingComplete = true
	err = m.Update()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !m.waitingForInput {
		t.Error("Expected waitingForInput to be true after spelling complete")
	}

	// We can't easily test inpututil.IsKeyJustPressed(ebiten.KeyEnter) without a real game loop or mocking input.
	// But we can manually set the state if we want to test the logic after enter is pressed.
}

func TestManager_Draw(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")

	// Not speaking, Draw should do nothing
	m.Draw(nil)
	if s.drawCalled != 0 {
		t.Error("Expected Draw not to be called on speech when not speaking")
	}

	m.ShowMessages([]string{"Hello"}, "bottom", 0)
	m.Draw(nil)
	if s.drawCalled != 1 {
		t.Errorf("Expected Draw to be called on speech, got %d", s.drawCalled)
	}
}

type typingSpeech struct {
	id               string
	visible          bool
	spellingComplete bool
	spelled          int
}

func (t *typingSpeech) ID() string { return t.id }
func (t *typingSpeech) Show() { t.visible = true }
func (t *typingSpeech) Hide() { t.visible = false }
func (t *typingSpeech) Visible() bool { return t.visible }
func (t *typingSpeech) Text(msg string) string {
	if t.spelled <= 0 {
		return ""
	}
	if t.spelled >= len(msg) {
		return msg
	}
	return msg[:t.spelled]
}
func (t *typingSpeech) ResetText() { t.spelled = 0; t.spellingComplete = false }
func (t *typingSpeech) SetID(id string) { t.id = id }
func (t *typingSpeech) SetSpellingDelay(d int) {}
func (t *typingSpeech) IsSpellingComplete() bool { return t.spellingComplete }
func (t *typingSpeech) CompleteSpelling() { t.spellingComplete = true }
func (t *typingSpeech) Count() int { return 0 }
func (t *typingSpeech) Update() error {
	t.spelled++
	return nil
}
func (t *typingSpeech) Draw(screen *ebiten.Image, text string) {}
func (t *typingSpeech) SetPosition(pos string) {}
func (t *typingSpeech) SetSpeed(speed int) {}
func (t *typingSpeech) SetColor(c color.Color) {}
func (t *typingSpeech) Color() color.Color { return color.Black }
func (t *typingSpeech) SetSkipFlash(frames int) {}
func (t *typingSpeech) IsAccumulative() bool { return false }
func (t *typingSpeech) SetAccumulative(bool) {}

func TestManager_ApplyDefaultSpeechAudio_Rotates(t *testing.T) {
	config.Set(&config.AppConfig{})
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetDefaultSpeechAudio([]string{"a", "b"})

	m.ApplyDefaultSpeechAudio(3)

	want := []string{"a", "b", "a"}
	if len(m.speechAudioByLine) != len(want) {
		t.Fatalf("expected %d entries, got %d", len(want), len(m.speechAudioByLine))
	}
	for i := range want {
		if m.speechAudioByLine[i] != want[i] {
			t.Fatalf("expected %s at %d, got %s", want[i], i, m.speechAudioByLine[i])
		}
	}
}

func TestManager_StartSpeechAudio_UsesLineAudio(t *testing.T) {
	config.Set(&config.AppConfig{})
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")
	m.SetAudioManager(&audio.AudioManager{})
	m.SetSpeechAudioQueue([]string{"assets/audio/bleeps/bleep001.ogg"})

	m.ShowMessages([]string{"Line 1"}, "bottom", 0)

	if m.speechAudioPlayingKey != "assets/audio/bleeps/bleep001.ogg" {
		t.Fatalf("expected speech audio key to be set, got %s", m.speechAudioPlayingKey)
	}
}

func TestManager_TypingSound_Rotates(t *testing.T) {
	config.Set(&config.AppConfig{
		EnableTypingSounds:        true,
		TypingSoundCooldownFrames: 1,
		TypingSoundVolume:         1,
	})
	s := &typingSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")
	m.SetAudioManager(&audio.AudioManager{})
	m.SetTypingSounds([]string{"a", "b", "c"})

	m.ShowMessages([]string{"abcd"}, "bottom", 0)

	_ = m.Update()
	if m.typingSoundIndex != 1 {
		t.Fatalf("expected typingSoundIndex 1, got %d", m.typingSoundIndex)
	}

	_ = m.Update()
	_ = m.Update()
	if m.typingSoundIndex != 2 {
		t.Fatalf("expected typingSoundIndex 2, got %d", m.typingSoundIndex)
	}
}
