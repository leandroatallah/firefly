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

// Jumping
type JumpState struct {
	BaseState
}

func (s *JumpState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if j, ok := s.GetRootOwner().(Jumpable); ok {
		j.OnJump()
	}
}

// Falling
type FallState struct {
	BaseState
}

func (s *FallState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if f, ok := s.GetRootOwner().(Fallable); ok {
		f.OnFall()
	}
}

// Landing
type LandingState struct {
	BaseState
}

func (s *LandingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if l, ok := s.GetRootOwner().(Landable); ok {
		l.OnLand()
	}
}

// Hurt
type HurtState struct {
	BaseState
	count         int
	durationLimit int
}

func (s *HurtState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
	s.durationLimit = 30 // 0.5 sec, duration of the hurt animation
}
