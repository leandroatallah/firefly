package kitstates

type aimLockSubState struct{}

func (s *aimLockSubState) OnStart(_ int) {}
func (s *aimLockSubState) OnFinish()     {}
func (s *aimLockSubState) TransitionTo(input GroundedInput) GroundedSubStateEnum {
	if !input.AimLockHeld() {
		return SubStateIdle
	}
	return SubStateAimLock
}
