package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func createPlayer(ctx *app.AppContext, playerType gameentitytypes.PlayerType) (platformer.PlatformerActorEntity, error) {
	var f func(*app.AppContext) (platformer.PlatformerActorEntity, error)

	switch playerType {
	case gameentitytypes.ClimberPlayerType:
		f = gameplayer.NewClimberPlayer
	}

	p, err := f(ctx)
	if err != nil {
		return nil, err
	}

	p.GetCharacter().AddSkill(skill.NewJumpSkill())
	p.GetCharacter().AddSkill(skill.NewHorizontalMovementSkill())

	return p, nil
}
