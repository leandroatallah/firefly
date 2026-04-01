package scene

type Freezable interface {
	FreezeFrame(durationFrames int)
	IsFrozen() bool
}
