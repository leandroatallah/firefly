package gamescenelevels

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/actors"
)

func (s *LevelsScene) DrawHUD(screen *ebiten.Image) {
	coinCount := 0

	if p, ok := s.player.(*actors.PlayerPlatform); ok {
		coinCount = p.CoinCount()
	}

	hud := ebiten.NewImage(74, 12)
	hud.Fill(color.White)
	hudOp := &ebiten.DrawImageOptions{}
	hudOp.GeoM.Translate(4, 5)
	textOp := &text.DrawOptions{}
	textOp.ColorScale.Scale(0, 0, 0, 255)
	textOp.GeoM.Translate(2, 2)
	s.mainText.Draw(hud, fmt.Sprintf("Score: %d", coinCount), 8, textOp)

	// Draw simple HUD score
	// HUD need to be drawed on screen and not on the camera.
	screen.DrawImage(hud, hudOp)
}
