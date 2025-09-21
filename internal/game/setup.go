package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/physics"
)

const (
	wallWidth = 20
)

func Setup() {
	ebiten.SetWindowSize(config.ScreenWidth*2, config.ScreenHeight*2)
	ebiten.SetWindowTitle("2D Boilerplate")

	// Boundaries
	wallTop := physics.NewObstacleRect(
		physics.NewRect(0, 0, config.ScreenWidth, wallWidth),
		[]*physics.CollisionArea{physics.NewCollisionArea(physics.NewRect(0, 0, config.ScreenWidth, wallWidth))},
	)
	wallLeft := physics.NewObstacleRect(
		physics.NewRect(0, 0, wallWidth, config.ScreenHeight), nil,
	)
	wallRight := physics.NewObstacleRect(
		physics.NewRect(config.ScreenWidth-wallWidth, 0, wallWidth, config.ScreenHeight), nil,
	)
	wallDown := physics.NewObstacleRect(
		physics.NewRect(0, config.ScreenHeight-wallWidth, config.ScreenWidth, wallWidth), nil,
	)

	// Enemies
	enemyRect := physics.NewObstacleRect(physics.NewRect(100, 100, 32, 32), nil)

	// Player
	p := physics.NewPlayer()

	game := NewGame(p)
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
