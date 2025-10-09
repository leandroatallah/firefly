package actors

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors/movement"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type Character struct {
	physics.PhysicsBody
	SpriteEntity
	count          int
	state          ActorState
	movementState  movement.MovementState
	animationCount int
	// TODO: Rename this
	op *ebiten.DrawImageOptions
}

func NewCharacter(sprites SpriteMap) *Character {
	spriteEntity := NewSpriteEntity(sprites)
	c := &Character{
		SpriteEntity: spriteEntity,
		op:           &ebiten.DrawImageOptions{},
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

func (c *Character) State() ActorStateEnum {
	return c.state.State()
}

func (c *Character) SetState(state ActorState) {
	c.state = state
	c.state.OnStart()
}

func (c *Character) SetMovementState(
	state movement.MovementStateEnum,
	target physics.Body,
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

func (c *Character) Update(space *physics.Space) error {
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
	c.op.GeoM.Reset()

	pos := c.Position()
	minX, minY := pos.Min.X, pos.Min.Y
	width := pos.Dx()

	fDirection := c.FaceDirection()
	if fDirection == physics.FaceDirectionRight {
		c.op.GeoM.Scale(-1, 1)
		c.op.GeoM.Translate(float64(width), 0)
	}

	// Apply character position
	c.op.GeoM.Translate(
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

func (c *Character) Draw(screen *ebiten.Image) {
	pos := c.Position()
	minX, minY := pos.Min.X, pos.Min.Y
	maxX, maxY := pos.Max.X, pos.Max.Y
	width := maxX - minX
	height := maxY - minY

	c.op.GeoM.Reset()

	fDirection := c.FaceDirection()
	if fDirection == physics.FaceDirectionRight {
		c.op.GeoM.Scale(-1, 1)
		c.op.GeoM.Translate(float64(width), 0)
	}

	// Apply character movement
	c.op.GeoM.Translate(
		float64(minX*config.Unit)/config.Unit,
		float64(minY*config.Unit)/config.Unit,
	)

	img := c.sprites[c.state.State()]
	characterWidth := img.Bounds().Dx()
	frameCount := characterWidth / width
	i := (c.count / frameRate) % frameCount
	sx, sy := frameOX+i*width, frameOY

	screen.DrawImage(
		img.SubImage(
			image.Rect(sx, sy, sx+width, sy+height),
		).(*ebiten.Image),
		c.op,
	)
}

func (c *Character) OnTouch(other physics.Body) {}

func (c *Character) OnBlock(other physics.Body) {}

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

func (c *Character) SetTouchable(t physics.Touchable) {
	c.PhysicsBody.Touchable = t
}

func (c *Character) Image() *ebiten.Image {
	img := ebiten.NewImage(c.Position().Dx(), c.Position().Dy())
	img.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	return img

	// pos := c.Position()
	// width := pos.Dx()
	// height := pos.Dy()
	//
	// // Fallback in case of missing sprite
	// if c.sprites == nil || c.sprites[c.state.State()] == nil {
	// 	if width == 0 || height == 0 {
	// 		// Avoid panic with zero size image
	// 		img := ebiten.NewImage(1, 1)
	// 		img.Fill(color.White)
	// 		return img
	// 	}
	// 	img := ebiten.NewImage(width, height)
	// 	img.Fill(color.RGBA{R: 255, A: 255})
	// 	return img
	// }
	//
	// img := c.sprites[c.state.State()]
	// characterWidth := img.Bounds().Dx()
	// if width == 0 {
	// 	// Avoid division by zero
	// 	img := ebiten.NewImage(1, 1)
	// 	img.Fill(color.White)
	// 	return img
	// }
	// frameCount := characterWidth / width
	// i := (c.count / frameRate) % frameCount
	// sx, sy := frameOX+i*width, frameOY
	//
	// return img.SubImage(
	// 	image.Rect(sx, sy, sx+width, sy+height),
	// ).(*ebiten.Image)
}

func (c *Character) ImageOptions() *ebiten.DrawImageOptions {
	return c.op
}
