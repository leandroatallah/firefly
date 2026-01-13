package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/audio"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
)

type SceneManager struct {
	app.AppContextHolder

	current      navigation.Scene
	factory      SceneFactory
	nextScene    navigation.Scene
	transitioner navigation.Transition
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

func (m *SceneManager) SwitchTo(scene navigation.Scene) {
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

func (m *SceneManager) NavigateTo(
	sceneType navigation.SceneType, sceneTransition navigation.Transition, freshInstance bool,
) {
	scene, err := m.factory.Create(sceneType, freshInstance)
	if err != nil {
		log.Fatalf("Error creating scene: %v", err)
	}

	if sceneTransition != nil {
		m.transitioner = sceneTransition
		m.nextScene = scene
		m.transitioner.StartTransition(func() {
			m.SwitchTo(m.nextScene)
			m.nextScene = nil
		})
	} else {
		m.SwitchTo(scene)
	}
}

func (m *SceneManager) AudioManager() *audio.AudioManager {
	return m.AppContext().AudioManager
}
