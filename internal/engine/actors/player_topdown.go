package actors

import (
	"github.com/leandroatallah/firefly/internal/engine/core/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type PlayerTopDown struct {
	Player
}

func NewPlayerTopDown(playerMovementBlocker physics.PlayerMovementBlocker) (PlayerEntity, error) {
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets sprites.SpriteAssets
	assets = assets.
		AddSprite(Idle, "assets/images/default-idle.png").
		AddSprite(Walk, "assets/images/default-walk.png").
		AddSprite(Hurted, "assets/images/default-hurt.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		return nil, err
	}

	character := NewCharacter(sprites)

	x, y := screenutil.GetCenterOfScreenPosition(frameWidth, frameHeight)
	bodyRect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := physics.NewRect(x+2, y+3, frameWidth-5, frameHeight-6)

	// TODO: Create a builder with director to automate this process
	player := &PlayerTopDown{
		Player: Player{
			Character: *character,
		},
	}
	player.SetBody(bodyRect)
	player.SetMaxHealth(5)
	player.SetCollisionArea(collisionRect)
	player.SetTouchable(player)
	player.SetSpeedAndMaxSpeed(4, 4)

	movementModel, err := physics.NewMovementModel(physics.TopDown, playerMovementBlocker)
	if err != nil {
		return nil, err
	}

	player.SetMovementModel(movementModel)

	return player, nil
}
