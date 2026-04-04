package actors

import "fmt"

// StateConstructor is a factory function that builds an ActorState from a BaseState.
type StateConstructor func(base BaseState) ActorState

// Singleton registry: intentional package-level state
//
//nolint:gochecknoglobals
var (
	stateConstructors = make(map[ActorStateEnum]StateConstructor)
	stateEnums        = make(map[string]ActorStateEnum)
	nextEnumValue     ActorStateEnum
)

// RegisterState adds a named state to the global registry and returns its unique enum value.
// If the name is already registered, the constructor is updated and the existing enum is returned.
func RegisterState(name string, constructor StateConstructor) ActorStateEnum {
	if val, ok := stateEnums[name]; ok {
		stateConstructors[val] = constructor
		return val
	}

	enumValue := nextEnumValue
	nextEnumValue++

	stateEnums[name] = enumValue
	stateConstructors[enumValue] = constructor
	return enumValue
}

// NewState constructs an ActorState for the given actor and state enum using the registered constructor.
func NewState(actor ActorEntity, state ActorStateEnum) (ActorState, error) {
	constructor, ok := stateConstructors[state]
	if !ok {
		return nil, fmt.Errorf("unregistered state: %d", state)
	}
	base := NewBaseState(actor, state)
	return constructor(base), nil
}

// GetStateEnum looks up the enum value for a registered state by name.
func GetStateEnum(name string) (ActorStateEnum, bool) {
	val, ok := stateEnums[name]
	return val, ok
}
