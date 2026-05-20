package kitactors

// healthSetter is the minimal interface PlayerDeathBehavior needs.
type healthSetter interface {
	SetHealth(int)
}

// PlayerDeathBehavior defines the behavior when a player dies.
type PlayerDeathBehavior struct {
	player healthSetter
}

// NewPlayerDeathBehavior creates a new PlayerDeathBehavior for the given player.
func NewPlayerDeathBehavior(p healthSetter) *PlayerDeathBehavior {
	return &PlayerDeathBehavior{player: p}
}

// OnDie sets the player's health to 0.
func (tm *PlayerDeathBehavior) OnDie() {
	tm.player.SetHealth(0)
}
