package gamescene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/utils"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

type MenuScene struct {
	scene.BaseScene

	fontText *font.FontText

	count              int
	isNavigating       bool
	navigationTrigger  utils.DelayTrigger
	shouldFadeOutSound bool
	isFadingOutSound   bool
}

func NewMenuScene(context *app.AppContext) *MenuScene {
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	scene := MenuScene{fontText: fontText}
	scene.SetAppContext(context)
	return &scene
}

func (s *MenuScene) OnStart() {
	am := s.AppContext().SceneManager.AudioManager()
	am.SetVolume(1)
	am.PlayMusic(TitleSound, true)  // Loop menu music
}

func (s *MenuScene) Update() error {
	canSkipDelay := s.count > timing.FromDuration(time.Second)
	if canSkipDelay && !s.isNavigating && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.isNavigating = true
		s.navigationTrigger.Enable(timing.FromDuration(time.Second))
	}

	s.navigationTrigger.Update()
	if s.navigationTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.SceneStory, transition.NewFader(0, config.Get().FadeVisibleDuration), true,
		)
	}

	if s.isNavigating && s.shouldFadeOutSound && !s.isFadingOutSound {
		s.AppContext().AudioManager.FadeOutAll(time.Second)
		s.isFadingOutSound = true
	}

	if s.isNavigating {
		s.shouldFadeOutSound = true
	}

	s.count++

	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xCC, 0x24, 0x40, 0xff})

	textOp := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
			LineSpacing:    0,
		},
	}
	textOp.GeoM.Translate(
		float64(config.Get().ScreenWidth/2),
		float64(config.Get().ScreenHeight/2),
	)
	textOp.ColorScale.Scale(1, 1, 1, float32(120))
	s.fontText.Draw(screen, "Press Enter to start", 8, textOp)
}

func (s *MenuScene) OnFinish() {}
