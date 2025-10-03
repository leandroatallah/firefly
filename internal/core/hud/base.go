package hud

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/assets/font"
	"github.com/leandroatallah/firefly/internal/config"
)

const (
	heartLimit = 5
	gap        = 4
)

type HUD interface {
	Draw(screen *ebiten.Image)
	Update() error
}

type StatusBar struct {
	heartImg *ebiten.Image
	player   *actors.Player
	score    int
	mainText *font.FontText
}

func NewStatusBar(player *actors.Player, score int) (*StatusBar, error) {
	heart, _, err := ebitenutil.NewImageFromFile("assets/heart.png")
	if err != nil {
		return nil, err
	}

	mainText, err := font.NewFontText(config.MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	return &StatusBar{
		heartImg: heart,
		player:   player,
		score:    score,
		mainText: mainText,
	}, nil
}

func (h *StatusBar) Update() error {
	return nil
}

func (h *StatusBar) Draw(screen *ebiten.Image) {
	if h.player == nil {
		return
	}

	hearthCount := min(h.player.Health(), heartLimit)
	width := h.heartImg.Bounds().Dx()
	height := h.heartImg.Bounds().Dy()

	bg := ebiten.NewImage(width*5+gap*4, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 32)
	textOp.ColorScale.ScaleWithColor(color.Black)
	h.mainText.Draw(screen, fmt.Sprintf("Score: %d", h.score), 12, textOp)

	for i := 0; i < hearthCount; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(width*i+gap*i), 0)
		bg.DrawImage(h.heartImg, op)
	}

	screen.DrawImage(bg, op)

}
