package gamestates

type duckingSubState struct{}

func (s *duckingSubState) OnStart(_ int) {}
func (s *duckingSubState) OnFinish()     {}
func (s *duckingSubState) TransitionTo(input GroundedInput) GroundedSubStateEnum {
	if !input.DuckHeld() && input.HasCeilingClearance() {
		return SubStateIdle
	}
	return SubStateDucking
}
