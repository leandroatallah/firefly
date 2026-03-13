package gameplayermethods

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

// StandardStateTransitionLogic handles common state transitions like dying.
func StandardStateTransitionLogic(c *actors.Character) bool {
	state := c.State()

	// When the character dies or is exiting, the state no longer changes.
	if state == gamestates.Dying || state == gamestates.Exiting {
		return true
	}

	if c.Health() <= 0 {
		state, err := c.NewState(gamestates.Dying)
		if err != nil {
			log.Printf("Failed to create new state %v: %v", gamestates.Dying, err)
			return false
		}
		c.SetState(state)
		return true
	}

	return false
}
