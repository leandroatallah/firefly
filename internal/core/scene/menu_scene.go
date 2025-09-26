package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/transition"
)

const (
	kickBackBG = "assets/kick_backOGG.ogg"
)

type MenuScene struct {
	BaseScene
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	bg := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	bg.Fill(color.RGBA{0x44, 0x65, 0x99, 0xff})
	screen.DrawImage(bg, nil)
}

func (s *MenuScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.Manager.NavigateTo(SceneSandbox, transition.NewFader())
	}
	return nil
}

func (s *MenuScene) OnStart() {
	s.audiomanager = s.Manager.AudioManager()
	s.audiomanager.SetVolume(0.1)
	s.audiomanager.PlayMusic(kickBackBG)
}

func (s *MenuScene) OnFinish() {
	s.audiomanager.PauseMusic(kickBackBG)
}
