package gamespeech

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/i18n"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
)

type StorySpeech struct {
	*baseSpeech
}

func NewStorySpeech(fontSource *speech.SpeechFont, i18nManager *i18n.I18nManager) *StorySpeech {
	// Create indicator image (a simple white square)
	indicatorImg := ebiten.NewImage(8, 8)
	indicatorImg.Fill(color.White)

	s := &StorySpeech{
		baseSpeech: newBaseSpeech(fontSource, i18nManager),
	}
	s.SetAccumulative(true)
	s.indicator = indicatorImg
	s.SetID(speech.StorySpeechID)
	s.SetColor(color.White) // Default to white for storytelling
	return s
}

func (s *StorySpeech) Show() {
	s.baseSpeech.Show()
	s.SpeechBase.SetSpellingDelay(0)
}

func (s *StorySpeech) TypingSoundEnabled() bool {
	return false
}

func (s *StorySpeech) Draw(screen *ebiten.Image, msg string) {
	if !s.Visible() && s.removed {
		return
	}

	// For storytelling, we use more space and no bubble
	w := config.Get().ScreenWidth - minMargin*2
	h := 100 // Increased height for multiple lines
	x := float64(minMargin)
	var y float64

	if s.GetPosition() == "top" {
		y = float64(minMargin)
	} else {
		// Default to bottom
		y = float64(config.Get().ScreenHeight) - float64(h) - float64(minMargin)
	}

	s.DrawText(screen, msg, x, y, w, h)

	s.DrawIndicator(screen, x, y, w, h)
}
