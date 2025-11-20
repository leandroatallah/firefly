package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

// TODO: Move to the right place
func getSprites(assets map[string]actors.AssetData) (sprites.SpriteMap, error) {
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
		s = s.AddSprite(state, value.Path)
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

	c := actors.NewCharacter(assets)
	c.SetFaceDirection(data.FacingDirection)
	c.SetFrameRate(data.FrameRate)

	return c, nil
}

type collisionRectSetter interface {
	AddCollisionRect(state actors.ActorStateEnum, rect *physics.Rect)
}

func SetPlayerBodies(player actors.ActorEntity, data actors.SpriteData) {
	bodyRect := physics.NewRect(data.BodyRect.Rect())
	// Set initial collision area for idle state
	collisionRect := physics.NewRect(data.Assets["idle"].CollisionRect.Rect())

	player.SetBody(bodyRect)
	player.SetCollisionArea(collisionRect)
	player.SetTouchable(player)

	setter, ok := player.(collisionRectSetter)
	if !ok {
		// Log error or just return
		return
	}

	for key, assetData := range data.Assets {
		var state actors.ActorStateEnum
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
		rect := physics.NewRect(assetData.CollisionRect.Rect())
		setter.AddCollisionRect(state, rect)
	}
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
