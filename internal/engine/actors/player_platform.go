package actors

import (
	"github.com/leandroatallah/firefly/internal/engine/core/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type PlayerPlatform struct {
	Player

	coinCount int
}

func NewPlayerPlatform() (PlayerEntity, error) {
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets sprites.SpriteAssets
	assets = assets.
		AddSprite(Idle, "assets/images/leandro-idle.png").
		AddSprite(Walk, "assets/images/leandro-walk.png").
		AddSprite(Hurted, "assets/images/default-hurt.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		return nil, err
	}

	character := NewCharacter(sprites, 6)

	x := 0
	_, y := screenutil.GetCenterOfScreenPosition(frameWidth, frameHeight)
	bodyRect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := physics.NewRect(x+2, y+3, frameWidth-5, frameHeight-6)

	// TODO: Create a builder with director to automate this process
	player := &PlayerPlatform{
		Player: Player{
			Character: *character,
		},
	}
	player.SetBody(bodyRect)
	player.SetMaxHealth(5)
	player.SetCollisionArea(collisionRect)
	player.SetTouchable(player)
	player.SetSpeedAndMaxSpeed(1, 1)

	movementModel, err := physics.NewMovementModel(physics.Platform, player)
	if err != nil {
		return nil, err
	}

	player.SetMovementModel(movementModel)

	return player, nil
}

func (p *PlayerPlatform) AddCoinCount(amount int) {
	p.coinCount += amount
}
func (p *PlayerPlatform) CoinCount() int {
	return p.coinCount
}
