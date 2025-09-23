package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneManager struct {
	current Scene
	factory SceneFactory
}

func NewSceneManager() *SceneManager {
	return &SceneManager{}
}

func (m *SceneManager) Update() error {
	if m.current == nil {
		return nil
	}
	if err := m.current.Update(); err != nil {
		return err
	}

	return nil
}
func (m *SceneManager) Draw(screen *ebiten.Image) {
	if m.current == nil {
		return
	}
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

func (m *SceneManager) SetFactory(factory SceneFactory) {
	m.factory = factory
}

func (m *SceneManager) GoToScene(sceneType SceneType) {
	scene, err := m.factory.Create(sceneType)
	if err != nil {
		log.Fatalf("Error creating scene: %v", err)
	}
	m.GoTo(scene)
}