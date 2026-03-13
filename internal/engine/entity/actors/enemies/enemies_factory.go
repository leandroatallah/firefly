package enemies

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
)

type EnemyFactory[T actors.ActorEntity] struct {
	enemyMap EnemyMap[T]
}

func NewEnemyFactory[T actors.ActorEntity](enemyMap EnemyMap[T]) *EnemyFactory[T] {
	return &EnemyFactory[T]{enemyMap: enemyMap}
}

func (f *EnemyFactory[T]) Create(enemyType EnemyType, x, y int, id string) (T, error) {
	enemyFunc, ok := f.enemyMap[enemyType]
	if !ok {
		var zero T
		return zero, fmt.Errorf("unknown enemy type: %s", enemyType)
	}

	enemy := enemyFunc(x, y, id)

	return enemy, nil
}
