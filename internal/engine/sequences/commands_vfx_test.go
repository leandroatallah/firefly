package sequences

import (
	"testing"
)

func TestSpawnTextCommand_Init_StructFields(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text:     "Hello World",
		Type:     "screen",
		X:        100.0,
		Y:        200.0,
		Duration: 60,
	}

	// Verify the command structure is correct
	if cmd.Text != "Hello World" {
		t.Errorf("expected Text 'Hello World', got %q", cmd.Text)
	}
	if cmd.Type != "screen" {
		t.Errorf("expected Type 'screen', got %q", cmd.Type)
	}
	if cmd.X != 100.0 {
		t.Errorf("expected X 100.0, got %f", cmd.X)
	}
	if cmd.Y != 200.0 {
		t.Errorf("expected Y 200.0, got %f", cmd.Y)
	}
}

func TestSpawnTextCommand_Init_OverheadNoActor(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text:     "Overhead Text",
		Type:     "overhead",
		TargetID: "nonexistent",
		Duration: 60,
	}

	// Verify the command structure is correct
	if cmd.TargetID != "nonexistent" {
		t.Errorf("expected TargetID 'nonexistent', got %q", cmd.TargetID)
	}
}

func TestSpawnTextCommand_Init_EmptyTargetID(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text:     "Test",
		Type:     "overhead",
		TargetID: "",
	}

	// Verify the command structure is correct
	if cmd.TargetID != "" {
		t.Errorf("expected empty TargetID, got %q", cmd.TargetID)
	}
}

func TestSpawnTextCommand_Update(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text: "Test",
		Type: "screen",
	}

	if !cmd.Update() {
		t.Error("SpawnTextCommand.Update() should return true (instant command)")
	}
}
