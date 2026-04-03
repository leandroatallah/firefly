package gamestates

import "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"

var StateGrounded actors.ActorStateEnum

func init() {
	StateGrounded = actors.RegisterState("grounded", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b} // placeholder; GroundedState is constructed directly
	})
}

// GroundedState is a composite state that owns and delegates to a grounded sub-state.
type GroundedState struct {
	deps      GroundedDeps
	activeSub groundedSubState
	activeKey GroundedSubStateEnum
	count     int
}

func NewGroundedState(deps GroundedDeps) *GroundedState {
	return &GroundedState{deps: deps}
}

func (g *GroundedState) State() actors.ActorStateEnum { return StateGrounded }

func (g *GroundedState) OnStart(currentCount int) {
	g.count = currentCount
	g.setSub(SubStateIdle, currentCount)
}

func (g *GroundedState) OnFinish() {
	if g.activeSub != nil {
		g.activeSub.OnFinish()
	}
}

// Update evaluates input, transitions sub-states, and returns the next ActorStateEnum.
// Returns a non-StateGrounded value when the parent state machine must exit Grounded.
func (g *GroundedState) Update() actors.ActorStateEnum {
	input := g.deps.Input

	if input.JumpPressed() {
		return actors.Falling
	}
	if input.DashPressed() {
		return StateDashing
	}

	if g.deps.Shooting != nil && input.ShootHeld() {
		g.deps.Shooting.Update(g.deps.Body, g.deps.Model)
	}

	next := g.activeSub.transitionTo(input)
	if next != g.activeKey {
		g.activeSub.OnFinish()
		g.setSub(next, g.count)
	}

	return StateGrounded
}

// ActiveSubState returns the current sub-state enum (used by tests).
func (g *GroundedState) ActiveSubState() GroundedSubStateEnum { return g.activeKey }

// ForceSubState sets the active sub-state directly (used by tests to set up initial conditions).
func (g *GroundedState) ForceSubState(key GroundedSubStateEnum) {
	if g.activeSub != nil {
		g.activeSub.OnFinish()
	}
	g.setSub(key, g.count)
}

// GetAnimationCount satisfies actors.ActorState.
func (g *GroundedState) GetAnimationCount(currentCount int) int { return currentCount - g.count }

// IsAnimationFinished satisfies actors.ActorState.
func (g *GroundedState) IsAnimationFinished() bool { return false }

func (g *GroundedState) setSub(key GroundedSubStateEnum, count int) {
	g.activeKey = key
	g.activeSub = newSubState(key)
	g.activeSub.OnStart(count)
}

func newSubState(key GroundedSubStateEnum) groundedSubState {
	switch key {
	case SubStateWalking:
		return &walkingSubState{}
	case SubStateDucking:
		return &duckingSubState{}
	case SubStateAimLock:
		return &aimLockSubState{}
	default:
		return &idleSubState{}
	}
}
