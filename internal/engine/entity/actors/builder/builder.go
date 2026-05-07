package builder

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
)

type collisionRectSetter interface {
	AddCollisionRect(state actors.ActorStateEnum, rect body.Collidable)
	RefreshCollisions()
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

// ApplySkills adds the provided skills to the character.
func ApplySkills(
	character actors.ActorEntity,
	skills []skill.Skill,
) error {
	for _, s := range skills {
		character.GetCharacter().AddSkill(s)
	}
	return nil
}
