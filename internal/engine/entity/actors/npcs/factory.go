package npcs

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
)

type NpcFactory[T actors.ActorEntity] struct {
	npcMap NpcMap[T]
}

func NewNpcFactory[T actors.ActorEntity](npcMap NpcMap[T]) *NpcFactory[T] {
	return &NpcFactory[T]{npcMap: npcMap}
}

func (nf *NpcFactory[T]) Create(npcType NpcType, x, y int, id string) (T, error) {
	creator, found := nf.npcMap[npcType]
	if !found {
		return *new(T), fmt.Errorf("npc type %s not found", npcType)
	}
	return creator(x, y, id), nil
}
