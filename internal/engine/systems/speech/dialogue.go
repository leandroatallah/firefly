package speech

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
)

// Manager handles the display of dialogue and speech bubbles.
type Manager struct {
	speech          Speech
	isSpeaking      bool
	currentText     string
	lines           []string
	currentLine     int
	waitingForInput bool
}

// NewManager creates a new dialogue manager.
func NewManager(speech Speech) *Manager {
	return &Manager{speech: speech}
}

// ShowMessages displays a list of messages.
func (m *Manager) ShowMessages(lines []string) {
	if len(lines) == 0 {
		return
	}
	m.lines = lines
	m.currentLine = 0
	m.isSpeaking = true
	m.waitingForInput = false
	m.speech.ResetText()
	m.speech.Show()
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

	if err := m.speech.Update(); err != nil {
		return err
	}

	if m.speech.IsSpellingComplete() && !m.waitingForInput {
		m.waitingForInput = true
	}

	if m.waitingForInput {
		if input.IsSomeKeyPressed(ebiten.KeySpace, ebiten.KeyEnter) {
			m.currentLine++
			if m.currentLine >= len(m.lines) {
				m.speech.Hide()
				m.isSpeaking = false
			} else {
				m.speech.ResetText()
				m.waitingForInput = false
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

	if m.currentLine < len(m.lines) {
		m.speech.Draw(screen, m.lines[m.currentLine])
	}
}
