package body

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/tilemaplayer"
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
	SetFreeze(freeze bool)
	Freeze() bool

	FaceDirection() animation.FacingDirectionEnum
	SetFaceDirection(value animation.FacingDirectionEnum)
	IsIdle() bool
	IsWalking() bool
	IsFalling() bool
	IsGoingUp() bool
	CheckMovementDirectionX()

	// Platform methods
	TryJump(force int)
	SetJumpForceMultiplier(multiplier float64)
	JumpForceMultiplier() float64

	SetHorizontalInertia(inertia float64)
	HorizontalInertia() float64
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
	SetPosition16(x16, y16 int)
	GetPosition16() (x16, y16 int)
	SetTouchable(t Touchable)
	ApplyValidPosition(distance16 int, isXAxis bool, space BodiesSpace) (x, y int, wasBlocked bool)
}

type Obstacle interface {
	Body
	Collidable
	Drawable
	DrawCollisionBox(screen *ebiten.Image, position image.Rectangle)
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
	Ownable
	ID() string
	SetID(id string)
	Position() image.Rectangle
	SetPosition(x, y int)
	SetPosition16(x16, y16 int)
	SetSize(width, height int)
	Scale() float64
	SetScale(float64)
	GetPosition16() (x16, y16 int)
	GetPositionMin() (x, y int)
	GetShape() Shape
}

type BodiesSpace interface {
	AddBody(body Collidable)
	Bodies() []Collidable
	RemoveBody(body Collidable)
	QueueForRemoval(body Collidable)
	ProcessRemovals()
	Clear()
	ResolveCollisions(body Collidable) (touching bool, blocking bool)
	SetTilemapDimensionsProvider(provider tilemaplayer.TilemapDimensionsProvider)
	GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider
	Find(id string) Collidable
	Query(rect image.Rectangle) []Collidable
}

type Ownable interface {
	Owner() interface{}
	SetOwner(interface{})
	LastOwner() interface{}
}
