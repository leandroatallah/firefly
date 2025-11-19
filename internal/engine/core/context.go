package core

import (
	"io/fs"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core/levels"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/datamanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/imagemanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
)

// AppContext holds all major systems and services that are shared across the
// application. It's used for dependency injection, allowing different parts of
// the game to access systems like input, audio, and scene management without
// relying on global variables.
type AppContext struct {
	InputManager    *input.Manager
	AudioManager    *audiomanager.AudioManager
	ImageManager    *imagemanager.ImageManager
	DataManager     *datamanager.Manager
	DialogueManager *speech.Manager
	ActorManager    *actors.Manager
	SceneManager    navigation.SceneManager
	LevelManager    *levels.Manager
	Assets          fs.FS
}
