package scene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/assets/font"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/screenutil"
	"github.com/leandroatallah/firefly/internal/core/transition"
	"github.com/leandroatallah/firefly/internal/navigation"
)

type SummaryScene struct {
	BaseScene

	fontText *font.FontText
}

func NewSummaryScene() *SummaryScene {
	fontText, err := font.NewFontText(config.MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	overlay := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	overlay.Fill(color.Black)
	return &SummaryScene{fontText: fontText}
}

func (s *SummaryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 255})
	screenutil.DrawCenteredText(screen, s.fontText, "Summary screen", 10, color.White)
}

func (s *SummaryScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.appContext.LevelManager.AdvanceToNextLevel()
		s.appContext.SceneManager.NavigateTo(navigation.SceneLevels, transition.NewFader())
	}

	return nil
}

func (s *SummaryScene) OnStart() {}

func (s *SummaryScene) OnFinish() {}
