package gameplayermethods

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
)

type PlayerDeathBehavior struct {
	player platformer.PlatformerActorEntity
}

func NewPlayerDeathBehavior(p platformer.PlatformerActorEntity) *PlayerDeathBehavior {
	tm := &PlayerDeathBehavior{
		player: p,
	}
	return tm
}

func (tm *PlayerDeathBehavior) OnDie() {
	tm.player.SetHealth(0)
}
