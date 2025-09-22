package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/scene"
)

func Setup() {
	ebiten.SetWindowSize(config.ScreenWidth*2, config.ScreenHeight*2)
	ebiten.SetWindowTitle("2D Boilerplate")

	// Scenes
	sceneManager := scene.NewSceneManager()
	// TODO: Replace with a factory pattern
	menuScene := &scene.MenuScene{}

	game := NewGame(sceneManager)

	game.sceneManager.GoTo(menuScene)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
