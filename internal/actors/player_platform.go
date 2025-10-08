package actors

import (
	"github.com/leandroatallah/firefly/internal/core/screenutil"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type PlayerPlatform struct {
	Player
}

func NewPlayerPlatform() (PlayerEntity, error) {
	const (
		frameWidth  = 16
		frameHeight = 16
	)

	var assets SpriteAssets
	assets = assets.
		AddSprite(Idle, "assets/default-idle.png").
		AddSprite(Walk, "assets/default-walk.png").
		AddSprite(Hurted, "assets/default-hurt.png")

	sprites, err := LoadSprites(assets)
	if err != nil {
		return nil, err
	}

	character := NewCharacter(sprites)

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
	player.SetSpeedAndMaxSpeed(4, 4)

	movementModel, err := physics.NewMovementModel(physics.Platform)
	if err != nil {
		return nil, err
	}

	player.SetMovementModel(movementModel)

	return player, nil
}
