package navigation

import (
	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/hajimehoshi/ebiten/v2"
)

// SceneType is an integer identifier for a registered scene.
type SceneType int

// Scene represents a single screen or game state with a defined lifecycle.
type Scene interface {
	// Draw renders the scene to the given screen image.
	Draw(screen *ebiten.Image)
	// Update advances the scene by one game tick, returning any error.
	Update() error
	// OnStart is called once when the scene becomes active.
	OnStart()
	// OnFinish is called once when the scene is about to be replaced.
	OnFinish()

	// SetAppContext injects the shared application context into the scene.
	SetAppContext(appContext any)
}

// SceneFactory creates Scene instances by type.
type SceneFactory interface {
	// Create returns a Scene for the given type, optionally as a fresh instance.
	Create(sceneType SceneType, freshInstance bool) (Scene, error)

	// SetAppContext injects the shared application context into the factory.
	SetAppContext(appContext any)
}

// SceneMap maps each SceneType to a constructor function that produces a Scene.
type SceneMap map[SceneType]func() Scene

// SceneManager controls scene navigation, transitions, and the active scene lifecycle.
type SceneManager interface {
	// AudioManager returns the shared audio manager.
	AudioManager() audio.Manager
	// Draw renders the current scene (and any active transition) to the screen.
	Draw(screen *ebiten.Image)
	// NavigateTo transitions to the scene identified by sceneType.
	NavigateTo(sceneType SceneType, sceneTransition Transition, freshInstance bool)
	// NavigateBack transitions to the previously active scene.
	NavigateBack(sceneTransition Transition)
	// SetFactory(factory SceneFactory)
	// SwitchTo immediately replaces the current scene without a transition.
	SwitchTo(scene Scene)
	// Update advances the scene manager and the active scene by one tick.
	Update() error
	// CurrentScene returns the scene that is currently active.
	CurrentScene() Scene
	// IsTransitioning reports whether a scene transition is in progress.
	IsTransitioning() bool
}

// Transition defines an animated transition played between scene changes.
type Transition interface {
	// Update advances the transition animation by one tick.
	Update()
	// Draw renders the transition overlay to the screen.
	Draw(screen *ebiten.Image)
	// StartTransition begins the transition and calls onComplete when the midpoint is reached.
	StartTransition(func())
	// EndTransition begins the exit phase and calls onComplete when finished.
	EndTransition(func())
	// IsRunning reports whether the transition is currently playing.
	IsRunning() bool
}
