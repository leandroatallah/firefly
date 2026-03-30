package gamestates

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

var (
	Dying   actors.ActorStateEnum
	Dead    actors.ActorStateEnum
	Exiting actors.ActorStateEnum
)

func init() {
	Dying = actors.Dying
	Dead = actors.Dead
	Exiting = actors.Exiting
}
