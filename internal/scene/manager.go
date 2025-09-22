package scene

import "github.com/hajimehoshi/ebiten/v2"

type SceneManager struct {
	current Scene
}

func NewSceneManager() *SceneManager {
	return &SceneManager{}
}

func (m *SceneManager) Update() error {
	if err := m.current.Update(); err != nil {
		return err
	}

	if next := m.current.Next(); next != nil {
		m.GoTo(next)
	}
	return nil
}
func (m *SceneManager) Draw(screen *ebiten.Image) {
	m.current.Draw(screen)
}

func (m *SceneManager) GoTo(scene Scene) {
	if m.current != nil {
		m.current.OnFinish()
	}

	m.current = scene

	if m.current != nil {
		m.current.OnStart()
	}
}
