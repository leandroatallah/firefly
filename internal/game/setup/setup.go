package gamesetup

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/game"
	"github.com/leandroatallah/firefly/internal/engine/core/levels"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
	gamescene "github.com/leandroatallah/firefly/internal/game/scenes"
	gamespeech "github.com/leandroatallah/firefly/internal/game/speech"
)

func Setup() {
	// Basic Ebiten setup
	ebiten.SetWindowSize(config.Get().ScreenWidth*3, config.Get().ScreenHeight*3)
	ebiten.SetWindowTitle("Firefly")

	// Initialize all systems and managers
	inputManager := input.NewManager()
	audioManager := audiomanager.NewAudioManager()
	sceneManager := scene.NewSceneManager()
	levelManager := levels.NewManager()
	actorManager := actors.NewManager()

	// Initialize Dialogue Manager
	fontText, err := font.NewFontText("assets/fonts/pressstart2p.ttf")
	if err != nil {
		log.Fatal(err)
	}
	speechFont := speech.NewSpeechFont(fontText, 8, 14)
	speechBubble := gamespeech.NewSpeechBubble(speechFont)
	dialogueManager := speech.NewManager(speechBubble)

	// Load audio assets
	loadAudioAssets(audioManager)

	// Load levels
	level1 := levels.Level{ID: 1, Name: "Level 1", TilemapPath: "assets/tilemap/sample-level-1.tmj", NextLevelID: 2}
	level2 := levels.Level{ID: 2, Name: "Level 2", TilemapPath: "assets/tilemap/sample-level-2.tmj", NextLevelID: 0} // 0 means no next level
	levelManager.AddLevel(level1)
	levelManager.AddLevel(level2)
	levelManager.SetCurrentLevel(1)

	appContext := &core.AppContext{
		InputManager:    inputManager,
		AudioManager:    audioManager,
		DialogueManager: dialogueManager,
		ActorManager:    actorManager,
		SceneManager:    sceneManager,
		LevelManager:    levelManager,
	}

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := game.NewGame(appContext)

	// Set initial game scene
	game.AppContext.SceneManager.NavigateTo(gamescene.SceneLevels, nil)

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
			audioItem, err := am.Load("assets/audio/" + file.Name())
			if err != nil {
				log.Printf("error loading audio file %s: %v", file.Name(), err)
				continue
			}
			am.Add(audioItem.Name(), audioItem.Data())
		}
	}
}
