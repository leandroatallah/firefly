package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/object"
)

type Game struct {
	player     *object.Player
	boundaries []object.Object
}

func NewGame(player *object.Player) *Game {
	return &Game{player: player}
}

func (g *Game) Update() error {
	g.player.Update(g.boundaries)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xcc, 0xcc, 0xdd, 0xff})

	for _, b := range g.boundaries {
		b.(object.Obstacle).Draw(screen)
	}

	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) AddBoundaries(boundaries ...object.Object) {
	for _, o := range boundaries {
		g.boundaries = append(g.boundaries, o)
	}
}
