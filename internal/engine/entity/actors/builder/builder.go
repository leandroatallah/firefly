package builder

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/jsonutil"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
)

type collisionRectSetter interface {
	AddCollisionRect(state actors.ActorStateEnum, rect body.Collidable)
	RefreshCollisions()
}

// PreparePlatformer loads sprite and stat data, builds the state map, and initializes a PlatformerCharacter.
func PreparePlatformer(
	ctx *app.AppContext,
	jsonPath string,
) (*platformer.PlatformerCharacter, schemas.SpriteData, actors.StatData, map[string]animation.SpriteState, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData](ctx.Assets, jsonPath)
	if err != nil {
		return nil, schemas.SpriteData{}, actors.StatData{}, nil, err
	}

	stateMap, err := BuildStateMap(spriteData)
	if err != nil {
		return nil, schemas.SpriteData{}, actors.StatData{}, nil, err
	}

	rect := BodyRectFromSpriteData(spriteData)
	character := platformer.NewPlatformerCharacter(stateMap, spriteData, rect)
	if character == nil {
		return nil, schemas.SpriteData{}, actors.StatData{}, nil, fmt.Errorf("failed to create platformer character")
	}
	character.SetAppContext(ctx)

	return character, spriteData, statData, stateMap, nil
}

// ApplyPlatformerPhysics sets up the movement model and touchable interface for a platformer actor.
func ApplyPlatformerPhysics(actor actors.ActorEntity, blocker physicsmovement.PlayerMovementBlocker) error {
	model, err := physicsmovement.NewMovementModel(physicsmovement.Platform, blocker)
	if err != nil {
		return err
	}
	actor.SetMovementModel(model)
	actor.GetCharacter().SetTouchable(actor)
	return nil
}

func SetCharacterBodies(
	character actors.ActorEntity,
	data schemas.SpriteData,
	stateMap map[string]animation.SpriteState,
	idPrefix string,
) error {
	character.SetID(idPrefix)

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
		rect.SetOwner(character)
		setter.AddCollisionRect(actorState, rect)
		setter.RefreshCollisions()
	}

	bodyphysics.SetCollisionBodies(character, data, stateMap, idProvider, addCollisionRect)
	return nil
}

func SetCharacterStats(character actors.ActorEntity, data actors.StatData) error {
	character.SetMaxHealth(data.Health)
	if err := character.SetSpeed(data.Speed); err != nil {
		return err
	}
	if err := character.SetMaxSpeed(data.MaxSpeed); err != nil {
		return err
	}
	return nil
}

func BuildStateMap(data schemas.SpriteData) (map[string]animation.SpriteState, error) {
	stateMap := make(map[string]animation.SpriteState)
	for stateName := range data.Assets {
		enum, ok := actors.GetStateEnum(stateName)
		if !ok {
			return nil, fmt.Errorf("state '%s' not registered", stateName)
		}
		stateMap[stateName] = enum
	}
	return stateMap, nil
}

func BodyRectFromSpriteData(data schemas.SpriteData) *bodyphysics.Rect {
	return bodyphysics.NewRect(data.BodyRect.Rect())
}

func ConfigureCharacter(
	character actors.ActorEntity,
	spriteData schemas.SpriteData,
	statData actors.StatData,
	stateMap map[string]animation.SpriteState,
	idPrefix string,
) error {
	if err := SetCharacterBodies(character, spriteData, stateMap, idPrefix); err != nil {
		return err
	}
	if err := SetCharacterStats(character, statData); err != nil {
		return err
	}
	return nil
}
