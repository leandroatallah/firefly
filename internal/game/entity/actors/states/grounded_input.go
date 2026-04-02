package gamestates

import contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

// GroundedInput is the input contract consumed by GroundedState and its sub-states.
type GroundedInput interface {
	HorizontalInput() int
	DuckHeld() bool
	HasCeilingClearance() bool
	JumpPressed() bool
	DashPressed() bool
	AimLockHeld() bool
	ShootHeld() bool
}

// GroundedDeps holds the dependencies injected into GroundedState.
type GroundedDeps struct {
	Input    GroundedInput
	Shooting *ShootingSkill
	Body     contractsbody.Movable
}
