package object

import (
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/config"
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
	BaseObject
	count   int
	sprites map[PlayerState]*ebiten.Image
	state   PlayerState
}

// TODO: Move this function to elsewhere
func getCenterOfScreenPosition(width, height int) (int, int) {
	x := config.ScreenWidth/2 - width/2
	y := config.ScreenHeight/2 - height/2
	return x, y
}

// TODO: Move this function to elsewhere
// TODO: Search for an alternative approach to reduce exponential
func checkRectIntersect(obj1, obj2 Object) bool {
	rects1 := obj1.CollisionPosition()
	rects2 := obj2.CollisionPosition()

	for _, r1 := range rects1 {
		for _, r2 := range rects2 {
			if r1.Overlaps(r2) {
				return true
			}
		}
	}

	return false
}

// TOOD: Move to a factory in Game module
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

	x, y := getCenterOfScreenPosition(frameWidth, frameHeight)

	playerElement := NewElement(x, y, frameWidth, frameHeight)
	collisionArea := NewElement(x+2, y+3, frameWidth-5, frameHeight-6)
	collisionList := []*CollisionArea{&CollisionArea{collisionArea}}

	return &Player{
		BaseObject: NewBaseObject(playerElement, collisionList),
		sprites:    sprites,
	}
}

// Object methods
func (p *Player) Position() (minX, minY, maxX, maxY int) {
	return p.BaseObject.Position()
}

func (p *Player) DrawCollisionBox(screen *ebiten.Image) {
	p.BaseObject.DrawCollisionBox(screen)
}

func (p *Player) CollisionPosition() []image.Rectangle {
	return p.BaseObject.CollisionPosition()
}

// updatePosition applies movement to player and collision areas
func (p *Player) updatePosition(velocity int, isXAxis bool) {
	if isXAxis {
		p.x16 += velocity
		for _, c := range p.collisionList {
			c.x16 += velocity
		}
	} else {
		p.y16 += velocity
		for _, c := range p.collisionList {
			c.y16 += velocity
		}
	}
}

// TODO: Move to the correct struct
func (p *Player) IsColliding(boundaries []Object) bool {
	for _, b := range boundaries {
		if checkRectIntersect(p, b.(Object)) {
			return true
		}
	}
	return false
}

// TODO: Move to the correct struct
func (p *Player) applyValidMovement(velocity int, isXAxis bool, boundaries []Object) {
	if velocity == 0 {
		return
	}

	p.updatePosition(velocity, isXAxis)

	isValid := !p.IsColliding(boundaries)
	if !isValid {
		p.updatePosition(-velocity, isXAxis)
	}
}

func (p *Player) Update(boundaries []Object) error {
	p.count++

	p.HandleInput()

	p.applyValidMovement(p.vx16, true, boundaries)
	p.applyValidMovement(p.vy16, false, boundaries)

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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.x16)/config.Unit, float64(p.y16)/config.Unit)

	// Animation frame rate
	img := p.sprites[p.state]
	playerWidth := img.Bounds().Dx()
	frameCount := playerWidth / p.width
	i := (p.count / frameRate) % frameCount
	sx, sy := frameOX+i*p.width, frameOY

	screen.DrawImage(img.SubImage(image.Rect(sx, sy, sx+p.width, sy+p.height)).(*ebiten.Image), op)
}

// TODO: Move it to elsewhere
func isSomeKeyPressed(keys ...ebiten.Key) bool {
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

func normalizeMoveOffset(move int, normalize bool) int {
	if normalize {
		return int(float64(move*config.Unit) / math.Sqrt2 / config.Unit)
	}
	return move
}

func (p *Player) HandleInput() {
	xMove, yMove := 0, 0
	if isSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		xMove = -playerXMove
	}
	if isSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		xMove = playerXMove
	}
	if isSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		yMove = -playerYMove
	}
	if isSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		yMove = playerYMove
	}

	isDiagonal := xMove != 0 && yMove != 0
	p.MoveY(normalizeMoveOffset(yMove, isDiagonal))
	p.MoveX(normalizeMoveOffset(xMove, isDiagonal))
}
