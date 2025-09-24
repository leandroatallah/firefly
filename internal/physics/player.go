package physics

import (
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/input"
	"github.com/leandroatallah/firefly/internal/screenutil"
)

const (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 32
	frameRate   = 8

	playerXMove = 2
	playerYMove = 2
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walk
)

type Player struct {
	PhysicsBody
	count      int
	sprites    map[PlayerState]*ebiten.Image
	state      PlayerState
	isMirrored bool
}

func NewPlayer() *Player {
	sprites := make(map[PlayerState]*ebiten.Image)
	var err error

	sprites[Idle], _, err = ebitenutil.NewImageFromFile("assets/default-idle.png")
	if err != nil {
		log.Fatal(err)
	}
	sprites[Walk], _, err = ebitenutil.NewImageFromFile("assets/default-walk.png")
	if err != nil {
		log.Fatal(err)
	}

	x, y := screenutil.GetCenterOfScreenPosition(frameWidth, frameHeight)

	playerElement := NewRect(x, y, frameWidth, frameHeight)
	collisionRect := NewRect(x+2, y+3, frameWidth-5, frameHeight-6)
	collisionArea := &CollisionArea{Shape: collisionRect}

	return &Player{
		PhysicsBody: *NewPhysicsBody(playerElement).AddCollision(collisionArea),
		sprites:     sprites,
	}
}

// Body methods
func (p *Player) Position() (minX, minY, maxX, maxY int) {
	return p.PhysicsBody.Position()
}

func (p *Player) DrawCollisionBox(screen *ebiten.Image) {
	p.PhysicsBody.DrawCollisionBox(screen)
}

func (p *Player) CollisionPosition() []image.Rectangle {
	return p.PhysicsBody.CollisionPosition()
}
func (p *Player) IsColliding(boundaries []Body) bool {
	return p.PhysicsBody.IsColliding(boundaries)
}

func (p *Player) ApplyValidMovement(distance int, isXAxis bool, boundaries []Body) {
	p.PhysicsBody.ApplyValidMovement(distance, isXAxis, boundaries)
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

func (p *Player) Update(boundaries []Body) error {
	p.count++

	// Get player input and set the raw acceleration for this frame.
	// p.accelerationX and p.accelerationY will be small integers (e.g., -2, 0, 2).
	p.HandleInput()

	if p.accelerationX > 0 {
		p.isMirrored = true
	}
	if p.accelerationX < 0 {
		p.isMirrored = false
	}

	// Apply physics to player's position based on the velocity from previous frame.
	// This is a simple Euler integration step: position += velocity * deltaTime (where deltaTime=1 frame).
	p.ApplyValidMovement(p.vx16, true, boundaries)
	p.ApplyValidMovement(p.vy16, false, boundaries)

	// Convert the raw input acceleration into a scaled and normalized vector.
	scaledAccX, scaledAccY := smoothDiagonalMovement(p.accelerationX, p.accelerationY)

	p.vx16 = increaseVelocity(p.vx16, scaledAccX)
	p.vy16 = increaseVelocity(p.vy16, scaledAccY)

	// Cap the magnitude of the velocity vector to enforce a maximum speed.
	// This is crucial for preventing faster movement on diagonals.
	// We need to check if the velocity magnitude `sqrt(vx² + vy²)` exceeds `speedMax16²`.
	// To avoid a costly square root, we can compare the squared values:
	speedMax16 := 3 * config.Unit
	// Use int64 for squared values to prevent potential overflow.
	velSq := int64(p.vx16)*int64(p.vx16) + int64(p.vy16)*int64(p.vy16)
	maxSq := int64(speedMax16) * int64(speedMax16)

	if velSq > maxSq {
		// If the speed is too high, we need to scale the velocity vector down.
		// The scaling factor is `scale = speedMax16 / current_speed`.
		// `current_speed` is `sqrt(velSq)`.
		// So, `scale = speedMax16 / sqrt(velSq)`.
		scale := float64(speedMax16) / math.Sqrt(float64(velSq))
		p.vx16 = int(float64(p.vx16) * scale)
		p.vy16 = int(float64(p.vy16) * scale)
	}

	// 5. Update Animation State
	// TODO: Improve Player state
	isWalking := p.vx16 != 0 || p.vy16 != 0
	if isWalking {
		p.state = Walk
	} else {
		p.state = Idle
	}

	// Reset frame-specific acceleration.
	// It will be recalculated on the next frame from input.
	p.accelerationX, p.accelerationY = 0, 0

	// Apply friction to slow the player down when there is no input.
	p.vx16 = reduceVelocity(p.vx16)
	p.vy16 = reduceVelocity(p.vy16)

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	body := p.Shape.(*Rect)
	op := &ebiten.DrawImageOptions{}

	if p.isMirrored {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(body.width), 0)
	}

	// Apply player movement
	op.GeoM.Translate(
		float64(body.x16)/config.Unit,
		float64(body.y16)/config.Unit,
	)

	img := p.sprites[p.state]
	playerWidth := img.Bounds().Dx()
	frameCount := playerWidth / body.width
	i := (p.count / frameRate) % frameCount
	sx, sy := frameOX+i*body.width, frameOY

	screen.DrawImage(
		img.SubImage(
			image.Rect(sx, sy, sx+body.width, sy+body.height),
		).(*ebiten.Image),
		op,
	)
}

func (p *Player) OnMoveLeft() {
	p.MoveX(-playerXMove)
}

func (p *Player) OnMoveRight() {
	p.MoveX(playerXMove)
}

func (p *Player) OnMoveUp() {
	p.MoveY(-playerYMove)
}

func (p *Player) OnMoveDown() {
	p.MoveY(playerYMove)
}

func (p *Player) HandleInput() {
	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		p.OnMoveLeft()
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		p.OnMoveRight()
	}
	if input.IsSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		p.OnMoveUp()
	}
	if input.IsSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		p.OnMoveDown()
	}
}
