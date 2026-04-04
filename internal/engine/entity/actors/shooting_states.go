package actors

// IdleShootingState plays the idle-shoot animation while the character fires standing still.
type IdleShootingState struct {
	BaseState
}

// WalkingShootingState plays the walk-shoot animation while the character fires while moving.
type WalkingShootingState struct {
	BaseState
}

// JumpingShootingState plays the jump-shoot animation while the character fires in the air going up.
type JumpingShootingState struct {
	BaseState
}

// FallingShootingState plays the fall-shoot animation while the character fires while descending.
type FallingShootingState struct {
	BaseState
}
