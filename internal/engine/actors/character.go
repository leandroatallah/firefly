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
	physics.PhysicsBody
	sprites.SpriteEntity
	count            int
	state            ActorState
	movementState    movement.MovementState
	movementBlockers int
	animationCount   int
	// TODO: Move to the right place
	frameRate int
	// TODO: Rename this
	imageOptions *ebiten.DrawImageOptions
}

func NewCharacter(s sprites.SpriteMap, frameRate int) *Character {
	spriteEntity := sprites.NewSpriteEntity(s)
	c := &Character{
		SpriteEntity: spriteEntity,
		frameRate:    frameRate,
		imageOptions: &ebiten.DrawImageOptions{},
	}
	state, err := NewActorState(c, Idle)
	if err != nil {
		log.Fatal(err)
	}
	c.SetState(state)
	return c
}

// Builder methods
func (c *Character) SetBody(rect *physics.Rect) ActorEntity {
	c.PhysicsBody = *physics.NewPhysicsBody(rect)
	c.PhysicsBody.SetTouchable(c)
	return c
}

func (c *Character) SetCollisionArea(rect *physics.Rect) ActorEntity {
	collisionArea := &physics.CollisionArea{Shape: rect}
	c.PhysicsBody.AddCollision(collisionArea)
	return c
}

func (c *Character) SetID(id string) {
	c.PhysicsBody.SetID(id)
}

func (c *Character) State() ActorStateEnum {
	return c.state.State()
}

func (c *Character) SetState(state ActorState) {
	c.state = state
	c.state.OnStart()
}

func (c *Character) SetMovementState(
	state movement.MovementStateEnum,
	target body.Body,
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

func (c *Character) SetMovementModel(model physics.MovementModel) {
	c.PhysicsBody.SetMovementModel(model)
}

func (c *Character) Update(space body.BodiesSpace) error {
	c.count++

	// Handle movement by Movement State - must happen BEFORE UpdateMovement
	if c.movementState != nil {
		c.movementState.Move()
	}

	// Update physics and apply movement
	c.UpdateMovement(space)

	// Check movement direction for sprite mirroring
	c.CheckMovementDirectionX()
	c.UpdateImageOptions()

	c.handleState()

	return nil
}

func (c *Character) UpdateImageOptions() {
	c.imageOptions.GeoM.Reset()

	pos := c.Position()
	minX, minY := pos.Min.X, pos.Min.Y
	width := pos.Dx()

	fDirection := c.FaceDirection()
	if fDirection == physics.FaceDirectionLeft {
		c.imageOptions.GeoM.Scale(-1, 1)
		c.imageOptions.GeoM.Translate(float64(width), 0)
	}

	// Apply character position
	c.imageOptions.GeoM.Translate(
		float64(minX),
		float64(minY),
	)
}

func (c *Character) handleState() {
	setNewState := func(s ActorStateEnum) {
		state, err := NewActorState(c, s)
		if err != nil {
			log.Fatal(err)
		}
		c.SetState(state)
	}

	state := c.state.State()

	switch {
	case state == Idle && c.IsWalking():
		setNewState(Walk)
	case state == Walk && !c.IsWalking():
		setNewState(Idle)
	case state == Hurted:
		// TODO: The player should be recover the mobility before becomes vulnerable again
		// TODO: Should add a panic checking here?
		isRecovered := c.state.(*HurtState).CheckRecovery()
		if isRecovered {
			setNewState(Idle)
			// TODO: Group this in a helper function or method
			c.SetImmobile(false)
			c.SetInvulnerable(false)
		}
	}
}

func (c *Character) OnTouch(other body.Body) {}

func (c *Character) OnBlock(other body.Body) {}

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
	// TODO: Group this in a helper function or method
	c.SetImmobile(true)
	c.SetInvulnerable(true)

	c.LoseHealth(damage)
}

func (c *Character) SetTouchable(t body.Touchable) {
	c.PhysicsBody.Touchable = t
}

func (c *Character) Image() *ebiten.Image {
	pos := c.Position()
	width := pos.Dx()
	height := pos.Dy()

	img := c.GetSpriteByState(c.state.State())
	if img == nil {
		// Try to fallback to idle sprite
		img = c.GetSpriteByState(Idle)
	}

	characterWidth := img.Bounds().Dx()

	if width <= 0 {
		return img
	}
	frameCount := characterWidth / width
	if frameCount <= 1 {
		return img
	}

	i := (c.count / c.frameRate) % frameCount

	sx, sy := frameOX+i*width, frameOY

	return img.SubImage(
		image.Rect(sx, sy, sx+width, sy+height),
	).(*ebiten.Image)
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
