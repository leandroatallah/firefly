package scene

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

func TestBaseScene_Lifecycle(t *testing.T) {
	ctx := &app.AppContext{
		Space:        space.NewSpace(),
		ActorManager: actors.NewManager(),
	}
	s := NewScene()
	s.SetAppContext(ctx)

	// Test OnStart clears space and actors
	ctx.Space.AddBody(&mocks.MockActor{Id: "test"})
	ctx.ActorManager.Register(&mocks.MockActor{Id: "test"})
	s.OnStart()
	if len(ctx.Space.Bodies()) != 0 {
		t.Error("OnStart did not clear physics space")
	}
	if _, found := ctx.ActorManager.Find("test"); found {
		t.Error("OnStart did not clear actor manager")
	}

	// Test PhysicsSpace
	if s.PhysicsSpace() != ctx.Space {
		t.Error("PhysicsSpace() returned wrong space")
	}

	// Test AddBoundaries
	actor := &mocks.MockActor{Id: "boundary"}
	s.AddBoundaries(actor)
	if ctx.Space.Find("boundary") == nil {
		t.Error("AddBoundaries did not add body to space")
	}

	// Test Key Toggling
	s.DisableKeys()
	if !s.IsKeysDisabled {
		t.Error("DisableKeys failed")
	}
	s.EnableKeys()
	if s.IsKeysDisabled {
		t.Error("EnableKeys failed")
	}

	// Test dummy methods for coverage
	s.Draw(ebiten.NewImage(1, 1))
	s.OnFinish()
	s.Exit()
	s.VFXManager()
}

func TestBaseScene_Schedule(t *testing.T) {
	ctx := &app.AppContext{FrameCount: 10}
	s := NewScene()
	s.SetAppContext(ctx)

	called := false
	// Schedule for +60 frames (1 second at 60 TPS)
	s.Schedule(time.Second, func() { called = true })

	if len(s.scheduledActions) != 1 {
		t.Fatal("action not scheduled")
	}

	s.Update()
	if called {
		t.Error("action called too early")
	}

	ctx.FrameCount = 70
	s.Update()
	if !called {
		t.Error("action not called after frame count reached")
	}
	if len(s.scheduledActions) != 0 {
		t.Error("action not removed after execution")
	}
}

func TestBaseScene_Music(t *testing.T) {
	// scene_test already has some audio manager mocking patterns, 
	// but BaseScene uses AppContext.AudioManager directly.
	// Since I haven't centralized MockAudioManager yet, I'll skip complex music tests
	// or just smoke test with nil checks if possible.
	
	s := NewScene()
	s.SetAppContext(&app.AppContext{}) // Nil AudioManager
	
	// Should not panic
	s.PauseAllMusic()
	s.PlayMusic("test.ogg", true)
	s.PlayMusicWithLoop("test.ogg", true, true)
}
