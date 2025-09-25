package physics

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
)

type ActorEntity interface {
	SetBody(rect *Rect) ActorEntity
	SetCollisionArea(rect *Rect) ActorEntity
}

type CharacterState int

const (
	Idle CharacterState = iota
	Walk
)

type Character struct {
	PhysicsBody
	SpriteEntity
	count int
	state CharacterState
}

func NewCharacter(sprites SpriteMap) Character {
	spriteEntity := NewSpriteEntity(sprites)
	return Character{SpriteEntity: spriteEntity}
}

// Builder methods
func (c *Character) SetBody(rect *Rect) ActorEntity {
	c.PhysicsBody = *NewPhysicsBody(rect)
	return c
}

func (c *Character) SetCollisionArea(rect *Rect) ActorEntity {
	collisionArea := &CollisionArea{Shape: rect}
	c.PhysicsBody.AddCollision(collisionArea)
	return c
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
func (c *Character) IsColliding(boundaries []Body) bool {
	return c.PhysicsBody.IsColliding(boundaries)
}

func (c *Character) ApplyValidMovement(distance int, isXAxis bool, boundaries []Body) {
	c.PhysicsBody.ApplyValidMovement(distance, isXAxis, boundaries)
}

func increaseVelocity(velocity, acceleration int) int {
	// increaseVelocity applies acceleration to the velocity for a single axis.
	// v_new = v_old + a
	// Capping is handled in the Update loop to correctly manage the 2D vector's magnitude.
	velocity += acceleration
	return velocity
}

func reduceVelocity(velocity int) int {
	// reduceVelocity applies friction to the velocity for a single axis, slowing it down.
	// It brings the velocity to zero if it's smaller than the friction value to prevent jitter.
	friction := config.Unit / 4
	if velocity > friction {
		return velocity - friction
	}
	if velocity < -friction {
		return velocity + friction
	}
	return 0
}

// smoothDiagonalMovement converts raw input acceleration into a scaled and normalized vector.
// This ensures that the player's acceleration is consistent in all directions.
//
// Math:
//  1. The base acceleration from input (e.g., 2) is scaled up to a value that can
//     overcome friction.
//  2. If moving diagonally, the acceleration vector's magnitude would be `sqrt(ax² + ay²)`.
//     To ensure the magnitude is the same as for cardinal movement, we normalize it by
//     dividing each component by `sqrt(2)`.
func smoothDiagonalMovement(accX, accY int) (int, int) {
	// This factor determines the player's acceleration strength.
	// It should be large enough to overcome the friction in `reduceVelocity`.
	// Friction is `config.Unit / 4`. The base input acceleration is 2.
	// We'll use a factor of `config.Unit / 6` so that the final acceleration
	// (2 * config.Unit / 6 = config.Unit / 3) is greater than friction.
	accelerationFactor := float64(config.Unit / 6)

	fAccX := float64(accX) * accelerationFactor
	fAccY := float64(accY) * accelerationFactor

	isDiagonal := accX != 0 && accY != 0
	if isDiagonal {
		fAccX /= math.Sqrt2
		fAccY /= math.Sqrt2
	}

	return int(fAccX), int(fAccY)
}

// TODO: Should it be splitted from Character to Movable?
func (c *Character) OnMoveLeft() {
	c.MoveX(-playerXMove)
}

func (c *Character) OnMoveRight() {
	c.MoveX(playerXMove)
}

func (c *Character) OnMoveUp() {
	c.MoveY(-playerYMove)
}

func (c *Character) OnMoveDown() {
	c.MoveY(playerYMove)
}

func (c *Character) Update(boundaries []Body, handleMovement func()) error {
	c.count++

	// Get player input and set the raw acceleration for this frame.
	// p.accelerationX and p.accelerationY will be small integers (e.g., -2, 0, 2).
	if handleMovement != nil {
		handleMovement()
	}

	if c.accelerationX > 0 {
		c.isMirrored = true
	}
	if c.accelerationX < 0 {
		c.isMirrored = false
	}

	// Apply physics to player's position based on the velocity from previous frame.
	// This is a simple Euler integration step: position += velocity * deltaTime (where deltaTime=1 frame).
	c.ApplyValidMovement(c.vx16, true, boundaries)
	c.ApplyValidMovement(c.vy16, false, boundaries)

	// Convert the raw input acceleration into a scaled and normalized vector.
	scaledAccX, scaledAccY := smoothDiagonalMovement(c.accelerationX, c.accelerationY)

	c.vx16 = increaseVelocity(c.vx16, scaledAccX)
	c.vy16 = increaseVelocity(c.vy16, scaledAccY)

	// Cap the magnitude of the velocity vector to enforce a maximum speed.
	// This is crucial for preventing faster movement on diagonals.
	// We need to check if the velocity magnitude `sqrt(vx² + vy²)` exceeds `speedMax16²`.
	// To avoid a costly square root, we can compare the squared values:
	speedMax16 := 3 * config.Unit
	// Use int64 for squared values to prevent potential overflow.
	velSq := int64(c.vx16)*int64(c.vx16) + int64(c.vy16)*int64(c.vy16)
	maxSq := int64(speedMax16) * int64(speedMax16)

	if velSq > maxSq {
		// If the speed is too high, we need to scale the velocity vector down.
		// The scaling factor is `scale = speedMax16 / current_speed`.
		// `current_speed` is `sqrt(velSq)`.
		// So, `scale = speedMax16 / sqrt(velSq)`.
		scale := float64(speedMax16) / math.Sqrt(float64(velSq))
		c.vx16 = int(float64(c.vx16) * scale)
		c.vy16 = int(float64(c.vy16) * scale)
	}

	// 5. Update Animation State
	// TODO: Improve Character state
	isWalking := c.vx16 != 0 || c.vy16 != 0
	if isWalking {
		c.state = Walk
	} else {
		c.state = Idle
	}

	// Reset frame-specific acceleration.
	// It will be recalculated on the next frame from input.
	c.accelerationX, c.accelerationY = 0, 0

	// Apply friction to slow the player down when there is no input.
	c.vx16 = reduceVelocity(c.vx16)
	c.vy16 = reduceVelocity(c.vy16)

	return nil
}

func (c *Character) Draw(screen *ebiten.Image) {
	body := c.Shape.(*Rect)
	op := &ebiten.DrawImageOptions{}

	if c.isMirrored {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(body.width), 0)
	}

	// Apply player movement
	op.GeoM.Translate(
		float64(body.x16)/config.Unit,
		float64(body.y16)/config.Unit,
	)

	img := c.sprites[c.state]
	playerWidth := img.Bounds().Dx()
	frameCount := playerWidth / body.width
	i := (c.count / frameRate) % frameCount
	sx, sy := frameOX+i*body.width, frameOY

	screen.DrawImage(
		img.SubImage(
			image.Rect(sx, sy, sx+body.width, sy+body.height),
		).(*ebiten.Image),
		op,
	)
}
