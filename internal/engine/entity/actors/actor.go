package actors

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
)

// Controllable defines methods for an actor that can be moved or have its movement blocked.
type Controllable interface {
	OnMoveLeft(force int)
	OnMoveRight(force int)
	BlockMovement()
	UnblockMovement()
	IsMovementBlocked() bool
}

// Stateful defines methods for an actor that has general and movement-specific states.
type Stateful interface {
	State() ActorStateEnum
	SetState(state ActorState)
	SetMovementState(
		state movement.MovementStateEnum,
		target body.MovableCollidable,
		options ...movement.MovementStateOption,
	)
	SwitchMovementState(state movement.MovementStateEnum)
	MovementState() movement.MovementState
}

// Damageable represents any actor that can take damage.
type Damageable interface {
	Hurt(damage int)
}

// ActorEntity is the master interface for all game actors.
// It is composed of smaller interfaces that define specific behaviors.
type ActorEntity interface {
	body.Drawable
	Controllable
	Stateful
	Damageable

	body.MovableCollidableAlive

	Update(space body.BodiesSpace) error
	MovementModel() physicsmovement.MovementModel
	SetMovementModel(model physicsmovement.MovementModel)
	GetCharacter() *Character
}
