package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

type CherryPlayer struct {
	actors.Player

	coinCount        int
	movementBlockers int
}

func NewCherryPlayer() (actors.ActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/actors/player/cherry.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		return nil, err
	}

	player := &CherryPlayer{
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

func (p *CherryPlayer) AddCoinCount(amount int) {
	p.coinCount += amount
}
func (p *CherryPlayer) CoinCount() int {
	return p.coinCount
}
