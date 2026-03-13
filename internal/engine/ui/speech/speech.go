package speech

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Speech interface {
	ID() string
	Show()
	Hide()
	Visible() bool
	Text(msg string) string
	ResetText()
	SetID(id string)
	SetSpellingDelay(d int)
	IsSpellingComplete() bool
	CompleteSpelling()
	Count() int
	Update() error
	Draw(screen *ebiten.Image, text string)
	SetPosition(pos string)
	SetSpeed(speed int)
	SetColor(c color.Color)
	Color() color.Color
	SetSkipFlash(frames int)
	IsAccumulative() bool
	SetAccumulative(bool)
}

type SpeechBase struct {
	id            string
	visible       bool
	count         int
	FontSource    *SpeechFont
	text          string
	spellingCount int
	spellingDelay int
	position      string
	speed         int
	color         color.Color
	skipFlash     int
	accumulative  bool
}

func NewSpeechBase(fontSource *SpeechFont) *SpeechBase {
	return &SpeechBase{
		FontSource: fontSource,
		position:   "bottom",
		speed:      4,
		color:      color.Black,
	}
}

func (s *SpeechBase) Update() error {
	s.count++
	if s.skipFlash > 0 {
		s.skipFlash--
	}
	if s.count > s.spellingDelay {
		effectiveCount := s.count - s.spellingDelay
		speed := s.speed
		if speed <= 0 {
			speed = 4
		}
		if effectiveCount > 0 && effectiveCount%speed == 0 {
			if s.spellingCount < len(s.text) {
				s.spellingCount++
			}
		}
	}
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

func (s *SpeechBase) Visible() bool {
	return s.visible
}

func (s *SpeechBase) SetID(id string) {
	s.id = id
}

func (s *SpeechBase) SetSpellingDelay(d int) {
	s.spellingDelay = d
}

func (s *SpeechBase) Text(msg string, speed int) string {
	s.text = msg // Store the full message

	if speed > 0 && speed != s.speed {
		s.speed = speed
	}

	if s.count < s.spellingDelay {
		return ""
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
	s.count = 0
	s.skipFlash = 0
}

func (s *SpeechBase) SetColor(c color.Color) {
	s.color = c
}

func (s *SpeechBase) Color() color.Color {
	return s.color
}

func (s *SpeechBase) SetSkipFlash(frames int) {
	s.skipFlash = frames
}

func (s *SpeechBase) SkipFlash() int {
	return s.skipFlash
}

func (s *SpeechBase) IsAccumulative() bool {
	return s.accumulative
}

func (s *SpeechBase) SetAccumulative(acc bool) {
	s.accumulative = acc
}

func (s *SpeechBase) Image(screen *ebiten.Image) *ebiten.Image {
	return nil
}
