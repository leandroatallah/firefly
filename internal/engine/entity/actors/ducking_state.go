package actors

const (
	duckHeightRatio = 0.5
)

var Ducking ActorStateEnum

func init() {
	Ducking = RegisterState("duck", func(b BaseState) ActorState { return &DuckingState{BaseState: b} })
}

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
