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

	count                int
	state                ActorState
	movementState        movement.MovementState
	movementModel        physics.MovementModel
	movementBlockers     int
	invulnerabilityTimer int
	imageOptions         *ebiten.DrawImageOptions
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

func (c *Character) GetCharacter() *Character {
	return c
}

// SetState set a new Character state and update current collision shapes.
func (c *Character) SetState(state ActorState) {
	if c.state == nil || c.state.State() != state.State() {
		c.state = state
		c.state.OnStart(c.count)
		c.StateCollisionManager.RefreshCollisions()
	}
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

	// Handle invulnerability timer
	if c.invulnerabilityTimer > 0 {
		c.invulnerabilityTimer--
		if c.invulnerabilityTimer == 0 {
			c.SetInvulnerability(false)
		}
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
	case state == Hurted:
		isAnimationOver := c.state.(*HurtState).IsAnimationFinished()
		if isAnimationOver {
			setNewState(Idle)
			c.SetImmobile(false)
		}
	case state != Falling && c.IsFalling():
		setNewState(Falling)
	case state != Walking && c.IsWalking():
		setNewState(Walking)
	case state != Idle && c.IsIdle():
		setNewState(Idle)
	}
}

func (c *Character) Hurt(damage int) {
	if c.Invulnerable() {
		return
	}

	c.LoseHealth(damage)

	// Switch to Hurt state
	state, err := NewActorState(c, Hurted)
	if err != nil {
		log.Fatal(err)
	}
	c.SetState(state)
	c.SetImmobile(true)
	c.SetInvulnerability(true)
	c.invulnerabilityTimer = 120 // 2 seconds at 60fps
}

func (c *Character) SetTouchable(t body.Touchable) {
	c.Touchable = t
}

func (c *Character) Image() *ebiten.Image {
	sprite := c.GetSpriteByState(c.state.State())
	if sprite == nil || sprite.Image == nil {
		// Try to fallback to idle sprite
		sprite = c.GetSpriteByState(Idle)
	}
	if sprite == nil || sprite.Image == nil {
		sprite = c.GetFirstSprite()
	}

	pos := c.Position()
	stateDurationCount := c.state.GetAnimationCount(c.count)
	return c.AnimatedSpriteImage(sprite, pos, stateDurationCount, c.SpriteEntity.FrameRate())
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
