package gameplayermethods

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

// StandardStateTransitionLogic handles common state transitions like dying.
func StandardStateTransitionLogic(c *actors.Character) bool {
	state := c.State()

	if state == gamestates.Dying && c.IsAnimationFinished() {
		c.SetNewStateFatal(gamestates.Dead)
		return true
	}

	// When the character is exiting, the state no longer changes.
	if state == gamestates.Exiting || state == gamestates.Dying || state == gamestates.Dead {
		return true
	}

	if c.Health() <= 0 && state != gamestates.Dying {
		c.SetNewStateFatal(gamestates.Dying)
		return true
	}

	return false
}
