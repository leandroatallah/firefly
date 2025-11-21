package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

type ZacPlayer struct {
	actors.Player

	coinCount int
}

func NewZacPlayer(
	movementBlocker physics.PlayerMovementBlocker,
) (actors.ActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/actors/player/zac.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		return nil, err
	}

	player := &ZacPlayer{
		Player: actors.Player{Character: *character},
	}
	if err = SetPlayerBodies(player, spriteData); err != nil {
		return nil, fmt.Errorf("SetPlayerBodies: %w", err)
	}
	if err = SetPlayerStats(player, statData); err != nil {
		return nil, fmt.Errorf("SetPlayerStats: %w", err)
	}
	// Pass player itself
	if err = SetMovementModel(player, physics.Platform, player); err != nil {
		return nil, fmt.Errorf("SetMovementModel: %w", err)
	}

	return player, nil
}
