package actors

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors/movement"
	"github.com/leandroatallah/firefly/internal/core/screenutil"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

const (
	frameOX   = 0
	frameOY   = 0
	frameRate = 8
)

type Player struct {
	Character
}

func NewPlayer() *Player {
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets SpriteAssets
	assets = assets.
		AddSprite(Idle, "assets/default-idle.png").
		AddSprite(Walk, "assets/default-walk.png").
		AddSprite(Hurted, "assets/default-idle.png")

	sprites, err := LoadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	character := NewCharacter(sprites)

	x, y := screenutil.GetCenterOfScreenPosition(frameWidth, frameHeight)
	bodyRect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := physics.NewRect(x+2, y+3, frameWidth-5, frameHeight-6)

	// TODO: Create a builder with director to automate this process
	player := &Player{Character: *character}
	player.SetBody(bodyRect)
	player.SetCollisionArea(collisionRect)
	player.PhysicsBody.SetTouchable(player)

	// TODO: Move it to the right place (builder)
	player.SetSpeedAndMaxSpeed(4, 4)
	player.SetMovementState(movement.Input, nil)

	return player
}

// Character Methods
func (p *Player) Update(space *physics.Space) error {
	return p.Character.Update(space)
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Character.Draw(screen)
}
