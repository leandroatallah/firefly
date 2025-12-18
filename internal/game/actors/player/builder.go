package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/schemas"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

func CreateAnimatedCharacter(data schemas.SpriteData) (*actors.Character, error) {
	stateMap := map[string]animation.SpriteState{
		"idle":   actors.Idle,
		"walk":   actors.Walking,
		"fall":   actors.Falling,
		"hurt":   actors.Hurted,
	}
	assets, err := sprites.GetSpritesFromAssets(data.Assets, stateMap)
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
func SetPlayerBodies(player actors.ActorEntity, data schemas.SpriteData) error {
	player.SetID("player")

	setter, ok := player.(collisionRectSetter)
	if !ok {
		return fmt.Errorf("player must implement collisionRectSetter")
	}

	stateMap := map[string]animation.SpriteState{
		"idle":   actors.Idle,
		"walk":   actors.Walking,
		"fall":   actors.Falling,
		"hurt":   actors.Hurted,
	}

	idProvider := func(assetKey string, index int) string {
		return fmt.Sprintf("PLAYER_COLLISION_RECT_%s_%d", assetKey, index)
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		actorState, ok := state.(actors.ActorStateEnum)
		if !ok {
			// This should not happen if the stateMap is correct
			return
		}
		setter.AddCollisionRect(actorState, rect)
	}

	physics.SetCollisionBodies(player, data, stateMap, idProvider, addCollisionRect)
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
