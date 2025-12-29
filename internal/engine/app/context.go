package app

import (
	"io/fs"

	"github.com/leandroatallah/firefly/internal/engine/assets/imagemanager"
	"github.com/leandroatallah/firefly/internal/engine/audio"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/data/datamanager"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
)

// AppContext holds all major systems and services that are shared across the
// application. It's used for dependency injection, allowing different parts of
// the game to access systems like input, audio, and scene management without
// relying on global variables.
type AppContext struct {
	AudioManager    *audio.AudioManager
	ImageManager    *imagemanager.ImageManager
	DataManager     *datamanager.Manager
	DialogueManager *speech.Manager
	ActorManager    *actors.Manager
	SceneManager    navigation.SceneManager
	PhaseManager    *phases.Manager
	Assets          fs.FS
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
