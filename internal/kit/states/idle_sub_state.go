package kitstates

// GroundedInputLike is the minimum input surface IdleSubState reads.
// Game-side GroundedInput satisfies this interface structurally.
type GroundedInputLike interface {
	AimLockHeld() bool
	DuckHeld() bool
	HorizontalInput() int
}

// IdleSubState is a no-op grounded sub-state that selects the next
// sub-state based on input. The concrete enum E is supplied by the
// caller so kit stays free of game-specific identifiers.
type IdleSubState[E comparable, I GroundedInputLike] struct {
	Idle    E
	Walking E
	Ducking E
	AimLock E
}

func (s *IdleSubState[E, I]) OnStart(_ int) {}
func (s *IdleSubState[E, I]) OnFinish()     {}

func (s *IdleSubState[E, I]) TransitionTo(input I) E {
	switch {
	case input.AimLockHeld():
		return s.AimLock
	case input.DuckHeld():
		return s.Ducking
	case input.HorizontalInput() != 0:
		return s.Walking
	default:
		return s.Idle
	}
}
