package actors

import (
	"testing"
)

func TestRegisterState_NewState(t *testing.T) {
	// Save current registry state
	savedConstructors := stateConstructors
	savedEnums := stateEnums
	savedNextValue := nextEnumValue
	
	// Clean up registry after test
	defer func() {
		stateConstructors = savedConstructors
		stateEnums = savedEnums
		nextEnumValue = savedNextValue
	}()

	// Reset for testing
	stateConstructors = make(map[ActorStateEnum]StateConstructor)
	stateEnums = make(map[string]ActorStateEnum)
	nextEnumValue = 0

	constructor := func(b BaseState) ActorState { return &IdleState{BaseState: b} }
	enum := RegisterState("test_state", constructor)

	if _, ok := stateEnums["test_state"]; !ok {
		t.Error("RegisterState() did not register state name")
	}

	if _, ok := stateConstructors[enum]; !ok {
		t.Error("RegisterState() did not register constructor")
	}
}

func TestRegisterState_DuplicateName(t *testing.T) {
	// Save current registry state
	savedConstructors := stateConstructors
	savedEnums := stateEnums
	savedNextValue := nextEnumValue
	
	// Clean up registry after test
	defer func() {
		stateConstructors = savedConstructors
		stateEnums = savedEnums
		nextEnumValue = savedNextValue
	}()

	// Reset for testing
	stateConstructors = make(map[ActorStateEnum]StateConstructor)
	stateEnums = make(map[string]ActorStateEnum)
	nextEnumValue = 0

	constructor1 := func(b BaseState) ActorState { return &IdleState{BaseState: b} }
	constructor2 := func(b BaseState) ActorState { return &WalkState{BaseState: b} }

	enum1 := RegisterState("duplicate_test", constructor1)
	enum2 := RegisterState("duplicate_test", constructor2)

	if enum1 != enum2 {
		t.Error("RegisterState() with duplicate name should return same enum")
	}

	// Constructor should be updated
	ctor, ok := stateConstructors[enum1]
	if !ok {
		t.Fatal("RegisterState() did not update constructor")
	}

	// Test that the new constructor works - just verify it doesn't return nil
	base := BaseState{state: enum1}
	state := ctor(base)
	if state == nil {
		t.Error("RegisterState() did not properly update constructor")
	}
}

func TestNewState_Success(t *testing.T) {
	tests := []struct {
		name    string
		state   ActorStateEnum
		wantErr bool
	}{
		{"idle state", Idle, false},
		{"walking state", Walking, false},
		{"jumping state", Jumping, false},
		{"falling state", Falling, false},
		{"landing state", Landing, false},
		{"hurted state", Hurted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := NewState(nil, tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewState() error = %v, wantErr %v", err, tt.wantErr)
			}
			if state == nil && !tt.wantErr {
				t.Error("NewState() returned nil state")
			}
			if state != nil && state.State() != tt.state {
				t.Errorf("NewState() state = %v, want %v", state.State(), tt.state)
			}
		})
	}
}

func TestNewState_Unregistered(t *testing.T) {
	// Use an unregistered state value
	unregisteredState := ActorStateEnum(9999)

	state, err := NewState(nil, unregisteredState)
	if err == nil {
		t.Error("NewState() with unregistered state should return error")
	}
	if state != nil {
		t.Error("NewState() with unregistered state should return nil state")
	}
}

func TestGetStateEnum_Existing(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		wantOk  bool
		wantVal ActorStateEnum
	}{
		{"idle", "idle", true, Idle},
		{"walk", "walk", true, Walking},
		{"jump", "jump", true, Jumping},
		{"fall", "fall", true, Falling},
		{"land", "land", true, Landing},
		{"hurt", "hurt", true, Hurted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := GetStateEnum(tt.state)
			if ok != tt.wantOk {
				t.Errorf("GetStateEnum(%q) ok = %v, want %v", tt.state, ok, tt.wantOk)
			}
			if tt.wantOk && val != tt.wantVal {
				t.Errorf("GetStateEnum(%q) = %v, want %v", tt.state, val, tt.wantVal)
			}
		})
	}
}

func TestGetStateEnum_NonExistent(t *testing.T) {
	val, ok := GetStateEnum("nonexistent_state")
	if ok {
		t.Errorf("GetStateEnum(%q) returned ok = true, want false", "nonexistent_state")
	}
	if val != 0 {
		t.Errorf("GetStateEnum(%q) returned val = %v, want 0", "nonexistent_state", val)
	}
}
