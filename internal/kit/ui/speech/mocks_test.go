package speech

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// mockSpeech is a hand-rolled fake satisfying internal/engine/ui/speech.Speech.
// It lives in this package to support white-box tests for Manager without
// exposing test helpers in production code or duplicating ebiten image
// boundaries elsewhere.
type mockSpeech struct {
	id                string
	visible           bool
	spellingComplete  bool
	updateCalled      int
	drawCalled        int
	resetTextCalled   int
	showCalled        int
	hideCalled        int
	setPositionCalled string
	setSpeedCalled    int
}

func (m *mockSpeech) ID() string                             { return m.id }
func (m *mockSpeech) Show()                                  { m.visible = true; m.showCalled++ }
func (m *mockSpeech) Hide()                                  { m.visible = false; m.hideCalled++ }
func (m *mockSpeech) Visible() bool                          { return m.visible }
func (m *mockSpeech) Text(msg string) string                 { return msg }
func (m *mockSpeech) ResetText()                             { m.resetTextCalled++; m.spellingComplete = false }
func (m *mockSpeech) SetID(id string)                        { m.id = id }
func (m *mockSpeech) SetSpellingDelay(d int)                 {}
func (m *mockSpeech) IsSpellingComplete() bool               { return m.spellingComplete }
func (m *mockSpeech) CompleteSpelling()                      { m.spellingComplete = true }
func (m *mockSpeech) Count() int                             { return 0 }
func (m *mockSpeech) Update() error                          { m.updateCalled++; return nil }
func (m *mockSpeech) Draw(screen *ebiten.Image, text string) { m.drawCalled++ }
func (m *mockSpeech) SetPosition(pos string)                 { m.setPositionCalled = pos }
func (m *mockSpeech) SetSpeed(speed int)                     { m.setSpeedCalled = speed }
func (m *mockSpeech) SetColor(c color.Color)                 {}
func (m *mockSpeech) Color() color.Color                     { return color.Black }
func (m *mockSpeech) SetSkipFlash(frames int)                {}
func (m *mockSpeech) IsAccumulative() bool                   { return false }
func (m *mockSpeech) SetAccumulative(bool)                   {}
