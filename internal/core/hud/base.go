package hud

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/actors"
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
}

func NewStatusBar(player *actors.Player) (*StatusBar, error) {
	heart, _, err := ebitenutil.NewImageFromFile("assets/heart.png")
	if err != nil {
		return nil, err
	}

	return &StatusBar{
		heartImg: heart,
		player:   player,
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

	for i := 0; i < hearthCount; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(width*i+gap*i), 0)
		bg.DrawImage(h.heartImg, op)
	}

	screen.DrawImage(bg, op)

}
