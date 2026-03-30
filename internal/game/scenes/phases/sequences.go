package gamescenephases

// deathSequence guards against multiple death triggers while the scene transition is in flight.
type deathSequence struct {
	active bool
}
