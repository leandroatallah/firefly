package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
)

type SceneManager struct {
	current      navigation.Scene
	factory      SceneFactory
	nextScene    navigation.Scene
	transitioner navigation.Transition
	appContext   *core.AppContext
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

func (m *SceneManager) SetAppContext(appContext *core.AppContext) {
	m.appContext = appContext
	m.factory.SetAppContext(appContext)
}

func (m *SceneManager) AudioManager() *audiomanager.AudioManager {
	return m.appContext.AudioManager
}
