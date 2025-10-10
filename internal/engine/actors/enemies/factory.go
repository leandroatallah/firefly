package enemies

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/actors"
)

type EnemyType int

const (
	BlueEnemyType EnemyType = iota
)

type EnemyFactory struct{}

func NewEnemyFactory() *EnemyFactory {
	return &EnemyFactory{}
}

func (f *EnemyFactory) Create(enemyType EnemyType, x, y int) (actors.ActorEntity, error) {
	switch enemyType {
	case BlueEnemyType:
		return NewBlueEnemy(x, y), nil
	default:
		return nil, fmt.Errorf("unknown enemy type")
	}
}
