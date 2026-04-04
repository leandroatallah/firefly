// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	"github.com/hajimehoshi/ebiten/v2"
)

type BaseScene struct {
	app.AppContextHolder

	count          int
	IsKeysDisabled bool

	scheduledActions []scheduledAction
}

type scheduledAction struct {
	targetFrame uint64
	action      func()
}

func NewScene() *BaseScene {
	return &BaseScene{}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {}

func (s *BaseScene) Update() error {
	// Update scheduled actions
	for i := 0; i < len(s.scheduledActions); i++ {
		if s.AppContext().FrameCount >= s.scheduledActions[i].targetFrame {
			s.scheduledActions[i].action()
			s.scheduledActions = append(s.scheduledActions[:i], s.scheduledActions[i+1:]...)
			i--
		}
	}
	return nil
}

func (s *BaseScene) Schedule(delay time.Duration, action func()) {
	target := s.AppContext().FrameCount + uint64(timing.FromDuration(delay))
	s.scheduledActions = append(s.scheduledActions, scheduledAction{
		targetFrame: target,
		action:      action,
	})
}

func (s *BaseScene) OnStart() {
	s.AppContext().Space.Clear()
	if s.AppContext().ActorManager != nil {
		s.AppContext().ActorManager.Clear()
	}
}

func (s *BaseScene) OnFinish() {}

func (s *BaseScene) Exit() {}

func (s *BaseScene) AddBoundaries(boundaries ...body.MovableCollidable) {
	space := s.PhysicsSpace()
	for _, o := range boundaries {
		space.AddBody(o)
	}
}

func (s *BaseScene) PhysicsSpace() body.BodiesSpace {
	return s.AppContext().Space
}

func (s *BaseScene) EnableKeys() {
	s.IsKeysDisabled = false
}

func (s *BaseScene) DisableKeys() {
	s.IsKeysDisabled = true
}

// PauseAllMusic pauses all music.
// Useful for scenes that need to control music manually.
func (s *BaseScene) PauseAllMusic() {
	ctx := s.AppContext()
	if ctx == nil || ctx.AudioManager == nil {
		return
	}
	ctx.AudioManager.PauseAll()
}

// PlayMusic plays music with rewind control.
// If rewind=false and music is already playing, does nothing.
// If rewind=true, restarts music from the beginning.
func (s *BaseScene) PlayMusic(path string, rewind bool) {
	s.PlayMusicWithLoop(path, false, rewind)
}

// PlayMusicWithLoop plays music with loop and rewind control.
func (s *BaseScene) PlayMusicWithLoop(path string, loop bool, rewind bool) {
	if path == "" {
		return
	}

	ctx := s.AppContext()
	if ctx == nil || ctx.AudioManager == nil {
		return
	}

	if !rewind && ctx.AudioManager.IsPlaying(path) {
		return
	}

	ctx.AudioManager.PlayMusic(path, loop)
}

// VFXManager returns the VFX manager. Override in subclasses.
func (s *BaseScene) VFXManager() interface{} {
	return s.count
}
