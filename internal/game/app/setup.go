package gamesetup

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/audio"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/data/i18n"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/render/particles/vfx"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
	gamescene "github.com/leandroatallah/firefly/internal/game/scenes"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
	gamespeech "github.com/leandroatallah/firefly/internal/game/ui/speech"
)

func Setup(assets fs.FS) error {
	cfg := config.Get()
	// Basic Ebiten setup
	ebiten.SetWindowSize(cfg.ScreenWidth*3, cfg.ScreenHeight*3)
	ebiten.SetFullscreen(cfg.Fullscreen)
	ebiten.SetWindowTitle("The Higher I Climb, The Farther You Are")

	// Initialize all systems and managers
	audioManager := audio.NewAudioManager()
	sceneManager := scene.NewSceneManager()
	phaseManager := phases.NewManager()
	actorManager := actors.NewManager()

	i18nManager := i18n.NewI18nManager(assets)
	if err := i18nManager.Load(cfg.Language); err != nil {
		return err
	}

	// Initialize Dialogue Manager
	fontMain, err := font.NewFontText(cfg.MainFontFace)
	if err != nil {
		return err
	}
	fontSmall, err := font.NewFontText(cfg.SmallFontFace)
	if err != nil {
		return err
	}

	speechFontMain := speech.NewSpeechFont(fontMain, 8, 14)
	speechFontSmall := speech.NewSpeechFont(fontSmall, 8, 12)

	speechBubble := gamespeech.NewSpeechBubble(speechFontMain)
	speechStory := gamespeech.NewStorySpeech(speechFontSmall)
	dialogueManager := speech.NewManager(speechBubble, speechStory)
	dialogueManager.SetActiveSpeech(speech.BubbleSpeechID)
	dialogueManager.SetAudioManager(audioManager)
	dialogueManager.SetTypingSounds(collectSpeechBleeps(assets))
	dialogueManager.SetDefaultSpeechAudio(collectSpeechBleeps(assets))

	// Load audio assets
	audio.LoadAudioAssetsFromFS(assets, audioManager)

	// Load VFX Manager (particles + floating text)
	vfxManager := vfx.NewManager("assets/particles/vfx.json")
	vfxManager.SetDefaultFont(fontMain)

	// Load phases
	for _, p := range GetPhases() {
		phaseManager.AddPhase(p)
	}
	phaseManager.SetCurrentPhase(1)

	appContext := &app.AppContext{
		AudioManager:    audioManager,
		DialogueManager: dialogueManager,
		EventManager:    event.NewManager(),
		ActorManager:    actorManager,
		SceneManager:    sceneManager,
		PhaseManager:    phaseManager,
		I18n:            i18nManager,
		ImageManager:    nil,
		DataManager:     nil,
		Assets:          assets,
		Config:          config.Get(),
		Space:           space.NewSpace(),
		VFX:             vfxManager,
		Font:            fontMain,
	}

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := app.NewGame(appContext)

	// Set initial game scene
	initialScene := scenestypes.SceneIntro
	if cfg.SkipIntro {
		// Catch the first scene phase
		for _, p := range GetPhases() {
			if p.SceneType == scenestypes.ScenePhases {
				initialScene = p.SceneType
				appContext.PhaseManager.SetCurrentPhase(p.ID)
				break
			}
		}
	}
	game.AppContext.SceneManager.NavigateTo(initialScene, nil, false)

	if err := ebiten.RunGame(game); err != nil {
		return err
	}

	return nil
}
