package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/actors/enemies"
)

const (
	BlueEnemyType enemies.EnemyType = iota
)

func InitEnemyMap() enemies.EnemyMap {
	enemyMap := map[enemies.EnemyType]actors.ActorEntity{
		BlueEnemyType: func() actors.ActorEntity {
			e, err := NewBlueEnemy()
			if err != nil {
				log.Fatal("InitEnemyMap: %w", err)
			}
			return e
		}(),
	}
	return enemyMap
}
