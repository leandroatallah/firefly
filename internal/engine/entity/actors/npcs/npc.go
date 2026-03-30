package npcs

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

type NpcType string

type NpcMap[T actors.ActorEntity] map[NpcType]func(x, y int, id string) T
