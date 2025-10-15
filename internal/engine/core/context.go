package core

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core/levels"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
)

// AppContext holds all major systems and services that are shared across the
// application. It's used for dependency injection, allowing different parts of
// the game to access systems like input, audio, and scene management without
// relying on global variables.
type AppContext struct {
	InputManager          *input.Manager
	AudioManager          *audiomanager.AudioManager
	DialogueManager       *speech.Manager
	ActorManager          *actors.Manager
	SceneManager          navigation.SceneManager
	LevelManager          *levels.Manager
	PlayerMovementBlocked bool
}

// SetPlayerMovementBlocked sets the flag to block or unblock player movement.
func (ac *AppContext) SetPlayerMovementBlocked(blocked bool) {
	ac.PlayerMovementBlocked = blocked
}

// IsPlayerMovementBlocked returns true if player movement is currently blocked.
func (ac *AppContext) IsPlayerMovementBlocked() bool {
	return ac.PlayerMovementBlocked
}
