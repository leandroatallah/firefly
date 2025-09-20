package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/object"
)

const (
	wallWidth = 20
)

func Setup() {
	ebiten.SetWindowSize(config.ScreenWidth*2, config.ScreenHeight*2)
	ebiten.SetWindowTitle("2D Boilerplate")

	// Boundaries
	wallTop := object.NewObstacleRect(
		object.NewElement(0, 0, config.ScreenWidth, wallWidth),
		[]*object.CollisionArea{object.NewCollisionArea(object.NewElement(0, 0, config.ScreenWidth, wallWidth))},
	)
	wallLeft := object.NewObstacleRect(
		object.NewElement(0, 0, wallWidth, config.ScreenHeight), nil,
	)
	wallRight := object.NewObstacleRect(
		object.NewElement(config.ScreenWidth-wallWidth, 0, wallWidth, config.ScreenHeight), nil,
	)
	wallDown := object.NewObstacleRect(
		object.NewElement(0, config.ScreenHeight-wallWidth, config.ScreenWidth, wallWidth), nil,
	)

	// Enemies
	enemyRect := object.NewObstacleRect(object.NewElement(100, 100, 32, 32), nil)

	// Player
	p := object.NewPlayer()

	game := NewGame(p)
	// TODO: Use a method to add (setter)
	game.AddBoundaries(
		wallTop,
		wallLeft,
		wallRight,
		wallDown,
		enemyRect,
	)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
