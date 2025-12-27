package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/app"
)

const (
	BlueEnemyType enemies.EnemyType = "BLUE"
)

func InitEnemyMap(ctx *app.AppContext) enemies.EnemyMap {
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
