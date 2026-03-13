package mocks

import (
	"time"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
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

func (m *MockAudioManager) IsPlaying(path string) bool {
	return m.PlayingPaths[path]
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

// MockVFXManager implements vfx.Manager for testing
type MockVFXManager struct {
	SpawnedTexts        []string
	SpawnedTextPositions []struct {
		X, Y float64
		Text string
	}
	SpawnedAboveTargets []string
}

func NewMockVFXManager() *MockVFXManager {
	return &MockVFXManager{
		SpawnedTexts:        make([]string, 0),
		SpawnedTextPositions: make([]struct {
			X, Y float64
			Text string
		}, 0),
		SpawnedAboveTargets: make([]string, 0),
	}
}

func (m *MockVFXManager) SpawnFloatingText(text string, x, y float64, duration int) {
	m.SpawnedTexts = append(m.SpawnedTexts, text)
	m.SpawnedTextPositions = append(m.SpawnedTextPositions, struct {
		X, Y float64
		Text string
	}{X: x, Y: y, Text: text})
}

func (m *MockVFXManager) SpawnFloatingTextAbove(actor actors.ActorEntity, text string, duration int) {
	m.SpawnedAboveTargets = append(m.SpawnedAboveTargets, actor.ID())
	m.SpawnedTexts = append(m.SpawnedTexts, text)
}

func (m *MockVFXManager) Update() error {
	return nil
}

func (m *MockVFXManager) Draw(screen interface{}) {}

func (m *MockVFXManager) Spawn(vfxType string, x, y float64) {}

func (m *MockVFXManager) SetCamera(cam interface{}) {}
