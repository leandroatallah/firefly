package gamestates

// GroundedSubStateEnum identifies the active sub-state inside GroundedState.
type GroundedSubStateEnum int

const (
	SubStateIdle    GroundedSubStateEnum = iota
	SubStateWalking
	SubStateDucking
	SubStateAimLock
)

// groundedSubState is the internal contract every sub-state must satisfy.
type groundedSubState interface {
	OnStart(currentCount int)
	OnFinish()
	transitionTo(input GroundedInput) GroundedSubStateEnum
}
