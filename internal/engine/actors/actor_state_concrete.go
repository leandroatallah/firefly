package actors

// Idle
type IdleState struct {
	BaseState
}

func (s *IdleState) OnStart() {}

// Walking
type WalkState struct {
	BaseState
}

func (s *WalkState) OnStart() {}

// Falling
type FallState struct {
	BaseState
}

func (s *FallState) OnStart() {}

// Hurt
type HurtState struct {
	BaseState
	count         int
	durationLimit int
}

func (s *HurtState) OnStart() {
	s.durationLimit = 30 // 0.5 sec
}

func (s *HurtState) CheckRecovery() bool {
	s.count++

	if s.count > s.durationLimit {
		return true
	}

	return false
}
