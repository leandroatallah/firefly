package gamenpcs

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/npcs"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
)

const (
	ClimberNpcType npcs.NpcType = "CLIMBER"
)

func InitNpcMap(ctx *app.AppContext) npcs.NpcMap[platformer.PlatformerActorEntity] {
	npcMap := map[npcs.NpcType]func(x, y int, id string) platformer.PlatformerActorEntity{
		ClimberNpcType: func(x, y int, id string) platformer.PlatformerActorEntity {
			npc, err := gameplayer.NewClimberPlayer(ctx)
			if err != nil {
				log.Fatal(err)
			}
			npc.SetPosition(x, y)
			npc.SetID(id)
			return npc
		},
	}
	return npcMap
}
