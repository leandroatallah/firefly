package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type CherryPlayer struct {
	gameentitytypes.PlatformerCharacter
}

func NewCherryPlayer() (gameentitytypes.PlatformerActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/player/cherry.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		return nil, err
	}
	character.AddSkill(skill.NewJumpSkill())
	character.AddSkill(skill.NewHorizontalMovementSkill())
	character.AddSkill(skill.NewDashSkill())

	player := &CherryPlayer{
		PlatformerCharacter: *character,
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
