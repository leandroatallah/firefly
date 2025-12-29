package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
)

type CoinCollector interface {
	AddCoinCount(amount int)
	CoinCount() int
}

type CherryPlayer struct {
	actors.Character

	coinCount        int
	movementBlockers int
}

func NewCherryPlayer() (actors.ActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/player/cherry.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		return nil, err
	}

	player := &CherryPlayer{
		Character: *character,
	}
	if err = SetPlayerBodies(player, spriteData); err != nil {
		return nil, fmt.Errorf("SetPlayerBodies: %w", err)
	}
	if err = SetPlayerStats(player, statData); err != nil {
		return nil, fmt.Errorf("SetPlayerStats: %w", err)
	}
	// Pass player itself
	if err = SetMovementModel(player, physicsmovement.Platform); err != nil {
		return nil, fmt.Errorf("SetMovementModel: %w", err)
	}

	character.StateCollisionManager.RefreshCollisions()

	return player, nil
}

func (p *CherryPlayer) GetCharacter() *actors.Character {
	return &p.Character
}

func (p *CherryPlayer) AddCoinCount(amount int) {
	p.coinCount += amount
}
func (p *CherryPlayer) CoinCount() int {
	return p.coinCount
}
