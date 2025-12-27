package builder

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/physics"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

func CreateAnimatedCharacter(data schemas.SpriteData, stateMap map[string]animation.SpriteState) (*actors.Character, error) {
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
	RefreshCollisions()
}

func SetCharacterBodies(
	character actors.ActorEntity,
	data schemas.SpriteData,
	stateMap map[string]animation.SpriteState,
	idPrefix string,
) error {
	setter, ok := character.(collisionRectSetter)
	if !ok {
		return fmt.Errorf("character must implement collisionRectSetter")
	}

	idProvider := func(assetKey string, index int) string {
		return fmt.Sprintf("%s_COLLISION_RECT_%s_%d", idPrefix, assetKey, index)
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		actorState, ok := state.(actors.ActorStateEnum)
		if !ok {
			return
		}
		setter.AddCollisionRect(actorState, rect)
		setter.RefreshCollisions()
	}

	physics.SetCollisionBodies(character, data, stateMap, idProvider, addCollisionRect)
	return nil
}

func SetCharacterStats(character actors.ActorEntity, data actors.StatData) error {
	character.SetMaxHealth(data.Health)
	var err error
	err = character.SetSpeed(data.Speed)
	if err != nil {
		return err
	}
	err = character.SetMaxSpeed(data.MaxSpeed)
	if err != nil {
		return err
	}
	return nil
}
