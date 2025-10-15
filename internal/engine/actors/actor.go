package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

type ActorEntity interface {
	body.Body
	ID() string
	SetID(id string)
	SetBody(rect *physics.Rect) ActorEntity
	SetCollisionArea(rect *physics.Rect) ActorEntity
	State() ActorStateEnum
	SetState(state ActorState)
	SetMovementState(
		state movement.MovementStateEnum,
		target body.Body,
		options ...movement.MovementStateOption,
	)
	SwitchMovementState(state movement.MovementStateEnum)
	MovementState() movement.MovementState
	MovementModel() physics.MovementModel
	SetMovementModel(model physics.MovementModel)
	Update(space body.BodiesSpace) error
	Hurt(damage int)
	OnMoveLeft(force int)
	OnMoveRight(force int)
	Image() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions
}

type ActorType int

type ActorMap map[ActorType]ActorEntity
