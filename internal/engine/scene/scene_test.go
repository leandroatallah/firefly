package scene

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

func TestSceneManagerNavigateWithTransition(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	factory := mocks.NewMockSceneFactory()
	factory.SetAppContext(ctx)
	manager.SetFactory(factory)

	sceneType := navigation.SceneType(1)
	transition := &mocks.MockTransition{}

	if err := manager.Update(); err != nil {
		t.Fatalf("initial Update error: %v", err)
	}

	manager.NavigateTo(sceneType, transition, false)

	created := factory.Scenes[sceneType]
	if created == nil {
		t.Fatalf("expected factory to create scene")
	}
	if !transition.StartCalled {
		t.Fatalf("expected transition StartTransition to be called")
	}
	if created.StartCount != 0 {
		t.Fatalf("scene should not start before transition completes")
	}

	if err := manager.Update(); err != nil {
		t.Fatalf("Update error while transition active: %v", err)
	}
	if transition.UpdateCount == 0 {
		t.Fatalf("expected transition Update to be called during manager.Update")
	}

	transition.Complete()

	if err := manager.Update(); err != nil {
		t.Fatalf("Update error after transition completion: %v", err)
	}

	if created.StartCount == 0 {
		t.Fatalf("expected scene OnStart to be called after transition completion")
	}

	screen := ebiten.NewImage(1, 1)
	manager.Draw(screen)
	if created.DrawCount == 0 {
		t.Fatalf("expected scene Draw to be called after switch")
	}
}

func TestSceneManagerNavigateWithoutTransition(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	factory := mocks.NewMockSceneFactory()
	factory.SetAppContext(ctx)
	manager.SetFactory(factory)

	sceneType1 := navigation.SceneType(1)
	sceneType2 := navigation.SceneType(2)

	manager.NavigateTo(sceneType1, nil, false)
	s1 := factory.Scenes[sceneType1]
	if s1 == nil {
		t.Fatalf("expected first scene to be created")
	}
	if s1.StartCount != 1 {
		t.Fatalf("expected first scene OnStart to be called once, got %d", s1.StartCount)
	}

	manager.NavigateTo(sceneType2, nil, false)
	s2 := factory.Scenes[sceneType2]
	if s2 == nil {
		t.Fatalf("expected second scene to be created")
	}
	if s1.FinishCount != 1 {
		t.Fatalf("expected first scene OnFinish to be called when switching, got %d", s1.FinishCount)
	}
	if s2.StartCount != 1 {
		t.Fatalf("expected second scene OnStart to be called once, got %d", s2.StartCount)
	}
}

func TestSceneManagerNavigateBack(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	factory := mocks.NewMockSceneFactory()
	factory.SetAppContext(ctx)
	manager.SetFactory(factory)

	sceneType1 := navigation.SceneType(1)
	sceneType2 := navigation.SceneType(2)

	manager.NavigateTo(sceneType1, nil, false)
	s1 := factory.Scenes[sceneType1]

	manager.NavigateTo(sceneType2, nil, false)
	s2 := factory.Scenes[sceneType2]

	manager.NavigateBack(nil)

	if s2.FinishCount == 0 {
		t.Fatalf("expected second scene OnFinish to be called when navigating back")
	}
	if s1.StartCount < 2 {
		t.Fatalf("expected first scene OnStart to be called again when navigated back, got %d", s1.StartCount)
	}
}

func TestSceneManager_Properties(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	// Test AudioManager
	if manager.AudioManager() != nil {
		t.Error("expected nil AudioManager")
	}

	// Test CurrentScene
	if manager.CurrentScene() != nil {
		t.Error("expected initial CurrentScene to be nil")
	}

	scene := &mocks.MockScene{}
	manager.SwitchTo(scene)
	if manager.CurrentScene() != scene {
		t.Error("CurrentScene() did not return switched scene")
	}
}

func TestSceneManager_NavigateBack_NoHistory(t *testing.T) {
	manager := NewSceneManager()
	// Should not panic
	manager.NavigateBack(nil)
}

func TestBaseSceneOnStartClearsPhysicsSpace(t *testing.T) {
	sp := space.NewSpace()
	rect := bodyphysics.NewRect(0, 0, 10, 10)
	obstacle := bodyphysics.NewObstacleRect(rect)
	obstacle.SetID("obstacle")
	sp.AddBody(obstacle)

	ctx := &app.AppContext{
		Space: sp,
	}

	base := NewScene()
	base.SetAppContext(ctx)

	if len(sp.Bodies()) == 0 {
		t.Fatalf("precondition: expected space to contain bodies before OnStart")
	}

	base.OnStart()

	if len(sp.Bodies()) != 0 {
		t.Fatalf("expected BaseScene.OnStart to clear physics space")
	}
}

func TestBaseSceneScheduleExecutesActionAfterDelay(t *testing.T) {
	ctx := &app.AppContext{}
	base := NewScene()
	base.SetAppContext(ctx)

	called := 0
	delay := time.Second

	base.Schedule(delay, func() {
		called++
	})

	frames := timing.FromDuration(delay)
	for i := 0; i < frames; i++ {
		ctx.FrameCount++
		if err := base.Update(); err != nil {
			t.Fatalf("Update error at frame %d: %v", i, err)
		}
	}

	if called != 1 {
		t.Fatalf("expected scheduled action to be called once after delay, got %d", called)
	}

	ctx.FrameCount++
	if err := base.Update(); err != nil {
		t.Fatalf("Update error after delay: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected scheduled action not to be called more than once, got %d", called)
	}
}
