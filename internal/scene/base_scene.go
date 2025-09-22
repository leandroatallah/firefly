// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/physics"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update() error
	OnStart()
	OnFinish()
	Next() Scene
}

type BaseScene struct {
	boundaries []physics.Body
	count      int
	nextScene  Scene
}

func NewScene() *BaseScene {
	return &BaseScene{}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) Update() error {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) OnStart() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) OnFinish() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) Exit() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) AddBoundaries(boundaries ...physics.Body) {
	for _, o := range boundaries {
		s.boundaries = append(s.boundaries, o)
	}
}

func (s *BaseScene) Next() Scene {
	return s.nextScene
}
