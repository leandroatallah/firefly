package gamespeech

import (
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/data/i18n"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
)

const (
	padding             = 8
	minMargin           = 6
	maxMargin           = 24
	delayBeforeRemove   = 12
	delayBeforeSpelling = 60
	speedText           = 4
	animDuration        = 15
)

type baseSpeech struct {
	*speech.SpeechBase

	I18n      *i18n.I18nManager
	delay     int
	ending    bool
	removed   bool
	speedText int
	indicator *ebiten.Image
}

func newBaseSpeech(fontSource *speech.SpeechFont, i18nManager *i18n.I18nManager) *baseSpeech {
	return &baseSpeech{
		SpeechBase: speech.NewSpeechBase(fontSource),
		I18n:       i18nManager,
		removed:    true,
		speedText:  speedText,
	}
}

func (s *baseSpeech) Update() error {
	if err := s.SpeechBase.Update(); err != nil {
		return err
	}

	s.delay++

	if !s.Visible() && s.delay > delayBeforeRemove {
		s.removed = true
	}

	return nil
}

func (s *baseSpeech) Show() {
	s.delay = 0
	s.ending = false
	s.removed = false
	s.SpeechBase.Show()
}

func (s *baseSpeech) Hide() {
	s.delay = 0
	s.ending = true
	s.SpeechBase.Hide()
}

func (s *baseSpeech) Visible() bool {
	return s.SpeechBase.Visible()
}

func (s *baseSpeech) Text(msg string) string {
	if s.I18n != nil {
		// Split by double newline (accumulative messages)
		lines := strings.Split(msg, "\n\n")
		translatedLines := make([]string, len(lines))
		for i, line := range lines {
			translatedLines[i] = s.I18n.T(line)
		}
		msg = strings.Join(translatedLines, "\n\n")
	}
	return s.SpeechBase.Text(msg, s.GetSpeed())
}

func (s *baseSpeech) SetSpeed(speed int) {
	s.SpeechBase.SetSpeed(speed)
	s.speedText = speed
}

func (s *baseSpeech) ResetText() {
	s.SpeechBase.ResetText()
}

func (s *baseSpeech) DrawText(screen *ebiten.Image, msg string, x, y float64, w, h int) {
	textStr := s.Text(msg)
	textX := x + padding
	textY := y + padding
	textW := w - padding*2
	textH := h - padding*2

	if textW > 0 && textH > 0 {
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: s.SpeechBase.FontSource.LineSpacing,
			},
		}
		op.ColorScale.ScaleWithColor(s.Color())

		textArea := ebiten.NewImage(textW, textH)
		s.FontSource.Draw(textArea, textStr, op)

		textAreaOp := &ebiten.DrawImageOptions{}
		textAreaOp.GeoM.Translate(math.Floor(textX), math.Floor(textY))
		screen.DrawImage(textArea, textAreaOp)
	}
}

func (s *baseSpeech) DrawIndicator(screen *ebiten.Image, x, y float64, w, h int) {
	if s.IsSpellingComplete() && !s.ending && s.indicator != nil {
		op := &ebiten.DrawImageOptions{}

		indX := x + float64(w) - float64(s.indicator.Bounds().Dx()) - float64(padding)
		indY := y + float64(h) - float64(s.indicator.Bounds().Dy()) - float64(padding)

		pulse := math.Sin(float64(s.Count()) / 15.0)
		alpha := 0.75 + (pulse * 0.25)
		op.ColorScale.ScaleAlpha(float32(alpha))

		op.GeoM.Translate(indX, indY)
		screen.DrawImage(s.indicator, op)
	}
}

func (s *baseSpeech) ImageOptions() *ebiten.DrawImageOptions {
	return nil
}
