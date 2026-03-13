package items

import "fmt"

type StateConstructor func(base BaseState) ItemState

var (
	stateConstructors = make(map[ItemStateEnum]StateConstructor)
	stateEnums        = make(map[string]ItemStateEnum)
	nextEnumValue     ItemStateEnum
)

func RegisterState(name string, constructor StateConstructor) ItemStateEnum {
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

func NewState(item Item, state ItemStateEnum) (ItemState, error) {
	constructor, ok := stateConstructors[state]
	if !ok {
		return nil, fmt.Errorf("unregistered item state: %d", state)
	}
	base := BaseState{item: item, state: state}
	return constructor(base), nil
}

func GetStateEnum(name string) (ItemStateEnum, bool) {
	val, ok := stateEnums[name]
	return val, ok
}
