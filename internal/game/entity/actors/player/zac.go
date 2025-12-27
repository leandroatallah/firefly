package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/physics"
)

type ZacPlayer struct {
	actors.Character

	coinCount int
}

func NewZacPlayer(
	movementBlocker physics.PlayerMovementBlocker,
) (actors.ActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/player/zac.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		return nil, err
	}

	player := &ZacPlayer{
		Character: *character,
	}
	if err = SetPlayerBodies(player, spriteData); err != nil {
		return nil, fmt.Errorf("SetPlayerBodies: %w", err)
	}
	if err = SetPlayerStats(player, statData); err != nil {
		return nil, fmt.Errorf("SetPlayerStats: %w", err)
	}
	// Pass player itself
	if err = SetMovementModel(player, physics.Platform); err != nil {
		return nil, fmt.Errorf("SetMovementModel: %w", err)
	}

	return player, nil
}

func (p *ZacPlayer) GetCharacter() *actors.Character {
	return &p.Character
}
