package actors

import "fmt"

type StateConstructor func(base BaseState) ActorState

var (
	stateConstructors = make(map[ActorStateEnum]StateConstructor)
	stateEnums        = make(map[string]ActorStateEnum)
	nextEnumValue     ActorStateEnum
)

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

func NewState(actor ActorEntity, state ActorStateEnum) (ActorState, error) {
	constructor, ok := stateConstructors[state]
	if !ok {
		return nil, fmt.Errorf("unregistered state: %d", state)
	}
	base := NewBaseState(actor, state)
	return constructor(base), nil
}

func GetStateEnum(name string) (ActorStateEnum, bool) {
	val, ok := stateEnums[name]
	return val, ok
}
