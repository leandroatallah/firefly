package gameenemies

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/schemas"
	"github.com/leandroatallah/firefly/internal/game/actors/builder"
)

func CreateAnimatedCharacter(data schemas.SpriteData) (*actors.Character, error) {
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
	}
	return builder.CreateAnimatedCharacter(data, stateMap)
}

// SetEnemyBodies
func SetEnemyBodies(enemy actors.ActorEntity, data schemas.SpriteData) error {
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
	}

	return builder.SetCharacterBodies(enemy, data, stateMap, "ENEMY")
}

func SetEnemyStats(enemy actors.ActorEntity, data actors.StatData) error {
	return builder.SetCharacterStats(enemy, data)
}
