package mocks

import (
	"time"
)

// MockAudioManager implements audio.Manager for testing
type MockAudioManager struct {
	PlayedPaths    []string
	PausedAllCount int
	VolumeSet      float64
	PlayingPaths   map[string]bool
	LoopSettings   map[string]bool
}

func NewMockAudioManager() *MockAudioManager {
	return &MockAudioManager{
		PlayedPaths:  make([]string, 0),
		PlayingPaths: make(map[string]bool),
		LoopSettings: make(map[string]bool),
	}
}

func (m *MockAudioManager) PlayMusic(path string, loop bool) {
	m.PlayedPaths = append(m.PlayedPaths, path)
	m.PlayingPaths[path] = true
	m.LoopSettings[path] = loop
}

func (m *MockAudioManager) PlaySound(path string) {
	m.PlayedPaths = append(m.PlayedPaths, path)
}

func (m *MockAudioManager) PlaySoundAtVolume(path string, _ float64) {
	m.PlayedPaths = append(m.PlayedPaths, path)
}

func (m *MockAudioManager) PauseCurrentMusic() {}

func (m *MockAudioManager) ResumeCurrentMusic() {}

func (m *MockAudioManager) FadeOutCurrentTrack(_ time.Duration) {}

func (m *MockAudioManager) IsPlaying(path string) bool {
	return m.PlayingPaths[path]
}

func (m *MockAudioManager) IsPaused(path string) bool {
	return false
}

func (m *MockAudioManager) SetVolume(volume float64) {
	m.VolumeSet = volume
}

func (m *MockAudioManager) PauseAll() {
	m.PausedAllCount++
}

func (m *MockAudioManager) FadeOutAll(duration time.Duration) {
	// Mock implementation - no-op
}

func (m *MockAudioManager) Stop(path string) {
	delete(m.PlayingPaths, path)
}

func (m *MockAudioManager) StopAll() {
	m.PlayingPaths = make(map[string]bool)
}

// MockDialogueManager implements speech.Manager for testing
type MockDialogueManager struct {
	ActiveSpeechID    string
	ShownLines        []string
	ShownPosition     string
	ShownSpeed        int
	IsSpeakingVal     bool
	MessagesDisplayed int
}

func NewMockDialogueManager() *MockDialogueManager {
	return &MockDialogueManager{
		IsSpeakingVal: false,
	}
}

func (m *MockDialogueManager) SetActiveSpeech(id string) {
	m.ActiveSpeechID = id
}

func (m *MockDialogueManager) ShowMessages(lines []string, position string, speed int) {
	m.ShownLines = lines
	m.ShownPosition = position
	m.ShownSpeed = speed
	m.MessagesDisplayed++
	m.IsSpeakingVal = true
}

func (m *MockDialogueManager) IsSpeaking() bool {
	return m.IsSpeakingVal
}

func (m *MockDialogueManager) Update() error {
	return nil
}

func (m *MockDialogueManager) Draw(screen interface{}) {}

func (m *MockDialogueManager) SetSpeaking(val bool) {
	m.IsSpeakingVal = val
}
