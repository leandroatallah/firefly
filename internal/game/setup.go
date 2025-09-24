package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
)

func Setup() {
	ebiten.SetWindowSize(config.ScreenWidth*2, config.ScreenHeight*2)
	ebiten.SetWindowTitle("2D Boilerplate")

	state, err := NewGameState(MainMenu)
	if err != nil {
		log.Fatal(err)
	}

	game := NewGame().
		SetSceneManager().
		SetAudioContext().
		SetAudioAssets().
		ChangeState(state)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
