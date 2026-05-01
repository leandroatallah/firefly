package app

import (
	"io/fs"
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/assets/imagemanager"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/datamanager"
	"github.com/boilerplate/ebiten-template/internal/engine/data/i18n"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/speech"
)

// AppContext holds all major systems and services that are shared across the
// application. It's used for dependency injection, allowing different parts of
// the game to access systems like input, audio, and scene management without
// relying on global variables.
type AppContext struct {
	AudioManager      audio.Manager
	ImageManager      *imagemanager.ImageManager
	DataManager       *datamanager.Manager
	DialogueManager   *speech.Manager
	EventManager      *event.Manager
	ActorManager      *actors.Manager
	SceneManager      navigation.SceneManager
	PhaseManager      *phases.Manager
	I18n              *i18n.I18nManager
	Assets            fs.FS
	Config            *config.AppConfig
	Space             body.BodiesSpace
	ProjectileManager contractscombat.ProjectileManager
	VFX               vfx.Manager
	Font              *font.FontText

	// Global frame counter
	FrameCount uint64
}

// AppContextHolder is a reusable component for embedding app context
type AppContextHolder struct {
	appContext *AppContext
}

func (c *AppContextHolder) SetAppContext(appContext any) {
	c.appContext = appContext.(*AppContext)
}

func (c *AppContextHolder) AppContext() *AppContext {
	return c.appContext
}

func (c *AppContext) GoToCurrentPhaseScene(t navigation.Transition, freshInstance bool) {
	if c.PhaseManager == nil || c.SceneManager == nil {
		return
	}

	phase, err := c.PhaseManager.GetCurrentPhase()
	if err != nil {
		log.Printf("failed to get current phase: %v", err)
		return
	}

	c.SceneManager.NavigateTo(phase.SceneType, t, freshInstance)
}

func (c *AppContext) CompleteCurrentPhase(t navigation.Transition, freshInstance bool) {
	if c.PhaseManager == nil {
		return
	}

	if err := c.PhaseManager.AdvanceToNextPhase(); err != nil {
		log.Printf("failed to advance to next phase: %v", err)
		return
	}

	c.GoToCurrentPhaseScene(t, freshInstance)
}
