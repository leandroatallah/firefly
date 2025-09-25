package physics

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/systems/input"
	"github.com/leandroatallah/firefly/internal/core/screenutil"
)

const (
	frameOX   = 0
	frameOY   = 0
	frameRate = 8

	playerXMove = 2
	playerYMove = 2
)

type Player struct {
	Character
}

func NewPlayer() *Player {
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets spriteAssets
	assets = assets.
		addSprite(Idle, "assets/default-idle.png").
		addSprite(Walk, "assets/default-walk.png")

	sprites, err := loadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	character := NewCharacter(sprites)

	x, y := screenutil.GetCenterOfScreenPosition(frameWidth, frameHeight)
	bodyRect := NewRect(x, y, frameWidth, frameHeight)
	collisionRect := NewRect(x+2, y+3, frameWidth-5, frameHeight-6)

	player := &Player{Character: character}
	player.SetBody(bodyRect)
	player.SetCollisionArea(collisionRect)

	return player
}

// Character Methods
func (p *Player) SetBody(rect *Rect) ActorEntity {
	return p.Character.SetBody(rect)
}

func (p *Player) SetCollisionArea(rect *Rect) ActorEntity {
	return p.Character.SetCollisionArea(rect)
}

func (p *Player) Update(boundaries []Body) error {
	return p.Character.Update(boundaries, p.HandleMovement)
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Character.Draw(screen)
}

func (p *Player) HandleMovement() {
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
