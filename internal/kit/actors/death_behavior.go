package kitactors

import "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"

// PlayerDeathBehavior defines the behavior when a player dies.
type PlayerDeathBehavior struct {
	player platformer.PlatformerActorEntity
}

// NewPlayerDeathBehavior creates a new PlayerDeathBehavior for the given player.
func NewPlayerDeathBehavior(p platformer.PlatformerActorEntity) *PlayerDeathBehavior {
	return &PlayerDeathBehavior{player: p}
}

// OnDie sets the player's health to 0.
func (tm *PlayerDeathBehavior) OnDie() {
	tm.player.SetHealth(0)
}
