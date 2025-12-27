package speech

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
)

type SpeechFont struct {
	source      *font.FontText
	size        float64
	LineSpacing float64
}

func NewSpeechFont(source *font.FontText, size, lineSpacing float64) *SpeechFont {
	return &SpeechFont{
		source:      source,
		size:        size,
		LineSpacing: lineSpacing,
	}
}

func (f *SpeechFont) Draw(screen *ebiten.Image, msg string, op *text.DrawOptions) {
	f.source.Draw(screen, msg, f.size, op)
}
