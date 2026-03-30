package gameentitytypes

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// EnemyActor is an interface that identifies an actor as an enemy.
type EnemyActor interface {
	actors.ActorEntity
	IsEnemy() bool
}
