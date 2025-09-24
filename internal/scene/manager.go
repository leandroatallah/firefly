package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/transition"
)

type SceneManager struct {
	current Scene
	factory SceneFactory

	nextScene    Scene
	transitioner transition.Transition
}

func NewSceneManager() *SceneManager {
	m := &SceneManager{}
	return m
}

func (m *SceneManager) Update() error {
	if m.transitioner != nil {
		m.transitioner.Update()
	}

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
	if m.transitioner != nil {
		m.transitioner.Draw(screen)
	}
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

func (m *SceneManager) GoToScene(sceneType SceneType, sceneTransition transition.Transition) {
	scene, err := m.factory.Create(sceneType)
	if err != nil {
		log.Fatalf("Error creating scene: %v", err)
	}

	if sceneTransition != nil {
		m.transitioner = sceneTransition
		m.nextScene = scene
		m.transitioner.StartTransition(func() {
			m.GoTo(m.nextScene)
			m.nextScene = nil
		})
	} else {
		m.GoTo(scene)
	}
}
