package gamesetup

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/config"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/game"
	"github.com/leandroatallah/firefly/internal/engine/core/levels"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
	gamescene "github.com/leandroatallah/firefly/internal/game/game_scene"
)

func Setup() {
	// Basic Ebiten setup
	ebiten.SetWindowSize(config.ScreenWidth*3, config.ScreenHeight*3)
	ebiten.SetWindowTitle("Firefly")

	// Initialize all systems and managers
	cfg := config.NewConfig()
	inputManager := input.NewManager()
	audioManager := audiomanager.NewAudioManager()
	sceneManager := scene.NewSceneManager()
	levelManager := levels.NewManager()

	// Load audio assets
	loadAudioAssets(audioManager)

	// Load levels
	level1 := levels.Level{ID: 1, Name: "Level 1", TilemapPath: "assets/sample-level-1.tmj", NextLevelID: 2}
	level2 := levels.Level{ID: 2, Name: "Level 2", TilemapPath: "assets/sample-level-2.tmj", NextLevelID: 0} // 0 means no next level
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

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := game.NewGame(appContext)

	// TODO: Game state is disabled. Check if it is necessary. The game is handled by scenes.
	// Set initial game state
	// game.SetState(state.Intro)
	game.AppContext.SceneManager.NavigateTo(gamescene.SceneIntro, nil)

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
