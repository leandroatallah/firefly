package gameenemies

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/enemies"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
)

const (
	WolfEnemyType  enemies.EnemyType = "WOLF"
	BatEnemyType   enemies.EnemyType = "BAT"
	SwarmEnemyType enemies.EnemyType = "SWARM"
)

func InitEnemyMap(ctx *app.AppContext) enemies.EnemyMap[platformer.PlatformerActorEntity] {
	enemyMap := map[enemies.EnemyType]func(x, y int, id string) platformer.PlatformerActorEntity{
		WolfEnemyType: func(x, y int, id string) platformer.PlatformerActorEntity {
			enemy, err := NewWolfEnemy(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			player, _ := ctx.ActorManager.GetPlayer()
			enemy.SetTarget(player)
			return enemy
		},
		BatEnemyType: func(x, y int, id string) platformer.PlatformerActorEntity {
			enemy, err := NewBatEnemy(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return enemy
		},
		SwarmEnemyType: func(x, y int, id string) platformer.PlatformerActorEntity {
			enemy, err := NewSwarmEnemy(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return enemy
		},
	}
	return enemyMap
}
