package gamestates

type idleSubState struct{}

func (s *idleSubState) OnStart(_ int)  {}
func (s *idleSubState) OnFinish()      {}
func (s *idleSubState) transitionTo(input GroundedInput) GroundedSubStateEnum {
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
