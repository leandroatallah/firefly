package sequences

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

// PlayMusicCommand plays background music.
// If music is already playing and Rewind=false, does nothing (continues from current position).
// If Rewind=true, restarts music from the beginning.
type PlayMusicCommand struct {
	Path   string  `json:"path"`
	Rewind bool    `json:"rewind"`
	Volume float64 `json:"volume"`
	Loop   bool    `json:"loop"`
}

func (c *PlayMusicCommand) Init(appContext any) {
	am := appContext.(*app.AppContext).AudioManager

	// If already playing and not rewinding, do nothing
	if !c.Rewind && am.IsPlaying(c.Path) {
		return
	}

	if c.Volume > 0 {
		am.SetVolume(c.Volume)
	}
	am.PlayMusic(c.Path, c.Loop)
}

func (c *PlayMusicCommand) Update() bool {
	return true // instant command
}

// PauseAllMusicCommand pauses all currently playing music.
// Useful for dramatic moments or transitions.
type PauseAllMusicCommand struct{}

func (c *PauseAllMusicCommand) Init(appContext any) {
	appContext.(*app.AppContext).AudioManager.PauseAll()
}

func (c *PauseAllMusicCommand) Update() bool {
	return true // instant command
}

type FadeOutAllMusicCommand struct {
	Duration int
}

func (c *FadeOutAllMusicCommand) Init(appContext any) {
	appContext.(*app.AppContext).AudioManager.FadeOutAll(timing.ToDuration(c.Duration))
}

func (c *FadeOutAllMusicCommand) Update() bool {
	return true
}
