package enemies

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/actors"
)

// To be initialized on game package.
type EnemyType int
type EnemyMap map[EnemyType]actors.ActorEntity

type EnemyFactory struct {
	enemyMap EnemyMap
}

func NewEnemyFactory(enemyMap EnemyMap) *EnemyFactory {
	return &EnemyFactory{enemyMap: enemyMap}
}

func (f *EnemyFactory) Create(enemyType EnemyType) (actors.ActorEntity, error) {
	enemy, ok := f.enemyMap[enemyType]
	if !ok {
		return nil, fmt.Errorf("unknown enemy type")
	}

	return enemy, nil
}
