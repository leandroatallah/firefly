package actors

// IdleState is active when the character is standing still with no input.
type IdleState struct {
	BaseState
}

func (s *IdleState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// WalkState is active while the character is moving horizontally on the ground.
type WalkState struct {
	BaseState
}

func (s *WalkState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// JumpState is active while the character is ascending after a jump.
type JumpState struct {
	BaseState
}

func (s *JumpState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if j, ok := s.GetRootOwner().(Jumpable); ok {
		j.OnJump()
	}
}

// FallState is active while the character is descending under gravity.
type FallState struct {
	BaseState
}

func (s *FallState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if f, ok := s.GetRootOwner().(Fallable); ok {
		f.OnFall()
	}
}

// LandingState plays the landing animation immediately after the character touches the ground.
type LandingState struct {
	BaseState
}

func (s *LandingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if l, ok := s.GetRootOwner().(Landable); ok {
		l.OnLand()
	}
}

// HurtState plays the hurt animation for a fixed duration after the character takes damage.
type HurtState struct {
	BaseState
	durationLimit int // frames the hurt animation lasts (~0.5 s at 60 fps)
}

func (s *HurtState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
	s.durationLimit = 30 // 0.5 sec, duration of the hurt animation
}
