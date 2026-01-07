package builder

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func CreateAnimatedCharacter(data schemas.SpriteData, stateMap map[string]animation.SpriteState) (*gameentitytypes.PlatformerCharacter, error) {
	assets, err := sprites.GetSpritesFromAssets(data.Assets, stateMap)
	if err != nil {
		return nil, err
	}

	rect := bodyphysics.NewRect(data.BodyRect.Rect())
	p := gameentitytypes.NewPlatformerCharacter(assets, rect)
	p.SetFaceDirection(data.FacingDirection)
	p.SetFrameRate(data.FrameRate)

	return p, nil
}

type collisionRectSetter interface {
	AddCollisionRect(state actors.ActorStateEnum, rect body.Collidable)
	RefreshCollisions()
}

func SetCharacterBodies(
	character gameentitytypes.PlatformerActorEntity,
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

	bodyphysics.SetCollisionBodies(character, data, stateMap, idProvider, addCollisionRect)
	return nil
}

func SetCharacterStats(character gameentitytypes.PlatformerActorEntity, data actors.StatData) error {
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
