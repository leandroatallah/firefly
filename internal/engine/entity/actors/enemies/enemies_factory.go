package enemies

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
)

type EnemyFactory struct {
	enemyMap EnemyMap
}

func NewEnemyFactory(enemyMap EnemyMap) *EnemyFactory {
	return &EnemyFactory{enemyMap: enemyMap}
}

func (f *EnemyFactory) Create(enemyType EnemyType, x, y int, id string) (actors.ActorEntity, error) {
	enemyFunc, ok := f.enemyMap[enemyType]
	if !ok {
		return nil, fmt.Errorf("unknown enemy type: %s", enemyType)
	}

	enemy := enemyFunc(x, y, id)

	return enemy, nil
}

