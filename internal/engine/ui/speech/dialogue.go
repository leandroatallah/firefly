package speech

import (
	"strings"

	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Manager handles the display of dialogue and speech bubbles.
type Manager struct {
	speeches        map[string]Speech
	activeSpeech    string
	isSpeaking      bool
	currentText     string
	lines           []string
	currentLine     int
	waitingForInput bool

	config                *config.AppConfig
	audioManager          *audio.AudioManager
	typingSounds          []string
	typingSoundIndex      int
	typingSoundLastCount  int
	typingSoundCooldown   int
	speechAudioByLine     []string
	speechAudioPlayingKey string
	defaultSpeechAudio    []string
	defaultSpeechIndex    int
	dialogueSkipEnabled   bool
}

type typingSoundPolicy interface {
	TypingSoundEnabled() bool
}

const (
	BubbleSpeechID = "bubble"
	StorySpeechID  = "story"
)

// NewManager creates a new dialogue manager.
func NewManager(s ...Speech) *Manager {
	m := &Manager{
		speeches: make(map[string]Speech),
		config:   config.Get(),
	}
	for _, s := range s {
		m.AddSpeech(s)
	}
	return m
}

func (m *Manager) AddSpeech(s Speech) {
	m.speeches[s.ID()] = s
}

func (m *Manager) SetSpeech(id string) {
	if _, ok := m.speeches[id]; ok {
		m.activeSpeech = id
	}
}

func (m *Manager) GetActiveSpeech() Speech {
	return m.speeches[m.activeSpeech]
}

func (m *Manager) SetActiveSpeech(id string) {
	m.SetSpeech(id)
}

func (m *Manager) SetAudioManager(manager *audio.AudioManager) {
	m.audioManager = manager
}

func (m *Manager) SetTypingSound(path string) {
	if path == "" {
		m.typingSounds = nil
		m.typingSoundIndex = 0
		return
	}
	m.typingSounds = []string{path}
	m.typingSoundIndex = 0
}

func (m *Manager) SetTypingSounds(paths []string) {
	m.typingSounds = append([]string{}, paths...)
	m.typingSoundIndex = 0
}

func (m *Manager) SetSpeechAudioQueue(paths []string) {
	m.speechAudioByLine = append([]string{}, paths...)
	m.speechAudioPlayingKey = ""
}

func (m *Manager) ClearSpeechAudioQueue() {
	m.speechAudioByLine = nil
	m.speechAudioPlayingKey = ""
}

func (m *Manager) SetDefaultSpeechAudio(paths []string) {
	m.defaultSpeechAudio = append([]string{}, paths...)
	m.defaultSpeechIndex = 0
}

func (m *Manager) ApplyDefaultSpeechAudio(lineCount int) {
	if len(m.defaultSpeechAudio) == 0 || lineCount <= 0 {
		m.ClearSpeechAudioQueue()
		return
	}
	queue := make([]string, 0, lineCount)
	for i := 0; i < lineCount; i++ {
		queue = append(queue, m.defaultSpeechAudio[m.defaultSpeechIndex])
		m.defaultSpeechIndex = (m.defaultSpeechIndex + 1) % len(m.defaultSpeechAudio)
	}
	m.SetSpeechAudioQueue(queue)
}

func (m *Manager) SetSpeechSkipEnabled(enabled bool) {
	m.dialogueSkipEnabled = enabled
}

// ShowMessages displays a list of messages.
func (m *Manager) ShowMessages(lines []string, position string, speed int) {
	if len(lines) == 0 {
		return
	}
	s := m.GetActiveSpeech()
	m.lines = lines
	m.currentLine = 0
	m.isSpeaking = true
	m.waitingForInput = false
	m.currentText = ""
	m.typingSoundLastCount = 0
	m.typingSoundCooldown = 0
	s.ResetText()
	s.SetPosition(position)
	if speed > 0 {
		s.SetSpeed(speed)
	} else {
		// Default speed if not specified
		s.SetSpeed(4)
	}
	s.Show()
	m.startSpeechAudioIfNeeded()
}

// IsSpeaking returns true if the dialogue manager is currently displaying a message.
func (m *Manager) IsSpeaking() bool {
	return m.isSpeaking
}

// Update updates the dialogue state. It handles input for proceeding.
func (m *Manager) Update() error {
	if !m.isSpeaking {
		return nil
	}

	s := m.GetActiveSpeech()
	if err := s.Update(); err != nil {
		return err
	}

	m.updateTypingSound(s)
	m.updateSpeechAudio()

	if !m.waitingForInput && !s.IsSpellingComplete() && m.shouldSkipTyping() {
		s.CompleteSpelling()
		m.stopSpeechAudio(false)
		m.waitingForInput = true
		return nil
	}

	if s.IsSpellingComplete() && !m.waitingForInput {
		m.waitingForInput = true
		m.stopSpeechAudio(false)
	}

	if m.waitingForInput {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.currentLine++
			if m.currentLine >= len(m.lines) {
				s.Hide()
				m.isSpeaking = false
				m.stopSpeechAudio(true)
			} else {
				if !s.IsAccumulative() {
					s.ResetText()
					m.currentText = ""
					m.typingSoundLastCount = 0
					m.typingSoundCooldown = 0
				}
				m.waitingForInput = false
				m.startSpeechAudioIfNeeded()
			}
		}
	}
	return nil
}

