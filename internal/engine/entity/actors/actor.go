package actors

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
)

type StatData struct {
	Health   int `json:"health"`
	Speed    int `json:"speed"`
	MaxSpeed int `json:"max_speed"`
}

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

type Jumpable interface {
	OnJump()
}

type Landable interface {
	OnLand()
}

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
