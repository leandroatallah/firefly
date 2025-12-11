package body

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type FacingDirectionEnum int

const (
	FaceDirectionRight FacingDirectionEnum = iota
	FaceDirectionLeft
)

type Shape interface {
	Width() int
	Height() int
}

// Movable is a Shape but with movement
type Movable interface {
	Body

	MoveX(distance int)
	MoveY(distance int)
	OnMoveLeft(distance int)
	OnMoveUpLeft(distance int)
	OnMoveDownLeft(distance int)
	OnMoveRight(distance int)
	OnMoveUpRight(distance int)
	OnMoveDownRight(distance int)
	OnMoveUp(distance int)
	OnMoveDown(distance int)

	Velocity() (vx16, vy16 int)
	SetVelocity(vx16, vy16 int)
	Acceleration() (accX, accY int)
	SetAcceleration(accX, accY int)

	SetSpeed(speed int) error
	SetMaxSpeed(maxSpeed int) error
	Speed() int
	MaxSpeed() int
	Immobile() bool
	SetImmobile(immobile bool)
	FaceDirection() FacingDirectionEnum
	SetFaceDirection(value FacingDirectionEnum)
	IsIdle() bool
	IsWalking() bool
	IsFalling() bool
	IsGoingUp() bool
	CheckMovementDirectionX()

	// Platform methods
	TryJump(force int)
}

type Collidable interface {
	Body
	Touchable

	GetTouchable() Touchable
	DrawCollisionBox(screen *ebiten.Image, position image.Rectangle)
	CollisionPosition() []image.Rectangle
	CollisionShapes() []Collidable
	IsObstructive() bool
	SetIsObstructive(value bool)
	AddCollision(list ...Collidable)
	ClearCollisions()
	SetPosition(x int, y int)
	SetTouchable(t Touchable)
	ApplyValidPosition(distance16 int, isXAxis bool, space BodiesSpace) (x, y int, wasBlocked bool)
}

// TODO: Should it be merge with Collidable?
type Obstacle interface {
	Body
	Collidable
	Drawable
	DrawCollisionBox(screen *ebiten.Image, position image.Rectangle)
	ImageCollisionBox() *ebiten.Image
}

// Drawable represents any object that can be drawn to the screen.
type Drawable interface {
	Image() *ebiten.Image
	ImageOptions() *ebiten.DrawImageOptions
	UpdateImageOptions()
}

type Touchable interface {
	OnTouch(other Collidable)
	OnBlock(other Collidable)
}

type Alive interface {
	Body
	Health() int
	MaxHealth() int
	SetHealth(health int)
	SetMaxHealth(health int)
	LoseHealth(damage int)
	RestoreHealth(heal int)
	Invulnerable() bool
	SetInvulnerability(value bool)
}

type Body interface {
	ID() string
	SetID(id string)
	Position() image.Rectangle
	SetPosition(x, y int)
	GetPositionMin() (x, y int)
	GetShape() Shape
}

type BodiesSpace interface {
	AddBody(body Collidable)
	Bodies() []Collidable
	RemoveBody(body Collidable)
	ResolveCollisions(body Collidable) (touching bool, blocking bool)
}
