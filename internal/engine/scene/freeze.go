package scene

type FreezeController struct {
	remaining int
}

func (f *FreezeController) FreezeFrame(durationFrames int) {
	if durationFrames <= 0 {
		return
	}
	f.remaining = durationFrames
}

func (f *FreezeController) IsFrozen() bool {
	return f.remaining > 0
}

func (f *FreezeController) Tick() {
	if f.remaining > 0 {
		f.remaining--
	}
}
