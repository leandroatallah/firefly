package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/builder"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func CreateAnimatedCharacter(data schemas.SpriteData) (*gameentitytypes.PlatformerCharacter, error) {
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
		"fall": actors.Falling,
		"hurt": actors.Hurted,
	}
	return builder.CreateAnimatedCharacter(data, stateMap)
}

// SetPlayerBodies
func SetPlayerBodies(player gameentitytypes.PlatformerActorEntity, data schemas.SpriteData) error {
	player.SetID("player")

	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
		"fall": actors.Falling,
		"hurt": actors.Hurted,
	}

	return builder.SetCharacterBodies(player, data, stateMap, "PLAYER")
}

func SetPlayerStats(player gameentitytypes.PlatformerActorEntity, data actors.StatData) error {
	return builder.SetCharacterStats(player, data)
}

func SetMovementModel(
	player gameentitytypes.PlatformerActorEntity,
	movementModel physicsmovement.MovementModelEnum,
) error {
	model, err := physicsmovement.NewMovementModel(movementModel, player)
	if err != nil {
		return err
	}
	player.SetMovementModel(model)
	return nil
}
