package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/audiomanager"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

type SummaryScene struct {
	scene.BaseScene

	audiomanager *audiomanager.AudioManager
	fontText     *font.FontText
}

func NewSummaryScene(context *app.AppContext) *SummaryScene {
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	overlay := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)
	overlay.Fill(color.Black)
	scene := SummaryScene{fontText: fontText}
	scene.SetAppContext(context)
	return &scene
}

func (s *SummaryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 255})
	screenutil.DrawCenteredText(screen, s.fontText, "Summary screen", 10, color.White)
}

func (s *SummaryScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.AppContext().LevelManager.AdvanceToNextLevel()
		s.AppContext().SceneManager.NavigateTo(scenestypes.SceneLevels, transition.NewFader(), true)
	}

	return nil
}

func (s *SummaryScene) OnStart() {}

func (s *SummaryScene) OnFinish() {}
