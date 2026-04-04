package actors

// Dieable is implemented by actors that have custom death behaviour.
type Dieable interface {
	OnDie()
}

// DyingState plays the death animation; transitions to DeadState when the animation finishes.
type DyingState struct {
	BaseState
}

func (s *DyingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if p, ok := s.GetRootOwner().(Dieable); ok {
		p.OnDie()
	}
}

// DeadState freezes and immobilises the actor after the death animation completes.
type DeadState struct {
	BaseState
}

func (s *DeadState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	actor := s.GetActor()
	actor.SetFreeze(true)
	actor.SetImmobile(true)
}

// ExitingState is a terminal state used when the actor leaves the scene (e.g. level exit).
type ExitingState struct {
	BaseState
}
