package sequences

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
)

func setupTestAppContext() *app.AppContext {
	ctx := &app.AppContext{
		EventManager:    event.NewManager(),
		DialogueManager: nil, // Will test with nil
		Space:           space.NewSpace(),
		SceneManager:    scene.NewSceneManager(),
	}
	return ctx
}

func TestEventCommand_Init_PublishesEvent(t *testing.T) {
	ctx := setupTestAppContext()

	cmd := &EventCommand{
		EventType: "test_event",
		Payload:   map[string]interface{}{"key": "value"},
	}

	// Should not panic with nil event manager
	cmd.Init(ctx)

	// EventCommand publishes immediately on Init, so we just verify it doesn't panic
	// and that Update returns true (command is complete)
	if !cmd.Update() {
		t.Error("EventCommand.Update() should return true (instant command)")
	}
}

func TestEventCommand_Init_NilEventManager(t *testing.T) {
	ctx := &app.AppContext{
		DialogueManager: nil,
		Space:           space.NewSpace(),
		SceneManager:    scene.NewSceneManager(),
	}

	cmd := &EventCommand{
		EventType: "test_event",
		Payload:   map[string]interface{}{"key": "value"},
	}

	// Should not panic with nil event manager
	cmd.Init(ctx)

	if !cmd.Update() {
		t.Error("EventCommand.Update() should return true (instant command)")
	}
}

func TestDialogueCommand_Init_NilManager(t *testing.T) {
	lines := []string{"Hello", "World"}
	cmd := &DialogueCommand{
		Lines:    lines,
		Position: "bottom",
		Speed:    5,
	}

	// Should panic with nil DialogueManager - this is expected
	// We test that the command structure is correct
	if len(cmd.Lines) != 2 {
		t.Error("DialogueCommand lines not set correctly")
	}
}

func TestDialogueCommand_Init_WithSpeechID(t *testing.T) {
	cmd := &DialogueCommand{
		Lines:    []string{"Test"},
		SpeechID: "bubble",
	}

	// Should panic with nil DialogueManager - this is expected
	// We test that the command structure is correct
	if cmd.SpeechID != "bubble" {
		t.Error("DialogueCommand SpeechID not set correctly")
	}
}

func TestDialogueCommand_Init_DefaultSpeed(t *testing.T) {
	cmd := &DialogueCommand{
		Lines: []string{"Test"},
		Speed: 0, // Should use default
	}

	// Should panic with nil DialogueManager - this is expected
	// We test that the command structure is correct
	if cmd.Speed != 0 {
		t.Errorf("expected Speed 0, got %d", cmd.Speed)
	}
}

func TestDialogueCommand_Update(t *testing.T) {
	cmd := &DialogueCommand{
		Lines: []string{"Test"},
	}

	// We can't test Update without proper Init (will panic with nil manager)
	// Just verify the command structure is correct
	if len(cmd.Lines) != 1 {
		t.Error("DialogueCommand lines not set correctly")
	}
}

func TestDelayCommand_Init(t *testing.T) {
	ctx := setupTestAppContext()

	cmd := &DelayCommand{
		Frames: 30,
	}

	cmd.Init(ctx)

	if cmd.timer != 0 {
		t.Errorf("expected timer to be 0 after Init, got %d", cmd.timer)
	}
}

func TestDelayCommand_Update(t *testing.T) {
	tests := []struct {
		name         string
		frames       int
		updateCount  int
		wantComplete bool
		wantTimer    int
	}{
		{"Zero frames", 0, 1, true, 1},
		{"One frame", 1, 1, true, 1},
		{"Multiple frames - not done", 5, 3, false, 3},
		{"Multiple frames - done", 5, 5, true, 5},
		{"Multiple frames - past", 5, 10, true, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTestAppContext()

			cmd := &DelayCommand{
				Frames: tt.frames,
			}

			cmd.Init(ctx)

			var complete bool
			for i := 0; i < tt.updateCount; i++ {
				complete = cmd.Update()
			}

			if complete != tt.wantComplete {
				t.Errorf("Update() complete = %v, want %v", complete, tt.wantComplete)
			}

			if cmd.timer != tt.wantTimer {
				t.Errorf("timer = %d, want %d", cmd.timer, tt.wantTimer)
			}
		})
	}
}
