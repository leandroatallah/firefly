package game

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core"
	"github.com/leandroatallah/firefly/internal/core/game/state"
	"github.com/leandroatallah/firefly/internal/core/scene"
	"github.com/leandroatallah/firefly/internal/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/systems/input"
)

func Setup() {
	// Basic Ebiten setup
	ebiten.SetWindowSize(config.ScreenWidth*2, config.ScreenHeight*2)
	ebiten.SetWindowTitle("Firefly")

	// Initialize all systems and managers
	cfg := config.NewConfig()
	inputManager := input.NewManager()
	audioManager := audiomanager.NewAudioManager()
	sceneManager := scene.NewSceneManager()

	// Setup SceneManager with its factory
	sceneFactory := scene.NewDefaultSceneFactory()
	sceneManager.SetFactory(sceneFactory)
	sceneFactory.SetManager(sceneManager)

	// Connect managers that depend on each other
	sceneManager.SetAudioManager(audioManager)

	// Load assets
	loadAudioAssets(audioManager)

	// Create the application context
	appContext := &core.AppContext{
		InputManager:  inputManager,
		AudioManager:  audioManager,
		SceneManager:  sceneManager,
		Configuration: cfg,
	}

	// Create and run the game
	game := NewGame(appContext)

	// Set initial game state
	game.SetState(state.MainMenu)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// loadAudioAssets is a helper function to load all audio files from the assets directory.
func loadAudioAssets(am *audiomanager.AudioManager) {
	files, err := os.ReadDir("assets")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".ogg") || strings.HasSuffix(file.Name(), ".wav")) {
			audioItem, err := am.Load("assets/" + file.Name())
			if err != nil {
				log.Printf("error loading audio file %s: %v", file.Name(), err)
				continue
			}
			am.Add(audioItem.Name(), audioItem.Data())
		}
	}
}
