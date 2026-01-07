package enemies

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// To be initialized on game package.
type EnemyType string
type EnemyMap[T actors.ActorEntity] map[EnemyType]func(x, y int, id string) T

type BaseEnemy struct {
	actors.Character
}

func NewBaseEnemy() *BaseEnemy {
	return &BaseEnemy{}
}

// Character Methods
func (e *BaseEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}
