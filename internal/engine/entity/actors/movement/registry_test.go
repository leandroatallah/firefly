package movement

import (
	"testing"
)

func TestRegisterMovementState_NewState(t *testing.T) {
	// Clean up registry after test
	defer func() {
		movementStateConstructors = make(map[MovementStateEnum]MovementStateConstructor)
		movementStateEnums = make(map[string]MovementStateEnum)
		nextMovementEnumValue = 0
	}()

	constructor := func(b BaseMovementState) MovementState { return &IdleMovementState{BaseMovementState: b} }
	enum := RegisterMovementState("test_state", constructor)

	if _, ok := movementStateEnums["test_state"]; !ok {
		t.Error("RegisterMovementState() did not register state name")
	}

	if _, ok := movementStateConstructors[enum]; !ok {
		t.Error("RegisterMovementState() did not register constructor")
	}

	// Verify enum value starts at 100 or higher (safe from hardcoded)
	// Note: First registration after cleanup returns 100
	// (nextMovementEnumValue starts at 0, gets set to 100, and returns that value)
	if enum < 100 {
		t.Errorf("RegisterMovementState() first enum = %d, want >= 100", enum)
	}
}

func TestRegisterMovementState_SecondState(t *testing.T) {
	// Clean up registry after test
	defer func() {
		movementStateConstructors = make(map[MovementStateEnum]MovementStateConstructor)
		movementStateEnums = make(map[string]MovementStateEnum)
		nextMovementEnumValue = 0
	}()

	constructor := func(b BaseMovementState) MovementState { return &IdleMovementState{BaseMovementState: b} }
	enum1 := RegisterMovementState("test_state1", constructor)
	enum2 := RegisterMovementState("test_state2", constructor)

	if enum2 != enum1+1 {
		t.Errorf("RegisterMovementState() second enum = %d, want %d", enum2, enum1+1)
	}
}

func TestRegisterMovementState_DuplicateName(t *testing.T) {
	// Clean up registry after test
	defer func() {
		movementStateConstructors = make(map[MovementStateEnum]MovementStateConstructor)
		movementStateEnums = make(map[string]MovementStateEnum)
		nextMovementEnumValue = 0
	}()

	constructor1 := func(b BaseMovementState) MovementState { return &IdleMovementState{BaseMovementState: b} }
	constructor2 := func(b BaseMovementState) MovementState { return &ChaseMovementState{BaseMovementState: b} }

	enum1 := RegisterMovementState("duplicate_test", constructor1)
	enum2 := RegisterMovementState("duplicate_test", constructor2)

	if enum1 != enum2 {
		t.Error("RegisterMovementState() with duplicate name should return same enum")
	}
}

func TestGetMovementStateConstructor_Existing(t *testing.T) {
	// Clean up registry after test
	defer func() {
		movementStateConstructors = make(map[MovementStateEnum]MovementStateConstructor)
		movementStateEnums = make(map[string]MovementStateEnum)
		nextMovementEnumValue = 0
	}()

	constructor := func(b BaseMovementState) MovementState { return &IdleMovementState{BaseMovementState: b} }
	enum := RegisterMovementState("test_ctor", constructor)

	retrieved, err := GetMovementStateConstructor(enum)
	if err != nil {
		t.Errorf("GetMovementStateConstructor() returned error: %v", err)
	}
	if retrieved == nil {
		t.Error("GetMovementStateConstructor() returned nil constructor")
	}
}

func TestGetMovementStateConstructor_Unregistered(t *testing.T) {
	// Use an unregistered state value
	unregisteredState := MovementStateEnum(9999)

	constructor, err := GetMovementStateConstructor(unregisteredState)
	if err == nil {
		t.Error("GetMovementStateConstructor() with unregistered state should return error")
	}
	if constructor != nil {
		t.Error("GetMovementStateConstructor() with unregistered state should return nil constructor")
	}
}
