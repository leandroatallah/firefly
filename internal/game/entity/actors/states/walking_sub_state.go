package gamestates

type walkingSubState struct{}

func (s *walkingSubState) OnStart(_ int) {}
func (s *walkingSubState) OnFinish()     {}
func (s *walkingSubState) TransitionTo(input GroundedInput) GroundedSubStateEnum {
	switch {
	case input.AimLockHeld():
		return SubStateAimLock
	case input.DuckHeld():
		return SubStateDucking
	case input.HorizontalInput() != 0:
		return SubStateWalking
	default:
		return SubStateIdle
	}
}
