package scene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/assets/font"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/transition"
	"github.com/leandroatallah/firefly/internal/navigation"
)

const (
	kickBackBG = "assets/kick_backOGG.ogg"
)

type MenuScene struct {
	BaseScene

	fontText *font.FontText
}

func NewMenuScene() *MenuScene {
	fontText, err := font.NewFontText(config.MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	return &MenuScene{fontText: fontText}
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
	textOp.GeoM.Translate(config.ScreenWidth/2, config.ScreenHeight/2)
	textOp.ColorScale.Scale(1, 1, 1, float32(120))
	s.fontText.Draw(screen, "Press Enter to start", 8, textOp)

}

func (s *MenuScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.Manager.NavigateTo(navigation.SceneLevels, transition.NewFader())
	}
	return nil
}

func (s *MenuScene) OnStart() {
	s.audiomanager = s.Manager.AudioManager()
	s.audiomanager.SetVolume(1)
	s.audiomanager.PlayMusic(kickBackBG)
}

func (s *MenuScene) OnFinish() {
	s.audiomanager.PauseMusic(kickBackBG)
}
