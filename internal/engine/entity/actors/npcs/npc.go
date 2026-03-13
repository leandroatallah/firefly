package npcs

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
)

type NpcType string

type NpcMap[T actors.ActorEntity] map[NpcType]func(x, y int, id string) T
