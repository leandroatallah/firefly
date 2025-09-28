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
	SetMovementFunc(func())
	Update(boundaries []physics.Body) error
}

type Character struct {
	physics.PhysicsBody
	SpriteEntity
	count         int
	state         ActorState
	movementState MovementStateEnum
	movementFunc  func()
}

func NewCharacter(sprites SpriteMap) Character {
	spriteEntity := NewSpriteEntity(sprites)
	c := Character{SpriteEntity: spriteEntity}
	state, err := NewActorState(&c, Idle)
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

func (c *Character) SetMovementFunc(cb func()) {
	c.movementFunc = cb
}

// Body methods
func (c *Character) Position() (minX, minY, maxX, maxY int) {
	return c.PhysicsBody.Position()
}

func (c *Character) DrawCollisionBox(screen *ebiten.Image) {
	c.PhysicsBody.DrawCollisionBox(screen)
}

func (c *Character) CollisionPosition() []image.Rectangle {
	return c.PhysicsBody.CollisionPosition()
}
func (c *Character) IsColliding(boundaries []physics.Body) bool {
	return c.PhysicsBody.IsColliding(boundaries)
}

func (c *Character) ApplyValidMovement(distance int, isXAxis bool, boundaries []physics.Body) {
	c.PhysicsBody.ApplyValidMovement(distance, isXAxis, boundaries)
}

// TODO: Should it be splitted from Character to Movable?
func (c *Character) OnMoveLeft() {
	c.MoveX(-playerXMove)
}
func (c *Character) OnMoveUpLeft() {
	c.MoveX(-playerXMove)
	c.MoveY(-playerYMove)
}
func (c *Character) OnMoveDownLeft() {
	c.MoveX(-playerXMove)
	c.MoveY(playerYMove)
}

func (c *Character) OnMoveRight() {
	c.MoveX(playerXMove)
}
func (c *Character) OnMoveUpRight() {
	c.MoveX(playerXMove)
	c.MoveY(-playerYMove)
}
func (c *Character) OnMoveDownRight() {
	c.MoveX(playerXMove)
	c.MoveY(playerYMove)
}

func (c *Character) OnMoveUp() {
	c.MoveY(-playerYMove)
}

func (c *Character) OnMoveDown() {
	c.MoveY(playerYMove)
}

var bodyToActorState = map[physics.BodyState]ActorStateEnum{
	physics.Idle: Idle,
	physics.Walk: Walk,
}

func (c *Character) Update(boundaries []physics.Body) error {
	c.count++

	// Sub class movement handler
	if c.movementFunc != nil {
		c.movementFunc()
	}

	isLeft, isRight := c.CheckMovementDirectionX()
	if isLeft {
		c.SetIsMirrored(false)
	}
	if isRight {
		c.SetIsMirrored(true)
	}

	c.UpdateMovement(boundaries)

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

	if c.isMirrored {
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
