package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/builder"
)

func CreateAnimatedCharacter(data schemas.SpriteData) (*actors.Character, error) {
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
		"fall": actors.Falling,
		"hurt": actors.Hurted,
	}
	return builder.CreateAnimatedCharacter(data, stateMap)
}

// SetPlayerBodies
func SetPlayerBodies(player actors.ActorEntity, data schemas.SpriteData) error {
	player.SetID("player")

	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
		"fall": actors.Falling,
		"hurt": actors.Hurted,
	}

	return builder.SetCharacterBodies(player, data, stateMap, "PLAYER")
}

func SetPlayerStats(player actors.ActorEntity, data actors.StatData) error {
	return builder.SetCharacterStats(player, data)
}

func SetMovementModel(
	player actors.ActorEntity,
	movementModel physicsmovement.MovementModelEnum,
) error {
	model, err := physicsmovement.NewMovementModel(movementModel, player)
	if err != nil {
		return err
	}
	player.SetMovementModel(model)
	return nil
}
