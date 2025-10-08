package scene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/assets/font"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/screenutil"
	"github.com/leandroatallah/firefly/internal/core/transition"
	"github.com/leandroatallah/firefly/internal/navigation"
)

const (
	animationDelay = 60
	fadeDelay      = 60
	maxDuration    = 60
	fadeAnimStep   = 2
)

type introAnimation int

const (
	idle introAnimation = iota
	fadeIn
	duration
	fadeOut
	over
	navigationStarted
)

type IntroScene struct {
	BaseScene

	count          int
	fontText       *font.FontText
	fadeAlpha      uint8
	duration       int
	introAnimation introAnimation
	fadeOverlay    *ebiten.Image
}

func NewIntroScene() *IntroScene {
	fontText, err := font.NewFontText(config.MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	overlay := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	overlay.Fill(color.Black)
	return &IntroScene{fontText: fontText, fadeOverlay: overlay}
}

func (s *IntroScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 255})

	screenutil.DrawCenteredText(screen, s.fontText, "Presented by", 10, color.White)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(1, 1, 1, float32(s.fadeAlpha)/255.0)
	screen.DrawImage(s.fadeOverlay, op)
}

func (s *IntroScene) Update() error {
	// TODO: REMOVE THIS
	// FORCE SKIP
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.NextScene()
	}

	s.count++

	// Allow user to skip
	if s.introAnimation == duration && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.duration = 0
		s.introAnimation = fadeOut
	}

	switch s.introAnimation {
	case idle:
		if s.count > animationDelay {
			s.introAnimation = fadeIn
		}
	case fadeIn:
		s.fadeAlpha -= fadeAnimStep
		if s.fadeAlpha <= fadeAnimStep {
			s.introAnimation = duration
		}
	case duration:
		s.duration++
		if s.duration > maxDuration {
			s.introAnimation = fadeOut
		}
	case fadeOut:
		s.fadeAlpha += 2
		if s.fadeAlpha > 255-fadeAnimStep {
			s.introAnimation = over
		}
	case over:
		s.NextScene()
	}

	return nil
}

func (s *IntroScene) NextScene() {
	s.appContext.SceneManager.NavigateTo(navigation.SceneMenu, transition.NewFader())
	s.introAnimation = navigationStarted
}

func (s *IntroScene) OnStart() {
	s.fadeAlpha = 255

}

func (s *IntroScene) OnFinish() {}
