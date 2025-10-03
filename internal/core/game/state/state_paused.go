package state

type PausedState struct {
	BaseState
}

func (s *PausedState) OnStart() {}
