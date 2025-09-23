package physics

import (
	"image"
	"log"

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

	playerXMove = 3
	playerYMove = 3
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walk
)

type Player struct {
	PhysicsBody
	count   int
	sprites map[PlayerState]*ebiten.Image
	state   PlayerState
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

func (p *Player) ApplyValidMovement(velocity int, isXAxis bool, boundaries []Body) {
	p.PhysicsBody.ApplyValidMovement(velocity, isXAxis, boundaries)
}

func (p *Player) Update(boundaries []Body) error {
	p.count++

	p.HandleInput()

	p.ApplyValidMovement(p.vx16, true, boundaries)
	p.ApplyValidMovement(p.vy16, false, boundaries)

	isWalking := p.vx16 != 0 || p.vy16 != 0
	if isWalking {
		p.state = Walk
	} else {
		p.state = Idle
	}

	isDiagonal := p.vx16 != 0 && p.vy16 != 0
	xMove := normalizeMoveOffset(playerXMove, isDiagonal)
	yMove := normalizeMoveOffset(playerYMove, isDiagonal)

	// Reduce velocity
	if p.vx16 > 0 {
		p.vx16 -= xMove * config.Unit
	} else if p.vx16 < 0 {
		p.vx16 += xMove * config.Unit
	}

	if p.vy16 > 0 {
		p.vy16 -= yMove * config.Unit
	} else if p.vy16 < 0 {
		p.vy16 += yMove * config.Unit
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	body := p.Shape.(*Rect)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(body.x16)/config.Unit,
		float64(body.y16)/config.Unit,
	)

	// Animation frame rate
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

func (p *Player) HandleInput() {
	xMove, yMove := 0, 0
	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		xMove = -playerXMove
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		xMove = playerXMove
	}
	if input.IsSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		yMove = -playerYMove
	}
	if input.IsSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		yMove = playerYMove
	}

	isDiagonal := xMove != 0 && yMove != 0
	p.MoveY(normalizeMoveOffset(yMove, isDiagonal))
	p.MoveX(normalizeMoveOffset(xMove, isDiagonal))
}
