package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/config"
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
			state = actors.Walking
		case "fall":
			state = actors.Falling
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

	rect := physics.NewRect(data.BodyRect.Rect())
	c := actors.NewCharacter(assets, rect)
	c.SetFaceDirection(data.FacingDirection)
	c.SetFrameRate(data.FrameRate)

	return c, nil
}

type collisionRectSetter interface {
	AddCollisionRect(state actors.ActorStateEnum, rect body.Collidable)
}

// SetPlayerBodies
func SetPlayerBodies(player actors.ActorEntity, data actors.SpriteData) error {
	player.SetID("player")
	cfg := config.Get()

	x16, y16 := player.GetPositionMin()
	x, y := x16/cfg.Unit, y16/cfg.Unit

	collisions := []body.Collidable{}
	for i, r := range data.Assets["idle"].CollisionRects {
		c := physics.NewCollidableBodyFromRect(physics.NewRect(r.Rect()))
		c.SetPosition(x+r.X, y+r.Y)
		c.SetID(fmt.Sprintf("%v_COLLISION_%d", player.ID(), i))
		collisions = append(collisions, c)
	}

	if len(collisions) > 0 {
		player.AddCollision(collisions...)
	}

	player.SetTouchable(player)

	setter, ok := player.(collisionRectSetter)
	if !ok {
		return fmt.Errorf("player must implement collisionRectSetter")
	}

	for key, assetData := range data.Assets {
		var state actors.ActorStateEnum
		switch key {
		case "idle":
			state = actors.Idle
		case "walk":
			state = actors.Walking
		case "fall":
			state = actors.Falling
		case "hurt":
			state = actors.Hurted
		default:
			continue
		}

		for i, r := range assetData.CollisionRects {
			rect := physics.NewCollidableBody(
				physics.NewBody(physics.NewRect(r.Rect())),
			)
			rect.SetPosition(r.X, r.Y)
			rect.SetID(fmt.Sprintf("PLAYER-COLLISION-RECT-%d", i))
			setter.AddCollisionRect(state, rect)
		}
	}

	return nil
}

func SetPlayerStats(player actors.ActorEntity, data actors.StatData) error {
	// TODO: Create set stats method
	// player.SetStats(statData)
	player.SetMaxHealth(data.Health)
	var err error
	err = player.SetSpeed(data.Speed)
	if err != nil {
		return err
	}
	err = player.SetMaxSpeed(data.MaxSpeed)
	if err != nil {
		return err
	}
	return nil
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
