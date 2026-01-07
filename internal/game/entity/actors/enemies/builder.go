package gameenemies

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/builder"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func CreateAnimatedCharacter(data schemas.SpriteData) (*gameentitytypes.PlatformerCharacter, error) {
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
	}
	return builder.CreateAnimatedCharacter(data, stateMap)
}

// SetEnemyBodies
func SetEnemyBodies(enemy gameentitytypes.PlatformerActorEntity, data schemas.SpriteData, id string) error {
	enemy.SetID(id)

	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
	}

	return builder.SetCharacterBodies(enemy, data, stateMap, "ENEMY")
}

func SetEnemyStats(enemy gameentitytypes.PlatformerActorEntity, data actors.StatData) error {
	return builder.SetCharacterStats(enemy, data)
}
