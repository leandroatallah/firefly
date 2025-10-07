package core

import (
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/levels"
	"github.com/leandroatallah/firefly/internal/navigation"
	"github.com/leandroatallah/firefly/internal/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/systems/input"
)

// AppContext holds all major systems and services that are shared across the
// application. It's used for dependency injection, allowing different parts of
// the game to access systems like input, audio, and scene management without
// relying on global variables.
type AppContext struct {
	InputManager  *input.Manager
	AudioManager  *audiomanager.AudioManager
	SceneManager  navigation.SceneManager
	LevelManager  *levels.Manager
	Configuration *config.Config
}
