package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

// TODO: Move to the right place
func getSprites(assets map[string]string) (sprites.SpriteMap, error) {
	var s sprites.SpriteAssets
	for key, value := range assets {
		var state animation.SpriteState
		switch key {
		case "idle":
			state = actors.Idle
		case "walk":
			state = actors.Walk
		case "hurt":
			state = actors.Hurted
		default:
			continue
		}
		s = s.AddSprite(state, value)
	}
	result, err := sprites.LoadSprites(s)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateAnimatedCharacter(data actors.SpriteData) (*actors.Character, error) {
	assets, err := getSprites(data.Assets)
	if err != nil {
		return nil, err
	}

	character := actors.NewCharacter(assets, data.FrameRate)
	character.SetFaceDirection(data.FacingDirection)

	return character, nil
}

func SetPlayerBodies(player actors.ActorEntity, data actors.SpriteData) {
	bodyRect := physics.NewRect(data.BodyRect.Rect())
	collisionRect := physics.NewRect(data.CollisionRect.Rect())

	player.SetBody(bodyRect)
	player.SetCollisionArea(collisionRect)
	player.SetTouchable(player)
}

func SetPlayerStats(player body.Body, data actors.StatData) {
	// TODO: Create set stats method
	// player.SetStats(statData)
	player.SetMaxHealth(data.Health)
	player.SetSpeedAndMaxSpeed(data.Speed, data.MaxSpeed)
}

func SetMovementModel(
	player actors.ActorEntity,
	movementModel physics.MovementModelEnum,
	movementBlocker physics.PlayerMovementBlocker,
) error {
	model, err := physics.NewMovementModel(movementModel, movementBlocker)
	if err != nil {
		return err
	}
	player.SetMovementModel(model)
	return nil
}
