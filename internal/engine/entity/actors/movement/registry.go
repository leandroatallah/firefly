package movement

import "fmt"

type MovementStateConstructor func(base BaseMovementState) MovementState

var (
	movementStateConstructors = make(map[MovementStateEnum]MovementStateConstructor)
	movementStateEnums        = make(map[string]MovementStateEnum)
	nextMovementEnumValue     MovementStateEnum
)

func RegisterMovementState(name string, constructor MovementStateConstructor) MovementStateEnum {
	if val, ok := movementStateEnums[name]; ok {
		return val
	}

	// If nextMovementEnumValue collides with existing hardcoded enums, skip them.
	// Hardcoded enums in state_base.go go up to 8 (SideToSide).
	// So we start from 100 to be safe, or we check if it is occupied.
	// A simpler approach is to rely on initialization order, but that is risky.
	// Since we are adding dynamic registration, we should ensure we don't conflict.
	// However, for now, let's assume dynamic states start after the hardcoded ones.

	if nextMovementEnumValue == 0 {
		// Initialize to a value safe from hardcoded conflicts.
		// Hardcoded values end at SideToSide (8).
		nextMovementEnumValue = 100
	}

	enumValue := nextMovementEnumValue
	nextMovementEnumValue++

	movementStateEnums[name] = enumValue
	movementStateConstructors[enumValue] = constructor
	return enumValue
}

func GetMovementStateConstructor(state MovementStateEnum) (MovementStateConstructor, error) {
	constructor, ok := movementStateConstructors[state]
	if !ok {
		return nil, fmt.Errorf("unregistered movement state: %d", state)
	}
	return constructor, nil
}
