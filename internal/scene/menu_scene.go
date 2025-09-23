package scene

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
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
		s.Manager.GoToScene(SceneSandbox)
	}
	return nil
}

func (s *MenuScene) OnStart() {}

func (s *MenuScene) OnFinish() {
	fmt.Println("Finish Menu Scence")
}
