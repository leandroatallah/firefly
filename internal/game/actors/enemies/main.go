package gameenemies

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/actors/enemies"
)

const (
	BlueEnemyType enemies.EnemyType = iota
)

func InitEnemyMap() enemies.EnemyMap {
	enemyMap := map[enemies.EnemyType]actors.ActorEntity{
		BlueEnemyType: NewBlueEnemy(),
	}
	return enemyMap
}
