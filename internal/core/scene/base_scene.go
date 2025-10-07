// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/core"
	"github.com/leandroatallah/firefly/internal/navigation"
	"github.com/leandroatallah/firefly/internal/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type BaseScene struct {
	count        int
	Manager      navigation.SceneManager
	audiomanager *audiomanager.AudioManager
	space        *physics.Space
	appContext   *core.AppContext
}

func NewScene() *BaseScene {
	return &BaseScene{space: physics.NewSpace()}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {}

func (s *BaseScene) Update() error {
	return nil
}

func (s *BaseScene) OnStart() {}

func (s *BaseScene) OnFinish() {}

func (s *BaseScene) Exit() {}

func (s *BaseScene) AddBoundaries(boundaries ...physics.Body) {
	space := s.PhysicsSpace()
	for _, o := range boundaries {
		space.AddBody(o)
	}
}

func (s *BaseScene) SetAppContext(appContext any) {
	s.appContext = appContext.(*core.AppContext)
	s.Manager = s.appContext.SceneManager
	s.audiomanager = s.appContext.AudioManager
}

func (s *BaseScene) PhysicsSpace() *physics.Space {
	if s.space == nil {
		s.space = physics.NewSpace()
	}
	return s.space
}
