package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/core/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/game/constants"
)

type SummaryScene struct {
	scene.BaseScene

	audiomanager *audiomanager.AudioManager
	fontText     *font.FontText
}

func NewSummaryScene(context *core.AppContext) *SummaryScene {
	fontText, err := font.NewFontText(constants.MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	overlay := ebiten.NewImage(constants.ScreenWidth, constants.ScreenHeight)
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
		s.AppContext.LevelManager.AdvanceToNextLevel()
		s.AppContext.SceneManager.NavigateTo(SceneLevels, transition.NewFader())
	}

	return nil
}

func (s *SummaryScene) OnStart() {}

func (s *SummaryScene) OnFinish() {}
