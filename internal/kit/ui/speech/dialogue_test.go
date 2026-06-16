package speech

import (
	"image/color"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/dialogue"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// TestManager_ImplementsDialogueContract is the canonical Red test for
// story 052-kit-ui-split. It does not compile until:
//   - the package internal/engine/contracts/dialogue exists with a
//     Manager interface;
//   - the package internal/kit/ui/speech exists with a Manager type
//     whose method set satisfies that interface.
//
// Once the relocation is performed, the compile-time assertion guarantees
// that the kit-side orchestrator continues to satisfy the engine-side
// contract.
func TestManager_ImplementsDialogueContract(t *testing.T) {
	var _ dialogue.Manager = (*Manager)(nil)
}

// TestManager mirrors the behavioural assertions previously held by
// internal/engine/ui/speech/dialogue_test.go. After the move, identical
// behaviour must be observable from the new package path.
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

	// Empty lines is a no-op — no Show, no state change.
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

	// Speed of 0 falls back to the engine-default of 4.
	m.ShowMessages(lines, "bottom", 0)
	if s.setSpeedCalled != 4 {
		t.Errorf("Expected default speed 4, got %d", s.setSpeedCalled)
	}

	if s.resetTextCalled != 2 {
		t.Errorf("Expected ResetText to be called twice, got %d", s.resetTextCalled)
	}

	if !s.visible {
		t.Error("Expected speech to be visible")
	}

	if s.showCalled < 1 {
		t.Errorf("Expected Show to be invoked on the active speech, got %d calls", s.showCalled)
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

// typingSpeech is a Speech implementation that advances one character per
// Update call, enabling deterministic typing-sound rotation tests.
type typingSpeech struct {
	id      string
	visible bool
	spelled int
}

func (t *typingSpeech) ID() string                             { return t.id }
func (t *typingSpeech) Show()                                  { t.visible = true }
func (t *typingSpeech) Hide()                                  { t.visible = false }
func (t *typingSpeech) Visible() bool                          { return t.visible }
func (t *typingSpeech) SetIndicator(*ebiten.Image)             {}
func (t *typingSpeech) ResetText()                             { t.spelled = 0 }
func (t *typingSpeech) SetID(id string)                        { t.id = id }
func (t *typingSpeech) SetSpellingDelay(d int)                 {}
func (t *typingSpeech) IsSpellingComplete() bool               { return false }
func (t *typingSpeech) CompleteSpelling()                      {}
func (t *typingSpeech) Count() int                             { return 0 }
func (t *typingSpeech) Update() error                          { t.spelled++; return nil }
func (t *typingSpeech) Draw(screen *ebiten.Image, text string) {}
func (t *typingSpeech) SetPosition(pos string)                 {}
func (t *typingSpeech) SetSpeed(speed int)                     {}
func (t *typingSpeech) SetColor(c color.Color)                 {}
func (t *typingSpeech) Color() color.Color                     { return color.Black }
func (t *typingSpeech) SetSkipFlash(frames int)                {}
func (t *typingSpeech) IsAccumulative() bool                   { return false }
func (t *typingSpeech) SetAccumulative(bool)                   {}
func (t *typingSpeech) Text(msg string) string {
	if t.spelled <= 0 {
		return ""
	}
	if t.spelled >= len(msg) {
		return msg
	}
	return msg[:t.spelled]
}

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

	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.typingSoundIndex != 1 {
		t.Fatalf("expected typingSoundIndex 1, got %d", m.typingSoundIndex)
	}

	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.typingSoundIndex != 2 {
		t.Fatalf("expected typingSoundIndex 2, got %d", m.typingSoundIndex)
	}
}

func TestManager_SetTypingSound(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)

	m.SetTypingSound("path/to/click.ogg")
	if len(m.typingSounds) != 1 || m.typingSounds[0] != "path/to/click.ogg" {
		t.Fatalf("expected single typing sound set, got %v", m.typingSounds)
	}

	// Empty path clears the configured sounds.
	m.SetTypingSound("")
	if len(m.typingSounds) != 0 {
		t.Fatalf("expected typing sounds to be cleared, got %v", m.typingSounds)
	}
}

func TestManager_ClearSpeechAudioQueue(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetSpeechAudioQueue([]string{"a", "b"})

	m.ClearSpeechAudioQueue()

	if len(m.speechAudioByLine) != 0 {
		t.Fatalf("expected queue to be empty, got %v", m.speechAudioByLine)
	}
	if m.speechAudioPlayingKey != "" {
		t.Fatalf("expected playing key to be cleared, got %q", m.speechAudioPlayingKey)
	}
}

