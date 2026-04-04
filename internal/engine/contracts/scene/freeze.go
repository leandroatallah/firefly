package scene

// Freezable allows a scene to be temporarily frozen for a hit-stop effect.
type Freezable interface {
	// FreezeFrame pauses scene updates for the given number of frames.
	FreezeFrame(durationFrames int)
	// IsFrozen reports whether the scene is currently frozen.
	IsFrozen() bool
}
