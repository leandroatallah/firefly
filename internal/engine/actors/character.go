package actors

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type Character struct {
	sprites.SpriteEntity

	*physics.MovableBody
	*physics.CollidableBody
	*physics.AliveBody
	*physics.StateCollisionManager[ActorStateEnum]

	Touchable body.Touchable

	count            int
	state            ActorState
	movementState    movement.MovementState
	movementModel    physics.MovementModel
	movementBlockers int
	animationCount   int
	imageOptions     *ebiten.DrawImageOptions
}

func NewCharacter(s sprites.SpriteMap, bodyRect *physics.Rect) *Character { // Modified signature
	spriteEntity := sprites.NewSpriteEntity(s)
	b := physics.NewBody(bodyRect)
	movable := physics.NewMovableBody(b)
	collidable := physics.NewCollidableBody(b)
	alive := physics.NewAliveBody(b)
	c := &Character{
		MovableBody:    movable,
		CollidableBody: collidable,
		AliveBody:      alive,

		SpriteEntity: spriteEntity,
		imageOptions: &ebiten.DrawImageOptions{},
	}
	c.StateCollisionManager = physics.NewStateCollisionManager[ActorStateEnum](c)

	state, err := NewActorState(c, Idle)
	if err != nil {
		log.Fatal(err)
	}
	c.SetState(state)
	return c
}

// Forwarding methods for Body to avoid ambiguous selector
// Always route via the MovableBody component
func (c *Character) ID() string {
	return c.MovableBody.ID()
}
func (c *Character) SetID(id string) {
	c.MovableBody.SetID(id)
}
func (c *Character) Position() image.Rectangle {
	return c.MovableBody.Position()
}
func (c *Character) SetPosition(x, y int) {
	c.CollidableBody.SetPosition(x, y)
}
func (c *Character) GetPositionMin() (int, int) {
	return c.MovableBody.GetPositionMin()
}
func (c *Character) GetShape() body.Shape {
	return c.MovableBody.GetShape()
}

// Builder methods
func (c *Character) State() ActorStateEnum {
	return c.state.State()
}

func (c *Character) AddCollisionRect(state ActorStateEnum, rect body.Collidable) {
	c.StateCollisionManager.AddCollisionRect(state, rect)
}

// SetState set a new Character state and update current collision shapes.
func (c *Character) SetState(state ActorState) {
	c.state = state
	c.StateCollisionManager.RefreshCollisions()
	c.state.OnStart()
}

func (c *Character) SetMovementState(
	state movement.MovementStateEnum,
	target body.MovableCollidable,
	options ...movement.MovementStateOption,
) {
	movementState, err := movement.NewMovementState(c, state, target, options...)
	if err != nil {
		log.Fatal(err)
	}

	c.movementState = movementState
	c.movementState.OnStart()
}
func (c *Character) SwitchMovementState(state movement.MovementStateEnum) {
	target := c.MovementState().Target()
	movementState, err := movement.NewMovementState(c, state, target)
	if err != nil {
		log.Fatal(err)
	}
	c.movementState = movementState
}

func (c *Character) MovementState() movement.MovementState {
	return c.movementState
}

func (c *Character) Update(space body.BodiesSpace) error {
	c.count++

	// Handle movement by Movement State - must happen BEFORE UpdateMovement
	if c.movementState != nil {
		c.movementState.Move()
	}

	// Update physics and apply movement
	c.UpdateMovement(space)

	c.handleState()

	return nil
}

func (c *Character) UpdateMovement(space body.BodiesSpace) {
	if c.movementModel != nil {
		c.movementModel.Update(c, space)
	}
}

func (c *Character) UpdateImageOptions() {
	if c.imageOptions == nil {
		return
	}
	c.imageOptions.GeoM.Reset()

	accX, _ := c.Acceleration()
	fDirection := c.FaceDirection()

	if accX > 0 {
		fDirection = body.FaceDirectionRight
	} else if accX < 0 {
		fDirection = body.FaceDirectionLeft
	}

	c.SetFaceDirection(fDirection)

	if fDirection == body.FaceDirectionLeft {
		width := c.Position().Dx()
		c.imageOptions.GeoM.Scale(-1, 1)
		c.imageOptions.GeoM.Translate(float64(width), 0)
	}

	// Apply character position
	x, y := c.GetPositionMin()
	c.imageOptions.GeoM.Translate(
		float64(x),
		float64(y),
	)
}

func (c *Character) handleState() {
	if c.state == nil {
		return
	}

	setNewState := func(s ActorStateEnum) {
		state, err := NewActorState(c, s)
		if err != nil {
			log.Fatal(err)
		}
		c.SetState(state)
	}

	state := c.state.State()

	switch {
	case state != Falling && c.IsFalling():
		setNewState(Falling)
	case state != Walking && c.IsWalking():
		setNewState(Walking)
	case state != Idle && c.IsIdle():
		setNewState(Idle)
	case state == Hurted:
		// TODO: The player should be recover the mobility before becomes vulnerable again
		isRecovered := c.state.(*HurtState).CheckRecovery()
		if isRecovered {
			setNewState(Idle)
			c.SetImmobile(false)
			c.SetInvulnerability(false)
		}
	}
}

func (c *Character) Hurt(damage int) {
	if c.Invulnerable() {
		return
	}

	// TODO: Check condition to react to damage 0
	// ...

	// Switch to Hurt state
	state, err := NewActorState(c, Hurted)
	if err != nil {
		log.Fatal(err)
	}
	c.SetState(state)
	c.SetImmobile(true)
	c.SetInvulnerability(true)

	c.LoseHealth(damage)
}

func (c *Character) SetTouchable(t body.Touchable) {
	c.Touchable = t
}

func (c *Character) Image() *ebiten.Image {
	img := c.GetSpriteByState(c.state.State())
	if img == nil {
		// Try to fallback to idle sprite
		img = c.GetSpriteByState(Idle)
	}
	if img == nil {
		img = c.GetFirstSprite()
	}

	pos := c.Position()
	return c.AnimatedSpriteImage(img, pos, c.count, c.SpriteEntity.FrameRate())
}

// WithCollisionBox extend Image method to show a rect with the collision area
func (c *Character) ImageCollisionBox() *ebiten.Image {
	img := c.Image()
	pos := c.Position()

	// Create a new image and copy the subimage to it
	res := ebiten.NewImage(img.Bounds().Dx(), img.Bounds().Dy())
	res.DrawImage(img, nil)

	c.DrawCollisionBox(res, pos)
	return res
}

func (c *Character) ImageOptions() *ebiten.DrawImageOptions {
	return c.imageOptions
}

// BlockMovement increases the count of systems blocking movement.
func (p *Character) BlockMovement() {
	p.movementBlockers++
}

// UnblockMovement decreases the count.
func (p *Character) UnblockMovement() {
	p.movementBlockers--
	if p.movementBlockers < 0 {
		p.movementBlockers = 0
	}
}

// IsPlayerMovementBlocked checks if any system is currently blocking movement.
func (p *Character) IsMovementBlocked() bool {
	return p.movementBlockers > 0
}

// Platform methods
func (c *Character) TryJump(force int) {
	c.MovableBody.TryJump(force)
}

// Movement Model methods
func (c *Character) SetMovementModel(model physics.MovementModel) {
	c.movementModel = model
}

func (c *Character) MovementModel() physics.MovementModel {
	return c.movementModel
}
