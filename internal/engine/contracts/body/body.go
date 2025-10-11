package body

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Shape interface {
	Position() image.Rectangle
}

// Movable is a Shape but with movement
type Movable interface {
	Shape
	ApplyValidMovement(velocity int, isXAxis bool, space BodiesSpace)

	SetSpeedAndMaxSpeed(speed, maxSpeed int)
	Speed() int
	Immobile() bool
	SetImmobile(immobile bool)

	OnMoveUp(distance int)
	OnMoveDown(distance int)
	OnMoveLeft(distance int)
	OnMoveRight(distance int)
	OnMoveUpLeft(distance int)
	OnMoveUpRight(distance int)
	OnMoveDownLeft(distance int)
	OnMoveDownRight(distance int)

	TryJump(force int)
}

type Collidable interface {
	Shape
	Touchable
	GetTouchable() Touchable
	DrawCollisionBox(screen *ebiten.Image)
	CollisionPosition() []image.Rectangle
	IsObstructive() bool
	SetIsObstructive(value bool)
}

// TODO: Should it be merge with Collidable?
type Obstacle interface {
	Body
	Draw(screen *ebiten.Image)
	DrawCollisionBox(screen *ebiten.Image)
	Image() *ebiten.Image
	ImageCollisionBox() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions
}

type Touchable interface {
	OnTouch(other Body)
	OnBlock(other Body)
}

type Alive interface {
	Health() int
	MaxHealth() int
	SetHealth(health int)
	SetMaxHealth(health int)
	LoseHealth(damage int)
	RestoreHealth(heal int)
	Invulnerable() bool
	SetInvulnerable(value bool)
}

// Body is a Shape with collision, movable and alive
type Body interface {
	Shape
	Movable
	Collidable
	Alive

	ID() string
	SetPosition(x, y int)

	// TODO: For now, I commented movement model methods to prevent import cycle error, by maybe it could be removed.
	// SetMovementModel(model MovementModel)
	// MovementModel() MovementModel
}

type BodiesSpace interface {
	AddBody(body Body)
	Bodies() []Body
	RemoveBody(body Body)
	ResolveCollisions(body Body) (touching bool, blocking bool)
	// TODO: For now, I commented tilemap methods to prevent import cycle error, by maybe it could be removed.
	// GetTilemapDimensionsProvider() TilemapDimensionsProvider
	// SetTilemapDimensionsProvider(provider TilemapDimensionsProvider)
}
