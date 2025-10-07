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
	"github.com/leandroatallah/firefly/internal/levels"
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
	levelManager := levels.NewManager()

	// Load audio assets
	loadAudioAssets(audioManager)

	// TODO: The factories and manager should be separated from the levels creation
	// Load levels
	level1 := levels.Level{ID: 1, Name: "Level 1", TilemapPath: "assets/mr-gimmick-stage-3.tmj", NextLevelID: 2}
	level2 := levels.Level{ID: 2, Name: "Level 2", TilemapPath: "assets/mr-gimmick-stage-3.tmj", NextLevelID: 0} // 0 means no next level
	levelManager.AddLevel(level1)
	levelManager.AddLevel(level2)
	levelManager.SetCurrentLevel(1)

	appContext := &core.AppContext{
		InputManager:  inputManager,
		AudioManager:  audioManager,
		SceneManager:  sceneManager,
		LevelManager:  levelManager,
		Configuration: cfg,
	}

	sceneFactory := scene.NewDefaultSceneFactory()
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

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
