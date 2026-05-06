package dialogue

import (
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
	"github.com/hajimehoshi/ebiten/v2"
)

// Manager is the engine-facing contract for a dialogue orchestrator.
// Concrete implementation lives in internal/kit/ui/speech.
type Manager interface {
	// Lifecycle / per-frame
	Update() error
	Draw(screen *ebiten.Image)
	IsSpeaking() bool

	// Speech registry
	AddSpeech(s speech.Speech)
	SetSpeech(id string)
	SetActiveSpeech(id string)
	GetActiveSpeech() speech.Speech

	// Audio wiring
	SetAudioManager(m *audio.AudioManager)
	SetTypingSound(path string)
	SetTypingSounds(paths []string)
	SetSpeechAudioQueue(paths []string)
	ClearSpeechAudioQueue()
	SetDefaultSpeechAudio(paths []string)
	ApplyDefaultSpeechAudio(lineCount int)

	// Behaviour flags
	SetSpeechSkipEnabled(enabled bool)

	// Display
	ShowMessages(lines []string, position string, speed int)
}

const (
	BubbleSpeechID = "bubble"
	StorySpeechID  = "story"
)
