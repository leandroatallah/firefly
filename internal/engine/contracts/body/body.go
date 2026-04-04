package body

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/hajimehoshi/ebiten/v2"
)

// Shape represents an object with a fixed width and height.
type Shape interface {
	// Width returns the width of the shape in pixels.
	Width() int
	// Height returns the height of the shape in pixels.
	Height() int
}

// Movable is a Shape but with movement
type Movable interface {
	Body

	// MoveX moves the body horizontally by the given distance.
	MoveX(distance int)
	// MoveY moves the body vertically by the given distance.
	MoveY(distance int)
	// OnMoveLeft is called when the body moves left.
	OnMoveLeft(distance int)
	// OnMoveUpLeft is called when the body moves up-left.
	OnMoveUpLeft(distance int)
	// OnMoveDownLeft is called when the body moves down-left.
	OnMoveDownLeft(distance int)
	// OnMoveRight is called when the body moves right.
	OnMoveRight(distance int)
	// OnMoveUpRight is called when the body moves up-right.
	OnMoveUpRight(distance int)
	// OnMoveDownRight is called when the body moves down-right.
	OnMoveDownRight(distance int)
	// OnMoveUp is called when the body moves up.
	OnMoveUp(distance int)
	// OnMoveDown is called when the body moves down.
	OnMoveDown(distance int)

	// Velocity returns the current velocity as fixed-point 16 values.
	Velocity() (vx16, vy16 int)
	// SetVelocity sets the current velocity as fixed-point 16 values.
	SetVelocity(vx16, vy16 int)
	// Acceleration returns the current acceleration components.
	Acceleration() (accX, accY int)
	// SetAcceleration sets the acceleration components.
	SetAcceleration(accX, accY int)

	// SetSpeed sets the movement speed, returning an error if invalid.
	SetSpeed(speed int) error
	// SetMaxSpeed sets the maximum allowed speed, returning an error if invalid.
	SetMaxSpeed(maxSpeed int) error
	// Speed returns the current movement speed.
	Speed() int
	// MaxSpeed returns the maximum allowed movement speed.
	MaxSpeed() int
	// Immobile reports whether the body is prevented from moving.
	Immobile() bool
	// SetImmobile sets whether the body is prevented from moving.
	SetImmobile(immobile bool)
	// SetFreeze sets whether the body's movement is frozen.
	SetFreeze(freeze bool)
	// Freeze reports whether the body's movement is currently frozen.
	Freeze() bool

	// FaceDirection returns the direction the body is currently facing.
	FaceDirection() animation.FacingDirectionEnum
	// SetFaceDirection sets the direction the body faces.
	SetFaceDirection(value animation.FacingDirectionEnum)
	// IsIdle reports whether the body is in an idle state.
	IsIdle() bool
	// IsWalking reports whether the body is currently walking.
	IsWalking() bool
	// IsFalling reports whether the body is currently falling.
	IsFalling() bool
	// IsGoingUp reports whether the body is currently moving upward.
	IsGoingUp() bool
	// CheckMovementDirectionX updates the facing direction based on horizontal velocity.
	CheckMovementDirectionX()

	// Platform methods

	// TryJump attempts to make the body jump with the given force.
	TryJump(force int)
	// SetJumpForceMultiplier sets a multiplier applied to the jump force.
	SetJumpForceMultiplier(multiplier float64)
	// JumpForceMultiplier returns the current jump force multiplier.
	JumpForceMultiplier() float64

	// SetHorizontalInertia sets the horizontal inertia factor.
	SetHorizontalInertia(inertia float64)
	// HorizontalInertia returns the current horizontal inertia factor.
	HorizontalInertia() float64
}

// Collidable represents a body that participates in collision detection.
type Collidable interface {
	Body
	Touchable

	// GetTouchable returns the Touchable associated with this body.
	GetTouchable() Touchable
	// DrawCollisionBox renders the collision box at the given position for debugging.
	DrawCollisionBox(screen *ebiten.Image, position image.Rectangle)
	// CollisionPosition returns the list of collision rectangles for this body.
	CollisionPosition() []image.Rectangle
	// CollisionShapes returns all collidable shapes attached to this body.
	CollisionShapes() []Collidable
	// IsObstructive reports whether this body blocks movement.
	IsObstructive() bool
	// SetIsObstructive sets whether this body blocks movement.
	SetIsObstructive(value bool)
	// AddCollision registers additional collidable shapes.
	AddCollision(list ...Collidable)
	// ClearCollisions removes all registered collision shapes.
	ClearCollisions()
	// SetPosition sets the body's position in pixel coordinates.
	SetPosition(x int, y int)
	// SetPosition16 sets the body's position as fixed-point 16 values.
	SetPosition16(x16, y16 int)
	// GetPosition16 returns the body's position as fixed-point 16 values.
	GetPosition16() (x16, y16 int)
	// SetTouchable assigns a Touchable handler to this body.
	SetTouchable(t Touchable)
	// ApplyValidPosition moves the body by distance16 along the given axis, resolving collisions.
	ApplyValidPosition(distance16 int, isXAxis bool, space BodiesSpace) (x, y int, wasBlocked bool)
}

