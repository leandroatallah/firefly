package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/schemas"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/game/actors/builder"
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
	movementModel physics.MovementModelEnum,
) error {
	model, err := physics.NewMovementModel(movementModel, player)
	if err != nil {
		return err
	}
	player.SetMovementModel(model)
	return nil
}
