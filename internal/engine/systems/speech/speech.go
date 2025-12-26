package speech

import "github.com/hajimehoshi/ebiten/v2"

type Speech interface {
	ID() string
	Show()
	Hide()
	Visible() bool
	Text(msg string) string
	ResetText()
	SetSpellingDelay(d int)
	IsSpellingComplete() bool
	CompleteSpelling()
	Count() int
	Update() error
	Draw(screen *ebiten.Image, text string)
}

type SpeechBase struct {
	id            string
	visible       bool
	count         int
	FontSource    *SpeechFont
	text          string
	spellingCount int
	spellingDelay int
}

func NewSpeechBase(fontSource *SpeechFont) *SpeechBase {
	return &SpeechBase{
		FontSource: fontSource,
	}
}

func (s *SpeechBase) Update() error {
	s.count++

	return nil
}

func (s *SpeechBase) Count() int {
	return s.count
}

func (s *SpeechBase) ID() string {
	return s.id
}

func (s *SpeechBase) Show() {
	s.visible = true
}

func (s *SpeechBase) Hide() {
	s.visible = false
}

func (s *SpeechBase) Visile() bool {
	return s.visible
}

func (s *SpeechBase) SetSpellingDelay(d int) {
	s.spellingDelay = d
}

func (s *SpeechBase) Text(msg string, speed int) string {
	s.text = msg // Store the full message

	if s.count < s.spellingDelay {
		s.spellingCount = 0
		return ""
	}

	// Adjust count to be relative to when spelling should start
	effectiveCount := s.count - s.spellingDelay
	if effectiveCount > 0 && effectiveCount%speed == 0 {
		s.spellingCount++
	}
	limit := min(s.spellingCount, len(s.text))
	return s.text[:limit]
}

func (s *SpeechBase) IsSpellingComplete() bool {
	// Check if the number of characters spelled is greater than or equal to the total message length.
	// Also check that the message is not empty.
	return s.spellingCount >= len(s.text) && len(s.text) > 0
}

func (s *SpeechBase) CompleteSpelling() {
	s.spellingCount = len(s.text)
}

func (s *SpeechBase) ResetText() {
	s.spellingCount = 0
}

func (s *SpeechBase) Image(screen *ebiten.Image) *ebiten.Image {
	return nil
}