// Obstacle is a fully collidable, drawable body that blocks movement.
type Obstacle interface {
	Body
	Collidable
	Drawable
	// DrawCollisionBox renders the collision box at the given position for debugging.
	DrawCollisionBox(screen *ebiten.Image, position image.Rectangle)
}

// Drawable represents any object that can be drawn to the screen.
type Drawable interface {
	// Image returns the sprite image for this object.
	Image() *ebiten.Image
	// ImageOptions returns the draw options used when rendering this object.
	ImageOptions() *ebiten.DrawImageOptions
	// UpdateImageOptions refreshes the draw options (e.g. position, scale, colour).
	UpdateImageOptions()
}

// Touchable handles touch and block events triggered by collision resolution.
type Touchable interface {
	// OnTouch is called when this body overlaps another collidable.
	OnTouch(other Collidable)
	// OnBlock is called when this body is blocked by another collidable.
	OnBlock(other Collidable)
}

// Alive represents a body with health and invulnerability state.
type Alive interface {
	Body
	// Health returns the current health value.
	Health() int
	// MaxHealth returns the maximum health value.
	MaxHealth() int
	// SetHealth sets the current health to the given value.
	SetHealth(health int)
	// SetMaxHealth sets the maximum health to the given value.
	SetMaxHealth(health int)
	// LoseHealth reduces health by the given damage amount.
	LoseHealth(damage int)
	// RestoreHealth increases health by the given heal amount.
	RestoreHealth(heal int)
	// Invulnerable reports whether the body is currently invulnerable.
	Invulnerable() bool
	// SetInvulnerability sets the invulnerability state.
	SetInvulnerability(value bool)
}

// Body is the base interface for all engine entities, providing identity, position, and size.
type Body interface {
	Ownable
	// ID returns the unique identifier of this body.
	ID() string
	// SetID sets the unique identifier of this body.
	SetID(id string)
	// Position returns the body's bounding rectangle in pixel coordinates.
	Position() image.Rectangle
	// SetPosition sets the body's position in pixel coordinates.
	SetPosition(x, y int)
	// SetPosition16 sets the body's position as fixed-point 16 values.
	SetPosition16(x16, y16 int)
	// SetSize sets the body's width and height in pixels.
	SetSize(width, height int)
	// Scale returns the current render scale factor.
	Scale() float64
	// SetScale sets the render scale factor.
	SetScale(float64)
	// GetPosition16 returns the body's position as fixed-point 16 values.
	GetPosition16() (x16, y16 int)
	// GetPositionMin returns the top-left corner of the body in pixel coordinates.
	GetPositionMin() (x, y int)
	// GetShape returns the Shape describing this body's dimensions.
	GetShape() Shape
}

// BodiesSpace manages a collection of collidable bodies and resolves collisions between them.
type BodiesSpace interface {
	// AddBody registers a collidable body in the space.
	AddBody(body Collidable)
	// Bodies returns all collidable bodies currently in the space.
	Bodies() []Collidable
	// RemoveBody immediately removes a body from the space.
	RemoveBody(body Collidable)
	// QueueForRemoval schedules a body to be removed at the end of the current frame.
	QueueForRemoval(body Collidable)
	// ProcessRemovals removes all bodies that were queued for removal.
	ProcessRemovals()
	// Clear removes all bodies from the space.
	Clear()
	// ResolveCollisions checks and resolves collisions for the given body.
	ResolveCollisions(body Collidable) (touching bool, blocking bool)
	// SetTilemapDimensionsProvider sets the provider used for tilemap boundary queries.
	SetTilemapDimensionsProvider(provider tilemaplayer.TilemapDimensionsProvider)
	// GetTilemapDimensionsProvider returns the current tilemap dimensions provider.
	GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider
	// Find returns the collidable body with the given ID, or nil if not found.
	Find(id string) Collidable
	// Query returns all collidable bodies whose collision shapes overlap the given rectangle.
	Query(rect image.Rectangle) []Collidable
}

// Ownable tracks ownership of a body, supporting handoff between owners.
type Ownable interface {
	// Owner returns the current owner of this body.
	Owner() interface{}
	// SetOwner sets the current owner of this body.
	SetOwner(interface{})
	// LastOwner returns the previous owner of this body.
	LastOwner() interface{}
}
