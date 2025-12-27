package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

const (
	kickBackBG = "assets/audio/kick_backOGG.ogg"
)

type MenuScene struct {
	scene.BaseScene

	fontText *font.FontText
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
	// Init audio
	am := s.AppContext().SceneManager.AudioManager()
	am.SetVolume(1)
	// s.audiomanager.PlayMusic(kickBackBG)
}

func (s *MenuScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.AppContext().SceneManager.NavigateTo(scenestypes.SceneLevels, transition.NewFader(), true)
	}

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

func (s *MenuScene) OnFinish() {
	s.AppContext().AudioManager.PauseMusic(kickBackBG)
}
