package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/audioplayer"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/scene"
)

func Setup() {
	ebiten.SetWindowSize(config.ScreenWidth*2, config.ScreenHeight*2)
	ebiten.SetWindowTitle("2D Boilerplate")

	// Scenes
	sceneManager := scene.NewSceneManager()
	sceneFactory := scene.NewDefaultSceneFactory()

	sceneManager.SetFactory(sceneFactory)
	sceneFactory.SetManager(sceneManager)

	// Audio Player
	audioContext := audioplayer.NewContext()
	audioItem, err := audioplayer.LoadAudio("assets/kick_backOGG.ogg")
	if err != nil {
		log.Fatal(err)
	}
	audioAssets := []*audioplayer.AudioItem{audioItem}

	state, err := NewGameState(MainMenu)
	if err != nil {
		log.Fatal(err)
	}

	game := NewGame().
		SetSceneManager(sceneManager).
		SetAudioContext(audioContext).
		SetAudioAssets(audioAssets).
		ChangeState(state)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
