package gamespeech

import (
	"image/color"
	"log"
	"math"

	"github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
)

const (
	padding             = 8 // Increased padding for the new bubble style
	minMargin           = 6
	maxMargin           = 24
	delayBeforeRemove   = 12
	delayBeforeSpelling = 60
	speedText           = 4
)

type SpeechBubble struct {
	speech.SpeechBase
	delay     int
	ending    bool
	removed   bool
	speedText int
	nineSlice *image.NineSlice
	indicator *ebiten.Image
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

	return &SpeechBubble{
		SpeechBase: *speech.NewSpeechBase(fontSource),
		removed:    true,
		speedText:  speedText,
		nineSlice:  ns,
		indicator:  indicatorImg,
	}
}

func (s *SpeechBubble) Update() error {
	if err := s.SpeechBase.Update(); err != nil {
		return err
	}

	s.delay++

	if !s.Visile() && s.delay > delayBeforeRemove {
		s.removed = true
	}

	return nil
}

func (s *SpeechBubble) Show() {
	s.delay = 0
	s.ending = false
	s.removed = false
	s.SpeechBase.Show()
}

func (s *SpeechBubble) Hide() {
	s.delay = 0
	s.ending = true
	s.SpeechBase.Hide()
}

func (s *SpeechBubble) Visible() bool {
	return s.SpeechBase.Visile()
}

func (s *SpeechBubble) Text(msg string) string {
	return s.SpeechBase.Text(msg, s.speedText)
}

func (s *SpeechBubble) ResetText() {
	s.SpeechBase.ResetText()
}

func (s *SpeechBubble) Draw(screen *ebiten.Image, msg string) {
	if !s.Visile() && s.removed {
		return
	}

	var x, y float64
	var w, h int

	// Resting state properties
	w_rest := float64(config.Get().ScreenWidth - minMargin*2)
	h_rest := float64(52)
	x_rest := float64(minMargin)
	y_rest := float64(config.Get().ScreenHeight) - h_rest - float64(minMargin)

	const animDuration = 15.0 // frames
	progress := float64(s.delay) / animDuration
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

	// --- Draw Text ---
	textStr := s.Text(msg) // Get current text
	textX := x + padding
	textY := y + padding
	textW := w - padding*2
	textH := h - padding*2

	if textW > 0 && textH > 0 {
		textArea := ebiten.NewImage(textW, textH)
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: s.SpeechBase.FontSource.LineSpacing,
			},
		}
		op.ColorScale.ScaleWithColor(color.Black)
		s.FontSource.Draw(textArea, textStr, op)

		textAreaOp := &ebiten.DrawImageOptions{}
		textAreaOp.GeoM.Translate(textX, textY)
		screen.DrawImage(textArea, textAreaOp)
	}

	// --- Draw Indicator ---
	if s.IsSpellingComplete() && !s.ending {
		op := &ebiten.DrawImageOptions{}

		// Position in bottom-right corner
		indX := x + float64(w) - float64(s.indicator.Bounds().Dx()) - float64(padding)
		indY := y + float64(h) - float64(s.indicator.Bounds().Dy()) - float64(padding)

		// Pulsating effect
		pulse := math.Sin(float64(s.Count()) / 15.0) // Use SpeechBase's count
		alpha := 0.75 + (pulse * 0.25)               // Varies alpha between 0.5 and 1.0
		op.ColorScale.ScaleAlpha(float32(alpha))

		op.GeoM.Translate(indX, indY)
		screen.DrawImage(s.indicator, op)
	}
}

func (s *SpeechBubble) ImageOptions() *ebiten.DrawImageOptions {
	return nil
}
