package actors

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type ActorEntity interface {
	SetBody(rect *physics.Rect) ActorEntity
	SetCollisionArea(rect *physics.Rect) ActorEntity
	SetState(state ActorState)
	SetMovementState(state MovementStateEnum, target physics.Body)
	SwitchMovementState(state MovementStateEnum)
	MovementState() MovementState
	Update(boundaries []physics.Body) error

	Position() (minX, minY, maxX, maxY int)
	Speed() int

	OnMoveUp(distance int)
	OnMoveDown(distance int)
	OnMoveLeft(distance int)
	OnMoveRight(distance int)
	OnMoveUpLeft(distance int)
	OnMoveUpRight(distance int)
	OnMoveDownLeft(distance int)
	OnMoveDownRight(distance int)
}

type Character struct {
	physics.PhysicsBody
	SpriteEntity
	count         int
	state         ActorState
	movementState MovementState
}

func NewCharacter(sprites SpriteMap) *Character {
	spriteEntity := NewSpriteEntity(sprites)
	c := &Character{SpriteEntity: spriteEntity}
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
	return c
}

func (c *Character) SetCollisionArea(rect *physics.Rect) ActorEntity {
	collisionArea := &physics.CollisionArea{Shape: rect}
	c.PhysicsBody.AddCollision(collisionArea)
	return c
}

func (c *Character) SetState(state ActorState) {
	c.state = state
	c.state.OnStart()
}

func (c *Character) SetMovementState(state MovementStateEnum, target physics.Body) {
	movementState, err := NewMovementState(c, state, target)
	if err != nil {
		log.Fatal(err)
	}

	c.movementState = movementState
}
func (c *Character) SwitchMovementState(state MovementStateEnum) {
	target := c.MovementState().Target()
	movementState, err := NewMovementState(c, state, target)
	if err != nil {
		log.Fatal(err)
	}
	c.movementState = movementState
}

func (c *Character) MovementState() MovementState {
	return c.movementState
}

// Movement methods
func (c *Character) OnMoveUp(distance int) {
	c.PhysicsBody.OnMoveUp(distance)
}
func (c *Character) OnMoveDown(distance int) {
	c.PhysicsBody.OnMoveDown(distance)
}
func (c *Character) OnMoveLeft(distance int) {
	c.PhysicsBody.OnMoveLeft(distance)
}
func (c *Character) OnMoveRight(distance int) {
	c.PhysicsBody.OnMoveRight(distance)
}
func (c *Character) OnMoveUpLeft(distance int) {
	c.PhysicsBody.OnMoveUp(distance)
	c.PhysicsBody.OnMoveLeft(distance)
}
func (c *Character) OnMoveUpRight(distance int) {
	c.PhysicsBody.OnMoveUp(distance)
	c.PhysicsBody.OnMoveRight(distance)
}
func (c *Character) OnMoveDownLeft(distance int) {
	c.PhysicsBody.OnMoveDown(distance)
	c.PhysicsBody.OnMoveLeft(distance)
}
func (c *Character) OnMoveDownRight(distance int) {
	c.PhysicsBody.OnMoveDown(distance)
	c.PhysicsBody.OnMoveRight(distance)
}

// Body methods
func (c *Character) Position() (minX, minY, maxX, maxY int) {
	return c.PhysicsBody.Position()
}
func (c *Character) Speed() int {
	return c.PhysicsBody.Speed()
}

func (c *Character) DrawCollisionBox(screen *ebiten.Image) {
	c.PhysicsBody.DrawCollisionBox(screen)
}

func (c *Character) CollisionPosition() []image.Rectangle {
	return c.PhysicsBody.CollisionPosition()
}
func (c *Character) IsColliding(boundaries []physics.Body) (isTouching, isBlocking bool) {
	return c.PhysicsBody.IsColliding(boundaries)
}

func (c *Character) ApplyValidMovement(distance int, isXAxis bool, boundaries []physics.Body) {
	c.PhysicsBody.ApplyValidMovement(distance, isXAxis, boundaries)
}

var bodyToActorState = map[physics.BodyState]ActorStateEnum{
	physics.Idle: Idle,
	physics.Walk: Walk,
}

func (c *Character) Update(boundaries []physics.Body) error {
	c.count++

	// Handle movement by Movement State - must happen BEFORE UpdateMovement
	c.movementState.Move()

	// Update physics and apply movement
	c.UpdateMovement(boundaries)

	// Check movement direction for sprite mirroring
	c.CheckMovementDirectionX()

	c.handleState()

	return nil
}

func (c *Character) handleState() {
	bodyState := c.CurrentBodyState()
	if state, exists := bodyToActorState[bodyState]; exists {
		if state == c.state.State() {
			return
		}

		s, err := NewActorState(c, state)
		if err != nil {
			log.Fatal(err)
		}
		c.SetState(s)
	}
}

func (c *Character) Draw(screen *ebiten.Image) {
	minX, minY, maxX, maxY := c.Position()
	width := maxX - minX
	height := maxY - minY

	op := &ebiten.DrawImageOptions{}

	fDirection := c.FaceDirection()
	if fDirection == physics.FaceDirectionRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(width), 0)
	}

	// Apply character movement
	op.GeoM.Translate(
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
		op,
	)
}

func (c *Character) HandleMovement() {}
