package movement

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

func TestNewMovementState_BuiltInStates(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0)
	target := &mockActor{}
	target.SetPosition(100, 100)

	tests := []struct {
		name    string
		state   MovementStateEnum
		wantErr bool
	}{
		{"Idle", Idle, false},
		{"Chase", Chase, false},
		{"DumbChase", DumbChase, false},
		{"Avoid", Avoid, false},
		{"Patrol", Patrol, false},
		{"SideToSide", SideToSide, false},
		{"Follow", Follow, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := NewMovementState(actor, tt.state, target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMovementState(%v) error = %v, wantErr %v", tt.state, err, tt.wantErr)
			}
			if state == nil && !tt.wantErr {
				t.Error("NewMovementState() returned nil state")
			}
			if state != nil && state.State() != tt.state {
				t.Errorf("NewMovementState() state = %v, want %v", state.State(), tt.state)
			}
		})
	}
}

func TestNewMovementState_RegisteredState(t *testing.T) {
	// Save and restore registry state
	savedConstructors := movementStateConstructors
	savedEnums := movementStateEnums
	savedNextValue := nextMovementEnumValue
	defer func() {
		movementStateConstructors = savedConstructors
		movementStateEnums = savedEnums
		nextMovementEnumValue = savedNextValue
	}()

	// Reset for testing
	movementStateConstructors = make(map[MovementStateEnum]MovementStateConstructor)
	movementStateEnums = make(map[string]MovementStateEnum)
	nextMovementEnumValue = 0

	// Register a custom state
	customState := MovementStateEnum(200) // Use a value that won't conflict
	movementStateConstructors[customState] = func(b BaseMovementState) MovementState {
		return &customTestState{BaseMovementState: b}
	}

	actor := &mockActor{speed: 5}
	target := &mockActor{}

	state, err := NewMovementState(actor, customState, target)
	if err != nil {
		t.Errorf("NewMovementState() with registered state returned error: %v", err)
	}
	if state == nil {
		t.Error("NewMovementState() with registered state returned nil")
	}
}

func TestNewMovementState_UnknownState(t *testing.T) {
	actor := &mockActor{speed: 5}
	target := &mockActor{}

	// Use an unregistered state value
	unknownState := MovementStateEnum(9999)

	state, err := NewMovementState(actor, unknownState, target)
	if err == nil {
		t.Error("NewMovementState() with unknown state should return error")
	}
	if state != nil {
		t.Error("NewMovementState() with unknown state should return nil")
	}
}

func TestNewMovementState_WithOptions(t *testing.T) {
	actor := &mockActor{speed: 5}
	target := &mockActor{}

	optionApplied := false
	testOption := func(s MovementState) {
		optionApplied = true
	}

	state, err := NewMovementState(actor, Idle, target, testOption)
	if err != nil {
		t.Errorf("NewMovementState() returned error: %v", err)
	}
	if state == nil {
		t.Error("NewMovementState() returned nil")
	}
	if !optionApplied {
		t.Error("NewMovementState() did not apply options")
	}
}

func TestNewMovementState_WithNilOption(t *testing.T) {
	actor := &mockActor{speed: 5}
	target := &mockActor{}

	// Should not panic with nil option
	state, err := NewMovementState(actor, Idle, target, nil)
	if err != nil {
		t.Errorf("NewMovementState() with nil option returned error: %v", err)
	}
	if state == nil {
		t.Error("NewMovementState() with nil option returned nil")
	}
}

func TestNewMovementState_WithMultipleOptions(t *testing.T) {
	actor := &mockActor{speed: 5}
	target := &mockActor{}

	option1Applied := false
	option2Applied := false

	option1 := func(s MovementState) {
		option1Applied = true
	}
	option2 := func(s MovementState) {
		option2Applied = true
	}

	state, err := NewMovementState(actor, Chase, target, option1, option2)
	if err != nil {
		t.Errorf("NewMovementState() returned error: %v", err)
	}
	if state == nil {
		t.Error("NewMovementState() returned nil")
	}
	if !option1Applied || !option2Applied {
		t.Error("NewMovementState() did not apply all options")
	}
}

// customTestState is a test implementation of MovementState
type customTestState struct {
	BaseMovementState
}

func (s *customTestState) Move(space body.BodiesSpace) {
	// No-op for testing
}
