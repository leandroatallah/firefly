package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

const (
	BlueEnemyType enemies.EnemyType = "BLUE"
)

func InitEnemyMap(ctx *core.AppContext) enemies.EnemyMap {
	enemyMap := map[enemies.EnemyType]func(x, y int, id string) actors.ActorEntity{
		BlueEnemyType: func(x, y int, id string) actors.ActorEntity {
			enemy, err := NewBlueEnemy(x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			player, _ := ctx.ActorManager.GetPlayer()
			enemy.SetTarget(player)
			return enemy
		},
	}
	return enemyMap
}

