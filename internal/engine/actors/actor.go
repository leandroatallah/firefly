package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

// Identifiable represents any object with a string ID.
type Identifiable interface {
	ID() string
	SetID(id string)
}

// Drawable represents any object that can be drawn to the screen.
type Drawable interface {
	Image() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions
}

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
		target body.Body,
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
	body.Body
	Identifiable
	Drawable
	Controllable
	Stateful
	Damageable

	Update(space body.BodiesSpace) error
	MovementModel() physics.MovementModel
	SetMovementModel(model physics.MovementModel)
	SetBody(rect *physics.Rect) ActorEntity
	SetCollisionArea(rect *physics.Rect) ActorEntity
	SetTouchable(t body.Touchable)
}