// Draw draws the speech bubble if it's active.
func (m *Manager) Draw(screen *ebiten.Image) {
	if !m.isSpeaking {
		return
	}

	s := m.GetActiveSpeech()
	if m.currentLine < len(m.lines) {
		s.Draw(screen, m.getCurrentMessage())
	}
}

func (m *Manager) getCurrentMessage() string {
	s := m.GetActiveSpeech()
	if s.IsAccumulative() {
		return strings.Join(m.lines[:m.currentLine+1], "\n\n")
	}
	return m.lines[m.currentLine]
}

func (m *Manager) shouldSkipTyping() bool {
	cfg := m.config
	if !m.dialogueSkipEnabled && (cfg == nil || !cfg.EnableSpeechSkip) {
		return false
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

func (m *Manager) startSpeechAudioIfNeeded() {
	if m.audioManager == nil || m.speechAudioPlayingKey != "" {
		return
	}
	if m.currentLine < 0 || m.currentLine >= len(m.speechAudioByLine) {
		return
	}
	key := m.speechAudioByLine[m.currentLine]
	if key == "" {
		return
	}
	// Check typing sound policy - if disabled, skip speech audio as well
	s := m.GetActiveSpeech()
	if policy, ok := s.(typingSoundPolicy); ok && !policy.TypingSoundEnabled() {
		return
	}
	m.speechAudioPlayingKey = key
	m.audioManager.PlaySound(m.speechAudioPlayingKey)
}

func (m *Manager) updateSpeechAudio() {
	if m.audioManager == nil || m.speechAudioPlayingKey == "" {
		return
	}
	if !m.audioManager.IsPlaying(m.speechAudioPlayingKey) {
		m.speechAudioPlayingKey = ""
	}
}

func (m *Manager) stopSpeechAudio(clearQueue bool) {
	if m.audioManager == nil {
		if clearQueue {
			m.speechAudioByLine = nil
		}
		m.speechAudioPlayingKey = ""
		return
	}
	if m.speechAudioPlayingKey != "" {
		m.audioManager.PauseMusic(m.speechAudioPlayingKey)
	}
	if clearQueue {
		m.speechAudioByLine = nil
	}
	m.speechAudioPlayingKey = ""
}

func (m *Manager) updateTypingSound(s Speech) {
	cfg := m.config
	if cfg == nil || !cfg.EnableTypingSounds || m.audioManager == nil || len(m.typingSounds) == 0 {
		return
	}
	if policy, ok := s.(typingSoundPolicy); ok && !policy.TypingSoundEnabled() {
		return
	}
	if s.IsSpellingComplete() {
		return
	}
	if m.currentLine >= len(m.lines) {
		return
	}
	if m.typingSoundCooldown > 0 {
		m.typingSoundCooldown--
		return
	}

	text := s.Text(m.getCurrentMessage())
	currentCount := len(text)
	if currentCount > m.typingSoundLastCount {
		if len(m.typingSounds) > 0 {
			path := m.typingSounds[m.typingSoundIndex%len(m.typingSounds)]
			if !m.audioManager.IsPlaying(path) {
				player := m.audioManager.PlaySound(path)
				if player != nil {
					player.SetVolume(m.audioManager.Volume() * cfg.TypingSoundVolume)
				}
			}
			m.typingSoundIndex = (m.typingSoundIndex + 1) % len(m.typingSounds)
		}
		m.typingSoundLastCount = currentCount
		cooldown := cfg.TypingSoundCooldownFrames
		if cooldown <= 0 {
			cooldown = 1
		}
		m.typingSoundCooldown = cooldown
	}
}
