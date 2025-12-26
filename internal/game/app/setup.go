package gamesetup

import (
	"io/fs"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core/game"
	"github.com/leandroatallah/firefly/internal/engine/core/levels"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
	gamescene "github.com/leandroatallah/firefly/internal/game/scenes"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
	gamespeech "github.com/leandroatallah/firefly/internal/game/speech"
)

func Setup(assets fs.FS) error {
	cfg := config.Get()
	// Basic Ebiten setup
	ebiten.SetWindowSize(cfg.ScreenWidth*3, cfg.ScreenHeight*3)
	ebiten.SetWindowTitle("Firefly")

	// Initialize all systems and managers
	audioManager := audiomanager.NewAudioManager()
	sceneManager := scene.NewSceneManager()
	levelManager := levels.NewManager()
	actorManager := actors.NewManager()

	// Initialize Dialogue Manager
	fontText, err := font.NewFontText(cfg.MainFontFace)
	if err != nil {
		return err
	}
	speechFont := speech.NewSpeechFont(fontText, 8, 14)
	speechBubble := gamespeech.NewSpeechBubble(speechFont)
	dialogueManager := speech.NewManager(speechBubble)

	// Load audio assets
	loadAudioAssetsFromFS(assets, audioManager)

	// Load levels
	level1 := levels.Level{ID: 1, Name: "Level 1", TilemapPath: "assets/tilemap/sample-level-1.tmj", NextLevelID: 2}
	level2 := levels.Level{ID: 2, Name: "Level 2", TilemapPath: "assets/tilemap/sample-level-2.tmj", NextLevelID: 0} // 0 means no next level
	levelManager.AddLevel(level1)
	levelManager.AddLevel(level2)
	levelManager.SetCurrentLevel(1)

	appContext := &app.AppContext{
		AudioManager:    audioManager,
		DialogueManager: dialogueManager,
		ActorManager:    actorManager,
		SceneManager:    sceneManager,
		LevelManager:    levelManager,
		ImageManager:    nil,
		DataManager:     nil,
		Assets:          assets,
	}

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := game.NewGame(appContext)

	// Set initial game scene
	game.AppContext.SceneManager.NavigateTo(scenestypes.SceneLevels, nil, false)

	if err := ebiten.RunGame(game); err != nil {
		return err
	}

	return nil
}

// loadAudioAssetsFromFS is a helper function to load all audio files from an fs.FS.
func loadAudioAssetsFromFS(assets fs.FS, am *audiomanager.AudioManager) {
	dir := "assets/audio"
	files, err := fs.ReadDir(assets, dir)
	if err != nil {
		log.Fatalf("error reading embedded audio dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		// Filter for supported audio types
		if !(strings.HasSuffix(fileName, ".ogg") || strings.HasSuffix(fileName, ".wav") || strings.HasSuffix(fileName, ".mp3")) {
			continue
		}

		fullPath := dir + "/" + fileName
		data, err := fs.ReadFile(assets, fullPath)
		if err != nil {
			log.Printf("failed to read embedded file %s: %v", fullPath, err)
			continue
		}

		// Use the existing Add method to process and store the player.
		am.Add(dir+"/"+fileName, data)
	}
}
