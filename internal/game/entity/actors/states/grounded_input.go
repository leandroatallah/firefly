package gamestates

import (
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
)

// GroundedInput is the input contract consumed by GroundedState and its sub-states.
type GroundedInput interface {
	HorizontalInput() int
	DuckHeld() bool
	HasCeilingClearance() bool
	JumpPressed() bool
	DashPressed() bool
	MeleePressed() bool
	AimLockHeld() bool
	ShootHeld() bool
}

// GroundedDeps holds the dependencies injected into GroundedState.
type GroundedDeps struct {
	Input GroundedInput
	Body  contractsbody.MovableCollidable
	Model *physicsmovement.PlatformMovementModel
}
