package mocks

import (
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/dialogue"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
	"github.com/hajimehoshi/ebiten/v2"
)

// MockDialogueManager is a mock implementation of dialogue.Manager for testing.
type MockDialogueManager struct {
	UpdateFunc                  func() error
	DrawFunc                    func(screen *ebiten.Image)
	IsSpeakingFunc              func() bool
	StopFunc                    func()
	AddSpeechFunc               func(s speech.Speech)
	SetSpeechFunc               func(id string)
	SetActiveSpeechFunc         func(id string)
	GetActiveSpeechFunc         func() speech.Speech
	SetAudioManagerFunc         func(m *audio.AudioManager)
	SetTypingSoundFunc          func(path string)
	SetTypingSoundsFunc         func(paths []string)
	SetSpeechAudioQueueFunc     func(paths []string)
	ClearSpeechAudioQueueFunc   func()
	SetDefaultSpeechAudioFunc   func(paths []string)
	ApplyDefaultSpeechAudioFunc func(lineCount int)
	SetSpeechSkipEnabledFunc    func(enabled bool)
	SetPlayerAdvanceEnabledFunc func(enabled bool)
	ShowMessagesFunc            func(lines []string, position string, speed int)

	UpdateCalls                  int
	DrawCalls                    int
	IsSpeakingCalls              int
	StopCalls                    int
	AddSpeechCalls               int
	SetSpeechCalls               int
	SetActiveSpeechCalls         int
	GetActiveSpeechCalls         int
	SetAudioManagerCalls         int
	SetTypingSoundCalls          int
	SetTypingSoundsCalls         int
	SetSpeechAudioQueueCalls     int
	ClearSpeechAudioQueueCalls   int
	SetDefaultSpeechAudioCalls   int
	ApplyDefaultSpeechAudioCalls int
	SetSpeechSkipEnabledCalls    int
	SetPlayerAdvanceEnabledCalls int
	ShowMessagesCalls            int

	LastAddedSpeech speech.Speech
	LastSpeechID    string
	LastPosition    string
	LastSpeed       int
	LastLines       []string
	LastAudioPaths  []string
	LastAudioMgr    *audio.AudioManager
	IsSpeakingValue bool
	ActiveSpeech    speech.Speech
}

func (m *MockDialogueManager) Update() error {
	m.UpdateCalls++
	if m.UpdateFunc != nil {
		return m.UpdateFunc()
	}
	return nil
}

func (m *MockDialogueManager) Draw(screen *ebiten.Image) {
	m.DrawCalls++
	if m.DrawFunc != nil {
		m.DrawFunc(screen)
	}
}

func (m *MockDialogueManager) IsSpeaking() bool {
	m.IsSpeakingCalls++
	if m.IsSpeakingFunc != nil {
		return m.IsSpeakingFunc()
	}
	return m.IsSpeakingValue
}

func (m *MockDialogueManager) Stop() {
	m.StopCalls++
	if m.StopFunc != nil {
		m.StopFunc()
	}
}

func (m *MockDialogueManager) AddSpeech(s speech.Speech) {
	m.AddSpeechCalls++
	m.LastAddedSpeech = s
	if m.AddSpeechFunc != nil {
		m.AddSpeechFunc(s)
	}
}

func (m *MockDialogueManager) SetSpeech(id string) {
	m.SetSpeechCalls++
	m.LastSpeechID = id
	if m.SetSpeechFunc != nil {
		m.SetSpeechFunc(id)
	}
}

func (m *MockDialogueManager) SetActiveSpeech(id string) {
	m.SetActiveSpeechCalls++
	m.LastSpeechID = id
	if m.SetActiveSpeechFunc != nil {
		m.SetActiveSpeechFunc(id)
	}
}

func (m *MockDialogueManager) GetActiveSpeech() speech.Speech {
	m.GetActiveSpeechCalls++
	if m.GetActiveSpeechFunc != nil {
		return m.GetActiveSpeechFunc()
	}
	return m.ActiveSpeech
}

func (m *MockDialogueManager) SetAudioManager(mgr *audio.AudioManager) {
	m.SetAudioManagerCalls++
	m.LastAudioMgr = mgr
	if m.SetAudioManagerFunc != nil {
		m.SetAudioManagerFunc(mgr)
	}
}

func (m *MockDialogueManager) SetTypingSound(path string) {
	m.SetTypingSoundCalls++
	if m.SetTypingSoundFunc != nil {
		m.SetTypingSoundFunc(path)
	}
}

func (m *MockDialogueManager) SetTypingSounds(paths []string) {
	m.SetTypingSoundsCalls++
	m.LastAudioPaths = paths
	if m.SetTypingSoundsFunc != nil {
		m.SetTypingSoundsFunc(paths)
	}
}

func (m *MockDialogueManager) SetSpeechAudioQueue(paths []string) {
	m.SetSpeechAudioQueueCalls++
	m.LastAudioPaths = paths
	if m.SetSpeechAudioQueueFunc != nil {
		m.SetSpeechAudioQueueFunc(paths)
	}
}

func (m *MockDialogueManager) ClearSpeechAudioQueue() {
	m.ClearSpeechAudioQueueCalls++
	if m.ClearSpeechAudioQueueFunc != nil {
		m.ClearSpeechAudioQueueFunc()
	}
}

func (m *MockDialogueManager) SetDefaultSpeechAudio(paths []string) {
	m.SetDefaultSpeechAudioCalls++
	m.LastAudioPaths = paths
	if m.SetDefaultSpeechAudioFunc != nil {
		m.SetDefaultSpeechAudioFunc(paths)
	}
}

func (m *MockDialogueManager) ApplyDefaultSpeechAudio(lineCount int) {
	m.ApplyDefaultSpeechAudioCalls++
	if m.ApplyDefaultSpeechAudioFunc != nil {
		m.ApplyDefaultSpeechAudioFunc(lineCount)
	}
}

func (m *MockDialogueManager) SetSpeechSkipEnabled(enabled bool) {
	m.SetSpeechSkipEnabledCalls++
	if m.SetSpeechSkipEnabledFunc != nil {
		m.SetSpeechSkipEnabledFunc(enabled)
	}
}

func (m *MockDialogueManager) SetPlayerAdvanceEnabled(enabled bool) {
	m.SetPlayerAdvanceEnabledCalls++
	if m.SetPlayerAdvanceEnabledFunc != nil {
		m.SetPlayerAdvanceEnabledFunc(enabled)
	}
}

func (m *MockDialogueManager) ShowMessages(lines []string, position string, speed int) {
	m.ShowMessagesCalls++
	m.LastLines = lines
	m.LastPosition = position
	m.LastSpeed = speed
	if m.ShowMessagesFunc != nil {
		m.ShowMessagesFunc(lines, position, speed)
	}
}

// Ensure MockDialogueManager implements dialogue.Manager
var _ dialogue.Manager = (*MockDialogueManager)(nil)
