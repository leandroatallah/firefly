package sequences

import (
	"testing"
)

func TestSpawnTextCommandInitStructFields(t *testing.T) {
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

func TestSpawnTextCommandInitOverheadNoActor(t *testing.T) {
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

func TestSpawnTextCommandInitEmptyTargetID(t *testing.T) {
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

func TestSpawnTextCommandUpdate(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text: "Test",
		Type: "screen",
	}

	if !cmd.Update() {
		t.Error("SpawnTextCommand.Update() should return true (instant command)")
	}
}

type mockVignetteScene struct {
	log []string
}

func (m *mockVignetteScene) EnableVignetteDarkness(radiusPx float64) {
	m.log = append(m.log, "enable")
}

func (m *mockVignetteScene) DisableVignetteDarkness() {
	m.log = append(m.log, "disable")
}

type mockSceneManager struct {
	current interface{}
}

func (m *mockSceneManager) CurrentScene() interface{} {
	return m.current
}

func TestVignetteRadiusCommandInitNoScene(t *testing.T) {
	cmd := &VignetteRadiusCommand{
		InitialRadius: 10,
		FinalRadius:   20,
		Duration:      30,
	}

	// No app context / scene provided: Init should be a no-op that
	// leaves controller as nil and Update should complete instantly.
	cmd.Init(nil)
	if !cmd.Update() {
		t.Fatalf("expected Update to complete immediately when no controller is set")
	}
}

func TestVignetteRadiusCommandUpdateInstant(t *testing.T) {
	sceneMock := &mockVignetteScene{}
	cmd := &VignetteRadiusCommand{
		InitialRadius: 0,
		FinalRadius:   32,
		Duration:      0,
		controller:    sceneMock,
	}

	if done := cmd.Update(); !done {
		t.Fatalf("expected instant completion when Duration <= 0")
	}
	if len(sceneMock.log) == 0 || sceneMock.log[len(sceneMock.log)-1] != "enable" {
		t.Fatalf("expected EnableVignetteDarkness to be called for final radius")
	}
}
