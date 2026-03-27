package actors

// Dieable interface for actors that have a death behavior.
type Dieable interface {
	OnDie()
}

// Dying
type DyingState struct {
	BaseState
}

func (s *DyingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if p, ok := s.GetRootOwner().(Dieable); ok {
		p.OnDie()
	}
}

// Dead
type DeadState struct {
	BaseState
}

func (s *DeadState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	actor := s.GetActor()
	actor.SetFreeze(true)
	actor.SetImmobile(true)
}

// Exiting
type ExitingState struct {
	BaseState
}
