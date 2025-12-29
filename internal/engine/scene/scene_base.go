// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

type BaseScene struct {
	app.AppContextHolder

	count          int
	space          *space.Space
	IsKeysDisabled bool
}

func NewScene() *BaseScene {
	return &BaseScene{space: space.NewSpace()}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {}

func (s *BaseScene) Update() error {
	return nil
}

func (s *BaseScene) OnStart() {}

func (s *BaseScene) OnFinish() {}

func (s *BaseScene) Exit() {}

func (s *BaseScene) AddBoundaries(boundaries ...body.MovableCollidable) {
	space := s.PhysicsSpace()
	for _, o := range boundaries {
		space.AddBody(o)
	}
}

func (s *BaseScene) PhysicsSpace() *space.Space {
	if s.space == nil {
		s.space = space.NewSpace()
	}
	return s.space
}

func (s *BaseScene) EnableKeys() {
	s.IsKeysDisabled = false
}

func (s *BaseScene) DisableKeys() {
	s.IsKeysDisabled = true
}
