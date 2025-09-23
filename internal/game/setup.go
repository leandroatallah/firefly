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
	sceneFactory := scene.NewDefaultSceneFactory()
	sceneManager := scene.NewSceneManager()
	menuScene, err := sceneFactory.Create(scene.SceneMenu)
	if err != nil {
		log.Fatal(err)
	}

	game := NewGame(sceneManager)

	// TODO: Should the possible scenes be create here, like routes address?
	game.sceneManager.GoTo(menuScene)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
