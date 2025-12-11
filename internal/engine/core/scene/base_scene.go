// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/imagemanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

type BaseScene struct {
	count          int
	Manager        navigation.SceneManager
	audiomanager   *audiomanager.AudioManager
	imagemanager   *imagemanager.ImageManager
	space          *physics.Space
	AppContext     *core.AppContext
	IsKeysDisabled bool
}

func NewScene() *BaseScene {
	return &BaseScene{space: physics.NewSpace()}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {}

func (s *BaseScene) Update() error {
	return nil
}

func (s *BaseScene) OnStart() {
	// Init asset managers
	s.audiomanager = s.AppContext.AudioManager
	s.imagemanager = s.AppContext.ImageManager
}

func (s *BaseScene) OnFinish() {}

func (s *BaseScene) Exit() {}

func (s *BaseScene) AddBoundaries(boundaries ...body.MovableCollidable) {
	space := s.PhysicsSpace()
	for _, o := range boundaries {
		space.AddBody(o)
	}
}

func (s *BaseScene) SetAppContext(appContext any) {
	s.AppContext = appContext.(*core.AppContext)
	s.Manager = s.AppContext.SceneManager
	s.audiomanager = s.AppContext.AudioManager
	s.imagemanager = s.AppContext.ImageManager
}

func (s *BaseScene) PhysicsSpace() *physics.Space {
	if s.space == nil {
		s.space = physics.NewSpace()
	}
	return s.space
}

func (s *BaseScene) AudioManager() *audiomanager.AudioManager {
	return s.audiomanager
}

func (s *BaseScene) ImageManager() *imagemanager.ImageManager {
	return s.imagemanager
}

func (s *BaseScene) EnableKeys() {
	s.IsKeysDisabled = false
}

func (s *BaseScene) DisableKeys() {
	s.IsKeysDisabled = true
}
