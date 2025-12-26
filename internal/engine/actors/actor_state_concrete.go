package actors

// Idle
type IdleState struct {
	BaseState
}

func (s *IdleState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Walking
type WalkState struct {
	BaseState
}

func (s *WalkState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Falling
type FallState struct {
	BaseState
}

func (s *FallState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// Hurt
type HurtState struct {
	BaseState
	count                int
	durationLimit        int
	invulnerabilityCount int
	invulnerabilityLimit int
}

func (s *HurtState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
	s.durationLimit = 30        // 0.5 sec
	s.invulnerabilityLimit = 120 // 2 sec
}

func (s *HurtState) CheckRecovery() bool {
	s.count++

	if s.count > s.durationLimit {
		return true
	}

	return false
}

func (s *HurtState) CheckInvulnerability() bool {
	s.invulnerabilityCount++

	if s.invulnerabilityCount > s.invulnerabilityLimit {
		return true
	}

	return false
}
