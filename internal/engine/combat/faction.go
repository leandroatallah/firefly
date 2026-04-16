package combat

// Faction identifies which side an actor or projectile belongs to.
// The zero value is FactionNeutral, which is safe by default.
type Faction int

const (
	FactionNeutral Faction = iota
	FactionPlayer
	FactionEnemy
)
