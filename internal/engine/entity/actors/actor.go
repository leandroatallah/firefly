package actors

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
)

// StatData holds the base numeric stats loaded from configuration for an actor.
type StatData struct {
	Health   int `json:"health"`
	Speed    int `json:"speed"`
	MaxSpeed int `json:"max_speed"`
}

// Controllable is implemented by actors that accept directional movement input.
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
	NewState(state ActorStateEnum) (ActorState, error)
}

// Damageable represents any actor that can take damage.
type Damageable interface {
	Hurt(damage int)
}

// Jumpable is implemented by actors that have custom on-jump behaviour.
type Jumpable interface {
	OnJump()
}

// Landable is implemented by actors that have custom on-land behaviour.
type Landable interface {
	OnLand()
}

// Fallable is implemented by actors that have custom on-fall behaviour.
type Fallable interface {
	OnFall()
}

// ActorEntity is the master interface for all game actors.
// It is composed of smaller interfaces that define specific behaviors.
type ActorEntity interface {
	body.Drawable
	Controllable
	Stateful
	Damageable
	body.Ownable

	body.MovableCollidableAlive

	Update(space body.BodiesSpace) error
	MovementModel() physicsmovement.MovementModel
	SetMovementModel(model physicsmovement.MovementModel)
	GetCharacter() *Character
}
