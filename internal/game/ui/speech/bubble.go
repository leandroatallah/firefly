package gamespeech

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
)

type SpeechBubble struct {
	*baseSpeech

	nineSlice *image.NineSlice
}

func NewSpeechBubble(fontSource *speech.SpeechFont) *SpeechBubble {
	// Load 9-slice bubble image
	img, _, err := ebitenutil.NewImageFromFile("assets/images/9-slice-speech.png")
	if err != nil {
		log.Fatal(err)
	}
	h := [3]int{4, 4, 4}
	v := [3]int{4, 4, 4}
	ns := image.NewNineSlice(img, h, v)

	// Create indicator image (a simple white square)
	indicatorImg := ebiten.NewImage(8, 8)
	indicatorImg.Fill(color.Black)

	s := &SpeechBubble{
		baseSpeech: newBaseSpeech(fontSource),
		nineSlice:  ns,
	}
	s.SetID(speech.BubbleSpeechID)
	s.indicator = indicatorImg
	return s
}

func (s *SpeechBubble) Show() {
	s.baseSpeech.Show()
	s.SpeechBase.SetSpellingDelay(animDuration)
}

func (s *SpeechBubble) ResetText() {
	s.baseSpeech.ResetText()
	if s.Visible() && s.delay >= animDuration {
		s.SpeechBase.SetSpellingDelay(0)
	}
}

func (s *SpeechBubble) TypingSoundEnabled() bool {
	return true
}

func (s *SpeechBubble) Draw(screen *ebiten.Image, msg string) {
	if !s.Visible() && s.removed {
		return
	}

	var x, y float64
	var w, h int

	// Resting state properties
	w_rest := float64(config.Get().ScreenWidth - minMargin*2)
	h_rest := float64(52)
	x_rest := float64(minMargin)
	var y_rest float64

	if s.GetPosition() == "top" {
		y_rest = float64(minMargin)
	} else {
		// Default to bottom
		y_rest = float64(config.Get().ScreenHeight) - h_rest - float64(minMargin)
	}

	const animDurationLocal = float64(animDuration)
	progress := float64(s.delay) / animDurationLocal
	if progress > 1.0 {
		progress = 1.0
	}

	var scale float64
	if s.ending {
		// Animate out: shrink to center
		scale = 1.0 - progress
	} else {
		// Animate in: grow from center
		scale = progress
	}

	w_anim := w_rest * scale
	h_anim := h_rest * scale

	x = x_rest + (w_rest-w_anim)/2
	y = y_rest + (h_rest-h_anim)/2
	w = int(w_anim)
	h = int(h_anim)

	if w <= 0 || h <= 0 {
		return
	}

	s.nineSlice.Draw(screen, w, h, func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(x, y)
	})

	s.DrawText(screen, msg, x, y, w, h)

	s.DrawIndicator(screen, x, y, w, h)
}
