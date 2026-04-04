package scene

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneManager struct {
	app.AppContextHolder

	current       navigation.Scene
	previousScene navigation.Scene
	factory       SceneFactory
	nextScene     navigation.Scene
	transitioner  navigation.Transition
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
	if m.current != nil {
		m.previousScene = m.current
	}

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

func (m *SceneManager) NavigateBack(sceneTransition navigation.Transition) {
	if m.previousScene == nil {
		log.Println("No previous scene to navigate back to.")
		return
	}

	sceneToLoad := m.previousScene
	m.previousScene = nil // Clear previous scene after navigating back

	if sceneTransition != nil {
		m.transitioner = sceneTransition
		m.nextScene = sceneToLoad
		m.transitioner.StartTransition(func() {
			m.SwitchTo(m.nextScene)
			m.nextScene = nil
		})
	} else {
		m.SwitchTo(sceneToLoad)
	}
}

func (m *SceneManager) AudioManager() *audio.AudioManager {
	if m.AppContext() == nil {
		return nil
	}
	return m.AppContext().AudioManager
}

// CurrentScene returns the currently active scene.
func (m *SceneManager) CurrentScene() navigation.Scene {
	return m.current
}

// IsTransitioning returns true if a scene transition is currently in progress.
func (m *SceneManager) IsTransitioning() bool {
	return m.transitioner != nil && m.transitioner.IsRunning()
}
