package gamesetup

import (
	"io/fs"
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/i18n"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles/vfx"
	enginestylevfx "github.com/boilerplate/ebiten-template/internal/engine/render/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/phaseoverlay"
	enginespeech "github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
	gamescene "github.com/boilerplate/ebiten-template/internal/game/scenes"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
	gamespeech "github.com/boilerplate/ebiten-template/internal/game/ui/speech"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/projectile"
	kitspeech "github.com/boilerplate/ebiten-template/internal/kit/ui/speech"
	"github.com/hajimehoshi/ebiten/v2"
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

	speechFontMain := enginespeech.NewSpeechFont(fontMain, 8, 14)
	speechFontSmall := enginespeech.NewSpeechFont(fontSmall, 8, 12)

	speechBubble := gamespeech.NewSpeechBubble(assets, speechFontMain, i18nManager)
	speechStory := gamespeech.NewStorySpeech(speechFontSmall, i18nManager)
	dialogueManager := kitspeech.NewManager(speechBubble, speechStory)
	dialogueManager.SetActiveSpeech(kitspeech.BubbleSpeechID)
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
		AudioManager:      audioManager,
		DialogueManager:   dialogueManager,
		EventManager:      event.NewManager(),
		ActorManager:      actorManager,
		SceneManager:      sceneManager,
		PhaseManager:      phaseManager,
		I18n:              i18nManager,
		ImageManager:      nil,
		DataManager:       nil,
		Assets:            assets,
		Config:            config.Get(),
		Space:             space.NewSpace(),
		VFX:               vfxManager,
		FadeOverlay:       enginestylevfx.NewFadeOverlay(),
		SolidColorOverlay: enginestylevfx.NewSolidColor(),

		Font: fontMain,
	}
	projManager := projectile.NewManager(appContext.Space)
	projManager.SetVFXManager(vfxManager)
	appContext.ProjectileManager = projManager

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	vfxManager.SetAppContext(appContext)

	// Create and run the game
	game := app.NewGame(appContext)
	game.DebugOverlay().SetFont(fontSmall.NewFace(8))
	game.ActorInspector().SetFont(fontSmall.NewFace(8))

	// F2 phase-jump overlay: list every phase and warp to the chosen one.
	phaseEntries := make([]phaseoverlay.Entry, 0, len(GetPhases()))
	for _, p := range GetPhases() {
		phaseEntries = append(phaseEntries, phaseoverlay.Entry{ID: p.ID, Name: p.Name})
	}
	game.PhaseOverlay().SetFont(fontSmall.NewFace(8))
	game.PhaseOverlay().SetEntries(phaseEntries)
	game.PhaseOverlay().SetOnSelect(func(id int) {
		if err := appContext.PhaseManager.SetCurrentPhase(id); err != nil {
			log.Printf("phase jump: %v", err)
			return
		}
		// Force a fresh start: the jump bypasses normal progression, so clear
		// any audio, dialogue, and VFX still lingering from the previous phase.
		appContext.ResetTransientState()
		appContext.GoToCurrentPhaseScene(nil, true)
	})

	// Set initial game scene
	initialScene := scenestypes.SceneMenu
	if cfg.SkipIntro {
		phase, _ := phaseManager.GetCurrentPhase()
		initialScene = phase.SceneType
	}
	game.AppContext.SceneManager.NavigateTo(initialScene, nil, false)

	if err := ebiten.RunGame(game); err != nil {
		return err
	}

	return nil
}
