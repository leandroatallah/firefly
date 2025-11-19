package gameplayer

import (
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
	SetPlayerBodies(player, spriteData)
	SetPlayerStats(player, statData)
	SetMovementModel(player, physics.Platform, movementBlocker)

	return player, nil
}
