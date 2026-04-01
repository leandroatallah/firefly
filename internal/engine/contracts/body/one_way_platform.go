package body

// OneWayPlatform is a Body that can be passed through from below or by
// explicit drop-through input, but acts as solid ground from above.
type OneWayPlatform interface {
	Body
	Collidable
	IsOneWay() bool
	// SetPassThrough registers actor for pass-through for the given number of frames.
	SetPassThrough(actor Collidable, frames int)
	// IsPassThrough reports whether actor is currently passing through this platform.
	IsPassThrough(actor Collidable) bool
	// Update decrements all pass-through countdowns and removes expired entries.
	Update()
}
