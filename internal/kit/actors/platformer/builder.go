package platformer

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/jsonutil"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
)

// PreparePlatformer loads sprite and stat data, builds the state map, and initializes a PlatformerCharacter.
func PreparePlatformer(
	ctx *app.AppContext,
	jsonPath string,
) (*PlatformerCharacter, schemas.SpriteData, actors.StatData, map[string]animation.SpriteState, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData](ctx.Assets, jsonPath)
	if err != nil {
		return nil, schemas.SpriteData{}, actors.StatData{}, nil, err
	}

	stateMap, err := builder.BuildStateMap(spriteData)
	if err != nil {
		return nil, schemas.SpriteData{}, actors.StatData{}, nil, err
	}

	rect := builder.BodyRectFromSpriteData(spriteData)
	character, err := NewPlatformerCharacter(ctx.Assets, stateMap, spriteData, rect)
	if err != nil {
		return nil, schemas.SpriteData{}, actors.StatData{}, nil, fmt.Errorf("failed to create platformer character: %w", err)
	}
	character.SetAppContext(ctx)

	return character, spriteData, statData, stateMap, nil
}
