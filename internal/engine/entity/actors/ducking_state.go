package actors

// State enum: part of engine public API
//
//nolint:gochecknoglobals
var Ducking ActorStateEnum

func init() {
	Ducking = RegisterState("duck", func(b BaseState) ActorState { return &DuckingState{BaseState: b} })
}

// DuckingState is active while the character is crouching.
// On entry it zeroes horizontal velocity; collision rects handle the reduced hitbox.
type DuckingState struct {
	BaseState
}

func (s *DuckingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	actor := s.GetActor()
	_, vy := actor.Velocity()
	actor.SetVelocity(0, vy)
}

func (s *DuckingState) OnFinish() {
	// No need to restore size - collision rects handle it
}
