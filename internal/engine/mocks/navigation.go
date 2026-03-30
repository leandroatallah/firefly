package mocks

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
)

// MockScene implements navigation.Scene
type MockScene struct {
	DrawCount     int
	UpdateCount   int
	StartCount    int
	FinishCount   int
	AppContextSet any
}

func (s *MockScene) Draw(screen *ebiten.Image) { s.DrawCount++ }
func (s *MockScene) Update() error             { s.UpdateCount++; return nil }
func (s *MockScene) OnStart()                  { s.StartCount++ }
func (s *MockScene) OnFinish()                 { s.FinishCount++ }
func (s *MockScene) SetAppContext(ctx any)     { s.AppContextSet = ctx }

// MockSceneFactory implements navigation.SceneFactory
type MockSceneFactory struct {
	Scenes        map[navigation.SceneType]*MockScene
	AppContextSet any
}

func NewMockSceneFactory() *MockSceneFactory {
	return &MockSceneFactory{
		Scenes: make(map[navigation.SceneType]*MockScene),
	}
}

func (f *MockSceneFactory) Create(sceneType navigation.SceneType, freshInstance bool) (navigation.Scene, error) {
	if !freshInstance {
		if s, ok := f.Scenes[sceneType]; ok {
			return s, nil
		}
	}
	s := &MockScene{}
	s.SetAppContext(f.AppContextSet)
	f.Scenes[sceneType] = s
	return s, nil
}

func (f *MockSceneFactory) SetAppContext(ctx any) { f.AppContextSet = ctx }

// MockTransition implements navigation.Transition
type MockTransition struct {
	StartCalled   bool
	UpdateCount   int
	DrawCount     int
	StartCallback func()
}

func (t *MockTransition) Update()                   { t.UpdateCount++ }
func (t *MockTransition) Draw(screen *ebiten.Image) { t.DrawCount++ }
func (t *MockTransition) StartTransition(cb func()) {
	t.StartCalled = true
	t.StartCallback = cb
}
func (t *MockTransition) EndTransition(cb func()) { cb() }
func (t *MockTransition) IsRunning() bool         { return t.StartCalled }

// Complete simulates the transition completion
func (t *MockTransition) Complete() {
	if t.StartCallback != nil {
		t.StartCallback()
	}
}

// MockSceneManager implements navigation.SceneManager
type MockSceneManager struct {
	UpdateCalled   bool
	DrawCalled     bool
	LastSceneType  navigation.SceneType
	LastFresh      bool
	NavigateCalls  int
	CurrentSceneTo navigation.Scene
	Transitioning  bool
}

func (m *MockSceneManager) AudioManager() *audio.AudioManager { return nil }
func (m *MockSceneManager) Draw(screen *ebiten.Image)         { m.DrawCalled = true }
func (m *MockSceneManager) NavigateTo(sceneType navigation.SceneType, trans navigation.Transition, fresh bool) {
	m.NavigateCalls++
	m.LastSceneType = sceneType
	m.LastFresh = fresh
}
func (m *MockSceneManager) NavigateBack(trans navigation.Transition) {}
func (m *MockSceneManager) SwitchTo(scene navigation.Scene)          {}
func (m *MockSceneManager) Update() error {
	m.UpdateCalled = true
	return nil
}
func (m *MockSceneManager) CurrentScene() navigation.Scene { return m.CurrentSceneTo }
func (m *MockSceneManager) IsTransitioning() bool          { return m.Transitioning }
