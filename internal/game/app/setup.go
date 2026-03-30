package gamesetup

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/i18n"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
	gamescene "github.com/boilerplate/ebiten-template/internal/game/scenes"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
	gamespeech "github.com/boilerplate/ebiten-template/internal/game/ui/speech"
)

func Setup(assets fs.FS) error {
	cfg := config.Get()
	// Basic Ebiten setup
	ebiten.SetWindowSize(cfg.ScreenWidth*3, cfg.ScreenHeight*3)
	ebiten.SetFullscreen(cfg.Fullscreen)
	ebiten.SetWindowTitle("Ebitengine Boilerplate")

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
	fontMain, err := font.NewFontText(assets, cfg.MainFontFace)
	if err != nil {
		return err
	}
	fontSmall, err := font.NewFontText(assets, cfg.SmallFontFace)
	if err != nil {
		return err
	}

	speechFontMain := speech.NewSpeechFont(fontMain, 8, 14)
	speechFontSmall := speech.NewSpeechFont(fontSmall, 8, 12)

	speechBubble := gamespeech.NewSpeechBubble(assets, speechFontMain, i18nManager)
	speechStory := gamespeech.NewStorySpeech(speechFontSmall, i18nManager)
	dialogueManager := speech.NewManager(speechBubble, speechStory)
	dialogueManager.SetActiveSpeech(speech.BubbleSpeechID)
	dialogueManager.SetAudioManager(audioManager)
	dialogueManager.SetTypingSounds(collectSpeechBleeps(assets))
	dialogueManager.SetDefaultSpeechAudio(collectSpeechBleeps(assets))

	// Load audio assets
	audio.LoadAudioAssetsFromFS(assets, audioManager)

	// Load VFX Manager (particles + floating text)
	vfxManager := vfx.NewManager(assets, "assets/particles/vfx.json")
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

	vfxManager.SetAppContext(appContext)

	// Create and run the game
	game := app.NewGame(appContext)

	// Set initial game scene
	initialScene := scenestypes.SceneMenu
	game.AppContext.SceneManager.NavigateTo(initialScene, nil, false)

	if err := ebiten.RunGame(game); err != nil {
		return err
	}

	return nil
}
