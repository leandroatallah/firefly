package gamehud

import (
	"bytes"
	"fmt"
	"image/color"
	"io/fs"

	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	heartLimit = 5
	gap        = 4
)

type StatusBar struct {
	heartImg *ebiten.Image
	player   body.Alive
	score    int
	mainText *font.FontText
}

func NewStatusBar(player body.Alive, score int, mainText *font.FontText, fsys fs.FS) (*StatusBar, error) {
	heartData, err := fs.ReadFile(fsys, "assets/images/heart.png")
	if err != nil {
		return nil, err
	}
	heart, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(heartData))
	if err != nil {
		return nil, err
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