func TestManager_ApplyDefaultSpeechAudio_NoDefault_ClearsQueue(t *testing.T) {
	config.Set(&config.AppConfig{})
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetSpeechAudioQueue([]string{"x"})

	// With no default audio configured, applying it must clear the queue.
	m.ApplyDefaultSpeechAudio(3)

	if len(m.speechAudioByLine) != 0 {
		t.Fatalf("expected queue cleared when no default audio, got %v", m.speechAudioByLine)
	}
}

func TestManager_SetSpeechSkipEnabled(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)

	m.SetSpeechSkipEnabled(true)
	if !m.dialogueSkipEnabled {
		t.Fatal("expected dialogueSkipEnabled to be true after SetSpeechSkipEnabled(true)")
	}

	m.SetSpeechSkipEnabled(false)
	if m.dialogueSkipEnabled {
		t.Fatal("expected dialogueSkipEnabled to be false after SetSpeechSkipEnabled(false)")
	}
}

func TestManager_StopSpeechAudio_NoManager_ClearsQueue(t *testing.T) {
	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")
	m.SetSpeechAudioQueue([]string{"a", "b"})

	// With no audio manager, completing the dialogue still clears the queue
	// when stopSpeechAudio is invoked via the spelling-complete + Update path.
	m.ShowMessages([]string{"hi"}, "bottom", 0)
	s.spellingComplete = true
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Spelling-complete only clears the playing key, not the queue.
	if m.speechAudioPlayingKey != "" {
		t.Fatalf("expected playing key cleared, got %q", m.speechAudioPlayingKey)
	}
}

func TestManager_GetActiveSpeech_ReturnsActive(t *testing.T) {
	s1 := &mockSpeech{id: "a"}
	s2 := &mockSpeech{id: "b"}
	m := NewManager(s1, s2)
	m.SetActiveSpeech("b")

	if got := m.GetActiveSpeech(); got != s2 {
		t.Fatalf("expected GetActiveSpeech to return s2, got %v", got)
	}
}

func TestManager_Update_ConfirmAdvancesLine(t *testing.T) {
	saved := input.CommandsReader
	t.Cleanup(func() { input.CommandsReader = saved })

	confirm := false
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Confirm: confirm}
	}

	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")

	m.ShowMessages([]string{"line1", "line2"}, "bottom", 0)
	s.spellingComplete = true

	// First Update transitions to waitingForInput and stops audio.
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.waitingForInput {
		t.Fatal("expected waitingForInput after spelling complete")
	}

	// Confirm pressed: should advance to next line and exit waitingForInput.
	confirm = true
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.currentLine != 1 {
		t.Fatalf("expected currentLine to advance to 1, got %d", m.currentLine)
	}
	if m.waitingForInput {
		t.Fatal("expected waitingForInput to be false after confirm advance")
	}
	if !m.IsSpeaking() {
		t.Fatal("expected manager to still be speaking on remaining line")
	}
}

func TestManager_Update_ConfirmHidesAtEnd(t *testing.T) {
	saved := input.CommandsReader
	t.Cleanup(func() { input.CommandsReader = saved })

	confirm := false
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Confirm: confirm}
	}

	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetActiveSpeech("test")

	m.ShowMessages([]string{"only"}, "bottom", 0)
	s.spellingComplete = true

	// Move to waitingForInput.
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Confirm pressed on the last line: must hide and stop speaking.
	confirm = true
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.IsSpeaking() {
		t.Fatal("expected IsSpeaking to be false after confirming final line")
	}
	if s.hideCalled == 0 {
		t.Fatal("expected Hide to be called on active speech at end of dialogue")
	}
}

func TestManager_Update_SkipTyping(t *testing.T) {
	saved := input.CommandsReader
	t.Cleanup(func() { input.CommandsReader = saved })
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Confirm: true}
	}

	s := &mockSpeech{id: "test"}
	m := NewManager(s)
	m.SetSpeechSkipEnabled(true)
	m.SetActiveSpeech("test")
	m.ShowMessages([]string{"long line"}, "bottom", 0)

	// Spelling not complete yet, but Confirm is held and skip is enabled.
	if err := m.Update(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.spellingComplete {
		t.Fatal("expected CompleteSpelling to be invoked when skip is enabled and Confirm is held")
	}
	if !m.waitingForInput {
		t.Fatal("expected waitingForInput after skip-completing the line")
	}
}
